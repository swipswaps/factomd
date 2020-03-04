// Keeps a map of messages without acks, a map of acks without messages and a list of Msg/Ack pairs,
// forwarding the Msg/Ack pairs to the VM in order.
//
// Also notifies the MissingMessage Request module about any missing messages.
// On minute change send the unack’d messages to the leader module if it is now leader for the VM
package msgorder

import (
	"context"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/modules/events"
	"github.com/FactomProject/factomd/modules/logging"
	"github.com/FactomProject/factomd/pubsub"
	"github.com/FactomProject/factomd/worker"
)

type LogData = logging.LogData

type Handler struct {
	Pub
	Sub
	*Events
	ctx     context.Context    // manage thread context
	cancel  context.CancelFunc // thread cancel
	holding *OrderedMessageList
	log     func(data LogData) bool //logger hook
}

func newLogger(nodeName string) *logging.ModuleLogger {
	log := logging.NewModuleLoggerLogger(
		logging.NewLayerLogger(
			logging.NewSequenceLogger(
				logging.NewFileLogger(".")),
			map[string]string{"thread": nodeName},
		), "msgorder.txt")

	log.AddNameField("logname", logging.Formatter("%s"), "unknown_log")
	log.AddPrintField("Msg", logging.Formatter("%s"), "MSG")
	return log
}

func New(nodeName string) *Handler {
	h := new(Handler)
	h.log = newLogger(nodeName).Log
	h.holding = NewOrderedMessageList()

	h.Events = &Events{
		DBHT: &events.DBHT{
			DBHeight: 0,
			Minute:   0,
		},
		Config: &events.LeaderConfig{
			NodeName: nodeName,
		}, // FIXME should use pubsub.Config
	}
	return h
}

type Pub struct {
	UnAck  pubsub.IPublisher
	Leader pubsub.IPublisher
}

// create and start all publishers
func (p *Pub) Init(nodeName string) {
	p.UnAck = pubsub.PubFactory.Threaded(100).Publish(
		pubsub.GetPath(nodeName, events.Path.UnAckMsgs),
	)
	go p.UnAck.Start()
}

type Sub struct {
	MsgInput      *pubsub.SubChannel
	MovedToHeight *pubsub.SubChannel
}

// Create all subscribers
func (s *Sub) Init() {
	s.MovedToHeight = pubsub.SubFactory.Channel(1000)
	s.MsgInput = pubsub.SubFactory.Channel(1000)
}

// start subscriptions
func (s *Sub) Start(nodeName string) {
	s.MovedToHeight.Subscribe(pubsub.GetPath(nodeName, events.Path.DBHT))
	s.MsgInput.Subscribe(pubsub.GetPath(nodeName, events.Path.BMV))
}

type Events struct {
	*events.DBHT                      // from move-to-ht
	Config       *events.LeaderConfig // FIXME: use pubsub.Config obj
}

func (h *Handler) Start(w *worker.Thread) {
	w.Spawn("MsgOrderThread", func(w *worker.Thread) {
		w.OnReady(func() {
			h.Sub.Start(h.Config.NodeName)
		})
		w.OnRun(h.Run)
		w.OnExit(func() {
			h.Pub.UnAck.Close()
			h.cancel()
		})
		h.Pub.Init(h.Config.NodeName)
		h.Sub.Init()
	})
}

// Detect minute changes on receiving new DBHT event
func (h *Handler) minuteChanged(evt *events.DBHT) bool {
	if !h.DBHT.MinuteChanged(evt) {
		return false
	}
	h.DBHT = evt // save the new db height event
	return true
}

func (h *Handler) Run() {
	h.ctx, h.cancel = context.WithCancel(context.Background())

runLoop:
	for {
		select {
		case v := <-h.MsgInput.Updates:
			h.HandleMsg(v.(interfaces.IMsg))
		case v := <-h.MovedToHeight.Updates:
			if !h.minuteChanged(v.(*events.DBHT)) {
				continue runLoop
			}
			h.sendUnAckedMessages()
		case <-h.ctx.Done():
			return
		}
	}
}

func (h *Handler) HandleMsg(m interfaces.IMsg) {
	pair, ok := h.holding.Add(m)
	if !ok {
		return
	}

	if pair.Complete() {
		h.log(LogData{"Msg": pair.Msg})
		h.log(LogData{"Msg": pair.Ack})
		// FIXME Send both to VM
	}
	// REVIEW: do we want to report to MMR module right away?
	// or perhaps just out-of-sequence Ack without Msg?

}

// send unmatched messages to the leader for processing
func (h *Handler) sendUnAckedMessages() {
	// FIXME: send to Leader
	for h, msg := range h.holding.MsgList {
		_ = h
		_ = msg
	}
}

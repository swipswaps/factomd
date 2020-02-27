package leader

import (
	"context"
	"time"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/log"
	"github.com/FactomProject/factomd/modules/event"
	"github.com/FactomProject/factomd/pubsub"
	"github.com/FactomProject/factomd/state"
	"github.com/FactomProject/factomd/worker"
)

type Handler struct {
	Pub
	Sub
	*Leader
	ctx    context.Context    // manage thread context
	cancel context.CancelFunc // thread cancel
}

type Pub struct {
	MsgOut pubsub.IPublisher
}

// isolate deps on state package - eventually functions will be relocated
var GetFedServerIndexHash = state.GetFedServerIndexHash

// create and start all publishers
func (p *Pub) Init(nodeName string) {
	// REVIEW: will need to spawn/stop leader thread
	// based on federated set membership
	p.MsgOut = pubsub.PubFactory.Threaded(100).Publish(
		pubsub.GetPath(nodeName, event.Path.LeaderMsgOut),
	)
	go p.MsgOut.Start()
}

type role = int

const (
	FederatedRole role = iota + 1
	AuditRole
	FollowerRole
)

var _ = AuditRole // REVIEW: if Audit responsibilities are different from normal Fed then use this

type Sub struct {
	role
	MsgInput       *pubsub.SubChannel
	MovedToHeight  *pubsub.SubChannel
	BalanceChanged *pubsub.SubChannel
	DBlockCreated  *pubsub.SubChannel
	LeaderConfig   *pubsub.SubChannel
	AuthoritySet   *pubsub.SubChannel
}

func (*Sub) mkChan() *pubsub.SubChannel {
	return pubsub.SubFactory.Channel(1000) // FIXME: should calibrate channel depths
}

// Create all subscribers
func (s *Sub) Init() {
	s.MovedToHeight = s.mkChan()
	s.MsgInput = s.mkChan()
	s.BalanceChanged = s.mkChan()
	s.DBlockCreated = s.mkChan()
	s.LeaderConfig = s.mkChan()
	s.AuthoritySet = s.mkChan()
}

// start subscriptions
func (s *Sub) Start(nodeName string) {
	s.LeaderConfig.Subscribe(pubsub.GetPath(nodeName, event.Path.LeaderConfig))
	s.AuthoritySet.Subscribe(pubsub.GetPath(nodeName, event.Path.AuthoritySet))
	{
		s.SetLeaderMode(nodeName) //  create initial subscriptions
	}
}

// start listening to subscriptions for leader duties
func (s *Sub) SetLeaderMode(nodeName string) {
	if s.role == FederatedRole {
		return
	}
	s.role = FederatedRole
	s.MsgInput.Subscribe(pubsub.GetPath(nodeName, event.Path.BMV))
	s.MovedToHeight.Subscribe(pubsub.GetPath(nodeName, event.Path.DBHT))
	s.DBlockCreated.Subscribe(pubsub.GetPath(nodeName, event.Path.Directory))
	s.BalanceChanged.Subscribe(pubsub.GetPath(nodeName, event.Path.Bank))
}

// stop subscribers that we do not need as a follower
func (s *Sub) SetFollowerMode() {
	if s.role == FollowerRole {
		return
	}
	s.role = FollowerRole
	s.MsgInput.Unsubscribe()
	s.MovedToHeight.Unsubscribe()
	s.BalanceChanged.Unsubscribe()
	s.DBlockCreated.Unsubscribe()
}

type Events struct {
	Config              *event.LeaderConfig //
	*event.DBHT                             // from move-to-ht
	*event.Balance                          // REVIEW: does this relate to a specific VM
	*event.Directory                        //
	*event.Ack                              // record of last sent ack by leader
	*event.AuthoritySet                     //
}

func (h *Handler) sendOut(msg interfaces.IMsg) {
	log.LogMessage(h.logfile, "sendout", msg)
	h.Pub.MsgOut.Write(msg)
}

func (h *Handler) Start(w *worker.Thread) {
	if !state.EnableLeaderThread {
		panic("LeaderThreadDisabled")
	}

	w.Spawn("LeaderThread", func(w *worker.Thread) {
		h.ctx, h.cancel = context.WithCancel(context.Background())
		w.OnReady(func() {
			h.Sub.Start(h.Config.NodeName)
		})
		w.OnRun(h.Run)
		w.OnExit(func() {
			h.Pub.MsgOut.Close()
			h.cancel()
		})
		h.Pub.Init(h.Config.NodeName)
		h.Sub.Init()
	})
}

func (h *Handler) processMin() (ok bool) {
	go func() {
		time.Sleep(time.Second * time.Duration(h.Config.BlocktimeInSeconds/10))
		h.eomTicker <- true
	}()

	for {
		select {
		case v := <-h.MsgInput.Updates:
			m := v.(interfaces.IMsg)
			if constants.NeedsAck(m.Type()) {
				log.LogMessage(h.logfile, "msgIn ", m)
				h.sendAck(m)
			}
		case <-h.eomTicker:
			log.LogPrintf(h.logfile, "Ticker:")
			return true
		case <-h.ctx.Done():
			return false
		}
	}
}

func (h *Handler) waitForNextMinute() (min int, ok bool) {
	for {
		select {
		case v := <-h.MovedToHeight.Updates:
			evt := v.(*event.DBHT)
			log.LogPrintf(h.logfile, "DBHT: %v", evt)

			if !h.DBHT.MinuteChanged(evt) {
				continue
			}

			h.DBHT = evt
			return h.DBHT.Minute, true
		case <-h.ctx.Done():
			return -1, false
		}
	}
}

// TODO: refactor to only get a single Directory event
func (h *Handler) WaitForDBlockCreated() (ok bool) {
	for { // wait on a new (unique) directory event
		select {
		case v := <-h.Sub.DBlockCreated.Updates:
			evt := v.(*event.Directory)
			if h.Directory != nil && evt.DBHeight == h.Directory.DBHeight {
				log.LogPrintf(h.logfile, "DUP Directory: %v", v)
				continue
			} else {
				log.LogPrintf(h.logfile, "Directory: %v", v)
			}
			h.Directory = v.(*event.Directory)
			return true
		case <-h.ctx.Done():
			return false
		}
	}
}

func (h *Handler) WaitForBalanceChanged() (ok bool) {
	select {
	case v := <-h.Sub.BalanceChanged.Updates:
		h.Balance = v.(*event.Balance)
		log.LogPrintf(h.logfile, "BalChange: %v", v)
		return true
	case <-h.ctx.Done():
		return false
	}
}

// get latest AuthoritySet event data
// and compare w/ leader config
func (h *Handler) currentAuthority() (isLeader bool, index int) {
	evt := h.Events.AuthoritySet

readLatestAuthSet:
	for {
		select {
		case v := <-h.Sub.AuthoritySet.Updates:
			{
				evt = v.(*event.AuthoritySet)
			}
		default:
			{
				h.Events.AuthoritySet = evt
				break readLatestAuthSet
			}
		}
	}

	return GetFedServerIndexHash(h.Events.AuthoritySet.FedServers, h.Config.IdentityChainID)
}

// wait to become leader (possibly forever for followers)
func (h *Handler) WaitForAuthority() (isLeader bool) {
	// REVIEW: do we need to check block ht?
	log.LogPrintf(h.logfile, "WaitForAuthority %v ", h.Events.AuthoritySet.LeaderHeight)

	defer func() {
		if isLeader {
			h.Sub.SetLeaderMode(h.Config.NodeName)
			log.LogPrintf(h.logfile, "GotAuthority %v ", h.Events.AuthoritySet.LeaderHeight)
		}
	}()

	for {
		select {
		case v := <-h.Sub.LeaderConfig.Updates:
			h.Config = v.(*event.LeaderConfig)
		case v := <-h.Sub.AuthoritySet.Updates:
			h.Events.AuthoritySet = v.(*event.AuthoritySet)
		case <-h.ctx.Done():
			return false
		}
		if isAuthority, index := h.currentAuthority(); isAuthority {
			h.VMIndex = index
			return true
		}
	}
}

func (h *Handler) waitForNewBlock() (ok bool) {
	if min, done := h.waitForNextMinute(); !done {
		return false
	} else {
		return min != 0
	}
}

func (h *Handler) Run() {
	// TODO: wait until after boot height
	// ignore these events during DB loading
	h.waitForNextMinute()
	h.ctx, h.cancel = context.WithCancel(context.Background())

blockLoop:
	for { //blockLoop
		ok := worker.RunSteps(
			h.WaitForAuthority,
			h.WaitForBalanceChanged,
			h.WaitForDBlockCreated,
		)
		if !ok {
			break blockLoop
		} else {
			h.sendDBSig()
		}
		log.LogPrintf(h.logfile, "MinLoopStart: %v", true)
	minLoop:
		for { // could be counted 1..9 to account for min
			ok := worker.RunSteps(
				h.processMin,
				h.sendEOM,
				h.waitForNewBlock,
			)
			if !ok {
				break minLoop
			}
		}
		log.LogPrintf(h.logfile, "MinLoopEnd: %v", true)
	}
}

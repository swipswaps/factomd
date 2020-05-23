package livefeed

import (
	"fmt"
	"github.com/FactomProject/factomd/modules/livefeed/eventconfig"
	"os"
	"time"

	"github.com/FactomProject/factomd/common/globals"
	"github.com/FactomProject/factomd/modules/events"
	"github.com/FactomProject/factomd/modules/livefeed/eventmessages/generated/eventmessages"
	"github.com/FactomProject/factomd/modules/livefeed/eventservices"
	"github.com/FactomProject/factomd/modules/pubsub"
	"github.com/FactomProject/factomd/modules/pubsub/subregistry"
	"github.com/FactomProject/factomd/util"
	"github.com/FactomProject/logrustash"
	log "github.com/sirupsen/logrus"
)

func (*liveFeedService) ConfigSender(_ ...interface{}) {
}

type LiveFeedService interface {
	Start(state StateEventServices, config *util.FactomdConfig, factomParams *globals.FactomParams)
	StartEventLogs(state StateEventServices, config *util.FactomdConfig, factomParams *globals.FactomParams)
	ConfigSender(_ ...interface{})
}

type liveFeedService struct {
	parentState StateEventServices
	//eventSender          eventservices.EventSender
	//factomEventPublisher pubsub.IPublisher
	subCommitChain   *pubsub.SubChannel
	subCommitEntry   *pubsub.SubChannel
	subRevealEntry   *pubsub.SubChannel
	subCommitDBState *pubsub.SubChannel
	subDBAnchored    *pubsub.SubChannel
	subBlkSeq        *pubsub.SubChannel
	subNodeMessage   *pubsub.SubChannel
	eventLog         *log.Logger
}

func NewLiveFeedService() LiveFeedService {
	liveFeedService := new(liveFeedService)
	//liveFeedService.factomEventPublisher = pubsub.PubFactory.Threaded(5000).Publish("/live-feed", pubsub.PubMultiWrap())
	pubsub.SubFactory.PrometheusCounter("factomd_livefeed_total_events_published", "Number of events published by the factomd backend")
	return liveFeedService
}

func (liveFeedService *liveFeedService) Start(serviceState StateEventServices, config *util.FactomdConfig, factomParams *globals.FactomParams) {
	liveFeedService.parentState = serviceState
	/*
		if liveFeedService.eventSender == nil {
			liveFeedService.eventSender = eventservices.NewEventSender(config, factomParams)
		}
	*/

	subRegistry := subregistry.New(serviceState.GetFactomNodeName())
	liveFeedService.subNodeMessage = subRegistry.NodeMessageChannel()
	liveFeedService.subCommitDBState = subRegistry.CommitDBStateChannel()
	liveFeedService.subDBAnchored = subRegistry.DBAnchoredChannel()
	liveFeedService.subCommitChain = subRegistry.CommitChainChannel()
	liveFeedService.subCommitEntry = subRegistry.CommitEntryChannel()
	liveFeedService.subRevealEntry = subRegistry.RevealEntryChannel()
	liveFeedService.subBlkSeq = subRegistry.BlkSeqChannel()
	go liveFeedService.processSubChannels()
	//go liveFeedService.factomEventPublisher.Start()
}

func (liveFeedService *liveFeedService) processSubChannels() {
	broadcastContent := eventconfig.BroadcastOnce // liveFeedService.eventSender.GetBroadcastContent()

	// Don't track different levels of persistence
	sendStateChangeEvents := false //liveFeedService.eventSender.IsSendStateChangeEvents()

	for {
		select {
		case v := <-liveFeedService.subBlkSeq.Updates:
			liveFeedService.Send(eventservices.MapDBHT(v.(*events.DBHT), liveFeedService.GetStreamSource()))
		case v := <-liveFeedService.subCommitChain.Updates:
			commitChainEvent := v.(*events.CommitChain)
			if !sendStateChangeEvents || commitChainEvent.RequestState == events.RequestState_HOLDING {
				liveFeedService.Send(eventservices.MapCommitChain(commitChainEvent, liveFeedService.GetStreamSource()))
			} else {
				liveFeedService.Send(eventservices.MapCommitChainState(commitChainEvent, liveFeedService.GetStreamSource()))
			}
		case v := <-liveFeedService.subCommitEntry.Updates:
			commitEntryEvent := v.(*events.CommitEntry)
			if !sendStateChangeEvents || commitEntryEvent.RequestState == events.RequestState_HOLDING {
				liveFeedService.Send(eventservices.MapCommitEntry(commitEntryEvent, liveFeedService.GetStreamSource()))
			} else {
				liveFeedService.Send(eventservices.MapCommitEntryState(commitEntryEvent, liveFeedService.GetStreamSource()))
			}
		case v := <-liveFeedService.subRevealEntry.Updates:
			revealEntryEvent := v.(*events.RevealEntry)
			if !sendStateChangeEvents || revealEntryEvent.RequestState == events.RequestState_HOLDING {
				liveFeedService.Send(eventservices.MapRevealEntry(revealEntryEvent, liveFeedService.GetStreamSource(), broadcastContent))
			} else {
				liveFeedService.Send(eventservices.MapRevealEntryState(revealEntryEvent, liveFeedService.GetStreamSource()))
			}
		case v := <-liveFeedService.subCommitDBState.Updates:
			liveFeedService.Send(eventservices.MapCommitDBState(v.(*events.DBStateCommit), liveFeedService.GetStreamSource(), broadcastContent))
		case v := <-liveFeedService.subDBAnchored.Updates:
			liveFeedService.Send(eventservices.MapCommitDBAnchored(v.(*events.DBAnchored), liveFeedService.GetStreamSource()))
		case v := <-liveFeedService.subNodeMessage.Updates:
			nodeMessageEvent := v.(*events.NodeMessage)
			liveFeedService.Send(eventservices.MapNodeMessage(nodeMessageEvent, liveFeedService.GetStreamSource()))
			if nodeMessageEvent.MessageCode == events.NodeMessageCode_SHUTDOWN {
				break
			}
		}
	}
}

func (liveFeedService *liveFeedService) GetStreamSource() eventmessages.EventSource {
	if liveFeedService.parentState == nil {
		return -1
	}

	if liveFeedService.parentState.GetRunLeader() {
		return eventmessages.EventSource_LIVE
	} else {
		return eventmessages.EventSource_REPLAY_BOOT
	}
}

func (liveFeedService *liveFeedService) Send(factomEvent *eventmessages.FactomEvent) {
	/*
		defer func() {
			if r := recover(); r != nil {
				log.WithField("FactomEvent", factomEvent).Error("A panic was caught while sending FactomEvent")
			}
		}()
	*/

	if factomEvent == nil {
		return
	}

	factomEvent.IdentityChainID = liveFeedService.parentState.GetIdentityChainID().Bytes()
	factomEvent.FactomNodeName = liveFeedService.parentState.GetFactomNodeName()
	// REVIEW: this is a global setup
	//liveFeedService.factomEventPublisher.Write(factomEvent)
	liveFeedService.handleEvent(factomEvent)
}

// Forward only live feed events to Logstash
func (liveFeedService *liveFeedService) StartEventLogs(serviceState StateEventServices, config *util.FactomdConfig, params *globals.FactomParams) {
	// KLUDGE: rather than expose via command line options - set an ENV var
	if _, enabled := os.LookupEnv("EVENTLOG"); !enabled {
		return
	}

	if params.LogstashURL == "" {
		panic("must set live feed url")
	}

	liveFeedService.eventLog = serviceState.GetLogger().(*log.Logger)
	hookLogstashLogger(liveFeedService.eventLog, params.LogstashURL)
	liveFeedService.Start(serviceState, config, params)
}

func extIDtoString(extIDS [][]byte) []string {
	out := make([]string, 0)
	for _, xid := range extIDS {
		out = append(out, fmtHex(xid))
	}
	return out
}

func fmtHex(d interface{}) string {
	return fmt.Sprintf("%x", d)
}

// Filter and dispatch Events to logger
func (liveFeedService *liveFeedService) handleEvent(evt *eventmessages.FactomEvent) {
	// TODO: handle stateChange events
	if e := evt.GetEntryCommit(); e != nil {
		// skip rejections - they seem to always happen right before acceptance
		if e.EntityState == eventmessages.EntityState_REJECTED {
			return
		}

		liveFeedService.eventLog.WithFields(
			log.Fields{
				"EventType":            "EntryCommit",
				"EntityState":          e.EntityState,
				"EntryHash":            fmtHex(e.EntryHash),
				"Timestamp":            e.Timestamp,
				"Credits":              e.Credits,
				"EntryCreditPublicKey": fmtHex(e.EntryCreditPublicKey),
				"Signature":            fmtHex(e.Signature),
				"Version":              e.Version,
			},
		).Info(evt.EventSource)
		return
	}
	if e := evt.GetChainCommit(); e != nil {
		// skip rejections - they seem to always happen right before acceptance
		if e.EntityState == eventmessages.EntityState_REJECTED {
			return
		}

		liveFeedService.eventLog.WithFields(
			log.Fields{
				"EventType":            "ChainCommit",
				"EntityState":          e.EntityState,
				"EntryHash":            fmtHex(e.EntryHash),
				"Weld":                 fmtHex(e.Weld),
				"Timestamp":            e.Timestamp,
				"Credits":              e.Credits,
				"EntryCreditPublicKey": fmtHex(e.EntryCreditPublicKey),
				"Signature":            fmtHex(e.Signature),
				"Version":              e.Version,
			},
		).Info(evt.EventSource)
		return
	}
	if e := evt.GetEntryReveal(); e != nil {
		// skip rejections - they seem to always happen right before acceptance
		if e.EntityState == eventmessages.EntityState_REJECTED {
			return
		}

		liveFeedService.eventLog.WithFields(
			log.Fields{
				"EventType":   "EntryReveal",
				"EntityState": e.EntityState,
				"Timestamp":   e.Timestamp,
				"ExternalIds": extIDtoString(e.Entry.ExternalIDs),
				"ChainID":     fmtHex(e.Entry.ChainID),
				"Content":     e.Entry.Content,
				"Hash":        fmtHex(e.Entry.Hash),
				"Version.":    e.Entry.Version,
			},
		).Info(evt.EventSource)
		return
	}
	if e := evt.GetDirectoryBlockCommit(); e != nil {
		blk := e.DirectoryBlock
		hdr := blk.Header
		liveFeedService.eventLog.WithFields(
			log.Fields{
				"EventType":                    "DirectoryBlockCommit",
				"Header.BodyMerkleRoot":        fmtHex(hdr.BodyMerkleRoot),
				"Header.PreviousKeyMerkleRoot": fmtHex(hdr.PreviousKeyMerkleRoot),
				"Header.PreviousFullHash":      fmtHex(hdr.PreviousFullHash),
				"Header.Timestamp ":            hdr.Timestamp,
				"Header.BlockHeight ":          hdr.BlockHeight,
				"Header.BlockCount":            hdr.BlockCount,
				"Header.Version":               hdr.Version,
				"Header.NetworkID":             hdr.NetworkID,
				//"Entries": blk.Entries, // REVIEW: should we send entries?
				"Hash":          fmtHex(blk.Hash),
				"ChainID":       fmtHex(blk.ChainID),
				"KeyMerkleRoot": fmtHex(blk.KeyMerkleRoot),
			},
		).Info(evt.EventSource)
		return
	}
	if e := evt.GetNodeMessage(); e != nil {
		liveFeedService.eventLog.WithFields(
			log.Fields{
				"EventType":           "NodeMessage",
				"Level":               e.Level,
				"MessageCode":         e.MessageCode,
				"Message.MessageText": e.MessageText,
			},
		).Info(evt.EventSource)
		return
	}

	if pe := evt.GetProcessListEvent(); pe != nil {
		if e := pe.GetNewBlockEvent(); e != nil {
			liveFeedService.eventLog.WithFields(
				log.Fields{
					"EventType": "NewBlock",
					"Height":    e.NewBlockHeight,
				},
			).Info(evt.EventSource)
			return
		}
		if e := pe.GetNewMinuteEvent(); e != nil {
			liveFeedService.eventLog.WithFields(
				log.Fields{
					"EventType": "NewMinute",
					"Minute":    e.NewMinute,
				},
			).Info(evt.EventSource)
		}
		return
	}

	return // Just drop everything else

	/*
		data, err := json.Marshal(evt.Event)
		if err != nil {
			panic(err)
		}

		evtString := string(data)
		liveFeedService.eventLog.WithFields(
			log.Fields{"Event": evtString, "Prefix": evtString[2:10]},
		).Info(evt.EventSource)
	*/
}

func hookLogstashLogger(logger *log.Logger, logStashURL string) error {
	hook, err := logrustash.NewAsyncHook("tcp", logStashURL, "factomdLogs")
	if err != nil {
		fmt.Printf("Failed to connect to logstash %v", err)
		return err
	}

	hook.ReconnectBaseDelay = time.Second // Wait for one second before first reconnect.
	hook.ReconnectDelayMultiplier = 2
	hook.MaxReconnectRetries = 10

	logger.Hooks.Add(hook)
	return nil

}

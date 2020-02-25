package msgorder

import (
	"github.com/FactomProject/factomd/common"
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
)

type ackPair struct {
	msg interfaces.IMsg
	ack interfaces.IMsg
}

type OrderedMessageList struct {
	common.Name
	PairList []*ackPair
	AckList  map[[32]byte]interfaces.IMsg
	MsgList  map[[32]byte]interfaces.IMsg
}

func NewOrderedMessageList() *OrderedMessageList {
	l := new(OrderedMessageList)
	// TODO: we likely want this data structure to show up in the Global Object Hierarchy
	// so eventually will need to accept a parent NamedObject
	//l.NameInit(parent, "OrderedMessageList", reflect.TypeOf(l).String())
	l.AckList = make(map[[32]byte]interfaces.IMsg)
	l.MsgList = make(map[[32]byte]interfaces.IMsg)
	l.PairList = make([]*ackPair,0)
	return l
}

// match msg/ack pairs as they arrive
func (ml *OrderedMessageList) Add(msg interfaces.IMsg) (matchedPair *ackPair) {
	var h [32]byte

	// REVIEW: do we need to account for duplicates here?
	// or is that covered as part of BSV - basic message validation
	if msg.Type() == constants.ACK_MSG {
		h = msg.(*messages.Ack).MessageHash.Fixed()
		_, foundAck := ml.AckList[h]
		targetMsg, foundMsg := ml.MsgList[h]
		if !foundAck {
			ml.AckList[h] = msg
			if foundMsg {
				matchedPair = &ackPair{msg: targetMsg, ack: msg}
				ml.PairList = append(ml.PairList, matchedPair)
			}
		}
	} else if constants.NeedsAck(msg.Type()) {
		h = msg.GetMsgHash().Fixed()
		targetAck, foundAck := ml.AckList[h]
		_, foundMsg := ml.MsgList[h]
		if ! foundMsg {
			ml.MsgList[h] = msg
			if foundAck {
				matchedPair = &ackPair{msg: msg, ack: targetAck}
				ml.PairList = append(ml.PairList, matchedPair)
			}
		}
	}
	return matchedPair
}

/*
// get and remove the list of dependent message for a hash
func (ml *OrderedMessageList) Get(h [32]byte) []interfaces.IMsg {
	rval := ml.list[h]
	delete(ml.list, h)

	// delete all the individual inMessages from the list
	for _, msg := range rval {
		if msg == nil {
			continue
		} else {
			//ml.s.LogMessage("DependentHolding", fmt.Sprintf("delete[%x]", h[:6]), msg)
			//ml.metric(msg).Dec()
			delete(ml.dependents, msg.GetMsgHash().Fixed())
		}
	}
	return rval
}
 */

package msgorder

import (
	"github.com/FactomProject/factomd/common"
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
)

type ackPair struct {
	Msg interfaces.IMsg
	Ack interfaces.IMsg
}

func (ap *ackPair) Complete() bool {
	return ap.Ack != nil && ap.Msg != nil
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
	l.PairList = make([]*ackPair, 0)
	return l
}

func (ml *OrderedMessageList) addPair(hash [32]byte, pair *ackPair) {
	ml.PairList = append(ml.PairList, pair)
	// TODO add prometheus metric
	delete(ml.AckList, hash)
	delete(ml.MsgList, hash)
}

func (ml *OrderedMessageList) lookup(h [32]byte) (pair *ackPair, foundAck bool, foundMsg bool) {
	pair = new(ackPair)
	pair.Ack, foundAck = ml.AckList[h]
	pair.Msg, foundMsg = ml.MsgList[h]
	return pair, foundAck, foundMsg
}

// match Msg/Ack pairs as they arrive
func (ml *OrderedMessageList) Add(msg interfaces.IMsg) (matchedPair *ackPair, ok bool) {
	var h [32]byte
	var foundAck bool
	var foundMsg bool

	if msg.Type() == constants.ACK_MSG {
		h = msg.(*messages.Ack).MessageHash.Fixed()
		matchedPair, foundAck, foundMsg = ml.lookup(h)
		if !foundAck {
			matchedPair.Ack = msg
			if foundMsg {
				ml.addPair(h, matchedPair)
			} else {
				ml.AckList[h] = msg
			}
		}
		ok = true
	} else if constants.NeedsAck(msg.Type()) {
		h = msg.GetMsgHash().Fixed()
		matchedPair, foundAck, foundMsg = ml.lookup(h)
		if !foundMsg {
			matchedPair.Msg = msg
			if foundAck {
				ml.addPair(h, matchedPair)
			} else {
				ml.MsgList[h] = msg
			}
		}
		ok = true
	}

	return matchedPair, ok
}

package leader

import (
	"encoding/binary"

	"strings"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/modules/events"
	"github.com/FactomProject/factomd/state"
)

type Leader struct {
	*Events                             // event aggregate
	VMIndex   int                       // vm this leader is responsible for
	eomTicker chan interface{}          // fires on calculated EOM
	logfile   string                    // hardcoded log for now
	SendOut   func(msg interfaces.IMsg) // function to broadcast messages
}

// initialize the leader event handler
func New(s state.LeaderConfig) (h *Handler) {
	h = new(Handler)
	l := new(Leader)

	h.Leader = l
	l.SendOut = h.sendOut // inject publish method

	l.VMIndex = s.LeaderVMIndex
	l.logfile = strings.ToLower(s.FactomNodeName) + "_leader"
	l.eomTicker = make(chan interface{})

	l.Events = &Events{
		Config: &events.LeaderConfig{
			NodeName:           s.FactomNodeName,
			Salt:               s.Salt,
			IdentityChainID:    s.IdentityChainID,
			ServerPrivKey:      s.ServerPrivKey,
			BlocktimeInSeconds: s.DirectoryBlockInSeconds,
		},
		DBHT: &events.DBHT{ // moved to new height/min
			DBHeight: s.DBHeightAtBoot,
			Minute:   0,
		},
		Balance:   nil, // last perm balance computed
		Directory: nil, // last dblock created
		Ack:       nil, // last ack
		AuthoritySet: &events.AuthoritySet{
			LeaderHeight: s.LLeaderHeight,
			FedServers:   make([]interfaces.IServer, 0),
			AuditServers: make([]interfaces.IServer, 0),
		},
	}

	return h
}

func (l *Leader) Sign(b []byte) interfaces.IFullSignature {
	return l.Config.ServerPrivKey.Sign(b)
}

// Use method injected by handler to broadcast message
func (l *Leader) sendOut(msg interfaces.IMsg) {
	l.SendOut(msg)
}

// Returns a millisecond timestamp
func (l *Leader) getTimestamp() interfaces.Timestamp {
	return primitives.NewTimestampNow()
}

func (l *Leader) getSalt(ts interfaces.Timestamp) uint32 {
	var b [32]byte
	copy(b[:], l.Config.Salt.Bytes())
	binary.BigEndian.PutUint64(b[:], uint64(ts.GetTimeMilli()))
	c := primitives.Sha(b[:])
	return binary.BigEndian.Uint32(c.Bytes())
}

func (l *Leader) sendAck(m interfaces.IMsg) {
	// TODO: if message cannot be ack'd send to Dependent Holding
	ack := l.NewAck(m, l.BalanceHash).(*messages.Ack) // LeaderExecute
	l.sendOut(ack)
}

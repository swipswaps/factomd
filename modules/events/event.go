package events

import (
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"path"
)

type pubSubPaths struct {
	EOM               string
	DBHT              string
	Directory         string
	Bank              string
	LeaderConfig      string
	LeaderMsgIn       string
	LeaderMsgOut      string
	CommitChain       string
	CommitEntry       string
	RevealEntry       string
	CommitDBState     string
	DBAnchored        string
	NodeMessage       string
	AuthoritySet      string
	ConnectionMetrics string
	ProcessListInfo   string
	StateUpdate       string
	BMV               string
	UnAckMsgs         string
}

var Path = pubSubPaths{
	EOM:               "EOM",
	DBHT:              "DBHT",
	Directory:         "directory",
	Bank:              "bank",
	LeaderConfig:      "leader-config",
	LeaderMsgIn:       "leader-msg-in",
	LeaderMsgOut:      "leader-msg-out",
	ConnectionMetrics: path.Join("connection", "metrics"),
	ProcessListInfo:   "process-list",
	StateUpdate:       "state-update",
	CommitChain:       "commit-chain",
	CommitEntry:       "commit-entry",
	RevealEntry:       "reveal-entry",
	CommitDBState:     "commit-dbstate",
	NodeMessage:       "node-message",
	AuthoritySet:      "authority-set",
	BMV:               path.Join("bmv", "rest"),
	UnAckMsgs:         "unack-msg",
}

type Balance struct {
	DBHeight    uint32
	BalanceHash interfaces.IHash
}

type Directory struct {
	DBHeight             uint32
	VMIndex              int
	DirectoryBlockHeader interfaces.IDirectoryBlockHeader
	Timestamp            interfaces.Timestamp
}

type DBHT struct {
	DBHeight uint32
	Minute   int
}

// Detect minute change excluding minute 10 boundary
func (evt *DBHT) MinuteChanged(newEvt *DBHT) bool {

	if newEvt.Minute == 10 { // ignore min 10 - we want to trigger on move from 9->0 instead
		return false
	}

	if newEvt.Minute == evt.Minute && newEvt.DBHeight == evt.DBHeight {
		return false // no change
	} else {
		return true
	}
}

// event created when Ack is actually sent out
type Ack struct {
	Height      uint32
	MessageHash interfaces.IHash
}

type LeaderConfig struct {
	NodeName           string
	IdentityChainID    interfaces.IHash
	Salt               interfaces.IHash // only change on boot
	ServerPrivKey      *primitives.PrivateKey
	BlocktimeInSeconds int
}

type EOM struct {
	Timestamp interfaces.Timestamp
}

type AuthoritySet struct {
	LeaderHeight uint32
	FedServers   []interfaces.IServer
	AuditServers []interfaces.IServer
}

type ProcessListInfo struct {
	ProcessTime interfaces.Timestamp
	Dump        string
	PrintMap    string
}

type StateUpdate struct {
	NodeTime           interfaces.Timestamp
	LeaderHeight       uint32
	Summary            string
	IdentitiesDetails  string
	AuthoritiesDetails string
}

type CommitChain struct {
	RequestState RequestState
	DBHeight     uint32
	CommitChain  ICommitChain
}

type CommitEntry struct {
	RequestState RequestState
	DBHeight     uint32
	CommitEntry  ICommitEntry
}

type RevealEntry struct {
	RequestState RequestState
	DBHeight     uint32
	RevealEntry  IRevealEntry
	MsgTimestamp interfaces.Timestamp
}

type DBStateCommit struct {
	DBHeight uint32
	DBState  IDBState
}

type DBAnchored struct {
	DBHeight     uint32
	DirBlockInfo interfaces.IDirBlockInfo
}

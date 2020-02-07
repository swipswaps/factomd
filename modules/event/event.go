package event

import (
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"path"
)

type pubSubPaths struct {
	EOM               string
	Seq               string
	Directory         string
	Bank              string
	LeaderConfig      string
	LeaderMsgIn       string
	LeaderMsgOut      string
	AuthoritySet string
	ConnectionMetrics string
	ConnectionAdded   string
	ConnectionRemoved string
}

var Path = pubSubPaths{
	EOM:               "EOM",
	Seq:               "seq",
	Directory:         "directory",
	Bank:              "bank",
	LeaderConfig:      "leader-config",
	LeaderMsgIn:       "leader-msg-in",
	LeaderMsgOut:      "leader-msg-out",
	AuthoritySet: "authority-set",
	ConnectionMetrics: path.Join("connection", "metrics"),
	ConnectionAdded:   "connection-added",
	ConnectionRemoved: "connection-removed",
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

type ConnectionChanged struct {
	IP       string
	Status   string
	Duration string
	Send     string
	Received string
	IsOnline bool
	State    string
}

type ConnectionAdded struct {
	ConnectionChanged
}
type ConnectionRemoved struct {
	ConnectionChanged
}

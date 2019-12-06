package leader

import (
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
)

func (l *Leader) CreateDBSig() (interfaces.IMsg, interfaces.IMsg) {

	dbs := new(messages.DirectoryBlockSignature)
	dbs.DirectoryBlockHeader = l.Directory.DirectoryBlockHeader
	dbs.ServerIdentityChainID = l.Config.IdentityChainID
	dbs.DBHeight = l.DBHT.DBHeight
	dbs.Timestamp = l.GetTimestamp()
	dbs.SetVMHash(nil)
	dbs.SetVMIndex(l.VMIndex)
	dbs.SetLocal(true)
	dbs.Sign(l)
	err := dbs.Sign(l)
	if err != nil {
		panic(err)
	}
	ack := l.NewAck(dbs, l.Balance.BalanceHash).(*messages.Ack)

	//l.LogMessage("dbstateprocess", "CreateDBSig", dbs)
	return dbs, ack
}

func (l *Leader) SendDBSig() {
	l.Ack = nil // this allows the DBState for block 1 to be written, but then stalls at 2-:-0
	dbs, ack := l.CreateDBSig()
	l.SendOut(dbs)
	l.SendOut(ack)
}
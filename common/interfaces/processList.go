package interfaces

type IProcessList interface {
	//Clear()
	GetKeysNewEntries() (keys [][32]byte)
	GetNewEntry(key [32]byte) IEntry
	LenNewEntries() int
	Complete() bool
	VMIndexFor(hash []byte) int
	GetVMStatsForFedServer(index int) (vmIndex int, listHeight int, listLength int, nextNil int)
	SortFedServers()
	SortAuditServers()
	SortDBSigs()
	FedServerFor(minute int, hash []byte) IServer
	GetVirtualServers(minute int, identityChainID *HashS) (found bool, index int)
	GetFedServerIndexHash(identityChainID *HashS) (bool, int)
	GetAuditServerIndexHash(identityChainID *HashS) (bool, int)
	MakeMap()
	PrintMap() string
	AddFedServer(identityChainID *HashS) int
	AddAuditServer(identityChainID *HashS) int
	RemoveFedServerHash(identityChainID *HashS)
	RemoveAuditServerHash(identityChainID *HashS)
	String() string
	GetDBHeight() uint32
}

type IRequest interface {
	//Key() (thekey [32]byte)
}

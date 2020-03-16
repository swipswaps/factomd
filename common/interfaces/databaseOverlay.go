// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package interfaces

//A simplified DBOverlay to make sure we are not calling functions that could cause problems
type DBOverlaySimple interface {
	Close() error
	DoesKeyExist(bucket, key []byte) (bool, error)
	ExecuteMultiBatch() error
	FetchABlock(*HashS) (IAdminBlock, error)
	FetchABlockByHeight(blockHeight uint32) (IAdminBlock, error)
	FetchDBKeyMRByHeight(dBlockHeight uint32) (dBlockKeyMR *HashS, err error)
	FetchDBlock(*HashS) (IDirectoryBlock, error)
	FetchDBlockByHeight(uint32) (IDirectoryBlock, error)
	FetchDBlockHead() (IDirectoryBlock, error)
	FetchEBlock(*HashS) (IEntryBlock, error)
	FetchEBlockHead(chainID *HashS) (IEntryBlock, error)
	FetchECBlock(*HashS) (IEntryCreditBlock, error)
	FetchECBlockByHeight(blockHeight uint32) (IEntryCreditBlock, error)
	FetchECTransaction(hash *HashS) (IECBlockEntry, error)
	FetchEntry(*HashS) (IEBEntry, error)
	FetchFBlock(*HashS) (IFBlock, error)
	FetchFBlockByHeight(blockHeight uint32) (IFBlock, error)
	FetchFactoidTransaction(hash *HashS) (ITransaction, error)
	FetchHeadIndexByChainID(chainID *HashS) (*HashS, error)
	FetchIncludedIn(hash *HashS) (*HashS, error)
	FetchPaidFor(hash *HashS) (*HashS, error)
	FetchAllEBlocksByChain(*HashS) ([]IEntryBlock, error)
	InsertEntryMultiBatch(entry IEBEntry) error
	InsertEntry(entry IEBEntry) error
	ProcessABlockMultiBatch(block DatabaseBatchable) error
	ProcessDBlockMultiBatch(block DatabaseBlockWithEntries) error
	ProcessEBlockBatch(eblock DatabaseBlockWithEntries, checkForDuplicateEntries bool) error
	ProcessEBlockMultiBatch(eblock DatabaseBlockWithEntries, checkForDuplicateEntries bool) error
	ProcessEBlockMultiBatchWithoutHead(eblock DatabaseBlockWithEntries, checkForDuplicateEntries bool) error
	ProcessECBlockMultiBatch(IEntryCreditBlock, bool) (err error)
	ProcessFBlockMultiBatch(DatabaseBlockWithEntries) error
	FetchDirBlockInfoByKeyMR(hash *HashS) (IDirBlockInfo, error)
	SetExportData(path string)
	StartMultiBatch()
	Trim()
	FetchAllEntriesByChainID(chainID *HashS) ([]IEBEntry, error)
	SaveKeyValueStore(kvs BinaryMarshallable, key []byte) error
	FetchKeyValueStore(key []byte, dst BinaryMarshallable) (BinaryMarshallable, error)
	SaveDatabaseEntryHeight(height uint32) error
	FetchDatabaseEntryHeight() (uint32, error)
}

// Db defines a generic interface that is used to request and insert data into db
type DBOverlay interface {
	// We let Database method calls flow through.
	IDatabase

	FetchHeadIndexByChainID(chainID *HashS) (*HashS, error)
	SetExportData(path string)

	StartMultiBatch()
	PutInMultiBatch(records []Record)
	ExecuteMultiBatch() error
	GetEntryType(hash *HashS) (*HashS, error)

	//**********************************Entry**********************************//

	// InsertEntry inserts an entry
	InsertEntry(entry IEBEntry) (err error)
	InsertEntryMultiBatch(entry IEBEntry) error

	// FetchEntry gets an entry by hash from the database.
	FetchEntry(*HashS) (IEBEntry, error)

	FetchAllEntriesByChainID(chainID *HashS) ([]IEBEntry, error)

	FetchAllEntryIDsByChainID(chainID *HashS) ([]*HashS, error)

	FetchAllEntryIDs() ([]*HashS, error)

	//**********************************EBlock**********************************//

	// ProcessEBlockBatche inserts the EBlock and update all it's ebentries in DB
	ProcessEBlockBatch(eblock DatabaseBlockWithEntries, checkForDuplicateEntries bool) error
	ProcessEBlockBatchWithoutHead(eblock DatabaseBlockWithEntries, checkForDuplicateEntries bool) error
	ProcessEBlockMultiBatchWithoutHead(eblock DatabaseBlockWithEntries, checkForDuplicateEntries bool) error
	ProcessEBlockMultiBatch(eblock DatabaseBlockWithEntries, checkForDuplicateEntries bool) error

	FetchEBlock(*HashS) (IEntryBlock, error)

	// FetchEBlockByHash gets an entry by hash from the database.
	FetchEBlockByPrimary(*HashS) (IEntryBlock, error)

	// FetchEBlockByKeyMR gets an entry by hash from the database.
	FetchEBlockBySecondary(hash *HashS) (IEntryBlock, error)

	// FetchEBKeyMRByHash gets an entry by hash from the database.
	FetchEBKeyMRByHash(hash *HashS) (*HashS, error)

	// FetchAllEBlocksByChain gets all of the blocks by chain id
	FetchAllEBlocksByChain(*HashS) ([]IEntryBlock, error)

	SaveEBlockHead(block DatabaseBlockWithEntries, checkForDuplicateEntries bool) error

	FetchEBlockHead(chainID *HashS) (IEntryBlock, error)

	FetchAllEBlockChainIDs() ([]*HashS, error)

	//**********************************DBlock**********************************//

	// ProcessDBlockBatche inserts the EBlock and update all it's ebentries in DB
	ProcessDBlockBatch(block DatabaseBlockWithEntries) error
	ProcessDBlockBatchWithoutHead(block DatabaseBlockWithEntries) error
	ProcessDBlockMultiBatch(block DatabaseBlockWithEntries) error

	// FetchHeightRange looks up a range of blocks by the start and ending
	// heights.  Fetch is inclusive of the start height and exclusive of the
	// ending height. To fetch all hashes from the start height until no
	// more are present, use -1 as endHeight.
	FetchDBlockHeightRange(startHeight, endHeight int64) ([]*HashS, error)

	// FetchBlockHeightByKeyMR returns the block height for the given hash.  This is
	// part of the database.Db interface implementation.
	FetchDBlockHeightByKeyMR(*HashS) (int64, error)

	FetchDBlock(*HashS) (IDirectoryBlock, error)

	// FetchDBlock gets an entry by hash from the database.
	FetchDBlockByPrimary(*HashS) (IDirectoryBlock, error)

	// FetchDBlock gets an entry by hash from the database.
	FetchDBlockBySecondary(*HashS) (IDirectoryBlock, error)

	// FetchDBlockByHeight gets an directory block by height from the database.
	FetchDBlockByHeight(uint32) (IDirectoryBlock, error)

	FetchDBlockHead() (IDirectoryBlock, error)

	// FetchDBKeyMRByHeight gets a dBlock KeyMR from the database.
	FetchDBKeyMRByHeight(dBlockHeight uint32) (dBlockKeyMR *HashS, err error)

	// FetchDBKeyMRByHash gets a DBlock KeyMR by hash.
	FetchDBKeyMRByHash(hash *HashS) (dBlockHash *HashS, err error)

	// FetchAllFBInfo gets all of the fbInfo
	FetchAllDBlocks() ([]IDirectoryBlock, error)
	FetchAllDBlockKeys() ([]*HashS, error)

	SaveDirectoryBlockHead(DatabaseBlockWithEntries) error

	FetchDirectoryBlockHead() (IDirectoryBlock, error)

	//**********************************ECBlock**********************************//

	// ProcessECBlockBatch inserts the ECBlock and update all it's ecbentries in DB
	ProcessECBlockBatch(IEntryCreditBlock, bool) (err error)
	ProcessECBlockBatchWithoutHead(IEntryCreditBlock, bool) (err error)
	ProcessECBlockMultiBatch(IEntryCreditBlock, bool) (err error)

	FetchECBlock(*HashS) (IEntryCreditBlock, error)

	// FetchECBlockByHash gets an Entry Credit block by hash from the database.
	FetchECBlockByPrimary(*HashS) (IEntryCreditBlock, error)

	// FetchECBlockByKeyMR gets an Entry Credit block by hash from the database.
	FetchECBlockBySecondary(hash *HashS) (IEntryCreditBlock, error)
	FetchECBlockByHeight(blockHeight uint32) (IEntryCreditBlock, error)

	// FetchAllECBlocks gets all of the entry credit blocks
	FetchAllECBlocks() ([]IEntryCreditBlock, error)
	FetchAllECBlockKeys() ([]*HashS, error)

	SaveECBlockHead(IEntryCreditBlock, bool) error

	FetchECBlockHead() (IEntryCreditBlock, error)

	//**********************************ABlock**********************************//

	// ProcessABlockBatch inserts the AdminBlock
	ProcessABlockBatch(block DatabaseBatchable) error
	ProcessABlockBatchWithoutHead(block DatabaseBatchable) error
	ProcessABlockMultiBatch(block DatabaseBatchable) error

	FetchABlock(*HashS) (IAdminBlock, error)

	// FetchABlockByHash gets an admin block by hash from the database.
	FetchABlockByPrimary(hash *HashS) (IAdminBlock, error)

	// FetchABlockByKeyMR gets an admin block by keyMR from the database.
	FetchABlockBySecondary(hash *HashS) (IAdminBlock, error)
	FetchABlockByHeight(blockHeight uint32) (IAdminBlock, error)

	// FetchAllABlocks gets all of the admin blocks
	FetchAllABlocks() ([]IAdminBlock, error)
	FetchAllABlockKeys() ([]*HashS, error)

	SaveABlockHead(DatabaseBatchable) error

	FetchABlockHead() (IAdminBlock, error)

	//**********************************FBlock**********************************//

	// ProcessFBlockBatch inserts the Factoid
	ProcessFBlockBatch(DatabaseBlockWithEntries) error
	ProcessFBlockBatchWithoutHead(DatabaseBlockWithEntries) error
	ProcessFBlockMultiBatch(DatabaseBlockWithEntries) error

	FetchFBlock(*HashS) (IFBlock, error)

	// FetchFBlockByHash gets a factoid block by hash from the database.
	FetchFBlockByPrimary(*HashS) (IFBlock, error)
	FetchFBlockBySecondary(*HashS) (IFBlock, error)
	FetchFBlockByHeight(blockHeight uint32) (IFBlock, error)

	// FetchAllFBlocks gets all of the factoid blocks
	FetchAllFBlocks() ([]IFBlock, error)
	FetchAllFBlockKeys() ([]*HashS, error)

	SaveFactoidBlockHead(fblock DatabaseBlockWithEntries) error

	FetchFactoidBlockHead() (IFBlock, error)
	FetchFBlockHead() (IFBlock, error)

	//******************************DirBlockInfo********************************//

	// ProcessDirBlockInfoBatch inserts the dirblock info block
	ProcessDirBlockInfoBatch(block IDirBlockInfo) error

	// FetchDirBlockInfoByHash gets a dirblock info block by hash from the database.
	FetchDirBlockInfoByHash(hash *HashS) (IDirBlockInfo, error)

	// FetchDirBlockInfoByKeyMR gets a dirblock info block by keyMR from the database.
	FetchDirBlockInfoByKeyMR(hash *HashS) (IDirBlockInfo, error)

	// FetchAllConfirmedDirBlockInfos gets all of the confirmed dirblock info blocks
	FetchAllConfirmedDirBlockInfos() ([]IDirBlockInfo, error)

	// FetchAllUnconfirmedDirBlockInfos gets all of the unconfirmed dirblock info blocks
	FetchAllUnconfirmedDirBlockInfos() ([]IDirBlockInfo, error)

	// FetchAllDirBlockInfos gets all of the dirblock info blocks
	FetchAllDirBlockInfos() ([]IDirBlockInfo, error)

	SaveDirBlockInfo(block IDirBlockInfo) error

	//******************************IncludedIn**********************************//

	SaveIncludedIn(entry, block *HashS) error
	SaveIncludedInMultiFromBlock(block DatabaseBlockWithEntries, checkForDuplicateEntries bool) error
	SaveIncludedInMulti(entries []*HashS, block *HashS, checkForDuplicateEntries bool) error
	FetchIncludedIn(hash *HashS) (*HashS, error)

	ReparseAnchorChains() error
	SetBitcoinAnchorRecordPublicKeysFromHex([]string) error
	SetEthereumAnchorRecordPublicKeysFromHex([]string) error

	FetchPaidFor(hash *HashS) (*HashS, error)

	FetchFactoidTransaction(hash *HashS) (ITransaction, error)
	FetchECTransaction(hash *HashS) (IECBlockEntry, error)

	//******************************KeyValueStore**********************************//
	SaveKeyValueStore(kvs BinaryMarshallable, key []byte) error
	FetchKeyValueStore(key []byte, dst BinaryMarshallable) (BinaryMarshallable, error)
	SaveDatabaseEntryHeight(height uint32) error
	FetchDatabaseEntryHeight() (uint32, error)
}

type ISCDatabaseOverlay interface {
	DBOverlay

	FetchWalletEntryByName(addr []byte) (IWalletEntry, error)
	FetchWalletEntryByPublicKey(addr []byte) (IWalletEntry, error)
	FetchAllWalletEntriesByName() ([]IWalletEntry, error)
	FetchAllWalletEntriesByPublicKey() ([]IWalletEntry, error)
	FetchAllAddressNameKeys() ([][]byte, error)
	FetchAllAddressPublicKeys() ([][]byte, error)
	FetchTransaction(key []byte) (ITransaction, error)
	SaveTransaction(key []byte, tx ITransaction) error
	DeleteTransaction(key []byte) error
	FetchAllTransactionKeys() ([][]byte, error)
	FetchAllTransactions() ([]ITransaction, error)
	SaveRCDAddress(key []byte, we IWalletEntry) error
	SaveAddressByPublicKey(key []byte, we IWalletEntry) error
	SaveAddressByName(key []byte, we IWalletEntry) error
}

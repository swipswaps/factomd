// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

type IDirectoryBlock interface {
	Printable
	DatabaseBlockWithEntries

	GetHeader() IDirectoryBlockHeader
	SetHeader(IDirectoryBlockHeader)
	GetDBEntries() []IDBEntry
	GetEBlockDBEntries() []IDBEntry
	SetDBEntries([]IDBEntry) error
	AddEntry(chainID *HashS, keyMR *HashS) error
	BuildKeyMerkleRoot() (*HashS, error)
	BuildBodyMR() (*HashS, error)
	GetKeyMR() *HashS
	GetHash() *HashS
	GetFullHash() *HashS
	GetHeaderHash() (*HashS, error)

	GetTimestamp() Timestamp
	BodyKeyMR() *HashS
	GetEntryHashesForBranch() []*HashS

	SetEntryHash(hash, chainID *HashS, index int)
	SetABlockHash(aBlock IAdminBlock) error
	SetECBlockHash(ecBlock IEntryCreditBlock) error
	SetFBlockHash(fBlock IFBlock) error
	IsSameAs(IDirectoryBlock) bool
}

type IDirectoryBlockHeader interface {
	Printable
	BinaryMarshallable

	GetVersion() byte
	SetVersion(byte)
	GetPrevFullHash() *HashS
	SetPrevFullHash(*HashS)
	GetBodyMR() *HashS
	SetBodyMR(*HashS)
	GetPrevKeyMR() *HashS
	SetPrevKeyMR(*HashS)
	GetHeaderHash() (*HashS, error)
	GetDBHeight() uint32
	SetDBHeight(uint32)
	GetBlockCount() uint32
	SetBlockCount(uint32)
	GetNetworkID() uint32
	SetNetworkID(uint32)
	GetTimestamp() Timestamp
	SetTimestamp(Timestamp)
	IsSameAs(IDirectoryBlockHeader) bool
}

type IDBEntry interface {
	Printable
	BinaryMarshallable
	GetChainID() *HashS
	SetChainID(*HashS)
	GetKeyMR() *HashS
	SetKeyMR(*HashS)
	IsSameAs(IDBEntry) bool
}

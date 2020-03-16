//  Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

type IEntryCreditBlock interface {
	Printable
	DatabaseBatchable

	GetHeader() IECBlockHeader
	GetBody() IECBlockBody
	GetHash() *HashS
	HeaderHash() (*HashS, error)
	GetFullHash() (*HashS, error)
	GetEntryHashes() []*HashS
	GetEntrySigHashes() []*HashS
	GetEntries() []IECBlockEntry
	GetEntryByHash(hash *HashS) IECBlockEntry

	UpdateState(IState) error
	IsSameAs(IEntryCreditBlock) bool
	BuildHeader() error
}

type IECBlockHeader interface {
	BinaryMarshallable

	String() string
	GetBodyHash() *HashS
	SetBodyHash(*HashS)
	GetPrevHeaderHash() *HashS
	SetPrevHeaderHash(*HashS)
	GetPrevFullHash() *HashS
	SetPrevFullHash(*HashS)
	GetDBHeight() uint32
	SetDBHeight(uint32)
	GetECChainID() *HashS
	SetHeaderExpansionArea([]byte)
	GetHeaderExpansionArea() []byte
	GetObjectCount() uint64
	SetObjectCount(uint64)
	GetBodySize() uint64
	SetBodySize(uint64)
	IsSameAs(IECBlockHeader) bool
}

type IECBlockBody interface {
	String() string
	GetEntries() []IECBlockEntry
	SetEntries([]IECBlockEntry)
	AddEntry(IECBlockEntry)
	IsSameAs(IECBlockBody) bool
}

type IECBlockEntry interface {
	Printable
	ShortInterpretable

	ECID() byte
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
	UnmarshalBinaryData(data []byte) ([]byte, error)
	Hash() *HashS
	GetHash() *HashS
	GetEntryHash() *HashS
	GetSigHash() *HashS
	GetTimestamp() Timestamp
	IsSameAs(IECBlockEntry) bool
}

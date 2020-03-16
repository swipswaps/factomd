// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

type IEntryBlock interface {
	Printable
	DatabaseBatchable

	GetHeader() IEntryBlockHeader

	// AddEBEntry creates a new Entry Block Entry from the provided Factom Entry
	// and adds it to the Entry Block Body.
	AddEBEntry(entry IEBEntry) error
	// AddEndOfMinuteMarker adds the End of Minute to the Entry Block. The End of
	// Minut byte becomes the last byte in a 32 byte slice that is added to the
	// Entry Block Body as an Entry Block Entry.
	AddEndOfMinuteMarker(m byte) error
	// BuildHeader updates the Entry Block Header to include information about the
	// Entry Block Body. BuildHeader should be run after the Entry Block Body has
	// included all of its EntryEntries.
	BuildHeader() error
	// Hash returns the simple Sha256 hash of the serialized Entry Block. Hash is
	// used to provide the PrevFullHash to the next Entry Block in a Chain.
	Hash() (*HashS, error)
	// KeyMR returns the hash of the hash of the Entry Block Header concatenated
	// with the Merkle Root of the Entry Block Body. The Body Merkle Root is
	// calculated by the func (e *EBlockBody) MR() which is called by the func
	// (e *EBlock) BuildHeader().
	KeyMR() (*HashS, error)

	GetBody() IEBlockBody

	BodyKeyMR() *HashS
	GetEntryHashes() []*HashS
	GetEntrySigHashes() []*HashS
	GetHash() *HashS
	HeaderHash() (*HashS, error)
	IsSameAs(IEntryBlock) bool
}

type IEntryBlockHeader interface {
	Printable
	BinaryMarshallable

	GetBodyMR() *HashS
	GetChainID() *HashS
	GetPrevFullHash() *HashS
	GetPrevKeyMR() *HashS
	SetBodyMR(bodyMR *HashS)
	SetChainID(*HashS)
	SetPrevFullHash(*HashS)
	SetPrevKeyMR(*HashS)

	GetDBHeight() uint32
	GetEBSequence() uint32
	GetEntryCount() uint32
	SetDBHeight(uint32)
	SetEBSequence(uint32)
	SetEntryCount(uint32)
	IsSameAs(IEntryBlockHeader) bool
}

type IEBlockBody interface {
	Printable

	AddEBEntry(*HashS)
	AddEndOfMinuteMarker(m byte)
	GetEBEntries() []*HashS
	MR() *HashS
	IsSameAs(IEBlockBody) bool
}

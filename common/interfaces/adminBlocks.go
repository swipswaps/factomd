// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

import "bytes"

// Administrative Block
// This is a special block which accompanies this Directory Block.
// It contains the signatures and organizational data needed to validate previous and future Directory Blocks.
// This block is included in the DB body. It appears there with a pair of the Admin AdminChainID:SHA256 of the block.
// For more details, please go to:
// https://github.com/FactomProject/FactomDocs/blob/master/factomDataStructureDetails.md#administrative-block
type IAdminBlock interface {
	Printable
	DatabaseBatchable

	IsSameAs(IAdminBlock) bool
	BackReferenceHash() (*HashS, error)
	GetABEntries() []IABEntry
	GetDBHeight() uint32
	GetDBSignature() IABEntry
	GetHash() *HashS
	GetHeader() IABlockHeader
	GetKeyMR() (*HashS, error)
	LookupHash() (*HashS, error)
	RemoveFederatedServer(*HashS) error
	SetABEntries([]IABEntry)
	SetHeader(IABlockHeader)
	AddEntry(IABEntry) error
	FetchCoinbaseDescriptor() IABEntry

	InsertIdentityABEntries() error
	AddABEntry(e IABEntry) error
	AddAuditServer(*HashS) error
	AddDBSig(serverIdentity *HashS, sig IFullSignature) error
	AddFedServer(*HashS) error
	AddFederatedServerBitcoinAnchorKey(*HashS, byte, byte, [20]byte) error
	AddFederatedServerSigningKey(*HashS, [32]byte) error
	AddFirstABEntry(e IABEntry) error
	AddMatryoshkaHash(*HashS, *HashS) error
	AddServerFault(IABEntry) error
	AddCoinbaseDescriptor(outputs []ITransAddress) error
	AddEfficiency(chain *HashS, efficiency uint16) error
	AddCoinbaseAddress(chain *HashS, add IAddress) error
	AddCancelCoinbaseDescriptor(descriptorHeight, index uint32) error

	UpdateState(IState) error
}

// Admin Block Header
type IABlockHeader interface {
	Printable
	BinaryMarshallable

	IsSameAs(IABlockHeader) bool
	GetAdminChainID() *HashS
	GetDBHeight() uint32
	GetPrevBackRefHash() *HashS
	SetDBHeight(uint32)
	SetPrevBackRefHash(*HashS)

	GetHeaderExpansionArea() []byte
	SetHeaderExpansionArea([]byte)
	GetHeaderExpansionSize() uint64

	GetBodySize() uint32
	GetMessageCount() uint32
	SetBodySize(uint32)
	SetMessageCount(uint32)
}

type IABEntry interface {
	Printable
	BinaryMarshallable
	ShortInterpretable

	UpdateState(IState) error // When loading Admin Blocks,

	Type() byte
	Hash() *HashS
}

type IIdentityABEntrySort []IIdentityABEntry

func (p IIdentityABEntrySort) Len() int {
	return len(p)
}
func (p IIdentityABEntrySort) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p IIdentityABEntrySort) Less(i, j int) bool {
	// Sort by Type
	if p[i].Type() != p[j].Type() {
		return p[i].Type() < p[j].Type()
	}

	// Sort if identities are the same
	if p[i].SortedIdentity().IsSameAs(p[j].SortedIdentity()) {
		return bytes.Compare(p[i].Hash().Bytes(), p[j].Hash().Bytes()) < 0
	}

	// Sort by identity
	return bytes.Compare(p[i].SortedIdentity().Bytes(), p[j].SortedIdentity().Bytes()) < 0
}

type IIdentityABEntry interface {
	IABEntry
	// Identity to sort by
	SortedIdentity() *HashS
}

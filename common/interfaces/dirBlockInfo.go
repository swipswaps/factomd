// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

type IDirBlockInfo interface {
	Printable
	DatabaseBatchable
	GetDBHeight() uint32
	GetBTCConfirmed() bool
	GetDBMerkleRoot() *HashS
	GetBTCTxHash() *HashS
	GetBTCTxOffset() int32
	GetTimestamp() Timestamp
	GetBTCBlockHeight() int32
	GetBTCBlockHash() *HashS
	GetEthereumAnchorRecordEntryHash() *HashS
	GetEthereumConfirmed() bool
}

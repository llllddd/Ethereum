package core

import (
	"xchain-go/common"
	"xchain-go/core/basic"
)

// NewTxsEvent is posted when a batch of transactions enter the transaction pool.
type NewTxsEvent struct{ Txs []*basic.Transaction }

// PendingLogsEvent is posted pre mining and notifies of pending logs.
//type PendingLogsEvent struct {
//	Logs []*types.Log
//}

// NewMinedBlockEvent is posted when a block has been imported.
//type NewMinedBlockEvent struct{ Block *types.Block }

// RemovedLogsEvent is posted when a reorg happens
//type RemovedLogsEvent struct{ Logs []*types.Log }

type ChainEvent struct {
	Block *basic.Block
	Hash  common.Hash
	// Logs  []*types.Log
}

type ChainSideEvent struct {
	Block *basic.Block
}

type ChainHeadEvent struct{ Block *basic.Block }

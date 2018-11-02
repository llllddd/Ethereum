package core

import (
	"xchain-go/core/basic"
)

type Validator interface {
	// ValidateBody validates the given block's content.
	ValidateBody(block *basic.Block) error

	// ValidateState validates the given statedb and optionally the receipts and
	// gas used.
	// ValidateState(block, parent *basic.Block, state *state.StateDB, receipts basic.Receipts, usedGas uint64) error
}

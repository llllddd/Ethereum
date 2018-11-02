// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"xchain-go/common"
	"xchain-go/core/basic"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader interface {
	// Config retrieves the blockchain's chain configuration.
	// Config() *params.ChainConfig

	// // CurrentHeader retrieves the current header from the local chain.
	// CurrentHeader() *basic.Header

	// GetHeader retrieves a block header from the database by hash and number.
	// GetHeader(hash common.Hash, number uint64) *basic.Header

	// // GetHeaderByNumber retrieves a block header from the database by number.
	// GetHeaderByNumber(number uint64) *basic.Header

	// // GetHeaderByHash retrieves a block header from the database by its hash.
	// GetHeaderByHash(hash common.Hash) *basic.Header

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *basic.Block
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(header *basic.Header) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainReader, header *basic.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	// VerifyHeaders与VerifyHeader类似，但同时验证一批标头。 该方法返回退出通道以中止操作，并返回结果通道以检索异步验证（顺序是输入切片的顺序）。
	VerifyHeaders(chain ChainReader, headers []*basic.Header, seals []bool) (chan<- struct{}, <-chan error)

	// // VerifySeal checks whether the crypto seal on a header is valid according to
	// // the consensus rules of the given engine.
	// VerifySeal(chain ChainReader, header *basic.Header) error

	// // Prepare initializes the consensus fields of a block header according to the
	// // rules of a particular engine. The changes are executed inline.
	// Prepare(chain ChainReader, header *basic.Header) error

	// // Finalize runs any post-transaction state modifications (e.g. block rewards)
	// // and assembles the final block.
	// // Note: The block header and state database might be updated to reflect any
	// // consensus rules that happen at finalization (e.g. block rewards).
	// Finalize(chain ChainReader, header *basic.Header, state *state.StateDB, txs []*basic.Transaction,
	// 	uncles []*basic.Header, receipts []*basic.Receipt) (*basic.Block, error)

	// // Seal generates a new block for the given input block with the local miner's
	// // seal place on top.
	// Seal(chain ChainReader, block *basic.Block, stop <-chan struct{}) (*basic.Block, error)

	// // CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
	// // that a new block should have.
	// CalcDifficulty(chain ChainReader, time uint64, parent *basic.Header) *big.Int

	// // APIs returns the RPC APIs this consensus engine provides.
	// APIs(chain ChainReader) []rpc.API
}

// PoW is a consensus engine based on proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
}

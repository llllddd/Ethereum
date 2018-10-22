package consensus

import (
	"xchain-go/common"
	"xchain-go/core/basic"
	// "xchain-go/rpc"
)

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// 区块验证着，矿工
	Author(header *basic.Header) (common.Address, error)

	// // 校验Header是否符合共识算法，seal可以选择是否使用VerifySeal 方法
	// VerifyHeader(chain ChainReader, header *basic.Header, seal bool) error

	// // VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// // concurrently. The method returns a quit channel to abort the operations and
	// // a results channel to retrieve the async verifications (the order is that of
	// // the input slice).
	// // VerifyHeaders和VerifyHeader类似，但是他同时执行批量Headers
	// VerifyHeaders(chain ChainReader, headers []*basic.Header, seals []bool) (chan<- struct{}, <-chan error)

	// // VerifySeal checks whether the crypto seal on a header is valid according to
	// // the consensus rules of the given engine.
	// // 根据共识算法来判断Header中的crypto seal是否符合共识算法
	// VerifySeal(chain ChainReader, header *basic.Header) error

	// // Prepare initializes the consensus fields of a block header according to the
	// // rules of a particular engine. The changes are executed inline.
	// // 初始化Header数据
	// Prepare(chain ChainReader, header *basic.Header) error

	// // 将 prepare 和 CommitNewWork 内容打包成新块，同时里面还有包含出块奖励、选举、更新打块计数等功能
	// // Finalize(chain ChainReader, header *basic.Header, state *state.StateDB, txs []*basic.Transaction,
	// // 	uncles []*basic.Header, receipts []*basic.Receipt, dposContext *basic.DposContext) (*basic.Block, error)

	// // Seal generates a new block for the given input block with the local miner's
	// // seal place on top.
	// // 对新块进行签名
	// Seal(chain ChainReader, block *basic.Block, stop <-chan struct{}) (*basic.Block, error)

	// // APIs returns the RPC APIs this consensus engine provides.
	// APIs(chain ChainReader) []rpc.API
}

type ChainReader interface {
	// // Config retrieves the blockchain's chain configuration.
	// Config() *params.ChainConfig

	// // CurrentHeader retrieves the current header from the local chain.
	// // 从本地区块链中加载当前区块头
	// CurrentHeader() *types.Header

	// GetHeader retrieves a block header from the database by hash and number.
	//通过hash和number从数据库中检索Header
	GetHeader(hash common.Hash, number uint64) *basic.Header

	// // GetHeaderByNumber retrieves a block header from the database by number.
	// // 通过number在数据库中检索Header
	// GetHeaderByNumber(number uint64) *types.Header

	// // GetHeaderByHash retrieves a block header from the database by its hash.
	// // 通过hash在数据库中检索Header
	// GetHeaderByHash(hash common.Hash) *types.Header

	// // GetBlock retrieves a block from the database by hash and number.
	// // 通过hash和number在数据库中检索区块
	// GetBlock(hash common.Hash, number uint64) *types.Block
}

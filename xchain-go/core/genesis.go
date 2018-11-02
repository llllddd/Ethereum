package core

import (
	"fmt"
	"math/big"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/core/rawdb"
	"xchain-go/ethdb"

	log "github.com/inconshreveable/log15"
)

// genesis 文件主要用于定义初始块，创建数据库，并将初始块信息和状态存入数据库

type Genesis struct {
	// qiqi-todo:其他需要补充的字段
	// Config *params.ChainConfig
	Timestamp  uint64         //时间戳
	ExtraData  []byte         //块额外信息
	Validator  common.Address //区块验证者地址
	Number     uint64         //块号
	ParentHash common.Hash    //父hash
}

// SetupGenesisBlock
// 如果存储的区块链配置不兼容那么会被更新().
// 为了避免发生冲突,会返回一个错误,并且新的配置和原来的配置会返回.
// 返回的链配置永远不会为空。
//                          genesis == nil       genesis != nil
//                       +------------------------------------------
//     db has no genesis |  default           |  genesis
//     db has genesis    |  from DB           |  genesis (if compatible)
func SetupGenesisBlock(db ethdb.Database, genesis *Genesis) (common.Hash, error) {
	// 1.从db查看genesis是否为空
	// 2.若为空，则将非空的genesis存入数据库
	// 3.若不为空，则比较传入的genesis block 与db中的是否一致，一致，则返回获取到的hash，否则
	stored := rawdb.ReadCanonicalHash(db, 0)
	if stored == (common.Hash{}) {
		if genesis == nil {
			//如果db中不存在，且传入的genesis为空，则设置genesis为默认的genesisblock
			log.Info("默认的genesis写入db")
			genesis = DefaultGenesisBlock()
		}
		block, err := genesis.Commit(db)
		return block.Header().Hash(), err
	}
	//如果genesis也不为空，则比较获取到的genesis block
	//一致则返回stored，不一致则重置genesis配置
	//todo:重置genesis
	if genesis != nil {
		block := genesis.ToBlock()
		hash := block.Header().Hash()
		if hash != stored {
			return hash, fmt.Errorf("db中存储的genesis的block的hash与传入的genesis不一致")
		}
	}
	return stored, nil
}

// DefaultGenesisBlock 返回构造好的genesis内容
func DefaultGenesisBlock() *Genesis {
	return &Genesis{
		// Config:     params.DposChainConfig,
		// Nonce:      66,
		Timestamp: 1522052340,
		ExtraData: []byte("Genesis Block"),
		// GasLimit:   4712388,
		// Difficulty: big.NewInt(17179869184),
		// Alloc:      decodePrealloc(mainnetAllocData),
	}
}

// ToBlock 方法使用genesis的内容构造block。使用genesis的数据，使用基于内存的数据库，然后创建了一个block并返回
func (genesis *Genesis) ToBlock() *basic.Block {
	// qiqi-todo:状态数据的构造
	// db := ethdb.NewMemDatabase()

	// qiqi-todo:root\txhash\与状态有关，目前为空hash
	// root := statedb.IntermediateRoot(false)
	root := common.Hash{}
	txhash := common.Hash{}
	dposContextProto := basic.MockDposProto()

	// header
	header := &basic.Header{
		ParentHash:  genesis.ParentHash,                        //父common.Hash
		Timestamp:   new(big.Int).SetUint64(genesis.Timestamp), //区块产生的时间戳
		Number:      new(big.Int).SetUint64(genesis.Number),    //区块号
		Extradata:   genesis.ExtraData,                         //额外信息
		Validator:   genesis.Validator,                         //区块验证者地址
		Root:        root,
		TxHash:      txhash,
		DposContext: dposContextProto,
	}
	//qiqi-todo:删掉body

	// body := &basic.Body{}
	block := basic.NewBlock(header, nil)
	return block
}

// MustCommit 将创世块和状态写入db
func (genesis *Genesis) MustCommit(db ethdb.Database) *basic.Block {
	block, err := genesis.Commit(db)
	if err != nil {
		panic(err)
	}
	return block
}

// Commit 方法把给定的genesis的block写入数据库
func (genesis *Genesis) Commit(db ethdb.Database) (*basic.Block, error) {
	// 使用genesis的数据构造block
	block := genesis.ToBlock()

	// 将genesisblock的内容存入数据库
	// 1. 判断block的number，需要为0
	// 2. 写入区块
	// 3. 写入 "LastBlock" -> hash
	// 4. 写入 number -> hash
	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("不能将区块号不为0的genesis block存入db")
	}
	if err := rawdb.WriteBlock(db, block); err != nil {
		log.Error("genesis db写入 block写入失败")
		return nil, err
	}
	if err := rawdb.WriteHeadBlockHash(db, block.Header().Hash()); err != nil {
		log.Error("genesis db写入 'LastBlock' -> hash写入失败")
		return nil, err
	}
	if err := rawdb.WriteCanonicalHash(db, block.Header().Hash(), block.NumberU64()); err != nil {
		log.Error("genesis db写入 'number -> hash写入失败")
		return nil, err
	}
	return block, nil
}

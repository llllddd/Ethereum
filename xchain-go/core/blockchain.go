package core

import (
	"errors"
	"fmt"
	"sync"
	"xchain-go/common"
	"xchain-go/consensus"
	"xchain-go/core/basic"
	"xchain-go/core/rawdb"
	"xchain-go/ethdb"
	"xchain-go/event"
	"xchain-go/metrics"

	log "github.com/inconshreveable/log15"
)

var (
	blockInsertTimer = metrics.NewRegisteredTimer("chain/inserts", nil)
	ErrNoGenesis     = errors.New("Genesis not found in chain")
	dbFile           = "xchain.db"
)

// WriteStatus status of write
type WriteStatus byte

const (
	NonStatTy WriteStatus = iota
	CanonStatTy
	SideStatTy
)

//blockchain结构
//是一个block数组
//补充db字段,存储的时候只存block的hash
type BlockChain struct {
	// Latestheadhash []byte
	db           ethdb.Database
	genesisBlock *basic.Block
	currentBlock *basic.Block            // Current head of the block chain
	chainFeed    event.Feed              // chain中的事件广播
	scope        event.SubscriptionScope // chain中的订阅
	wg           sync.WaitGroup          // chain processing wait group for shutting down
	chainmu      sync.RWMutex            // blockchain insertion lock
	procmu       sync.RWMutex            // block processor lock
	engine       consensus.Engine
	validator    Validator // block and state validator interface

}

// NewBlockChain 功能初始化一条链
func NewBlockChain(db ethdb.Database) (*BlockChain, error) {
	// 1，NewHeaderChain()初始化区块头部链
	// 2，bc.genesisBlock = bc.GetBlockByNumber(0)  拿到第0个区块，也就是创世区块
	// 3，bc.loadLastState() 加载最新的状态数据，目前为加载currentBlock
	// 4，查找本地区块链上时候有硬分叉的区块，如果有调用bc.SetHead回到硬分叉之前的区块头
	// 5，todo：go bc.update() 定时处理future block

	//qiqi-todo:加载各种缓存、状态
	bc := &BlockChain{
		db: db,
	}
	//qiqi-todo:bc.SetValidator(NewBlockValidator(chainConfig, bc, engine))
	// bc.SetProcessor(NewStateProcessor(chainConfig, bc, engine))
	// var err error
	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock == nil {
		return nil, ErrNoGenesis
	}
	_, bc.currentBlock = bc.loadLastState()
	fmt.Println("bc.currentBlock.Number()", bc.currentBlock.Number())
	//qiqi-todo:定时处理future block
	// go bc.update()
	return bc, nil
}

// loadLastState 加载currentBlock
// todo:加载其他的一些状态
func (bc *BlockChain) loadLastState() (error, *basic.Block) {
	// Log := log(prefix)
	//1.从db中读取到最新的区块"LastBlock"的hash
	hash := rawdb.ReadHeadBlockHash(bc.db)
	if hash == (common.Hash{}) {
		err := errors.New("从db中读取'LastBlock'失败，没有获取到hash")
		log.Error("db没有读取到hash", "err", err)
		//qiqi-todo:重置链，重新加载
		return err, nil
	}
	//2.根据hash获取currentBlock
	currentBlock := bc.GetBlockByHash(hash)
	return nil, currentBlock
}
func (bc *BlockChain) WriteBlockWithoutState(block *basic.Block) error {
	bc.wg.Add(1)
	defer bc.wg.Done()
	if err := rawdb.WriteBlock(bc.db, block); err != nil {
		log.Error("block写入db失败", "err：", err)
		return err
	}
	log.Debug("block写入db成功")
	bc.insert(block)

	return nil

	//todo:存储难度值td
	// if err := bc.hc.WriteTd(block.Hash(), block.NumberU64(), td); err != nil {
	// 	return err
	// }

}

// InsertChain 同步完成后，插入块至本地blockchain
// todo:定义结构blocks，参数传入blocks数组
func (bc *BlockChain) InsertChain(newblocks basic.Blocks) error {
	events, err := bc.insertChain(newblocks)
	//qiqi-todo:添加完区块后，开始广播
	// n, events, logs, err := bc.insertChain(newblock)
	bc.PostChainEvents(events)
	return err
}

// insertChain 将block写入db
// todo:校验
func (bc *BlockChain) insertChain(newblocks basic.Blocks) ([]interface{}, error) {
	// qiqi-todo:对传入的blocks数组进行健全性检查，确保这个链是有序链接的,即检查这些block的number及hash与上一个block之间的关系
	// qiqi-todo:调用共识引擎，验证这些区块的headers
	// qiqi-todo:写入区块前，获取状态
	// batch := bc.db.NewBatch()

	//1. 前置校验
	//1）传入的blocks长度判断，不为空则继续
	//2）判断传入的blocks是否有序和链接，正常则继续
	//2. 开始导入block
	//1) todo:进程和加锁的功能
	//2) 并行header校验
	//3）迭代块，验证完成块的body之后，插入块。

	// 检查传入的链中的blocks是有序和链接的
	// 检查的内容：number是按顺序的、父hash是上一个块的blockhash
	if len(newblocks) == 0 {
		return nil, nil
	}

	for i := 1; i < len(newblocks); i++ {
		if newblocks[i].NumberU64() != newblocks[i-1].NumberU64()+1 || newblocks[i].ParentHash() != newblocks[i-1].Hash() {
			log.Error("非连续块插入,", "number:", newblocks[i].Number(), ",hash:", newblocks[i].Hash(),
				",parent:", newblocks[i].ParentHash(), ",prevnumber:", newblocks[i-1].Number(), ",prevhash:", newblocks[i-1].Hash())
			return nil, fmt.Errorf("非连续块插入: item %d is #%d [%x…], item %d is #%d [%x…] (parent [%x…])", i-1, newblocks[i-1].NumberU64(),
				newblocks[i-1].Hash().Bytes()[:4], i, newblocks[i].NumberU64(), newblocks[i].Hash().Bytes()[:4], newblocks[i].ParentHash().Bytes()[:4])
		}
	}
	log.Debug("传入的blocks正确")

	// 前置条件检查完，开始插入块的流程
	bc.wg.Add(1)
	defer bc.wg.Done()

	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	var (
		events    = make([]interface{}, 0, len(newblocks))
		lastCanon *basic.Block
	)

	// todo:校验传入的所有的block的header
	// headers := make([]*basic.Header, len(newblocks))
	// seals := make([]bool, len(newblocks)) //共识相关
	// abort, results := bc.engine.VerifyHeaders(bc, headers, seals)
	// defer close(abort) //获取到返回的结束信道值后，结束header的校验

	//迭代每个块，并在验证者验证完成后，插入
	for _, block := range newblocks {
		//todo:校验传入的所有的block的body
		// bstart := time.Now()
		// err := <-results
		// if err == nil {
		// 	err = bc.Validator().ValidateBody(block)
		// }
		// switch {
		// case err == ErrKnownBlock:
		// 	if bc.CurrentBlock().NumberU64() >= block.NumberU64() {
		// 		//回滚
		// 		//stats.ignored++
		// 		continue
		// 	}
		// }

		//若校验body没有错误，使用父区块创建一个新的statedb,若遇错则抛错
		//qiqi-todo：根据父区块创建statedb，验证状态，处理状态

		//上述操作执行完后，将新区块(todo：和状态)写入db
		err := bc.WriteBlockWithoutState(block)
		if err != nil {
			return events, err
		}

		log.Debug("Inserted new block", "number", block.Number(), "hash", block.Hash(),
			"txs", len(block.Transactions()))

		// blockInsertTimer.UpdateSince(bstart)
		events = append(events, ChainEvent{block, block.Hash()})
		lastCanon = block

	}
	if lastCanon != nil && bc.CurrentBlock().Hash() == lastCanon.Hash() {
		events = append(events, ChainHeadEvent{lastCanon})
	}
	return events, nil
}

// insert 方法,将block根据number将对应的hash存储至db，并将block设置为bc的currentBlock
func (bc *BlockChain) insert(block *basic.Block) {

	// 1.根据number存储对应的block的hash
	// 2.将LastBlock的hash添加到db
	if err := rawdb.WriteCanonicalHash(bc.db, block.Header().Hash(), block.NumberU64()); err != nil {
		log.Error("db存储blocknumber-hash失败", "err", err)
	}
	if err := rawdb.WriteHeadBlockHash(bc.db, block.Header().Hash()); err != nil {
		log.Error("db存储'LastBlock'-hash失败", "err", err)
	}
	bc.currentBlock = block
}

// GetBlockByHash 根据blockhash获取block
func (bc *BlockChain) GetBlockByHash(hash common.Hash) *basic.Block {
	// 1. 先根据hash从db中获取number
	// 2. 根据hash和number获取block
	// qiqi-todo:从缓存中读取
	number := rawdb.ReadHeaderNumber(bc.db, hash)
	if number == nil {
		log.Error("从db中根据hash获取number失败,db中没有该hash对应的number")
		return nil
	}
	block := bc.GetBlock(hash, *number)
	return block
}

// GetBlockByNumber 根据blocknumber获取block
func (bc *BlockChain) GetBlockByNumber(number uint64) *basic.Block {
	// 1. 先根据number从db中获取hash
	// 2. 根据hash和number获取block

	hash := rawdb.ReadCanonicalHash(bc.db, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return bc.GetBlock(hash, number)
}

// GetBlock 目前功能：根据hash和number从db中读取到block
func (bc *BlockChain) GetBlock(hash common.Hash, number uint64) *basic.Block {
	//qiqi-todo:从cache中读取。如果是从db中读取到的，则读完之后，将读到的block存到cache
	block := rawdb.ReadBlock(bc.db, hash, number)
	if block == nil {
		log.Error("读取到的block为nil")
		return nil
	}
	return block
}

// Validator returns the current validator.
func (bc *BlockChain) Validator() Validator {
	bc.procmu.RLock()
	defer bc.procmu.RUnlock()
	return bc.validator
}

//循环遍历链中的数据
func (bc *BlockChain) PrintBlockchain() {
	log.Debug("打印区块链信息")

	db := bc.db
	defer db.Close()
	log.Debug("打开数据库成功")

	//根据最新的块hash，获取最大的块号
	latesthash := rawdb.ReadHeadBlockHash(db)
	blocknumber := rawdb.ReadHeaderNumber(db, latesthash)
	n := int(*blocknumber)
	log.Info("最新的区块号", "number", n)

	//根据最大的块号，查询数据库中的从0-number的block信息
	for i := 0; i <= n; i++ {
		hash := rawdb.ReadCanonicalHash(db, *blocknumber)
		// Log.Infoln("number:", blocknumber, "对应的hash:", hash.String())
		block := rawdb.ReadBlock(db, hash, *blocknumber)
		// Log.Infoln(block.String())
		block.PrintBlockstruct()
		*blocknumber--
	}

}

// CurrentBlock 返回链中的currentBlock
func (bc *BlockChain) CurrentBlock() *basic.Block {
	return bc.currentBlock
}

// CurrentBlock 返回链中的currentBlock
func (bc *BlockChain) DB() ethdb.Database {
	return bc.db
}

// PostChainEvents迭代链插入生成的事件并将它们发布到事件源中。
// TODO：不应公开PostChainEvents。 链事件应该在WriteBlock中发布。
// qiqi-todo:传入log，定义log
func (bc *BlockChain) PostChainEvents(events []interface{}) {
	for _, event := range events {
		//判断event事件的类型
		switch ev := event.(type) {
		case ChainEvent:
			bc.chainFeed.Send(ev)
		}
	}
}

// SubscribeChainEvent 注册ChainEvent的订阅CHANGEBYLIDIAN。
func (bc *BlockChain) SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription {
	return bc.scope.Track(bc.chainFeed.Subscribe(ch))
}

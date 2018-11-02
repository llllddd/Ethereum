package core

import (
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"
	"xchain-go/common"
	"xchain-go/utils"

	log "github.com/inconshreveable/log15"

	"xchain-go/core/basic"
	"xchain-go/event"
)

const (
	chainHeadChanSize = 10
)

var (
	validDuration = big.NewInt(7200) //节点时间的有效范围
)

type blockChain interface {
	CurrentBlock() *basic.Block
	GetBlock(hash common.Hash, number uint64) *basic.Block
	//	StateAt(root common.Hash) (*state.StateDB, error)

	SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription
}

type TxPool struct {
	//TODO: config,chainconfig,event相关
	chain blockChain
	//currentState *state.StateDB //用来从区块链中查询当前状态TODO
	signer   basic.Signer //获取交易的发送者址地
	gasPrice *big.Int

	pending map[common.Address]*txList //存储对应账户的交易列表，由一个堆结构进行维护
	all     *txLookup

	chainHeadCh  chan ChainHeadEvent //TODO:
	chainHeadSub event.Subscription
	//设置事件订阅，以是的交易可以广播至订阅此事件类型的信道
	txFeed event.Feed              //订阅的入口
	Scope  event.SubscriptionScope //追踪多个订阅者提供集中取消订阅的功能
	wg     sync.WaitGroup          //定义同步等待的组
	//读写锁
	mu sync.RWMutex
}

//初始化txpool
func NewTxPool(chain blockChain) *TxPool {
	pool := &TxPool{
		chain:       chain,
		gasPrice:    utils.GasPrice, //TODO:gasprice 为定值.
		signer:      basic.NewEIP155Signer(utils.ChainID),
		pending:     make(map[common.Address]*txList),
		all:         newTxLookup(),
		chainHeadCh: make(chan ChainHeadEvent, chainHeadChanSize), //TODO:定义事件信息
	}
	//TODO: pool.locals相关
	//pool.reset(nil, chain.CurrentBlock().Header()) //初始化txpool

	//pool.chainHeadSub = chain.SubscribeChainHeadEvent(pool.chainHeadCh)

	//开始事件循环
	pool.wg.Add(1)
	//	go pool.loop()

	return pool
}

//初始化txpool中的statedb，并且去除已经在链上的交易。
/*
func (pool *TxPool) reset() {
	statedb, err := pool.chain.StateAt(newHead.Root)
	if err != nil {
		log.Error("Failed to reset txpool state", "err", err)
		return
	}
	pool.currentState = statedb
}
*/
//获取节点时间TODO:
func (pool *TxPool) getTime() (*big.Int, *big.Int) {
	t := time.Now().UnixNano()
	t_ := big.NewInt(t)
	up := new(big.Int).Add(t_, validDuration)
	down := new(big.Int).Sub(t_, validDuration)
	return up, down
}

//TODO:TxPool的主事件循环
func (pool *TxPool) loop() {
	defer pool.wg.Done()

	head := pool.chain.CurrentBlock()
	log.Debug("head", "head", head)

	for {
		select {
		case ev := <-pool.chainHeadCh:
			if ev.Block != nil {
				pool.mu.Lock()
			}
			head = ev.Block
			pool.mu.Unlock()
		case <-pool.chainHeadSub.Err():
			return
		}

	}
}

//验证交易,通过验证后的交易可放入pending
func (pool *TxPool) validateTx(tx *basic.Transaction, local bool) error {
	//判断交易的尺寸
	if tx.Size() > 32*1024 {
		return ErrOversizedData
	}
	if tx.Value().Sign() < 0 {
		return ErrNegativeValue
	}
	//
	//获取交易的发送者地址TODO:
	from, err := basic.Sender(pool.signer, tx)
	if err != nil {
		return ErrUnderpriced
	}
	log.Info("from", "from", from)
	/*
		if pool.currentState.GetLatestTime(from).Cmp(tx.Timestamp()) > -1 {
			return ErrTimestampTooLow
		}
	*/
	//验证当前交易的时间戳与状态树中时间戳的大小
	//if pool.chain.currentBlock.Time().Cmp(tx.Timestamp()) != 1 {
	//		return ErrTimestampTooLow
	//	}

	//TODO: Balance相关
	//if pool.CurrentState.GetBalance(from).Cmp(tx.Cost()) < 0 {
	//		return ErrInsufficientFunds
	//	}

	//up, down := pool.getTime()
	//if err != nil {
	//	return ErrTimeError
	//}
	//确认交易的时间戳在时间范围内TODO:
	//if tx.Timestamp().Cmp(down) == -1 || tx.Timestamp().Cmp(up) == 1 {
	//	return ErrTimestampOutBound
	//	}
	return nil
}

func (pool *TxPool) add(tx *basic.Transaction, local bool) (bool, error) {
	//	defer pool.wg.Done()
	//如果一笔交易是已知的,丢弃它
	hash := tx.Hash()
	if pool.all.Get(hash) != nil {
		log.Info("忽略已知交易", "hash", hash)
		return false, fmt.Errorf("已知的交易: %x", hash)
	}
	//如果一笔交易验证失败,丢弃它
	if err := pool.validateTx(tx, local); err != nil {
		log.Info("丢弃无效交易", "hash", hash, "err", err)
		return false, err
	}
	//TODO:交易池满的情况
	from, _ := basic.Sender(pool.signer, tx)
	if list := pool.pending[from]; list != nil {
		re := list.Add(tx)
		pool.all.Add(tx)
		go pool.txFeed.Send(tx)
		return re, nil
	}
	list := newtxList()
	pool.pending[from] = list
	insert := pool.pending[from].Add(tx)
	pool.all.Add(tx)

	go pool.txFeed.Send(tx)
	return insert, nil

}

func (pool *TxPool) AddLocal(tx *basic.Transaction) error {
	return pool.addTx(tx, true)
}

//AddLocals 将远程的一组交易添加到本地交易池中,若它们都有效
//func (pool *TxPool) AddRemotes(txs []*basic.Transaction) []error {
//	return pool.addTxs(txs, false)
//}

//AddRemote
func (pool *TxPool) AddRemote(tx *basic.Transaction) error {
	return pool.addTx(tx, false)
}

//TODO:
func (pool *TxPool) addTx(tx *basic.Transaction, local bool) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	replace, err := pool.add(tx, local)
	if err != nil {
		return err
	}

	if !replace {
		from, _ := basic.Sender(pool.signer, tx)
		log.Debug("from", "from", from)
		//TODO:
	}
	return nil
}

//TODO:
//func (pool *TxPool) addTxs(txs []*basic.Transaction, local bool) []error {

//}

/*
//向交易池中添加交易
func (pool *TxPool) Add(tx *basic.Transaction) error {
	//	from, err := types.Sender(pool.signer, tx)
	//	if err != nil {
	//		fmt.Println("Invalid Address")
	//	}
	from := *tx.To()
	if pool.pending[from] == nil {
		list := newtxList()
		pool.pending[from] = list
	}
	if err := pool.validateTx(tx); err != nil {
		fmt.Println("This transaction is invalid")
		//	return error
	}
	pool.pending[from].Add(tx)
	//向订阅了交易池的信道发送交易
	pool.txFeed.Send(*tx)
	return nil
}
*/

//TODO:对pending中的交易进行过滤，去除无效的交易
/*
func (pool *TxPool) AddPending() error {
	for _, list := range pool.pending {
		//TODO:
		local := true //用于标记是否为local
		list.Forward(pool.chain.CurrentBlock().Time())
		for _, tx := range list.items {
			if err := pool.validateTx(tx, local); err != nil {
				list.Remove(tx.Timestamp())
				fmt.Println("invalid transaction")
			}
		}
	}
	return nil
}
*/
//对交易池中的交易按照timestamp升序排序TODO: 所有交易构成一个交易列表一起排列
//func (pool *TxPool) Pending() (map[common.Address]basic.Transactions, error) {
//	pool.mu.Lock()
//	defer pool.mu.Unlock()

//	pending := make(map[common.Address]basic.Transactions)
//	for addr, list := range pool.pending {
//		pending[addr] = list.Flatten()
//	}
//
//	return pending, nil
//}

func (pool *TxPool) Pending() (basic.Transactions, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	var txs basic.Transactions
	for _, list := range pool.pending {
		txs = append(txs, list.Ready()...)
	}
	sort.Sort(basic.TransactionsByTimestamp(txs))

	return txs, nil
}

//验证pending中的交易是否已经存在于区块链上
/*
func (pool *TxPool) demonteUnSave() error {
	for addr, list := range pool.pending {
		timestamp := pool.currentState.GetTimestamp(addr)
		list.Forward(timestamp)
	}
	return nil
}
*/
//返回有效的交易列表
func (pool *TxPool) Ready() basic.Transactions {
	//	if err := pool.demonteUnSave(); err != nil {
	//		fmt.Println("Invalid tx in pending")
	//	}

	if _, err := pool.Pending(); err != nil {
		fmt.Println("Invalid tx in pending")
	}

	var ready basic.Transactions
	for _, list := range pool.pending {
		for _, tx := range list.Ready() {
			ready = append(ready, tx)
		}

	}
	return ready
}

//TODO:
func (pool *TxPool) removeTx(hash common.Hash) {}

//停止所有信道的订阅
func (pool *TxPool) Stop() {
	pool.Scope.Close()

	//	pool.chainHeadSub.Unsubscribe()
	pool.wg.Wait()
}

//订阅交易
func (pool *TxPool) SubscribeNewTxsEvent(ch chan<- *basic.Transaction) event.Subscription {
	return pool.Scope.Track(pool.txFeed.Subscribe(ch))
}

type txLookup struct {
	all  map[common.Hash]*basic.Transaction
	lock sync.RWMutex
}

// newTxLookup returns a new txLookup structure.
func newTxLookup() *txLookup {
	return &txLookup{
		all: make(map[common.Hash]*basic.Transaction),
	}
}

// Range calls f on each key and value present in the map.
func (t *txLookup) Range(f func(hash common.Hash, tx *basic.Transaction) bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for key, value := range t.all {
		if !f(key, value) {
			break
		}
	}
}

// Get returns a transaction if it exists in the lookup, or nil if not found.
func (t *txLookup) Get(hash common.Hash) *basic.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.all[hash]
}

// Count returns the current number of items in the lookup.
func (t *txLookup) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.all)
}

// Add adds a transaction to the lookup.
func (t *txLookup) Add(tx *basic.Transaction) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.all[tx.Hash()] = tx
}

// Remove removes a transaction from the lookup.
func (t *txLookup) Remove(hash common.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	delete(t.all, hash)

}

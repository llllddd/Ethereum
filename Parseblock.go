package main
import (
"fmt"
"log"
"time"
"sync"
"github.com/ethereum/go-ethereum/ethdb"
"github.com/ethereum/go-ethereum/common"
"github.com/ethereum/go-ethereum/core/types"
"github.com/ethereum/go-ethereum/core/rawdb"
)
const (
	//mainnet数据库路径
	dbpath = "/media/lzt/Elements/.ethereum/geth/chaindata"
	//并发最⼤协程数
	maxRoutine = 1000
	//最早区块的时间戳
	earliestTS = 1438269988
	//现在同步到的最新区块
	latestSyncedBlock = 5000000
)

type Database struct {
	leveldb *ethdb.LDBDatabase
}

func NewDatabase (path string) (*Database) {
	db, err := ethdb.NewLDBDatabase(path, 128, 1024)
	if err != nil {
		log.Fatal(err)
	}
	return &Database {db}
}

func (db *Database) getBlock(hash common.Hash, number uint64) *types.Block {
	block := rawdb.ReadBlock(db.leveldb, hash, number)
	if block == nil {
		return nil
	}
	return block
}

func (db *Database) GetBlockByNumber(number uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(db.leveldb, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return db.getBlock(hash, number)
}

func (db *Database) GetBlockByHash(hash common.Hash) *types.Block {
	num := rawdb.ReadHeaderNumber(db.leveldb, hash)
	if num == nil {
		return nil
	}
	return db.getBlock(hash, *num)
}

//在同步完成之前，这个⽅法是有问题的，返回的number总是0
func (db *Database) GetTailBlock() *types.Block {
	hash := rawdb.ReadHeadBlockHash(db.leveldb)
	if hash == (common.Hash{}) {
		return nil
	}
	return db.GetBlockByHash(hash)
}

//这个⽅法来代替上⾯的
func (db *Database) GetLatestSyncedBlock() *types.Block {
	return db.GetBlockByNumber(uint64(latestSyncedBlock))
}

//⼆分法快速定位⼀个区块
func (db *Database) numberInPeriod(ts, te time.Time, hs, he uint64) (uint64, types.Transactions, error) {
	hm := (hs + he) / 2
	block := db.GetBlockByNumber(hm)
	if block == nil {
		return 0, nil, fmt.Errorf("number %d block is not found!", hm)
	}
	bt := block.Time().Int64()
	if bt < ts.Unix() {
		return db.numberInPeriod(ts, te, hm+1, he)
	}
	if bt > te.Unix() {
		return db.numberInPeriod(ts, te, hs, hm-1)
	}
	return hm, block.Transactions(), nil
}

//获取某天的全部交易单
func (db *Database) GetTXsByDay(y,m,d int) (types.Transactions, uint64, error) {
	ts := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	te := ts.Add(24*time.Hour)
	block := db.GetLatestSyncedBlock()
	if te.Unix() < earliestTS || ts.Unix() > block.Time().Int64() {
		return nil, 0, fmt.Errorf("The date is not in interval")
	}
	hm, TX, err := db.numberInPeriod(ts, te, 1, block.NumberU64())
	if err != nil {
		return nil, 0, err
	}
	TXChan := make(chan types.Transactions)
	last := make(chan uint64, 1)
	var wg sync.WaitGroup
	wg.Add(2)
	go func(hh uint64) {
		for i:=hh; ; i++ {
			block := db.GetBlockByNumber(i)
			if block == nil {
				log.Printf("number %d block is not found!\n", i)
				continue
			}
			if block.Time().Int64() > te.Unix() {
				last <- i-1
				break
			}
			TXChan <- block.Transactions()
		}
		wg.Done()
	}(hm+1)
	go func(hh uint64) {
		for i:=hh; ; i-- {
			block := db.GetBlockByNumber(i)
			if block == nil {
				log.Printf("number %d block is not found!\n", i)
				continue
			}
			if block.Time().Int64() < ts.Unix() {
				break
			}
			TXChan <- block.Transactions()
		}
		wg.Done()
	}(hm-1)
	go func() {
		wg.Wait()
		close(TXChan)
	}()
	for tx := range TXChan {
		TX = append(TX, tx...)
	}
	return TX, <- last, nil
}

func main() {
	db := NewDatabase(dbpath)
	defer db.leveldb.Close()
	block := db.GetBlockByNumber(uint64(1000000))
	hb, err := block.Header().MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", hb)
	if txs, _, err := db.GetTXsByDay(2017, 8, 1); err != nil {
		log.Fatal(err)
	} else {
		for _, tx := range txs {
			txj, err := tx.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", txj)
		}
	}
}

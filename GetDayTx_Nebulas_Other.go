package main

import (
        "fmt"
        "log"
	"time"
        "sync"

        "github.com/nebulasio/go-nebulas/storage"
        "github.com/nebulasio/go-nebulas/core"
        "github.com/nebulasio/go-nebulas/util"
        "github.com/nebulasio/go-nebulas/util/byteutils"
        "github.com/gogo/protobuf/proto"
        "github.com/nebulasio/go-nebulas/core/pb"
        "github.com/nebulasio/go-nebulas/core/state"
 )

const (
        //mainnet数据库路径
        dbpath = "/home/lzt/Documents/mygo/src/github.com/alexmiaomiao/neb/data.db"
	//并发最大协程数
	maxRoutine = 5000
	//最早区块的时间戳
	earliestTS = 1522377345
)

type Database struct {
	rocksdb *storage.RocksStorage
}

func NewDatabase (path string) (*Database) {
	db, err := storage.NewRocksStorage(path)
	if err != nil {
		log.Fatal(err)
	}
	return &Database {db}
}

//通过hash获取某个区块
func (db *Database) GetBlockByHash(hash []byte) (*core.Block, error) {
        value, err := db.rocksdb.Get(hash)
	if err != nil {
                return nil, err
        } 
	pbBlock := new(corepb.Block)
	block := new(core.Block)
	if err = proto.Unmarshal(value, pbBlock); err != nil {
		return nil, err
	}
	if err = block.FromProto(pbBlock); err != nil {
		return nil, err
	}
	return block, nil
}

//通过高度得到某个区块
func (db *Database) GetBlockByHeight(h uint64) (*core.Block, error) {
	hash, err := db.rocksdb.Get(byteutils.FromUint64(h))
	if err != nil {
                return nil, err
        }
	return db.GetBlockByHash(hash)
}

//得到最后一个区块
func (db *Database) GetTailBlock() (*core.Block, error) {
        hash, err := db.rocksdb.Get([]byte(core.Tail))
	if err != nil {
                return nil, err
        }
	return db.GetBlockByHash(hash)
}

//获取整个区块链上所有交易单数量
func (db *Database) GetTotalTXs() (uint64, error) {
        var currentH uint64
        txChan := make(chan int)
	sema := make(chan struct{}, maxRoutine)
	res := make(chan uint64)
        var wg sync.WaitGroup

        if block, err := db.GetTailBlock(); err != nil {
                return 0, err
        } else {
                currentH = block.Height()
		fmt.Println(currentH)
        }

	wg.Add(int(currentH))
	go func() {
                wg.Wait()
                close(txChan)
        }()

	go func() {
		var total uint64
		for num := range txChan {
			total += uint64(num)
		}
		res <- total
	}()

        for i := uint64(1); i <= currentH; i++ {
		sema <- struct{}{}
                go func(h uint64) {
                        defer wg.Done()

			b, err := db.GetBlockByHeight(h)
			if err != nil {
				log.Println(err)
				log.Println(h)
				txChan <- 0
			} else {
				txChan <- len(b.Transactions())

			}
			<- sema
                }(i)
        }

        return  <- res, nil
}

//二分法快速定位一个时间段的某个区块高度
func (db *Database) heightInPeriod(ts, te time.Time, hs, he uint64) (uint64, core.Transactions, error) {
	hm := (hs + he) / 2
        block, err := db.GetBlockByHeight(hm)
	if err != nil {
                return 0, nil, err
        }
	bt := block.Timestamp()
	if bt < ts.Unix() {
		return db.heightInPeriod(ts, te, hm+1, he)
	}
	if bt > te.Unix() {
		return db.heightInPeriod(ts, te, hs, hm-1)
	}
	return hm, block.Transactions(), nil
}

//获取某天的全部交易单
func (db *Database) GetTXsByDay(y,m,d int) (core.Transactions, uint64, error) {
	ts := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	te := ts.Add(24*time.Hour)
	
        block, err := db.GetTailBlock()
	if err != nil {
                return nil, 0, err
        }

	if te.Unix() < earliestTS || ts.Unix() > block.Timestamp() {
		return nil, 0, fmt.Errorf("The date is not in interval")
	}
	
	hm, TX, err := db.heightInPeriod(ts, te, 2, block.Height())
	if err != nil {
		return nil, 0, err
	}

	TXChan := make(chan core.Transactions)
	last := make(chan uint64, 1)
        var wg sync.WaitGroup
	wg.Add(2)

	go func(hh uint64) {
		for i:=hh; ; i++ {
			block, err := db.GetBlockByHeight(i)
			if err != nil {
				log.Println(err)
				continue
			}
			if block.Timestamp() > te.Unix() {
				last <- i-1
				break
			}
			TXChan <- block.Transactions()
		}
		wg.Done()
	}(hm+1)

	go func(hh uint64) {
		for i:=hh; ; i-- {
			block, err := db.GetBlockByHeight(i)
			if err != nil {
				log.Println(err)
				continue
			}
			if block.Timestamp() < ts.Unix() {
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

//整理交易单
func ArrangeTXs(TXs core.Transactions) (binary, call, deploy core.Transactions) {
	for _, tx := range TXs {
		switch tx.Type() {
		case core.TxPayloadBinaryType:
			binary = append(binary, tx)
		case core.TxPayloadCallType:
			call = append(call, tx)
		case core.TxPayloadDeployType:
			deploy = append(deploy, tx)
		}
	}
	return
}

//统计合约调用次数
func CountContractCall(TXs core.Transactions) (map[string]int) {
	m := make(map[string]int)
	for _, tx := range TXs {
		m[tx.To().String()] += 1
	}
	return m
}

// //统计合约入账出账
// func countContractInOut(addr string, TXs core.Transactions) ([2]*util.Uint128, error) {
// 	inout := [2]*util.Uint128{util.NewUint128(), util.NewUint128()}

// 	Addr, err := core.AddressParse(addr)
// 	if err != nil {
// 		return inout, err
// 	}	
// 	if err != nil {
// 		return inout, err
// 	}
	
// 	b, c, _ := ArrangeTXs(TXs)
// 	for _, tx := range c {
// 		if tx.To().Equals(Addr) {
// 			inout[0], err = inout[0].Add(tx.Value())
// 			if err != nil {
// 				return inout, err
// 			}
// 		}
// 	}
// 	for _, tx := range b {
// 		if tx.From().Equals(Addr) {
// 			log.Fatal()
// 			inout[1], err = inout[1].Add(tx.Value())
// 			if err != nil {
// 				return inout, err
// 			}
// 		}
// 	}
// 	return inout, nil
// }

//获取账户在指定区块的余额
func (db *Database) GetBalanceAtHeight(addr string, h uint64) (*util.Uint128, error) {
	Addr, err := core.AddressParse(addr)
	if err != nil {
		return nil, err
	}
	
	block, err := db.GetBlockByHeight(h)
	if err != nil {
		return nil, err
	}

	accst, err := state.NewAccountState(block.StateRoot(), db.rocksdb)
	if err != nil {
		return nil, err
	}

	acc, err := accst.GetOrCreateUserAccount(Addr.Bytes())
	if err != nil {
		return nil, err
	}
	
	return acc.Balance(), nil
}

func main() {
	core.SetCompatibilityOptions(core.MainNetID)
	
        db := NewDatabase(dbpath)
        defer db.rocksdb.Close()

        // if block, err := GetBlockByHeight(uint64(410511), db); err != nil {
        //         log.Fatal(err)
        // } else {
        //         fmt.Println(block)
        //         for _, tx := range block.Transactions() {
        //                 fmt.Println(tx)
        //         }
        // }

        // if num, err := db.GetTotalTXs(); err != nil {
        //         log.Fatal(err)
	// } else {
        //         fmt.Println(num)
        // }

	if txs, l, err := db.GetTXsByDay(2018, 6, 2); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("2018-6-2\ntotal: %d\n\n", len(txs))
		b, c, d := ArrangeTXs(txs)
		fmt.Printf("binary: %d\n", len(b))
		for _, i := range b {
			fmt.Println(i)
		}
		fmt.Printf("call: %d\n", len(c))
		for _, i := range c {
			fmt.Println(i)
		}
		fmt.Printf("deploy: %d\n", len(d))
		for _, i := range d {
			fmt.Println(i)
		}
		fmt.Println()
		m := CountContractCall(c)
		fmt.Printf("Contract\tCalled count\tLastbalance\n")
		for k, v := range m {
			hb, err := db.GetBalanceAtHeight(k, l)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\t%d\t%v\n", k, v, hb)
			// cio, err := countContractInOut(k, txs)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Println(cio)
		}
	}
}

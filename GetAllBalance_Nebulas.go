package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/nebulasio/go-nebulas/core"
	"github.com/nebulasio/go-nebulas/core/pb"
	"github.com/nebulasio/go-nebulas/storage"
	"github.com/nebulasio/go-nebulas/util/byteutils"
)

const (
	//mainnet数据库路径
	dbpath = "/path/to/workspace/src/github.com/nebulasio/go-nebulas/mainnet/data.db"
)

//获取整个区块链上所有交易单数量
func getTotalTXs(db *storage.RocksStorage) (uint64, error) {
	var currentH uint64
	txChan := make(chan uint64)
	var wg sync.WaitGroup

	if block, err := getTailBlock(db); err != nil {
		return 0, err
	} else {
		currentH = block.Height()
	}

	//虚拟机性能限制，所以每个go协程负责500个区块，不然会栈溢出
	for i := uint64(1); i <= currentH; i += 500 {
		wg.Add(1)
		go func(h uint64) {
			defer wg.Done()

			t := uint64(0)
			for i := uint64(0); i < 500; i++ {
				if h+i > currentH {
					break
				}
				b, err := getBlockByHeight(h+i, db)
				if err != nil {
					log.Println(err)
				} else {
					t += uint64(len(b.Transactions()))
				}
			}
			txChan <- t
		}(i)
	}

	go func() {
		wg.Wait()
		close(txChan)
	}()

	var total uint64
	for num := range txChan {
		total += num
	}

	return total, nil
}

//得到区块上的所有交易地址

func getAddressFromBlock(db *storage.RocksStorage, height uint64) error {
	b, err := getBlockByHeight(aa)
	for _, v := range b.Transactions() {
		fmt.Println("from :" + v.from().address + " to: " + v.to().address)
	}
	return nil
}

//给定账户检索所有的账户交易单
func getTotalTXsOfAccount(db *storage.RocksStorage, aa uint64, addr byteutils.Hash) error {
	var totalFrom *util.Uint128
	var totalTo *util.Uint128
	for i = uint64(1); i <= aa; i += 1 {
		b, err := getBlockByHeight(i, db)
		for _, v := range b.Transactions() {
			if v.from().Equals(addr) {
				fmt.Println("支出" + v.to().address)
				totalFrom += v.value()
			} else {
				if v.to().Equals(addr) {
					fmt.Println("获得" + v.from().address)
					totalTo += v.value()
				}
			}

		}
	}
	return nil
}

//或取到指定高度所有交易单
func getTotalTXsDay(db *storage.RocksStorage, aa uint64, bb uint64) (uint64, error) {
	var t uint64
	var i uint64
	for i = aa; i <= bb; i += 1 {
		b, err := getBlockByHeight(i, db)
		for _, v := range b.Transactions() {
			fmt.Println(v.data.Type)
			fmt.Println(v.from)
		}
		if err != nil {
			log.Println(err)
		} else {
			t += uint64(len(b.Transactions()))
		}
	}

	return t, nil
}

//通过hash获取某个区块
func getBlockByHash(hash []byte, db *storage.RocksStorage) (*core.Block, error) {
	if b, err := db.Get(hash); err != nil {
		return nil, err
	} else {
		pbBlock := new(corepb.Block)
		block := new(core.Block)
		if err := proto.Unmarshal(b, pbBlock); err != nil {
			return nil, err
		}
		if err := block.FromProto(pbBlock); err != nil {
			return nil, err
		}
		return block, nil
	}
}

//通过高度得到某个区块
func getBlockByHeight(h uint64, db *storage.RocksStorage) (*core.Block, error) {
	if hash, err := db.Get(byteutils.FromUint64(h)); err != nil {
		return nil, err
	} else {
		if b, err := getBlockByHash(hash, db); err != nil {
			return nil, err
		} else {
			return b, nil
		}
	}
}

//得到最后一个区块
func getTailBlock(db *storage.RocksStorage) (*core.Block, error) {
	if hash, err := db.Get([]byte(core.Tail)); err != nil {
		return nil, err
	} else {
		if b, err := getBlockByHash(hash, db); err != nil {
			return nil, err
		} else {
			return b, nil
		}
	}
}

func main() {
	db, err := storage.NewRocksStorage(dbpath)

	if err != nil {
		fmt.Println(dbpath)
		log.Fatal(err)
	}
	defer db.Close()

	if block, err := getBlockByHeight(uint64(256977), db); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(block)
		for _, tx := range block.Transactions() {
			fmt.Println(tx)
		}
	}

	if num, err := getTotalTXs(db); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(num)
	}

}

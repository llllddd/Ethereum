package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/nebulasio/go-nebulas/core"
	"github.com/nebulasio/go-nebulas/core/pb"
	"github.com/nebulasio/go-nebulas/storage"
	"github.com/nebulasio/go-nebulas/util/byteutils"
	"log"
	"sync"
)

const (
	//mainnet数据库路径
	dbpath     = "/path/to/workspace/src/github.com/nebulasio/go-nebulas/mainnet/data.db"
	maxRoutine = 5000
)

/*
//获取整个区块链上所有交易单数量
func getTotalTXs(tt unit64 ,db *storage.RocksStorage) (uint64, error) {
        var currentT uint64
        // txChan := make(chan uint64)
        //var wg sync.WaitGroup
        var tranT unit64
        if block, err := getBlockByTimeStamp(tt,db) ; err != nil {
                return 0, err
        } else {
                currentT = block.Timestamp()
        }

        //虚拟机性能限制，所以每个go协程负责500个区块，不然会栈溢出
        for i := currentT; i >= currenT-86400; i-=15 {
          b,err := getBlockByTimeStamp(i,db)
          if err != nil{
            log.Println(err)
          }else{
            tranT+=unit64(len(b.Transactions()))
            fmt.Println(b.Transactions().data.Type)
         }
        return total, nil
}*/
//计算指定区块间的所有交易数目
func getTotalTXsInHeight(db *storage.RocksStorage, hs, he uint64) (uint64, error) {
	txChan := make(chan int)                //存储交易的信道
	sema := make(chan struct{}, maxRoutine) //
	res := make(chan uint64)
	var wg sync.WaitGroup

	if block, err := getTailBlock(db); err != nil {
		return 0, err
	} else {
		if he >= block.Height() {
			he = block.Height()
		}
	}

	temp := int(he - hs)
	wg.Add(temp) //工作中goroutines的个数
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

	for i := hs; i <= he; i++ {
		sema <- struct{}{}
		go func(h uint64) {
			defer wg.Done()
			b, err := getBlockByHeight(h, db)
			if err != nil {
				log.Println(err)
				log.Println(h)
				txChan <- 0
			} else {
				txChan <- len(b.Transactions())
			}
			<-sema
		}(i)
	}
	return <-res, nil
}

//遍历读取一天的区块数目
/*
func getTotalTXsDay(db *storage.RocksStorage, aa uint64, bb uint64) (uint64, error) {
	var t uint64
	var i uint64
	for i = aa; i <= bb; i += 1 {
		b, err := getBlockByHeight(i, db)
		for _, v := range b.Transactions() {
			fmt.Println(v.data.Type)
		}
		if err != nil {
			log.Println(err)
		} else {
			t += uint64(len(b.Transactions()))
		}
	}

	return t, nil
}
*/

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

//通过时间戳得到某个区块的高度
func getBlockHeightBytime(t int64, db *storage.RocksStorage, op int64) (uint64, error) {
	var h uint64
	b, _ := getBlockByTimeStamp(t, db, op)
	h = b.Height()
	for {
		c, _ := getBlockByHeight(h, db)
		if c.Timestamp() < t {
			h += 1
		} else {
			if op == 1 {
				break
			} else {
				if c.Timestamp() == t {
					break
				} else {
					h -= 1
					break
				}
			}
		}
	}

	for {
		c, _ := getBlockByHeight(h, db)
		if c.Timestamp() > t {
			h -= 1
		} else {
			if op == 1 {
				if c.Timestamp() == t {
					break
				} else {
					h += 1
					break
				}
			} else {
				break
			}
		}
	}
	//fmt.Println(h)
	return h, nil
}

//如果15s产生一个区块则理论上得到的区块
func getBlockByTimeStamp(t int64, db *storage.RocksStorage, op int64) (*core.Block, error) {
	var temp, temp1 int64
	temp = t % 15
	if temp != 0 {
		if op == 1 {
			temp1 = 15 - temp
			t += temp1
		} else {
			t -= temp1
		}
	}
	var h uint64
	h = uint64(((t - 1522377345) / 15) + 2) //通过时间戳判断区块高度
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

//通过高度得到区块
func getBlockByHeight(h uint64, db *storage.RocksStorage) (*core.Block, error) {
	if hash, err := db.Get(byteutils.FromUint64(h)); err != nil {
		log.Fatal(err)
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

//计算指定天数之后的区块高度
func getBlockHeightByday(day int64, day0 int64, db *storage.RocksStorage) (uint64, error) {
	var time int64
	time = day*86400 + day0
	if height, err := getBlockHeightBytime(time, db, 0); err != nil { //输入一个合法的时间区块
		log.Fatal(err)
		return 0, err
	} else {
		//fmt.Println(height)
		return height, err
	}

}

//整理交易单
func ArrangeTxs(tx core.Transactions) (binary, call, deploy core.Transactions) {
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
}

//统计合约调用次数
func CountContractCall(TXs core.Transactions) map[string]int {
	m := make(map[string]int)
	for _, tx := range TXs {
		m[tx.To().String()] += 1
	}
	return m
}
func main() {
	var fir, las int64
	las = 1526579745
	fir = 1526579655 //初始时间

	db, err := storage.NewRocksStorage(dbpath) //获取区块的数据
	if err != nil {
		fmt.Println(dbpath)
		log.Fatal(err)
	}
	defer db.Close()

	aa, _ := getBlockHeightBytime(fir, db, 1)
	end, _ := getBlockHeightByday(1, fir, db)
	end1, _ := getBlockHeightBytime(las, db, 0)
	fmt.Println("起始高度:", aa)
	fmt.Println("终止高度:", end)
	db.Close()
	d1b, err := storage.NewRocksStorage(dbpath) //获取区块的数据

	num1, _ := getTotalTXsInHeight(d1b, aa, end1)

	fmt.Println("交易的总数:", num1)
}

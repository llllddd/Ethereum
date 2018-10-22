package core

import (
	"fmt"
	"testing"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/core/rawdb"
	"xchain-go/crypto/sha3"
	"xchain-go/ethdb"
	"xchain-go/rlp"
)

// var Logger *mylog.SimpleLogger

// func log(prefix string) *mylog.SimpleLogger {
// 	common.InitLog(prefix)
// 	return common.Logger.NewSessionLogger()
// }

func TestAddBlock(t *testing.T) {
	// Log := common.Logger.NewSessionLogger()

	db, err := basic.OpenDatabase(dbFile, 512, 512)
	if err != nil {
		fmt.Printf("failed to open database,the err info ：%v ", err)
	}
	// db := ethdb.NewMemDatabase()

	// defer db.Close()
	genesis := basic.NewGenesisBlock()
	rawdb.WriteBlock(db, genesis)
	rawdb.WriteHeadBlockHash(db, genesis.Header().Hash())
	rawdb.WriteCanonicalHash(db, genesis.Header().Hash(), genesis.NumberU64())
	// Log.Infoln("写入初始块，块hash：", genesis.Header().Hash())
	//blockchain结构中的latestheadhash，存储当前的hash为块中的最新的hash
	latestheadhash, _ := rlp.EncodeToBytes(genesis.Header().Hash())

	//返回初始的blockchain结构
	bc := &BlockChain{latestheadhash, db}

	mockheader := &basic.Header{}
	mockbody := &basic.Body{}

	mockheader.Extradata = []byte("testdata")
	db.Close()
	bc.AddBlock(mockheader, mockbody)

}
func mockHash(x interface{}) common.Hash {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hash := common.BytesToHash(hw.Sum(nil))
	return hash
}

func TestPrintblock(t *testing.T) {
	// common.InitLog("TEST")
	// Log := common.Logger.NewSessionLogger()

	Log := log("core")

	Log.Infoln("循环遍历数据")

	//打开db
	db := ethdb.NewMemDatabase()
	genesis := basic.NewGenesisBlock()
	rawdb.WriteBlock(db, genesis)
	rawdb.WriteHeadBlockHash(db, genesis.Header().Hash())
	rawdb.WriteCanonicalHash(db, genesis.Header().Hash(), genesis.NumberU64())
	// Log.Infoln("写入初始块，块hash：", genesis.Header().Hash())
	//blockchain结构中的latestheadhash，存储当前的hash为块中的最新的hash
	latestheadhash, _ := rlp.EncodeToBytes(genesis.Header().Hash())

	//返回初始的blockchain结构
	bc := &BlockChain{latestheadhash, db}
	Log.Infoln("构造完成初始区块链，bc:", bc)
	// //新加一个区块
	header := &basic.Header{}
	body := &basic.Body{}
	header.Extradata = []byte("testdata")
	db.Close()
	bc.AddBlock(header, body)
	bc.PrintBlockchain()
}

package core

import (
	"fmt"
	"math/big"
	"testing"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/core/rawdb"

	log "github.com/inconshreveable/log15"
)

func TestAddBlock(t *testing.T) {
	//初始化blockchain
	blockchain := mockBlockchain()
	//获取到currentblock
	block := blockchain.CurrentBlock()
	// mockbody := &basic.Body{}
	txs := basic.Mocktxs()
	//获得链中的最大区块的区块号，然后在此基础上增加4个块
	latesthash := rawdb.ReadHeadBlockHash(blockchain.db)
	latestnumber := rawdb.ReadHeaderNumber(blockchain.db, latesthash)
	fmt.Println("latesthash", latesthash, "latestnumber", *latestnumber)
	fmt.Println("blockhash", block.Hash())

	number := int(*latestnumber)
	var blocks []*basic.Block
	//构造header,开始连续添加区块

	for i := number + 1; i <= number+4; i++ {
		mockheader := makeHeader(block)
		fmt.Println("mockheader.parenthash,", mockheader.ParentHash)
		extradata := fmt.Sprintf("添加的第%v个区块", i)
		mockheader.Extradata = []byte(extradata)
		block = basic.NewBlockWithHeader(mockheader).WithBody(txs)
		block.Header().BlockHash = block.Hash()

		blocks = append(blocks, block)

		for i := 0; i < len(blocks); i++ {
			fmt.Printf("-------------block[%v]----------", i)
			fmt.Println("parenthash:", blocks[i].ParentHash())
			fmt.Println("blockhash:", blocks[i].Hash())
		}
	}
	blockchain.InsertChain(blocks)

	fmt.Println("区块添加构造结束，打印所有区块信息")
	blockchain.PrintBlockchain()

}

func TestPrintblock(t *testing.T) {

	log.Debug("循环遍历数据")

	//初始化db、genesis
	blockchain := mockBlockchain()
	blockchain.PrintBlockchain()
}
func mockBlockchain() *BlockChain {
	//初始化db、genesis
	// db := ethdb.NewMemDatabase()
	db, err := basic.OpenDatabase(dbFile, 512, 512)
	if err != nil {
		fmt.Printf("failed to open database,the err info ：%v ", err)
	}
	// defer db.Close() //关闭数据库

	gspec := DefaultGenesisBlock()
	// genesis := gspec.MustCommit(db)
	log.Info("genesis信息已存入数据库", "genesis", gspec)
	_, err = SetupGenesisBlock(db, gspec)
	if err != nil {
		log.Error("genesis块加载失败")
	}
	blockchain, _ := NewBlockChain(db)
	return blockchain
}

// makeHeader 构造block的header结构
func makeHeader(preblock *basic.Block) *basic.Header {
	root := common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017")
	txhash := common.HexToHash("0a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9e")
	dposContextProto := basic.MockDposProto()
	// fmt.Println("preblock.blockhash", preblock.Hash())
	// header
	return &basic.Header{
		ParentHash:  preblock.Hash(),                                                 //父common.Hash
		Timestamp:   big.NewInt(1426516743),                                          //区块产生的时间戳
		Number:      new(big.Int).Add(preblock.Number(), common.Big1),                //区块号
		Extradata:   []byte("testdata"),                                              //额外信息
		Validator:   common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"), //区块验证者地址
		Root:        root,
		TxHash:      txhash,
		DposContext: dposContextProto,
	}
}

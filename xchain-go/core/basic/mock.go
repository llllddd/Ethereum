package basic

import (
	"fmt"
	"math/big"
	"xchain-go/common"
)

func MockBlockWithtxs() *Block {
	header := Mockheader()
	body := &Body{}
	// tx1 := NewTransaction(big.NewInt(0), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), []byte("test"))
	body.Transactions = Mocktxs()
	fmt.Println("len:", len(body.Transactions))
	// for _, tx := range body.Transactions {
	// 	fmt.Println("tx.data:", tx.data)
	// }

	block := NewBlockWithHeader(header).WithBody(body.Transactions)
	// for _, tx := range block.Transactions() {
	// 	fmt.Println("block-tx.data:", tx.data)
	// }
	// fmt.Println("mockblock:", block.String())
	return block
}

func Mockheader() *Header {
	// time := time.Now().Unix()
	root := common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017")
	txhash := common.HexToHash("0a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9e")
	dposContextProto := MockDposProto()

	// header
	return &Header{
		ParentHash:  common.HexToHash("bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff498"), //父common.Hash
		Timestamp:   big.NewInt(1426516743),                                                               //区块产生的时间戳
		Number:      big.NewInt(12),                                                                       //区块号
		Extradata:   []byte("testdata"),                                                                   //额外信息
		Validator:   common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"),                      //区块验证者地址
		Root:        root,
		TxHash:      txhash,
		DposContext: dposContextProto,
	}
}

// Mocktx 测试用的构造tx
func Mocktxs() Transactions {
	var txs Transactions
	tx1 := NewTransaction(big.NewInt(0), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), big.NewInt(50000000), []byte("test"))
	// hash := tx1.Hash()
	// tx1.data.Hash = &hash
	txs = append(txs, tx1)
	tx2 := NewTransaction(big.NewInt(12), common.HexToAddress("095e7baea6a6c7c4c2dfe1237efac326af552d8"), big.NewInt(6), big.NewInt(50000000), []byte("test2"))
	// hash = tx2.Hash()
	// tx2.data.Hash = &hash
	txs = append(txs, tx2)
	for _, tx := range txs {
		// fmt.Println("tx.data:", tx.data)
		fmt.Println("tx.hash:", tx.Hash())
	}
	return txs
}
func Mocktx() *Transaction {

	tx := NewTransaction(big.NewInt(0), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), big.NewInt(50000000), []byte("test"))
	// hash := tx1.Hash()
	// tx1.data.Hash = &hash
	// txs = append(txs, tx1)

	return tx
}

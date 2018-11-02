package main

import (
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"
	"xchain-go/common"
	"xchain-go/core"
	"xchain-go/core/basic"
	"xchain-go/core/rawdb"
	"xchain-go/p2p"

	log "github.com/inconshreveable/log15"
)

//负责处理命令行参数的CLI
type CLI struct {
	//对象是链
	// bc *core.BlockChain
}

//方法说明
func (cli *CLI) printusage() {
	fmt.Println("usage:")
	fmt.Println(" ./xchain createblockchain                  -- create a blockchain")
	fmt.Println(" ./xchain createwallet -passwd 'passward'   -- create a createwallet")
	fmt.Println(" ./xchain printblockchain                   -- print all block in the blockchain")
	fmt.Println(" ./xchain startServer                       -- start the blockchain")
	fmt.Println(" ./xchain addblock -data 'Block data'       -- add a new block into the blockchain")

}

//新建区块链信息
func (cli *CLI) createBlockchain() {
	// Log := log(prefix)
	// 1. 初始化genesis
	// 2. 根据genesis创建blockchain
	gspec := core.DefaultGenesisBlock()
	db, err := basic.OpenDatabase(dbFile, 0, 0)
	if err != nil {
		log.Error("打开数据库失败")
	}
	log.Info("打开数据库成功")
	defer db.Close()
	genesis := gspec.MustCommit(db)
	log.Info("genesis信息已存入数据库", "genesis", genesis)
	// hash, err := core.SetupGenesisBlock(db, genesis)
	// if err != nil {
	// 	Log.Error("")
	// }
	blockchain, err := core.NewBlockChain(db)
	if err != nil {
		log.Error("初始化区块链失败。", "err", err)
	}
	log.Debug("blockchain:", "blockchain", blockchain)
}

//产生交易
func (cli *CLI) createTx() {
	// key, _ := crypto.GenerateKey()
	// addr := crypto.PubkeyToAddress(key.PublicKey)

	// signer := basic.NewEIP155Signer(big.NewInt(18))
	// tx := basic.Mocktx()
	// signedTx, err := basic.SignTx(tx, signer, key)
	// if err != nil {
	// 	log.Error("err", "err", err)
	// }
	// submitTransactionWithoutBackend(signedTx)

}

//添加区块
func (cli *CLI) addBlock(extradata string) {
	//初始化区块链
	db, err := basic.OpenDatabase(dbFile, 0, 0)
	if err != nil {
		log.Error("打开数据库失败")
	}
	log.Debug("打开数据库成功")
	blockchain, err := core.NewBlockChain(db)
	if err != nil {
		log.Error("初始化区块链失败", "err", err)
	}
	//构造块
	currentBlock := blockchain.CurrentBlock()
	latesthash := rawdb.ReadHeadBlockHash(db)
	latestnumber := rawdb.ReadHeaderNumber(db, latesthash)
	number := int(*latestnumber)
	log.Debug("构造的block的number信息", "number", number)
	header := makeHeader(currentBlock)
	body := &basic.Body{}
	body.Transactions = basic.Mocktxs()
	header.Extradata = []byte(extradata)
	log.Debug("构造的block的header信息", "header", header)
	header.Number = big.NewInt(int64(number + 1))
	block := basic.NewBlockWithHeader(header).WithBody(body.Transactions)
	var blocks []*basic.Block
	blocks = append(blocks, block)
	blockchain.InsertChain(blocks)
	blockchain.PrintBlockchain()

}

//打印区块链内容
func (cli *CLI) printblockchain() {
	//初始化区块链
	db, err := basic.OpenDatabase(dbFile, 0, 0)
	if err != nil {
		log.Error("打开数据库失败")
	}
	log.Debug("打开数据库成功")
	blockchain, err := core.NewBlockChain(db)
	if err != nil {
		log.Error("初始化区块链失败", "err", err)
	}
	blockchain.PrintBlockchain()
}

//启动dpos的链
func (cli *CLI) startServer() {
	//开启监听服务
	go func() {
		p2p.UdpListen()
	}()

	CreateNode()
	log.Debug("NodeArr", "NodeArr", NodeArr)
	Vote()
	nodes := SortNodes()
	log.Debug("nodes", "nodes", nodes)
	//初始化一条链
	db, err := basic.OpenDatabase(dbFile, 0, 0)
	if err != nil {
		log.Error("打开数据库失败")
	}
	log.Debug("打开数据库成功")
	blockchain, err := core.NewBlockChain(db)
	if err != nil {
		log.Error("初始化区块链失败", "err：", err)
	}
	// 添加区块
	nodesNum := len(nodes)
	currentBlock := blockchain.CurrentBlock()

	for i := 0; i < nodesNum; i++ {

		//构造区块
		body := &basic.Body{}
		body.Transactions = basic.Mocktxs()
		header := makeHeader(currentBlock)
		extradata := fmt.Sprintf("sender 1btc to alice")
		header.Extradata = []byte(extradata)
		header.Validator = nodes[i].Name
		block := basic.NewBlockWithHeader(header).WithBody(body.Transactions)
		// cli.bc.AddBlock(header, body)
		var blocks []*basic.Block
		blocks = append(blocks, block)
		blockchain.InsertChain(blocks)
		p2p.UdpDial()
		time.Sleep(time.Second)
		currentBlock = block
	}

	blockchain.PrintBlockchain()

}

//创建钱包
// func (cli *CLI) createWallet(walletPasswd string) {
// 	// Log := log(prefix)
// 	// 1. 初始化genesis
// 	// 2. 根据genesis创建blockchain
// 	// ctx := *new(context.Context)
// 	var b Backend
// 	private := NewPrivateAccountAPI(b)
// 	accaddr, err := private.NewAccount(walletPasswd)
// 	if err != nil {
// 		fmt.Println("创建账户失败")
// 	}
// 	fmt.Println("新建账户的地址为：", accaddr)
// }

//构造交易
// func (cli *CLI) createTx(txJSON string) {
// 	// 1. 获得传入的tx参数
// 	// 2. 根据tx中的from获得wallet
// 	// 3. 调用相关方法进行

// 	//解析tx参数
// 	var sendTxArgs SendTxArgs
// 	if err := json.Unmarshal([]byte(txJSON), &sendTxArgs); err != nil {
// 		log.Info("获取到的参数的地址", "from", sendTxArgs.From)
// 		log.Error("ERR:", "ERR", err)
// 	}
// 	// Log.Infoln("from地址：", sendTxArgs.From)
// 	// Log.Infoln("To地址：", sendTxArgs.To)
// 	// Log.Infoln("Timestamp地址：", sendTxArgs.Timestamp)
// 	// from := sendTxArgs.From
// 	//获取from的wallet
// 	// account := accounts.Account{Address: sendTxArgs.From}
// 	passwd := "111"
// 	ctx := *new(context.Context)
// 	var b Backend
// 	private := NewPrivateAccountAPI(b)
// 	private.SendTransaction(ctx, sendTxArgs, passwd)
// }

//校验输入的参数是否合法
//如果参数的长度小于2，即参数中只有二进制项目名，没有后面相应的执行的命令及参数，则退出
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		// Log.Infoln(os.Args[:])
		cli.printusage()
		os.Exit(1)
	}
	// Log.Infoln(os.Args[:])
}

func (cli *CLI) Run() {
	cli.validateArgs()
	//flag命令进行定义
	printusageCmd := flag.NewFlagSet("printusage", flag.ExitOnError)             //打印用法
	printBlockchainCmd := flag.NewFlagSet("printblockchain", flag.ExitOnError)   //打印区块链
	addblockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)                 //添加区块
	addblockData := addblockCmd.String("data", "", "Block ExtraData")            //addblockData定义一个string的flag,来对data的参数进行判断。得到的是一个指针对象
	startServerCmd := flag.NewFlagSet("startServer", flag.ExitOnError)           //启动区块链进行投票出块
	createblockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError) //创建区块链
	// createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)         //创建钱包
	// walletPasswd := createWalletCmd.String("passwd", "", "passward")             //pw
	createTxCmd := flag.NewFlagSet("createtx", flag.ExitOnError) //新增转账交易
	// txJSON := createWalletCmd.String("tx", "", "{\"from\":\"123\",\"to\":\"123\",\"value\":12}") //{"from":"123","to":"123","value":"12"}

	//判断传入的命令参数内容
	//参数内容中，第一个字符为二进制项目名，第二个字符开始内容为实际的输入。
	//第二个字符为实际想执行的命令字段，第三个字符为要传入的参数
	//判断第二个参数的内容，来执行printblockchain、addblock、printusage的对应操作
	switch os.Args[1] {
	case "printusage":
		err := printusageCmd.Parse(os.Args[2:])
		if err != nil {
			log.Error("err", "err", err)
		}

	case "printblockchain":
		err := printBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Error("err", "err", err)
		}

	case "addblock":
		err := addblockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Error("err", "err", err)
		}

	case "startServer":
		err := startServerCmd.Parse(os.Args[2:])
		if err != nil {
			log.Error("err", "err", err)
		}
	case "createblockchain":
		err := createblockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Error("err", "err", err)
		}
	// case "createwallet":
	// 	err := createWalletCmd.Parse(os.Args[2:])
	// 	if err != nil {
	// 		log.Error("err", "err", err)
	// 	}
	case "createtx":
		err := createTxCmd.Parse(os.Args[2:])
		if err != nil {
			log.Error("err", "err", err)
		}
	default:
		//非零为错误退出
		cli.printusage()
		os.Exit(1)
	}

	//初始化启动的时候，打印usage
	if printusageCmd.Parsed() {
		cli.printusage()
	}

	//命令行信息读取结束，则执行打印区块链的操作
	if printBlockchainCmd.Parsed() {
		cli.printblockchain()
	}

	//判断addblock命令传入的参数，添加新区块
	if addblockCmd.Parsed() {
		//判断输入的data参数的值是否为空
		if *addblockData == "" {
			// Log.Infoln(*addblockData)
			cli.printusage()
			os.Exit(1)
		}
		log.Debug("addblockData", "addblockData", *addblockData)
		cli.addBlock(*addblockData)
		// cli.printblockchain()
	}

	if startServerCmd.Parsed() {
		cli.startServer()
	}
	if createblockchainCmd.Parsed() {
		cli.createBlockchain()
	}
	//创建钱包
	// if createWalletCmd.Parsed() {
	// 	if *walletPasswd == "" {
	// 		// Log.Infoln(*addblockData)
	// 		cli.printusage()
	// 		os.Exit(1)
	// 	}
	// 	log.Info("walletPasswd", "walletPasswd", *walletPasswd)
	// 	cli.createWallet(*walletPasswd)
	// }

	//创建tx
	if createTxCmd.Parsed() {
		// if *txJSON == "" {
		// 	// Log.Infoln(*addblockData)
		// 	cli.printusage()
		// 	os.Exit(1)
		// }
		// log.Info("txJSON", "txJSON", *txJSON)
		// cli.createTx(*txJSON)
		cli.createTx()
	}
}

// 创建节点
func CreateNode() {
	for i := 0; i < 10; i++ {
		name := RandStringBytes(basic.AddressLength)
		NodeArr[i] = Node{name, 0}
	}
}

func RandStringBytes(n int) common.Address {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var b common.Address
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

//简单模拟投票
func Vote() {
	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().UnixNano())
		vote := rand.Intn(10) + 1
		NodeArr[i].Votes = vote
	}
}

//选出票数最多的前3位
func SortNodes() []Node {
	n := NodeArr
	for i := 0; i < len(n); i++ {
		for j := 0; j < len(n)-1; j++ {
			if n[j].Votes < n[j+1].Votes {
				n[j], n[j+1] = n[j+1], n[j]
			}
		}
	}
	return n[:3]
}

// makeHeader 构造block的header结构
func makeHeader(preblock *basic.Block) *basic.Header {
	// var time *big.Int
	// if preblock.Time() == nil {
	// 	time = big.NewInt(10)
	// } else {
	// 	time = new(big.Int).Add(preblock.Time(), big.NewInt(10)) // block time is fixed at 10 seconds
	// }
	time := time.Now().Unix()
	root := common.Hash{}
	txhash := common.Hash{}
	dposContextProto := basic.MockDposProto()

	// header
	return &basic.Header{
		ParentHash:  preblock.Header().BlockHash,                      //父common.Hash
		Timestamp:   big.NewInt(time),                                 //区块产生的时间戳
		Number:      new(big.Int).Add(preblock.Number(), common.Big1), //区块号
		Extradata:   []byte{},                                         //额外信息
		Validator:   preblock.Validator(),                             //区块验证者地址
		Root:        root,
		TxHash:      txhash,
		DposContext: dposContextProto,
	}
}

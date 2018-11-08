package main

import (
	"fmt"
	"math/big"
	"net"
	"os"
	"sync"
	tim "time"
	"xchain-go/accounts"
	"xchain-go/accounts/keystore"
	"xchain-go/core"
	"xchain-go/core/basic"
	"xchain-go/rlp"
)

var (
	gasLimit = big.NewInt(10000)
	time     = big.NewInt(111)
	amount   = big.NewInt(10)
)

func main() {

	//新建blockchain
	db, err := basic.OpenDatabase("测试数据库", 512, 512)
	if err != nil {
		fmt.Println("新建数据库失败")
	}
	//创建创世区块
	genesis := core.DefaultGenesisBlock()
	hash, err := core.SetupGenesisBlock(db, genesis)
	if err != nil {
		fmt.Println("创世区块存入失败")
	}
	fmt.Println("创世区块的哈希值为:", hash)
	blockchain, err := core.NewBlockChain(db)
	if err != nil {
		fmt.Println("创建区块链失败")
	}
	blockchain.PrintBlockchain()

	//生成密钥
	//key := keystore.NewKeyForDirectICAP(rand.Reader)
	//密钥文件存储
	//Account相关的api在internal/ethapi
	backends := []accounts.Backend{keystore.NewPlaintextKeyStore("测试钱包"), keystore.NewKeyStore("加密的钱包", keystore.StandardScrypN, keystore.StandardScrypP)}

	manager := accounts.NewManager(backends...)

	fmt.Println("bakend的个数:", len(manager.Backends(keystore.KeyStoreType)))
	plainkeystore := manager.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	//创建新的账户
	account, err := plainkeystore.NewAccount("test")
	if err != nil {
		fmt.Println("创建账户失败")
	}
	fmt.Printf("%x,%s", account.Address, account.URL)

	wallet := manager.Wallets()
	fmt.Println("钱包个数:", len(wallet))
	fmt.Println("账户个数", len(wallet[1].Accounts()))
	signaccount := wallet[0].Accounts()[0]
	txs1 := basic.NewTransaction(time, account.Address, amount, gasLimit, []byte{111})
	txs2 := basic.NewTransaction(big.NewInt(222), wallet[0].Accounts()[0].Address, big.NewInt(222), gasLimit, []byte{222})
	txs3 := basic.NewTransaction(big.NewInt(333), wallet[0].Accounts()[0].Address, big.NewInt(333), gasLimit, []byte{33})

	walletsign, _ := manager.Find(signaccount)

	signtxs1, _ := walletsign.SignTxWithPassphrase(signaccount, "test", txs1, big.NewInt(1))
	signtxs2, _ := walletsign.SignTxWithPassphrase(signaccount, "test", txs2, big.NewInt(1))
	signtxs3, _ := walletsign.SignTxWithPassphrase(signaccount, "test", txs3, big.NewInt(1))

	//创建交易池
	txs := basic.Transactions{signtxs1, signtxs2, signtxs3}
	txpool := core.NewTxPool(blockchain)

	//配置网络
	conn, err := net.Dial("udp", "192.168.82.67:8333")
	defer conn.Close()
	if err != nil {
		os.Exit(1)
	}

	//创建节点监听本地交易池的信息

	type Node struct {
		TxChan chan *basic.Transaction
	}

	node1 := Node{TxChan: make(chan *basic.Transaction, 5)}

	sub := txpool.SubscribeNewTxsEvent(node1.TxChan)
	defer sub.Unsubscribe()

	var Wg = sync.WaitGroup{}
	for _, sendtx := range txs {
		Wg.Add(1)
		if err := txpool.AddLocal(sendtx); err != nil {
			fmt.Println(err)
		}

		go func() {
			for {
				defer Wg.Done()
				select {
				case result := <-node1.TxChan:
					fmt.Println("节点从交易池中得到的交易是:", result)
					//将监听到的本地交易池中的交易进行RLP编码
					str, _ := rlp.EncodeToBytes(result)
					conn.Write(str)
					conn.SetWriteDeadline(tim.Now().Add(tim.Second * 1))
					fmt.Println("发送交易", result.Hash())
					return
				case <-sub.Err():
					fmt.Println("监听出错")
					return

				}
			}
		}()

		Wg.Wait()
	}

	//交易添加到交易池
	//if err := txpool.AddLocal(signtxs1); err != nil {
	//	fmt.Println(err)
	//}

	//if err := txpool.AddLocal(signtxs2); err != nil {
	//		fmt.Println(err)
	//}
	//	if err := txpool.AddLocal(signtxs3); err != nil {
	//		fmt.Println(err)
	//	}
	//	txpool.Stop()
	//	txss, _ := txpool.Pending()
	//	fmt.Println(txss)
}

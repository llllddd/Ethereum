package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"xchain-go/core"
	"xchain-go/core/basic"
	"xchain-go/rlp"
	//	"time"
)

func main() {
	fmt.Println("UdpListen Start")
	packetConn, err := net.ListenPacket("udp", "192.168.82.67:8333")

	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
		return
	}
	defer packetConn.Close()

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

	//创建交易池子
	txpool := core.NewTxPool(blockchain)

	var buf [512]byte

	packetConn.SetReadDeadline(time.Now().Add(time.Second * 15))
	for {

		n, addr, err := packetConn.ReadFrom(buf[0:])

		//fmt.Println("buf:", buf[0:n])
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
			break
		}
		fmt.Fprintf(os.Stdout, "listen recv: %x\n", buf[0:n])

		// 将数组反序列化
		var tx *basic.Transaction
		err = rlp.DecodeBytes(buf[0:n], &tx)
		//fmt.Println("number0:", blocks[0].header.Number)
		if err != nil {
			fmt.Println("decodeerr", err)
		}
		txpool.AddRemote(tx)
		fmt.Println(tx)

		_, err = packetConn.WriteTo(buf[0:n], addr)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
			break
		}
	}
	txs, _ := txpool.Pending()
	fmt.Println("交易池中交易的数量:", len(txs))

}

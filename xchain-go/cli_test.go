package main

import (
	"testing"
)

func TestCreateBlockchain(t *testing.T) {

	cli := CLI{}
	// cli.Run()
	cli.createBlockchain()
}

func TestAddblock(t *testing.T) {

	cli := CLI{}
	// cli.Run()
	cli.addBlock("data")
}

func TestStartServer(t *testing.T) {

	cli := CLI{}
	// cli.Run()
	cli.startServer()
}

func TestCreatTX(t *testing.T) {

	cli := CLI{}
	// cli.Run()
	// tx := "{\"from\":\"0x8888f1f195afa192cfee860698584c030f4c9db1\",\"to\":\"0x8888f1f195afa192cfee860698584c030f4c9db1\",\"value\":\"0x12\"}"
	cli.createTx()
}

// func TestCreateWallet(t *testing.T) {

// 	cli := CLI{}
// 	// cli.Run()
// 	walpasswd := "111"
// 	cli.createWallet(walpasswd)
// }

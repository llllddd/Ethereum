package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

//keystoreWallet{}实现了钱包接口，结构体中包含账户地址，以及与秘钥有关的成员KeyStore{}

func main() {
	var testSigData = make([]byte, 32)
	keydir, _ := ioutil.TempDir("", "Eth-wallet")
	key_store := keystore.NewKeyStore(keydir, 4096, 6)

	account1, _ := key_store.NewAccount("Dli")
	fmt.Println(account1)
	account2, _ := key_store.NewAccount("Dli")
	fmt.Println(account2)

	fmt.Println(key_store)
	sign, _ := key_store.SignHashWithPassphrase(account1, "Dli", testSigData)
	fmt.Println(sign)
	a := key_store.Wallets()
	fmt.Println(a)
	b := key_store.Accounts()
	fmt.Println(len(b))
	fmt.Println(sk)
}

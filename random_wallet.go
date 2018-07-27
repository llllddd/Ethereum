package main

import (
	"fmt"
	"io/ioutil"

	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

//keystoreWallet{}实现了钱包接口，结构体中包含账户地址，以及与秘钥有关的成员KeyStore{}
const (
	ScryptN = 4096
	ScryptP = 6
)

func main() {
	var testSigData = make([]byte, 32)
	keydir, _ := ioutil.TempDir("", "Eth-wallet")
	key_store := keystore.NewKeyStore(keydir, 4096, 6)

	account1, _ := key_store.NewAccount("aaa")
	fmt.Printf("账户地址为: %x\n", account1.Address)
	fmt.Printf("账户路径协议为: %s\n账户存储路径为: %s\n", account1.URL.Scheme, account1.URL.Path)

	account2, _ := key_store.NewAccount("bbb")
	fmt.Printf("账户地址为: %x\n", account2.Address)
	fmt.Printf("账户路径协议为: %s\n账户存储路径为: %s\n", account2.URL.Scheme, account2.URL.Path)

	Wallet_list := key_store.Wallets()

	for _, wallet := range Wallet_list {
		fmt.Println(wallet.URL().Path)
	}

	sign, _ := key_store.SignHashWithPassphrase(account1, "aaa", testSigData)
	fmt.Printf("由账户1私钥得到的签名:%x", sign)
	b := key_store.Accounts()
	fmt.Println(len(b))
	//解析本地存储的KeyJson文件
	storage := key_store.Storage()
	keyjson, _ := ioutil.ReadFile(account1.URL.Path)
	//首先定义一个简单的映射来确定密钥文件的版本,键string的值可以使任意类型所以用空接口。
	m := make(map[string]interface{})
	//对json文件解码
	if err := json.Unmarshal(keyjson, &m); err != nil {
		fmt.Println("解码错误")
	}
	//从json文件中读取cypto参数
	fmt.Println("账户地址:", m["address"], "id:", m["id"], "密钥文件类型:", m["version"])

	key, _ := storage.GetKey(account1.Address, account1.URL.Path, "aaa")
	fmt.Println("私钥为:%x", key)
}

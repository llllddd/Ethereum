package accounts

import (
	"math/big"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/event"
)

type Account struct {
	Address common.Address
	URL     URL `json:"url"`
}

type Wallet interface {
	//查询密钥文件位置
	URL()
	//根据用户口令打开私钥
	Open(passphrase string) error
	//关闭钱包
	Close() error
	//返回所有已知的可签名账户
	Accounts() []Account
	//验证账户是否存在钱包中
	Contains(account Account) bool
	//TODO: 下面两个函数时针对HD钱包
	//Derive(path DerivationPath,pin bool)(Account, error)
	//SelfDerive(base DerivationPath, chain ethereun.ChainStateReader)
	SignHash(account Account, hash []byte) ([]byte, error)
	//由用户输入的口令进行签名
	SignHashWithPassphrase(account Account, passphrase string, hash []byte) ([]byte, error)

	SignTxWithPassphrase(account Account, passphrase string, tx *basic.Transaction, chainID *big.Int)
}

type Backend interface {
	//返回所有已知的钱包列表
	Wallets() []Wallet
	//订阅通道
	Subscribe(sink chan<- WalletEvent) event.Subscription
}

type WalletEventType int

const (
	WalletArrived WalletEventType = iota

	WalletOpend

	WalletDropped
)

type WalletEvent struct {
	Wallet Wallet
	Kind   WalletEventType
}

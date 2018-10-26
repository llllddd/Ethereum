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
	// SignTx请求钱包将给定的交易进行签名。
	//它仅通过其中包含的地址查找指定的帐户，或者可选地借助嵌入的URL字段中的任何位置元数据。
	//如果钱包需要额外的身份验证来签署请求（例如，解密帐户的密码，或PIN码验证交易），
	//将返回AuthNeededError实例，其中包含用户的信息，说明需要哪些字段或操作。
	//用户可以通过SignTxWithPassphrase或其他方式（例如，在密钥库中解锁帐户）提供所需的详细信息来重试。
	SignTx(account Account, tx *basic.Transaction, chainID *big.Int) (*basic.Transaction, error)
	// SignTxWithPassphrase请求钱包签署给定的事务，并将给定的密码作为额外的身份验证信息。
	//它仅通过其中包含的地址查找指定的帐户，或者可选地借助嵌入的URL字段中的任何位置元数据。
	SignTxWithPassphrase(account Account, passphrase string, tx *basic.Transaction, chainID *big.Int) (*basic.Transaction, error)
	//查询密钥文件位置
	URL() URL
	//返回钱包的状态
	Status() (string, error)
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

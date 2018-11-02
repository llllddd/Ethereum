package main

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"time"
	"xchain-go/accounts"
	"xchain-go/accounts/keystore"
	"xchain-go/common"
	"xchain-go/common/hexutil"
	"xchain-go/core/basic"
	"xchain-go/utils"

	log "github.com/inconshreveable/log15"
)

// EthAPIBackend implements ethapi.Backend for full nodes
// type EthAPIBackend struct {
// 	eth *Ethereum
// }
type Backend interface {
	// General Ethereum API
	AccountManager() *accounts.Manager

	// TxPool API
	SendTx(ctx context.Context, signedTx *basic.Transaction) error
}

// PrivateAccountAPI 提供了一个API来访问此节点管理的帐户。
//它提供了创建（un）锁定列表帐户的方法。有些方法接受密码，因此默认情况下被视为私有。
type PrivateAccountAPI struct {
	am *accounts.Manager //账户管理，获得wallet
	b  Backend
	// nonceLock *AddrLocker
	// qiqi-todo:增加一个关于时间戳的lock
}

// NewPrivateAccountAPI创建一个新的PrivateAccountAPI对象
func NewPrivateAccountAPI(b Backend) *PrivateAccountAPI {
	return &PrivateAccountAPI{
		// am: b.AccountManager(),
		am: new(accounts.Manager),
		b:  b,
	}
}

// SendTransaction 使用节点上的from账户的私钥，给tx签名,并将交易提交到txpool
func (s *PrivateAccountAPI) SendTransaction(ctx context.Context, args SendTxArgs, passwd string) (common.Hash, error) {
	// qiqi-todo:判断传入的参数中的时间戳的，并对时间戳加锁。
	// 1. 传入的交易参数，构造成tx，并增加签名
	// 2. 将带有签名的交易提交至txpool
	signed, err := s.signTransaction(ctx, args, passwd)
	if err != nil {
		return common.Hash{}, err
	}
	return submitTransaction(ctx, s.b, signed)

}

// signTransaction 传入的交易参数，构造成tx，并增加签名
func (s *PrivateAccountAPI) signTransaction(ctx context.Context, args SendTxArgs, passwd string) (*basic.Transaction, error) {
	// 1. 查找交易发起者的account
	// 2. 设置一些健全性默认值并在失败时终止
	// 3. 组装交易并用钱包签名

	//查找account
	account := accounts.Account{Address: args.From}
	wallet, err := s.am.Find(account)
	if err != nil {
		return nil, err
	}

	//校验健全值
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}

	//组装交易并签名
	tx := args.toTransaction()
	var chainID *big.Int
	// qiqi-todo:从config中读取chainID
	// if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) {
	// 	chainID = config.ChainID
	// }
	chainID = utils.ChainID

	return wallet.SignTxWithPassphrase(account, passwd, tx, chainID)
}

// submitTransaction 是一个帮助函数，它将tx提交给txPool并记录消息。
func submitTransaction(ctx context.Context, b Backend, signedTx *basic.Transaction) (common.Hash, error) {
	// qiqi-todo:对交易进行简单的校验validate
	// 将交易提交至txpool,addLocal
	// if err := *core.TxPool.AddLocal(signedTx); err != nil {
	// 	return common.Hash{}, err
	// }
	// qiqi-todo:增加判断。如果tx.To()==nil，那么为创建合约的交易
	log.Info("提交交易至txpool", "txhash", signedTx.Hash().Hex(), "交易接收方", signedTx.To())
	return signedTx.Hash(), nil
}

// submitTransaction 是一个帮助函数，它将tx提交给txPool并记录消息。
func submitTransactionWithoutBackend(signedTx *basic.Transaction) (common.Hash, error) {
	// qiqi-todo:对交易进行简单的校验validate
	// 将交易提交至txpool,addLocal
	// if err := (*core.TxPool).AddLocal(signedTx); err != nil {
	// 	return common.Hash{}, err
	// }
	// qiqi-todo:增加判断。如果tx.To()==nil，那么为创建合约的交易
	log.Info("提交交易至txpool", "txhash", signedTx.Hash().Hex(), "交易接收方", signedTx.To())
	return signedTx.Hash(), nil
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
// SendTxArgs 表示将新事务提交到事务池的参数。
type SendTxArgs struct {
	From      common.Address  `json:"from"`
	To        *common.Address `json:"to"`
	GasLimit  *hexutil.Big    `json:"gas"` //gas的amount
	GasPrice  *hexutil.Big    `json:"gasPrice"`
	Value     *hexutil.Big    `json:"value"`
	Timestamp *hexutil.Big    `json:"timeStamp"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	// qiqi-todo:创建合约用到的两个参数，目前不需要
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`
}

//setDefaults 设置一些健全性默认值并在失败时终止
func (args *SendTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.GasLimit == nil {
		gasLimit := big.NewInt(90000)
		args.GasLimit = (*hexutil.Big)(gasLimit)
	}
	//gasprice参数，若没有传，设置为固定值
	if args.GasPrice == nil {
		price := utils.GasPrice
		args.GasPrice = (*hexutil.Big)(price)
	}
	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}
	//时间戳参数，若没有传，则置为系统当前时间
	if args.Timestamp == nil {
		time := big.NewInt(int64(time.Now().Unix()))
		args.Timestamp = (*hexutil.Big)(time)
	}
	// qiqi-todo：创建合约用到的参数data和input的判断
	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return errors.New(`Both "data" and "input" are set and not equal. Please use "input" to pass transaction call data.`)
	}
	if args.To == nil {
		// Contract creation
		var input []byte
		if args.Data != nil {
			input = *args.Data
		} else if args.Input != nil {
			input = *args.Input
		}
		if len(input) == 0 {
			return errors.New(`contract creation without any data provided`)
		}
	}
	return nil
}

// 根据传入的tx的参数，构造交易
func (args *SendTxArgs) toTransaction() *basic.Transaction {

	return basic.NewTransaction((*big.Int)(args.Timestamp), *args.To, (*big.Int)(args.Value), (*big.Int)(args.GasLimit), *args.Data)
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PrivateAccountAPI) NewAccount(password string) (common.Address, error) {
	acc, err := fetchKeystore(s.am).NewAccount(password)
	if err == nil {
		return acc.Address, nil
	}
	return common.Address{}, err
}

// fetchKeystore retrives the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) *keystore.KeyStore {
	return am.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

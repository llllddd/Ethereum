package keystore

import (
	"math/big"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"xchain-go/accounts"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/crypto"
	"xchain-go/event"
)

//对硬盘上加密密钥文件的管理
type KeyStore struct {
	storage keyStore      //对存储文件操作的接口
	cache   *accountCache //内存文件系统中的账户缓存

	changes  chan struct{}                //从账户缓存中接受账户变化的通知
	unlocked map[common.Address]*unlocked //当前已解锁的账户

	wallets     []accounts.Wallet
	updateFeed  event.Feed
	updateScope event.SubscriptionScope

	mu sync.RWMutex
}

type unlocked struct {
	*Key
	abort chan struct{}
}

//NewKeyStore 在给定目录下创建一个新的keystore文件
func NewKeyStore(keydir string, scryptN, scryptP int) *KeyStore {
	keydir, _ = filepath.Abs(keydir)
	ks := &KeyStore{storage: &keyStorePassphrase{keydir, scryptN, scryptP}}
	ks.init(keydir)
	return ks
}

func NewPlaintextKeyStore(keydir string) *KeyStore {
	keydir, _ = filepath.Abs(keydir)
	ks := &KeyStore{storage: &keyStorePlain{keydir}}
	ks.init(keydir)
	return ks
}

//init 函数初始化KeyStore
func (ks *KeyStore) init(keydir string) {
	//锁定互斥锁，因为帐户缓存可能会调用事件
	ks.mu.Lock()
	defer ks.mu.Unlock()
	//初始化
	ks.unlocked = make(map[common.Address]*unlocked)
	ks.cache, ks.changes = newAccountCache(keydir)

	runtime.SetFinalizer(ks, func(m *KeyStore) {
		m.cache.close()
	})
	//从账户缓存中创建初始的钱包列表
	accs := ks.cache.accounts()
	ks.wallets = make([]accounts.Wallet, len(accs))
	for i := 0; i < len(accs); i++ {
		ks.wallets[i] = &keystoreWallet{account: accs[i], keystore: ks}
	}
}

//Wallets 实现了accounts.Backend的接口,返回所有的密钥文件.
func (ks *KeyStore) Wallets() []accounts.Wallet {
	ks.refreshWallets()

	ks.mu.RLock()
	defer ks.mu.RUnlock()

	cpy := make([]accounts.Wallet, len(ks.wallets))
	copy(cpy, ks.wallets)
	return cpy
}

//refreshWallets 检索当前帐户列表并基于此做任何必要的钱包刷新操作。
func (ks *KeyStore) refreshWallets() {
	//检索当前的账户列表
	ks.cache.accounts()
	accs := ks.cache.accounts()
	// 将当前的账户列表转换为一个新的列表
	wallets := make([]accounts.Wallet, 0, len(accs))
	events := []accounts.WalletEvent{}

	for _, account := range accs {
		//在下一个账户前丢掉钱包
		for len(ks.wallets) > 0 && ks.wallets[0].URL().Cmp(account.URL) < 0 {
			events = append(events, accounts.WalletEvent{Wallet: wallet, Kind: accounts.WalletArrived})
			ks.wallets = ks.wallets[1:]
		}
		if len(ks.wallets) == 0 || ks.wallets[0].URL().Cmp(account.URL) > 0 {
			//TODO:keystoreWallet
			wallet := &keystoreWallet{account: account, keystore: ks}
			events = append(events, accounts.WalletEvent{Wallet: wallet, Kind: accounts.WalletArrived})
			wallets = append(wallets, wallet)
			continue
		}
		//如果账户和前一个钱包的一样保留
		if ks.wallets[0].Accounts()[0] == account {
			wallets = append(wallets, ks.wallets[0])
			ks.wallets = ks.wallets[1:]
			continue
		}
	}
	//
	for _, wallet := range ks.wallets {
		events = append(events, accounts.WalletEvent{Wallet: wallet, Kind: accounts.WalletDropped})
	}
	ks.wallets = wallets
	ks.mu.Unlock()
	//
	for _, event := range events {
		ks.updateFeed.Send(event)
	}
}

// Subscribe 实现了accounts.Backend的接口,生成了一个异步的订阅来通知钱包的增减状况
func (ks *KeyStore) Subscribe(sink chan<- accounts.WalletEvent) event.Subscription {
	//使用锁来使安全的更新循环
	ks.mu.Lock()
	defer ks.mu.Unlock()
	//订阅调用者并且追踪订阅者
	sub := ks.updateScope.Track(ks.updateFeed.Subscribe(sink))

	//订阅服务器需要一个激活的通知循环
	if !ks.updating {
		ks.updating = true
		go ks.updater()
	}
	return sub
}

//更新器负责维护存储在最新的钱包列表。密钥存储库，用于启动钱包添加/删除事件。它倾听

//来自基础帐户缓存的帐户更改事件，也周期性地强制手动刷新（仅对文件系统通知器的系统进行触发器）

//不运行。
func (ks *KeyStore) updater() {
	for {
		// 等待账户更新或者刷新超时
		select {
		case <-ks.changes:
		case <-time.After(walletRefreCycle):
		}
		// 运行钱包刷新器
		ks.refreshWallets()
		// 如果订阅者离开,停止更新
		ks.mu.Lock()
		if ks.updateScope.Count() == 0 {
			ks.updating = false
			ks.mu.Unlock()
			return
		}
		ks.mu.Unlock()
	}
}

func (ks *KeyStore) SignHash(a accounts.Account, hash []byte) ([]byte, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	unlockedKey, found := ks.unlocked[a.Address]
	if !found {
		return nil, Errlocked
	}
	return crypto.Sign(hash, unlocked.PrivateKey)
}

func (ks *KeyStore) SignTx(a accounts.Account, tx *basic.Transaction, chainId *big.Int) (*basic.Transaction, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	unlockedKey, found := ks.unlocked[a.Address]
	if !found {
		return nil, ErrLocked
	}
	return basic.SignTx(tx, basic.NewEIP155Signer(chainId), unlockedKey.PrivateKey)
}

func (ks *KeyStore) SignHashWithPassphrase(a accounts.Account, passphrase string, hash []byte) (signature []byte, err error) {
	_, key, err := ks.getDecryptedKey(a, passphrase)
	if err != nil {
		return nil, err
	}
	defer zeroKey(key.PrivateKey)
	return crypto.Sign(hash, key.PrivateKey)
}

func (ks *KeyStor) SignTxWithPassphrase(a accounts.Account, passphrase string, tx *basic.Transaction, chainID *big.Int) (*basic.Transaction, error) {
	_, key, err := ks.getDecryptedKey(a, passphrase)
	if err != nil {
		return nil, err
	}
	defer zeroKey(key.PrivateKey)

	return basic.SignTx(tx, basic.NewEIP155Signer(chainID), key.PrivateKey)
}

func (ks *KeyStore) getDecryptedKey(a accounts.Account, auth string) (accounts.Account, *Key, error) {
	a, err := ks.Find(a)
	if err != nil {
		return a, nil, err
	}
	key, err := ks.storage.GetKey(a.Address, a.URL.Path, auth)
	return a, key, err
}

func (ks *KeyStore) Find(a accounts.Account) (accounts.Account, error) {
	ks.cache.maybeReload()
	ks.cache.mu.Lock()
	a, err := ks.cache.find(a)
	ks.cache.mu.Unlock()
	return a, err
}

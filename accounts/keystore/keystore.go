package keystore

import (
	"path/filepath"
	"runtime"
	"sync"
	"xchain-go/accounts"
	"xchain-go/common"
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
		// acc
	}
}

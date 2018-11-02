package accounts

import (
	"errors"
	"reflect"
	"sort"
	"sync"
	"xchain-go/event"
)

var (
	ErrUnknownAccount = errors.New("未知账户错误")
)

// Manager 管理不同钱包中的账户
type Manager struct {
	backends map[reflect.Type][]Backend //Backend实现的功能实际上就是管理本地拥有的账户文件,并进行消息传递,进而根据传递来的信息由manager进行管理.
	updaters []event.Subscription       //订阅来自backend的信息.
	updates  chan WalletEvent
	wallets  []Wallet //缓存所有已知backend的钱包.

	feed event.Feed

	quit chan chan error
	lock sync.RWMutex
}

//NewManager 生成一个账户管理的Manager来通过钱包进行签名
func NewManager(backends ...Backend) *Manager {

	//从已知的Bckend中检索钱包信息并按照URL进行排序
	var wallets []Wallet
	for _, backend := range backends {
		wallets = merge(wallets, backend.Wallets()...)
	}
	//订阅钱包的变化的信息
	updates := make(chan WalletEvent, 4*len(backends))
	subs := make([]event.Subscription, len(backends))
	for i, backend := range backends {
		subs[i] = backend.Subscribe(updates)
	}
	//组装manager结构体并返回
	am := &Manager{
		backends: make(map[reflect.Type][]Backend),
		updaters: subs,
		updates:  updates,
		wallets:  wallets,
		quit:     make(chan chan error),
	}
	for _, backend := range backends {
		kind := reflect.TypeOf(backend)
		am.backends[kind] = append(am.backends[kind], backend)
	}
	go am.update()

	return am
}

//Wallets 返回管理器下的所有可以用来签名的钱包
func (am *Manager) Wallets() []Wallet {
	am.lock.RLock()
	defer am.lock.RUnlock()
	wallet := make([]Wallet, len(am.wallets))
	copy(wallet, am.wallets)
	return wallet
}

//Wallet由指定URL返回对应的钱包
func (am *Manager) Wallet(url string) (Wallet, error) {
	am.lock.Lock()
	defer am.lock.Unlock()

	parsed, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	for _, wallet := range am.Wallets() {
		if wallet.URL() == parsed {
			return wallet, nil
		}
	}
	return nil, ErrUnknownAccount
}

//Subscribe 订阅钱包事件
func (am *Manager) Subscribe(sink chan<- WalletEvent) event.Subscription {
	return am.feed.Subscribe(sink)
}

//Close 关闭管理器
func (am *Manager) Close() error {
	errc := make(chan error)
	am.quit <- errc
	return <-errc
}

//Find 从钱包列表中返回包含特定账户的钱包来进行签名
func (am *Manager) Find(account Account) (Wallet, error) {
	am.lock.RLock()
	defer am.lock.RUnlock()

	for _, wallet := range am.wallets {
		if wallet.Contains(account) {
			return wallet, nil
		}
	}
	return nil, ErrUnknownAccount
}

// Backends retrieves the backend(s) with the given type from the account manager.
func (am *Manager) Backends(kind reflect.Type) []Backend {
	return am.backends[kind]

}

//update 循环检测backends中钱包的信息及时更新钱包缓存
func (am *Manager) update() {
	defer func() {
		am.lock.Lock()
		for _, sub := range am.updaters {
			sub.Unsubscribe()
		}
		am.updaters = nil
		am.lock.Unlock()
	}()
	//开始循环监听
	for {
		select {
		case event := <-am.updates:
			am.lock.Lock()
			switch event.Kind {
			case WalletArrived:
				am.wallets = merge(am.wallets, event.Wallet)
			case WalletDropped:
				am.wallets = drop(am.wallets, event.Wallet)
			}
			am.lock.Unlock()

			am.feed.Send(event)
		case errc := <-am.quit:
			errc <- nil
			return
		}
	}
}

//merge 函数将钱包合并成一个列表
func merge(slice []Wallet, wallets ...Wallet) []Wallet {
	for _, wallet := range wallets {
		n := sort.Search(len(slice), func(i int) bool { return slice[i].URL().Cmp(wallet.URL()) >= 0 })
		if n == len(slice) {
			slice = append(slice, wallet)
			continue
		}
		slice = append(slice[:n], append([]Wallet{wallet}, slice[n:]...)...)
	}
	return slice
}

func drop(slice []Wallet, wallets ...Wallet) []Wallet {
	for _, wallet := range wallets {
		n := sort.Search(len(slice), func(i int) bool { return slice[i].URL().Cmp(wallet.URL()) >= 0 })
		if n == len(slice) || slice[n].URL().Cmp(wallet.URL()) != 0 {
			continue
		}
		slice = append(slice[:n], slice[n+1:]...)
	}
	return slice
}

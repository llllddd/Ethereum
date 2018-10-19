package keystore

import (
	"time"

	mylog "mylog2"

	"github.com/rjeczalik/notify"
)

type watcher struct {
	ac       *accountCache
	starting bool
	running  bool
	ev       chan notify.EventInfo
	quit     chan struct{}
}

func newWatcher(ac *accountCache) *watcher {
	return &watcher{
		ac:   ac,
		ev:   make(chan notify.EventInfo, 10),
		quit: make(chan struct{}),
	}
}

//在后台启动观察者循环.调用者必须由锁
func (w *watcher) start() {
	if w.starting || w.running {
		return
	}
	w.starting = true
	go w.loop()
}

func (w *watcher) loop() {
	defer func() {
		w.ac.mu.Lock()
		w.running = false
		w.starting = false
		w.ac.mu.Unlock()
	}()
	logger := mylog.NewLogger(w.ac.keydir)

	if err := notify.Watch(w.ac.keydir, w.ev, notify.All); err != nil {
		logger.Infoln("监控密钥存储文件夹失败")
		return
	}
	defer notify.Stop(w.ev)
	logger.Infoln("开始监控密钥存储文件夹")
	defer logger.Infoln("停止监控密钥存储文件夹")

	w.ac.mu.Lock()
	w.running = true
	w.ac.mu.Unlock()

	//等待文件系统的事件并且重载,当事件产生时,重载会稍微滞后因此当有许多事件到达时仅导致单个重新加载
	var (
		debounceDuration = 500 * time.Millisecond
		rescanTriggered  = false
		debounce         = time.NewTimer(0)
	)
	//忽略初始触发器
	if !debounce.Stop() {
		<-debounce.C
	}
	defer debounce.Stop()
	for {
		select {
		case <-w.quit:
			return
		case <-w.ev:
			if !rescanTriggered {
				debounce.Reset(debounceDuration)
				rescanTriggered = true
			}
		case <-debounce.C:
			w.ac.scanAccounts()
			rescanTriggered = false
		}
	}
}

func (w *wathcer) close() {
	close(w.quit)
}

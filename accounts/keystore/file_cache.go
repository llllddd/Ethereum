package keystore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set"
	log "github.com/inconshreveable/log15"
)

//fileCache是在密钥库扫描过程中得到的文件的缓存。
type fileCache struct {
	all     mapset.Set //密钥存储文件夹下所有文件的Set
	lastMod time.Time  //最后一次文件被修改的时间
	mu      sync.RWMutex
}

//scan 对给定的目录执行新的扫描,与已有的进行比较,缓存文件名并返回文件集：创建删除更新.
func (fc *fileCache) scan(keyDir string) (mapset.Set, mapset.Set, mapset.Set, error) {
	t0 := time.Now()

	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return nil, nil, nil, err
	}
	t1 := time.Now()

	fc.mu.Lock()
	defer fc.mu.Unlock()

	all := mapset.NewThreadUnsafeSet()
	mods := mapset.NewThreadUnsafeSet()

	var newLastMod time.Time
	for _, fi := range files {
		path := filepath.Join(keyDir, fi.Name())
		if nonKeyFile(fi) {
			log.Info("当前文件夹下没有密钥存储文件", "path", path)
			continue
		}
		all.Add(path)

		modified := fi.ModTime()
		if modified.After(fc.lastMod) {
			mods.Add(path)
		}
		if modified.After(newLastMod) {
			newLastMod = modified
		}
	}
	t2 := time.Now()

	deletes := fc.all.Difference(all)

	//更新追踪的文件返回文件集合
	deletes = fc.all.Difference(all)
	creates := all.Difference(fc.all)
	updates := mods.Difference(creates)

	fc.all, fc.lastMod = all, newLastMod
	t3 := time.Now()

	log.Debug("FS scan times", "list", t1.Sub(t0), "set", t2.Sub(t1), "diff", t3.Sub(t2))
	return creates, deletes, updates, nil
}

//nonKeyFile 忽视备份文件,隐藏文件和链接文件.
func nonKeyFile(fi os.FileInfo) bool {
	//跳过隐藏文件和备份文件
	if strings.HasSuffix(fi.Name(), "~") || strings.HasPrefix(fi.Name(), ".") {
		return true
	}
	//跳过特殊文件,目录和系统链接
	if fi.IsDir() || fi.Mode()&os.ModeType != 0 {
		return true
	}
	return false
}

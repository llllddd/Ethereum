package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"xchain-go/common"
)

type keyStorePlain struct {
	keysDirPath string
}

//GetKey 从文件中得到对应账户地址的密钥
func (ks keyStorePlain) GetKey(addr common.Address, filename, auth string) (*Key, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	key := new(Key)
	if err := json.NewDecoder(fd).Decode(key); err != nil {
		return nil, err
	}
	if key.Address != addr {
		return nil, fmt.Errorf("密钥和账户匹配错误: 错误地址 %x, 正确地址 %x", key.Address, addr)
	}
	return key, nil
}

//StoreKey 将密钥存储到对应文件中
func (ks keyStorePlain) StoreKey(filename string, key *Key, auth string) error {
	content, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return writeKeyFile(filename, content)
}

//JoinPath 连接文件名和密钥文件目录
func (ks keyStorePlain) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(ks.keysDirPath, filename)
}

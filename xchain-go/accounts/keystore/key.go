package keystore

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"xchain-go/accounts"
	"xchain-go/common"
	"xchain-go/crypto"

	"github.com/pborman/uuid"
)

const (
	version        = 1
	KeyStoreScheme = "XChain"
)

type Key struct {
	Id uuid.UUID //通用唯一识别码,作为密钥的ID

	Address common.Address //地址可以由私钥先导出公钥再由公钥得到

	PrivateKey *ecdsa.PrivateKey
}

//对密钥文件的基本管理功能,加密密钥的存储设施的功能接口
type keyStore interface {
	//从硬盘中加载并解密密钥文件
	GetKey(addr common.Address, filename string, auth string) (*Key, error)
	//加密并存储密钥文件
	StoreKey(filename string, k *Key, auth string) error
	//将文件名与密钥目录连接,除非已经是绝对路径
	JoinPath(filename string) string
}

//预留version字段用来以后替换抗量子攻击的签名
type plainKeyJSON struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privatekey"`
	Id         string `json:"id"`
	Version    int    `json:"version"`
}

type encryptedKeyJSON struct {
	Address string     `json:"address"`
	Crypto  cryptoJSON `json:"crypto"`
	Id      string     `json:"id"`
	Version int        `json:"version"`
}

type cryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"` //私钥加密之后的密文
	CipherParams cipherparamsJSON       `json:cipherparams`
	KDF          string                 `json:"kdf"` //密钥导出函数
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

//MarshalJSON 对密钥文件进行JSON编码
func (k *Key) MarshalJSON() ([]byte, error) {
	jStruct := plainKeyJSON{
		hex.EncodeToString(k.Address[:]),
		hex.EncodeToString(crypto.FromECDSA(k.PrivateKey)),
		k.Id.String(),
		version,
	}
	j, err := json.Marshal(jStruct)
	return j, err
}

//UnmarshalJSON 对密钥文件进行JSON解码
func (k *Key) UnmarshalJSON(j []byte) error {
	keyJSON := new(plainKeyJSON)
	err := json.Unmarshal(j, &keyJSON)
	if err != nil {
		return err
	}
	u := new(uuid.UUID)
	*u = uuid.Parse(keyJSON.Id)
	k.Id = *u
	addr, err := hex.DecodeString(keyJSON.Address)
	if err != nil {
		return err
	}
	privkey, err := crypto.HexToECDSA(keyJSON.PrivateKey)
	if err != nil {
		return err
	}
	k.Address = common.BytesToAddress(addr)
	k.PrivateKey = privkey
	return nil
}

func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *Key {
	id := uuid.NewRandom()
	key := &Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}

//newKey 使用椭圆曲线加密算法生成一个新的私钥
func newKey(rand io.Reader) (*Key, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return newKeyFromECDSA(privateKeyECDSA), nil
}

//storeNewKey 由用户口令加密生成密钥文件
func storeNewKey(ks keyStore, rand io.Reader, auth string) (*Key, accounts.Account, error) {
	key, err := newKey(rand)
	if err != nil {
		return nil, accounts.Account{}, err
	}
	a := accounts.Account{Address: key.Address, URL: accounts.URL{Scheme: KeyStoreScheme, Path: ks.JoinPath(keyFileName(key.Address))}}
	if err := ks.StoreKey(a.URL.Path, key, auth); err != nil {
		zeroKey(key.PrivateKey)
		return nil, a, err
	}
	return key, a, err
}

/*
存储逻辑: 1.新建存储目录
          2. 创建临时存储文件
          3. 写入内容
*/
func writeTemporaryKeyFile(file string, content []byte) (string, error) {

	const dirPerm = 0700

	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}

	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", nil
	}

	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil

}

//writeKeyFile 将密钥存储在文件中
func writeKeyFile(file string, content []byte) error {
	name, err := writeTemporaryKeyFile(file, content)
	if err != nil {
		return err
	}
	return os.Rename(name, file)
}

//keyFileName 定义密钥文件名称
func keyFileName(keyAddr common.Address) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", toISO8601(ts), hex.EncodeToString(keyAddr[:]))
}

func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

func zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}

/*
func NewKeyForDirectICAP(rand io.Reader) *Key {
	randBytes := make([]byte, 64)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic("key generation: could not read from random source: " + err.Error())
	}
	reader := bytes.NewReader(randBytes)
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), reader)
	if err != nil {
		panic("key generation: ecdsa.GenerateKey failed: " + err.Error())
	}
	key := newKeyFromECDSA(privateKeyECDSA)
	if !strings.HasPrefix(key.Address.Hex(), "0x00") {
		return NewKeyForDirectICAP(rand)
	}
	return key
}
*/

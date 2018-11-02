package keystore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"xchain-go/common"
	"xchain-go/common/math"
	"xchain-go/crypto"

	"github.com/pborman/uuid"
	"golang.org/x/crypto/scrypt"
)

var (
	ErrDecrypt = errors.New("解密密钥文件出错")
)

const (
	keyHeaderKDF   = "scrypt"
	StandardScrypN = 1 << 18
	StandardScrypP = 1
	scryptR        = 8
	scryptDKLen    = 32
)

type keyStorePassphrase struct {
	keysDirPath string
	//密钥导出算法的参数
	scryptN int
	scryptP int
}

/*
1.从密钥存储文件中获取内容
2.将密钥解密
3.验证账户地址是否一致
*/
func (ks keyStorePassphrase) GetKey(addr common.Address, filename, auth string) (*Key, error) {
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	key, err := DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}

	if key.Address != addr {
		return nil, fmt.Errorf("密钥和账户匹配错误: 错误地址 %x, 正确地址 %x", key.Address, addr)
	}
	return key, nil
}

func StoreKey(dir, auth string, scryptN, scryptP int) (common.Address, error) {
	_, a, err := storeNewKey(&keyStorePassphrase{dir, scryptN, scryptP}, rand.Reader, auth)
	if err != nil {
		return common.Address{}, err
	}
	return a.Address, nil
}

/*
存储私钥文件:1.用给定口令加密私钥
            2.生成临时存储文件
            3.写入加密后的私钥并存储
*/

func (ks keyStorePassphrase) StoreKey(filename string, key *Key, auth string) error {
	keyjson, err := EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	tmpName, err := writeTemporaryKeyFile(filename, keyjson)
	if err != nil {
		return err
	}
	return os.Rename(tmpName, filename)
}

func (ks keyStorePassphrase) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(ks.keysDirPath, filename)

}

/*
 加密私钥: 1.Scrypt算法得到AES加密私钥
           2.Hash得到mac
           3.使用AES进行加密
           4.将密文和相关参数编码为JSON
*/
func EncryptKey(key *Key, auth string, scryptN, scryptP int) ([]byte, error) {
	authArray := []byte(auth)

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic("从crypto/rand读取随机数出错: " + err.Error())
	}
	//Scrypt算法由参数N,P,Salt得到AES加密密钥
	deriveKey, err := scrypt.Key(authArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	encryptKey := deriveKey[:16]
	keyBytes := math.PaddedBigBytes(key.PrivateKey.D, 32)

	//采用AES的CTR模式进行加密
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic("从crypto/rand读取随机数出错: " + err.Error())
	}
	cipherText, err := aesCTRXOR(encryptKey, keyBytes, iv)
	if err != nil {
		return nil, err
	}
	//计算AES加密私钥的Mac
	mac := crypto.Keccak256(deriveKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherparamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}
	cryptoStruct := cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherparamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	encryptedKeyJSON := encryptedKeyJSON{
		hex.EncodeToString(key.Address[:]),
		cryptoStruct,
		key.Id.String(),
		version,
	}

	return json.Marshal(encryptedKeyJSON)
}

/*解密得到密钥
  1.将密钥Json解码
  2.根据用户口令解密
*/

func DecryptKey(keyjson []byte, auth string) (*Key, error) {
	//Json文件解码为一个map
	m := make(map[string]interface{})
	if err := json.Unmarshal(keyjson, &m); err != nil {
		return nil, err
	}
	var (
		keyBytes []byte
		err      error
	)
	k := new(encryptedKeyJSON)
	if err := json.Unmarshal(keyjson, &k); err != nil {
		return nil, err
	}
	keyBytes, keyId, err := decryptKey(k, auth)
	if err != nil {
		return nil, err
	}
	key := crypto.ToECDSAUnsafe(keyBytes)

	return &Key{
		Id:         uuid.UUID(keyId),
		Address:    crypto.PubkeyToAddress(key.PublicKey),
		PrivateKey: key,
	}, nil
}

func decryptKey(keyProtected *encryptedKeyJSON, auth string) (keyBytes []byte, keyId []byte, err error) {
	if keyProtected.Crypto.Cipher != "aes-128-ctr" {
		return nil, nil, fmt.Errorf("不支持的密文类型: %v", keyProtected.Crypto.Cipher)
	}

	keyId = uuid.Parse(keyProtected.Id)
	mac, err := hex.DecodeString(keyProtected.Crypto.MAC)
	if err != nil {
		return nil, nil, err
	}
	iv, err := hex.DecodeString(keyProtected.Crypto.CipherParams.IV)
	if err != nil {
		return nil, nil, err
	}

	cipherText, err := hex.DecodeString(keyProtected.Crypto.CipherText)
	if err != nil {
		return nil, nil, err
	}

	derivedKey, err := getKDFKey(keyProtected.Crypto, auth)
	if err != nil {
		return nil, nil, err
	}

	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, nil, ErrDecrypt
	}
	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, nil, err
	}
	return plainText, keyId, err
}

//getKDFKEY 由密钥导出函数,通过相关参数和用户口令得到AES加密的密钥.
func getKDFKey(cryptoJSON cryptoJSON, auth string) ([]byte, error) {
	authArray := []byte(auth)
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}

	if cryptoJSON.KDF == keyHeaderKDF {
		// dkLen := cryptoJSON.KDFParams["dklen"].(int)
		// n := cryptoJSON.KDFParams["n"].(int)
		// p := cryptoJSON.KDFParams["p"].(int)
		// r := cryptoJSON.KDFParams["r"].(int)
		// return scrypt.Key(authArray, salt, n, r, p, dkLen)
		dkLen := int(cryptoJSON.KDFParams["dklen"].(float64))
		n := int(cryptoJSON.KDFParams["n"].(float64))
		p := int(cryptoJSON.KDFParams["p"].(float64))
		r := int(cryptoJSON.KDFParams["r"].(float64))
		keybytes, _ := scrypt.Key(authArray, salt, n, r, p, dkLen)
		return keybytes, nil
	}
	return nil, fmt.Errorf("不支持的密钥导出函数: %s", cryptoJSON.KDF)
}

//aesCTRXOR 使用AES-128 CTR 模式对密文解密得到私钥.
func aesCTRXOR(aeskey, ciphertext, iv []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(ciphertext))
	stream.XORKeyStream(outText, ciphertext)
	return outText, nil
}

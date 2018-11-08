package crypto

import (
	"PQC/dilithium"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"io"
	"io/ioutil"
	"os"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h common.Hash) {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

// Keccak512 calculates and returns the Keccak512 hash of the input data.
func Keccak512(data ...[]byte) []byte {
	d := sha3.NewKeccak512()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

//LoadDilithum 从指定文件中加载一个Dilithium私钥
func LoadDilithum(file string) (dilithium.PrivateKey, error) {
	buf := make([]byte, 896)
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	if _, err := io.ReadFull(fd, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

//SaveDilithium 将生成的Dilithium私钥保存到一个给定的文件中
func SaveDilithium(file string, key dilithium.PrivateKey) error {
	return ioutil.WriteFile(string, key, 0600)
}

//GenerateKey 生成Dilithium算法私钥
func GenerateKey() (dilithium.PrivateKey, error) {
	_, sk := dilithium.GenerateSk()
	return sk, nil
}

func ValidateSignatureValues(sig []byte, pubkey dilithium.PublicKey) bool {
	if len(sig) != 1519 || len(pubkey) != 896 {
		return false
	}
	return true
}

//PubkeyToAddress 由Dilithium公钥通过Keccak256算法生成地址
func PubkeyToAddress(p dilithium.PublicKey) common.Address {
	return common.BytesToAddress(Keccak256(pubBytes[:])[12:])
}

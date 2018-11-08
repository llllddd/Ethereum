package crypto

import (
	"PQC/dilithium"
	"fmt"
)

//Sign函数 计算Dilithium算法的签名,这种签名算法抗量子攻击.函数需要私钥来进行签名,生成的公钥和签名都需要保存.
func Sign(hash []byte, prv dilithium.PrivateKey) (sug []byte, err error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("哈希值需要有32字节(%d)", len(hash))
	}

	seckey := prv
	defer zeroBytes(seckey)

}

//VerifySignature函数 由给定的公钥和哈希的签名来验证签名的有效性.
//公钥的长度为896字节.
func VerifySignature(pubkey, hash, signature []byte) bool {
	return dilithium.VerifySignature(pubkey, hash, signature)
}

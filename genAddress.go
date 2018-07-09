/*
	账户地址的生成：
	1.ECDSA算法生成 32bytes私钥 64bytes公钥(elliptic,x,y)
	2.对公钥[]byte进行hash运算得到地址 address = Keccak256(publickey)
	3.[]byte形式的公钥 04开头为(xy)形式公钥
	交易签名
	1.完整的交易的数据结构包含九个数据成员
	2.对交易进行RLP编码，生成序列化信息
	3.对编码后的信息进行hash运算Keccak256得到32bytes的摘要
	4.椭圆曲线加密算法对hash值进行签名算法，得到的签名格式为(R||S||V)，65字节V用来的得到公钥简化计算1bytes
	验证签名
	1.由交易的hash值和签名得到公钥pubkey
	2.pubkey和签名使用椭圆曲线加密算法对交易进行验证
*/
package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"io"
	"sdk/secp256k1"
	"sdk/sha3"
)

const (
	HashLength    = 32
	AddressLength = 20
)

type Address [AddressLength]byte
type Hash [HashLength]byte

//输入初始随机数生成私钥包含 加密的椭圆曲线，公钥， 私钥
func genAccountSek(rand io.Reader) (*ecdsa.PrivateKey, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		return nil, err
	}
	return privateKeyECDSA, nil
}

//直接生成并得到私钥信息
func genNormalSk(rand io.Reader) ([]byte, error) {
	privatekey, err := genAccountSek(rand)
	if err != nil {
		return nil, err
	} else {
		sk := privatekey.D.Bytes()
		return sk, nil
	}

}

//得到公钥信息
func getPubKey(sk *ecdsa.PrivateKey) *ecdsa.PublicKey {
	if sk == nil {
		return nil
	}
	return &sk.PublicKey
}

//得到字节串私钥信息
func getByteSk(sk *ecdsa.PrivateKey) []byte {
	sk1 := sk.D.Bytes()
	return sk1
}

//得到字节串公钥信息
func getBytePkFromSk(sk *ecdsa.PrivateKey) []byte {
	if sk == nil {
		return nil
	}
	pk := sk.PublicKey
	return elliptic.Marshal(secp256k1.S256(), pk.X, pk.Y)
}

//直接编组公钥
func getBytePk(pk *ecdsa.PublicKey) []byte {
	if pk == nil || pk.X == nil || pk.Y == nil {
		return nil
	}
	return elliptic.Marshal(secp256k1.S256(), pk.X, pk.Y)
}

//keaccak256哈希
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

//得到公钥的hash值
func genHashFromPub(pk *ecdsa.PublicKey) []byte {
	if pk == nil || pk.X == nil || pk.Y == nil {
		return nil
	}
	pubBytes := getBytePk(pk)
	tempBytes := Keccak256(pubBytes[1:])
	return tempBytes
}

//生成账户地址
func genAddressFromPub(pk *ecdsa.PublicKey) Address {
	temp := genHashFromPub(pk)
	var a Address
	if len(temp) > len(a) {
		temp = temp[len(temp)-AddressLength:]
	}
	copy(a[AddressLength-len(temp):], temp)
	return a
}

//字符串进制转换

/*
//对交易信息进行hash运算算法为Keccak256
func hashTxByKeccak256(tx *types.Transaction) ([]byte, error) {
	var v []byte
	v, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	} else {
		haTx := crypto.Keccak256(v)
		return haTx, nil
	}
}

//对交易的hash值用椭圆曲线算法签名
func signTxBySecp256k1(tx *types.Transaction, sk []byte) ([]byte, error) {
	Txhash, err := hashTxByKeccak256(tx)
	if err != nil {
		return nil, err
	}
	signedTx, err := secp256k1.Sign(Txhash, sk)
	if err != nil {
		return nil, err
	} else {
		return signedTx, nil
	}
}

//验证交易签名,签名为(R||S||V)65bytes
func verifySignedTx(txhash []byte, txsign []byte) bool {
	pubkey, err := secp256k1.RecoverPubkey(txhash, txsign)
	if err != nil || len(txhash) != 32 {
		return false
	}
	bb := secp256k1.VerifySignature(pubkey, txhash, txsign[0:64])
	return bb
}
*/
func main() {
	seed := []byte{110, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111}
	bytesBuffer := bytes.NewBuffer(seed)
	v, _ := genAccountSek(bytesBuffer)
	fmt.Println("PrivateKey结构体", v)
	pk := getPubKey(v)
	//pk1 := getBytePkFromSk(v)
	fmt.Println("公钥为: ", pk)
	address := genAddressFromPub(pk)
	fmt.Println("地址: ", address)
}

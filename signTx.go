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
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

//输入初始随机数生成私钥包含 加密的椭圆曲线，公钥， 私钥
func genAccountSek(rand io.Reader) (*ecdsa.PrivateKey, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return privateKeyECDSA, nil
}

func genNormalSk(rand io.Reader) ([]byte, error) {
	privatekey, err := genAccountSek(rand)
	if err != nil {
		return nil, err
	} else {
		sk := privatekey.D.Bytes()
		return sk, nil
	}

}

//得到私钥信息
func getNormalSk(sk *ecdsa.PrivateKey) []byte {
	sk1 := sk.D.Bytes()
	return sk1
}

//从公钥得到账户地址
func genAddressFromPub(pubkey ecdsa.PublicKey) common.Address {
	pubBytes := crypto.FromECDSAPub(&pubkey)
	tempBytes := crypto.Keccak256(pubBytes[1:])
	address := common.BytesToAddress(tempBytes[12:])
	return address
}

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

func main() {
	emptyTx := types.NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil)

	seed := []byte{111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111}
	bytesBuffer := bytes.NewBuffer(seed)
	v, _ := genAccountSek(bytesBuffer)
	fmt.Println("PrivateKey机构体", v)
	sk := getNormalSk(v)
	fmt.Println("私钥：", common.ToHex(sk))
	addr := genAddressFromPub(v.PublicKey)
	fmt.Println("地址：", addr)
	txhash, _ := hashTxByKeccak256(emptyTx)
	fmt.Println("交易hash值：", common.ToHex(txhash))
	signtx, _ := signTxBySecp256k1(emptyTx, sk)
	fmt.Println("交易的签名：", common.ToHex(signtx))
	b := verifySignedTx(txhash, signtx)
	fmt.Println("验证交易：", b)
}

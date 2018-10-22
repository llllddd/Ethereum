package basic

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"xchain-go/common"
	"xchain-go/crypto"
	//	"xchain-go/params"
)

//TODO: 定义关于ChainId相关的操作

//ChainId 对应不同的区块链网络类型
var (
	ChainID = big.NewInt(1)
)

var (
	ErrInvalidChainId = errors.New("错误的ChainId")
	ErrInvalidSig     = errors.New("无效的签名")
)

type sigCache struct {
	signer Signer
	from   common.Address
}

//SignTx 函数由给定的交易和私钥对交易进行签名
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

//Sender 函数返回交易对应的发送者地址
func Sender(signer Signer, tx *Transaction) (common.Address, error) {
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}
	addr, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.from.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

//Signer 定义了交易签名处理的接口
type Signer interface {
	Sender(tx *Transaction) (common.Address, error)

	SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error)

	Hash(tx *Transaction) common.Hash

	Equal(Signer) bool
}

type EIP155Signer struct {
	chainId, chainIdMul *big.Int
}

func NewEIP155Signer(chainId *big.Int) EIP155Signer {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return EIP155Signer{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}

//Equal 函数判段是否使用的是给定的签名算法
func (s EIP155Signer) Equal(s2 Signer) bool {
	eip155, ok := s2.(EIP155Signer)
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}

var big8 = big.NewInt(8)

//Sender 函数返回交易发送者的地址
func (s EIP155Signer) Sender(tx *Transaction) (common.Address, error) {
	//TODO:
	/*
		if tx.ChainId().Cmp(s.chainId) != 0 {
			return common.Address{}, ErrInvalidChainId
		}
	*/
	V := new(big.Int).Sub(tx.data.V, s.chainIdMul)
	V.Sub(V, big8)
	return recoverPlain(s.Hash(tx), tx.data.R, tx.data.S, V, true)
}

//SignatureVlues 函数将数字签名转换为 [R||S||V]格式,其中V等于0或1
func (s EIP155Signer) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	if len(sig) != 65 {
		panic(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}
	R = new(big.Int).SetBytes(sig[:32])
	S = new(big.Int).SetBytes(sig[32:64])
	V = new(big.Int).SetBytes([]byte{sig[64] + 27})

	if s.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, s.chainIdMul)
	}
	return R, S, V, nil
}

// Hash 函数计算交易的哈希值
func (s EIP155Signer) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.data.Timestamp,
		tx.data.FixPrice,
		tx.data.Recipient,
		tx.data.Amount,
		tx.data.Payload,
		s.chainId, uint(0), uint(0),
	})
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		return common.Address{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return common.Address{}, ErrInvalidSig
	}

	// 将big.Int类型的签名转换为[]byte
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// 由签名中恢复公钥
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}

	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

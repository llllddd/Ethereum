package basic

import (
	//	"container/heap"
	//	"errors"
	//	"xchain-go/hexutil"
	//	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"math/big"
	"sync/atomic"

	//	"time"
	"xchain-go/common"
	"xchain-go/rlp"
)

var (
	Price = big.NewInt(111) //定义交易中的固定交易费
)

type Transaction struct {
	data txdata //交易的主体
	//cache
	size atomic.Value //缓存交易的尺寸
	from atomic.Value //缓存交易的发送者地址
	hash atomic.Value //缓存交易的哈希值
}

type txdata struct {
	Timestamp *big.Int        //交易的时间戳
	FixPrice  *big.Int        //固定的交易费
	Recipient *common.Address //交易接收者的地址
	Amount    *big.Int        //交易的金额
	Payload   []byte          //交易中的数据
	//签名值
	V *big.Int
	R *big.Int
	S *big.Int

	Hash *common.Hash //交易的哈希值.
}

//TODO:关于ChainId相关的函数

//NewTransaction 函数由给定的参数创建一个新的交易
func NewTransaction(time *big.Int, to common.Address, amount *big.Int, data []byte) *Transaction {
	return newTransaction(time, &to, amount, data)
}

func newTransaction(time *big.Int, to *common.Address, amount *big.Int, data []byte) *Transaction {
	d := txdata{
		Timestamp: time,
		FixPrice:  Price,
		Recipient: to,
		Amount:    new(big.Int),
		Payload:   data,

		V: new(big.Int),
		R: new(big.Int),
		S: new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	return &Transaction{data: d}
}

//EncodeRLP 函数返回此交易的RLP编码
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.data)
}

//DecodeRLP 函数将RLP编码值解码为交易类型并存储交易的尺寸
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&tx.data)
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

func (tx *Transaction) Data() []byte        { return tx.data.Payload }
func (tx *Transaction) Timestamp() *big.Int { return new(big.Int).Set(tx.data.Timestamp) }
func (tx *Transaction) Value() *big.Int     { return new(big.Int).Set(tx.data.Amount) }
func (tx *Transaction) To() *common.Address {
	if tx.data.Recipient == nil {
		return nil
	}
	to := *tx.data.Recipient
	return &to
}

//Hash 函数计算交易的RLP编码哈希值,并存储.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

//Size 函数返回交易的存储大小
func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

//RawSignature 函数返回交易的签名
func (tx *Transaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return tx.data.V, tx.data.R, tx.data.S
}

//WithSignature 函数将签名转换为[R||S||V]形式存储到交易中
func (tx *Transaction) WithSignature(signer Signer, sign []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sign)
	if err != nil {
		return nil, err
	}
	cpy := &Transaction{data: tx.data}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	return cpy, nil
}

type Transactions []*Transaction

func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

//定义交易列表别名在TXPOOL中使用
type TransactionsByTimestamp Transactions

func (t TransactionsByTimestamp) Len() int { return len(t) }
func (t TransactionsByTimestamp) Less(i, j int) bool {
	return t[i].data.Timestamp.Cmp(t[j].data.Timestamp) < 0
}

func (t TransactionsByTimestamp) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

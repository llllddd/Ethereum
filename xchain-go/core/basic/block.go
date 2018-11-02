package basic

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"time"
	"xchain-go/common"
	"xchain-go/crypto"
	"xchain-go/crypto/sha3"
	"xchain-go/ethdb"
	"xchain-go/rlp"
)

var (
	EmptyRootHash = DeriveSha(Transactions{})
)

const AddressLength = 20

const HashLength = 32

type Header struct {
	ParentHash  common.Hash    //父common.Hash
	Timestamp   *big.Int       //区块产生的时间戳
	BlockHash   common.Hash    //区块common.Hash
	Number      *big.Int       //区块号
	Extradata   []byte         //额外信息
	Validator   common.Address //区块验证者地址
	Root        common.Hash
	TxHash      common.Hash
	DposContext *DposContextProto
	//todo:difficulty\gaslimit\gasused\nonce\totaldifficulty
	//todo:uncle\stateroot\Txcommon.Hash\receipthash\bloom\

}

//body先预留一个字段
type Body struct {
	// Bodydata     []byte
	Transactions []*Transaction
}

type Block struct {
	header *Header
	// body         *Body
	transactions Transactions
	// caches
	hash atomic.Value
	//qiqi-todo:补充DposContext *DposContext

}
type Blocks []*Block

//计算Blockhash
//将block中的已有字段拼接起来，转为byte，通过keccak256算法得到Hash值
//组装的参数为每个参数的string类型：ParentHash，Timestamp，Number，Extradata
func (block *Block) SetHash() common.Hash {
	header := block.header
	record := string(header.Extradata) + header.Timestamp.String() + header.Number.String() + header.ParentHash.String()
	recordbyte := []byte(record)
	h := crypto.Keccak256Hash(recordbyte) //转为commonhash
	return h
}

func (block *Block) PrintBlockstruct() {
	// Log := log(prefix)
	fmt.Println("====================== Block", block.Number(), "======================")
	fmt.Println("ParentHash: ", block.ParentHash().String())
	fmt.Println("Timestamp : ", TimeFormat(block.header.Timestamp.Int64()))
	fmt.Println("BlockHash : ", block.Header().Hash().String())
	fmt.Println("Number    : ", block.Number())
	fmt.Println("Extradata : ", string(block.Extra()))
	fmt.Println("Validator : ", block.Validator().String())
	fmt.Println("Root      : ", block.Root().String())
	fmt.Println("Validator : ", block.TxHash().String())
	fmt.Println("Candidate : ", block.DposContext().CandidateHash.String())
	fmt.Println("Epoch     : ", block.DposContext().EpochHash.String())
	fmt.Println()
}

//时间转换：把unix时间戳format后再打印
func TimeFormat(unixtimestamp int64) string {
	time := time.Unix(unixtimestamp, 0)
	timeformat := time.Format("Mon Jan _2 15:04:05 MST 2006")
	return timeformat
}

//Newblock
//qiqi-todo：补充block中其他字段的构造--替换成这个函数
func NewBlock(header *Header, txs []*Transaction) *Block {
	// func NewBlock(header *Header, body *Body) *Block {

	//先根据出传入的参数，构造Block结构
	b := &Block{header: CopyHeader(header)}
	// b.body = body
	//判断传入的交易列表中的交易数量，如果为空，则将txhash置为空hash
	//否则计算tx列表的hash
	if len(txs) == 0 {
		b.header.TxHash = EmptyRootHash
	} else {
		//计算tx列表的hash
		b.header.TxHash = DeriveSha(Transactions(txs))
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)

	}
	return b

}

//处理header,防止修改头变量出现其他影响
func CopyHeader(header *Header) *Header {
	copyhead := *header
	if copyhead.Timestamp = new(big.Int); header.Timestamp != nil {
		copyhead.Timestamp.Set(header.Timestamp)
	}
	if copyhead.Number = new(big.Int); header.Number != nil {
		copyhead.Number.Set(header.Number)
	}
	if len(copyhead.Extradata) > 0 {
		copyhead.Extradata = make([]byte, len(header.Extradata))
		copyhead.Extradata = header.Extradata
	}
	// add dposContextProto to header
	copyhead.DposContext = &DposContextProto{}
	if header.DposContext != nil {
		copyhead.DposContext = header.DposContext
	}
	return &copyhead
}

// NewBlockWithHeader 返回一个已经给定header的block，且其他对于block的header的改动不会影响该block
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

func MockHash(x interface{}) common.Hash {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hash := common.BytesToHash(hw.Sum(nil))
	return hash
}
func (block *Block) Number() *big.Int          { return new(big.Int).Set(block.header.Number) }
func (block *Block) Time() *big.Int            { return new(big.Int).Set(block.header.Timestamp) }
func (block *Block) NumberU64() uint64         { return block.header.Number.Uint64() }
func (block *Block) Validator() common.Address { return block.header.Validator }
func (block *Block) Root() common.Hash         { return block.header.Root }
func (block *Block) ParentHash() common.Hash   { return block.header.ParentHash }
func (block *Block) TxHash() common.Hash       { return block.header.TxHash }
func (block *Block) Extra() []byte             { return common.CopyBytes(block.header.Extradata) }

// Body returns the non-header content of the block.
func (block *Block) Body() *Body     { return &Body{block.transactions} }
func (block *Block) Header() *Header { return CopyHeader(block.header) }
func (block *Block) DposContext() *DposContextProto {
	return &DposContextProto{EpochHash: block.header.DposContext.EpochHash, CandidateHash: block.header.DposContext.CandidateHash}
}

func (block *Block) Transactions() Transactions { return block.transactions }

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (block *Block) Hash() common.Hash {
	// if hash := block.blockhash.Load(); hash != nil {
	// 	return hash.(common.Hash)
	// }
	v := block.header.Hash()
	// block.blockhash.Store(v)
	return v
}

// qiqi-todo:补充：
// func (block *Block) DposCtx() *DposContext { return block.DposContext }

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	return rlpHash(h)
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// WithSeal返回一个包含b数据的新块，但header被的标题替换。
// 主要是用在打包区块后，进行区块数据组装时的调用--
func (block *Block) WithSeal(header *Header) *Block {
	copy := *header
	return &Block{
		header: &copy,
		//qiqi-todo-补充block中其他字段的内容
		// transactions: block.transactions,
		// DposContext: block.DposContext,
	}
}

// string形式返回block的内容
func (block *Block) String() string {
	str := fmt.Sprintf(`Block(#%v):
	header:
	{
		%v
	},
	body:
	{
		Transactions:
		%v
	}`, block.Number(), block.Header().String(), block.Transactions())
	return str
}

// string形式返回transactions的内容
func (header *Header) String() string {
	str := fmt.Sprintf(`Header(%x):
[
	ParentHash: %x
	Timestamp : %x
	BlockHash : %x
	Number    : %x
	Extradata : %x
	Root      : %x
	Validator : %x
	Candidate : %x
	Txhash    : %x
	Epoch     : %x
]`, header.Hash(), header.ParentHash, header.Timestamp, header.BlockHash, header.Number, header.Extradata, header.Root, header.Validator, header.DposContext.CandidateHash, header.TxHash, header.DposContext.EpochHash)
	return str
}

// WithBody返回包含了给定的transaction内容的block
// qiqi-todo:增加了transaction的block结构改造
func (block *Block) WithBody(transactions []*Transaction) *Block {
	//先构造一个block结构，header就是原block中的header。然后把参数中给定的transactions内容补充到block中
	//qiqi-todo:考虑还有没有其他的block中的字段作为参数传入的
	b := &Block{
		header:       CopyHeader(block.header),
		transactions: make([]*Transaction, len(transactions)),
	}
	copy(b.transactions, transactions)
	// for _, tx := range transactions {
	// 	Log.Debugln("tx.data:", tx.data)
	// }
	return b
}

//打开或新建一个数据库
func OpenDatabase(name string, cache, handles int) (ethdb.Database, error) {
	if name == "" {
		return ethdb.NewMemDatabase(), nil
	}
	return ethdb.NewLDBDatabase(name, cache, handles)
}

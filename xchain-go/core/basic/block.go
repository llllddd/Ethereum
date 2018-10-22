package basic

import (
	"fmt"
	"math/big"
	"time"
	"xchain-go/common"
	"xchain-go/crypto"
	"xchain-go/crypto/sha3"
	"xchain-go/ethdb"
	"xchain-go/rlp"
	"xchain-go/utils"

	mylog "mylog2"
)

const AddressLength = 20

const HashLength = 32

var prefix = "basic"

func log(prefix string) *mylog.SimpleLogger {
	common.InitLog(prefix)
	return common.Logger.NewSessionLogger()
}

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
	Bodydata []byte
}

type Block struct {
	header *Header
	// transactions Transactions
	body *Body
}

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

//生成新区块
//TODO：body中的tx等作为参数传入
func NewBlock(header *Header, body *Body) *Block {
	Log := log(prefix)

	b := &Block{header: CopyHeader(header), body: body}
	Log.Infoln("开始构造块")
	b.header.Timestamp = big.NewInt(time.Now().Unix())
	b.header.BlockHash = b.Header().Hash()

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
	return &copyhead
}

// NewBlockWithHeader 返回一个已经给定header的block，且其他对于block的header的改动不会影响该block
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// WithBody返回一个给定body内容的block
func (block *Block) WithBody(blockbody *Body) *Block {
	b := &Block{
		header: CopyHeader(block.header),
		body:   blockbody,
	}

	return b
}

//genesisblock
//todo，通过文件形式传入genesisblock
func NewGenesisBlock() *Block {

	header := &Header{}
	body := &Body{}
	b := &Block{header: header, body: body}
	b.header.Extradata = []byte("Genesis Block")
	// b.header.Coinbase = nil
	b.header.Number = big.NewInt(0)
	b.header.Timestamp = big.NewInt(1535706356)
	//构造两个hash
	candidate := []byte{1, 2}
	epoch := []byte{1, 3}
	dposContext := &DposContextProto{}
	dposContext.CandidateHash = MockHash(candidate)
	// Log.Infoln("dposContext.CandidateHash", dposContext.CandidateHash.String())
	dposContext.EpochHash = MockHash(epoch)
	header.DposContext = dposContext
	b.header.BlockHash = b.Header().Hash()

	return b
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
func (block *Block) Body() *Body     { return &Body{block.body.Bodydata} }
func (block *Block) Header() *Header { return CopyHeader(block.header) }
func (block *Block) DposContext() *DposContextProto {
	return &DposContextProto{EpochHash: block.header.DposContext.EpochHash, CandidateHash: block.header.DposContext.CandidateHash}
}

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

//打开或新建一个数据库
func OpenDatabase(name string, cache, handles int) (ethdb.Database, error) {
	if name == "" {
		return ethdb.NewMemDatabase(), nil
	}
	return ethdb.NewLDBDatabase(name, cache, handles)
}

// //将新增的块存入db
func (block *Block) StoreAddedBlock(db ethdb.Database) ([]byte, error) {
	Log := log(prefix)
	// Log := common.Logger.NewSessionLogger()
	//DB存储结构1--key:[]byte("latesthash"),value:rlp(blockhash)
	//DB存储结构2--key:[]byte("latestnumber"),value:[]byte(blocknumber)
	//DB存储结构3--key:rlp(blockhash),value:rlp(block结构)，均为rlp的编码值
	//将编码后的key，value存入db
	storedBlockHash, _ := rlp.EncodeToBytes(block.header.BlockHash)
	storedBlock, _ := rlp.EncodeToBytes(*block)
	// defer db.Close()
	// err := db.Put([]byte("blockchain"), storedBlockHash)
	err := db.Put(storedBlockHash, storedBlock)
	if err != nil {
		Log.Infoln("新增的区块存入数据库出错")
		return nil, err
	}
	err = db.Put([]byte("latesthash"), storedBlockHash)
	if err != nil {
		Log.Infoln("lasthash存储出错，err：", err)
	}
	enc := utils.EncodeBlockNumber(block.header.Number.Uint64())
	if err != nil {
		Log.Infoln("新增的区块存入数据库出错")
		return nil, err
	}
	err = db.Put([]byte("latestnumber"), enc)
	if err != nil {
		Log.Infoln("新增的区块存入数据库出错")
		return nil, err
	}
	Log.Infoln("latesthash,", storedBlockHash)
	Log.Infoln("新区块写入db成功")

	return storedBlockHash, nil

}

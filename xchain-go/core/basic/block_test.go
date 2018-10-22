package basic

import (
	"math/big"
	"testing"
	"xchain-go/common"
	"xchain-go/ethdb"
)

// var prefix = "basic"

// func log(prefix string) *mylog.SimpleLogger {
// 	common.InitLog(prefix)
// 	return common.Logger.NewSessionLogger()
// }

func TestSetHash(t *testing.T) {
	// Log := log(prefix)

	header := &Header{}
	body := &Body{}
	tests := []struct {
		name       string
		Extradata  []byte
		Number     *big.Int
		Timestamp  *big.Int
		ParentHash common.Hash
		want       string
	}{
		{"block中各参数格式正确", []byte("Genesis Block"), big.NewInt(0), big.NewInt(1535706356), [32]byte{}, "0x1cc694fb2c1c03172445fd1fcb06d7130445a07826b5b808854ff391624d368d"},
		{"block中extradata为空值", []byte{}, big.NewInt(0), big.NewInt(1535706356), [32]byte{}, "0xa734d7dd44af909e1406bd31dbee1dbf053d521673449e5ec9fd6afd6d41b83a"},
		{"block中Number数值很大", []byte{}, big.NewInt(0), big.NewInt(999999999999999999), [32]byte{}, "0xaf642341bfdc99b717871efe0990f1bed0092462ee035f85333007826c97b178"},

		// {},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := &Block{header: header, body: body}
			header.Extradata, header.Number, header.Timestamp, header.ParentHash = tt.Extradata, tt.Number, tt.Timestamp, tt.ParentHash
			got := block.SetHash()
			if got.String() != tt.want {
				t.Error("测试结果:", got.String(), "!=预期结果want:", tt.want)
			}
			// else {
			// 	Log.Infoln("测试结果:", got.String(), "==预期结果want:", tt.want)
			// }
		})
	}
}

func TestCopyHeader(t *testing.T) {
	// Log := log(prefix)

	header := &Header{}
	header.Extradata, header.Number, header.Validator, header.ParentHash = []byte("Genesis Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 12})
	formerheaderhash := header.Hash().String()
	// Log.Infoln("formerheaderhash", formerheaderhash)
	tests := []struct {
		name       string
		Extradata  []byte
		Number     *big.Int
		Validator  common.Address
		ParentHash common.Hash
		want       string
	}{
		{"改变Number-*big.Int的值", []byte("Genesis Block"), big.NewInt(1), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 2, 3}), formerheaderhash},
		{"改变Extradata-[]byte类型的值", []byte("test Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 2, 3}), formerheaderhash},
		{"改变Validator-common.Address的值", []byte("Genesis Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 13}), common.BytesToHash([]byte{1, 2, 3}), formerheaderhash},
		{"改变ParentHash-common.Hash类型的值", []byte("Genesis Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 2, 4}), formerheaderhash},
	}
	testheader := CopyHeader(header)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Log.Infoln("")
			//改变header中的值，看会不会影响testheader
			header.Extradata, header.Number, header.Validator, header.ParentHash = tt.Extradata, tt.Number, tt.Validator, tt.ParentHash
			got := testheader.Hash().String()
			// Log.Infoln("testheaderhash:", testheader.Hash().String())
			if got != tt.want {
				t.Error("测试结果:", got, "!=预期结果want:", tt.want)
			}
		})
	}
}

func TestStoreAddedBlock(t *testing.T) {
	// Log := log(prefix)

	db := ethdb.NewMemDatabase()
	header := &Header{}
	body := &Body{}
	header.Extradata, header.Number, header.Timestamp, header.Validator, header.ParentHash = []byte("Genesis Block"), big.NewInt(0), big.NewInt(1535706356), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 12})
	block := &Block{header: header, body: body}

	if _, err := block.StoreAddedBlock(db); err != nil {
		t.Error("存储失败,err:", err)
	}
}

// func TestNewBlock(t *testing.T) {
// 	var block Block
// 	var block2 *Block
// 	header := &Header{}
// 	body := &Body{}
// 	header.Extradata = []byte("Genesis Block")
// 	// header.Coinbase = nil
// 	header.Number = big.NewInt(0)
// 	// header.ParentHash = ""
// 	// header.Timestamp = 1535706356
// 	// header.Blockhash = "0xefa03879003f4ce162de0463832f4b9d78d14bace3d8d3eac04278ee0d72fd50"
// 	Log.Infoln("-----------------------------------------")
// 	Log.Infoln("header,", header)

// 	block = Block{header: header, body: body}
// 	block2 = &Block{header: header, body: body}
// 	Log.Infoln("-----------------------------------------")
// 	Log.Infoln("block,", block)
// 	Log.Infoln("block2,", block2)
// 	Log.Infoln("-----------------------------------------")

// 	blockrlp, err := rlp.EncodeToBytes(block)
// 	if err != nil {
// 		Log.Infoln("block encode fail!err,", err)
// 	}
// 	block2rlp, err := rlp.EncodeToBytes(block2)
// 	if err != nil {
// 		Log.Infoln("block encode fail!err,", err)
// 	}
// 	Log.Infoln("-----------------------------------------")
// 	Log.Infoln("blockrlp,", blockrlp)
// 	Log.Infoln("block2rlp,", block2rlp)
// 	Log.Infoln("-----------------------------------------")

// 	var deblockrlp Block
// 	var deblock2rlp *Block
// 	err = rlp.DecodeBytes(blockrlp, &deblockrlp)
// 	if err != nil {
// 		Log.Infoln("block decode fail!err,", err)
// 	}
// 	err = rlp.DecodeBytes(block2rlp, &deblock2rlp)
// 	if err != nil {
// 		Log.Infoln("block2 decode fail!err,", err)
// 	}
// 	Log.Infoln("-----------------------------------------")
// 	Log.Infoln("deblockrlp,", deblockrlp)
// 	Log.Infoln("deblock2rlp,", deblock2rlp)
// 	Log.Infoln("-----------------------------------------")

// }

package basic

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"xchain-go/common"
	"xchain-go/rlp"

	log "github.com/inconshreveable/log15"
)

func TestSetHash(t *testing.T) {
	// Log := log(prefix)

	header := &Header{}
	// body := &Body{}
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
			block := &Block{header: header}
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

	header := &Header{}
	header.Extradata, header.Number, header.Validator, header.ParentHash = []byte("Genesis Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 12})
	header.DposContext = &DposContextProto{}
	formerheaderhash := header.Hash().String()
	log.Debug("formerheaderhash", "formerheaderhash", formerheaderhash)
	tests := []struct {
		name        string
		Extradata   []byte
		Number      *big.Int
		Validator   common.Address
		ParentHash  common.Hash
		DposContext *DposContextProto
		want        string
	}{
		{"改变Number-*big.Int的值", []byte("Genesis Block"), big.NewInt(1), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 2, 3}), &DposContextProto{}, formerheaderhash},
		{"改变Extradata-[]byte类型的值", []byte("test Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 2, 3}), &DposContextProto{}, formerheaderhash},
		{"改变Validator-common.Address的值", []byte("Genesis Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 13}), common.BytesToHash([]byte{1, 2, 3}), &DposContextProto{}, formerheaderhash},
		{"改变ParentHash-common.Hash类型的值", []byte("Genesis Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 2, 4}), &DposContextProto{}, formerheaderhash},
		// {"改变DposContext-*DposContextProto类型的值", []byte("Genesis Block"), big.NewInt(0), common.BytesToAddress([]byte{1, 12}), common.BytesToHash([]byte{1, 2, 4}), &DposContextProto{}, formerheaderhash},
	}
	testheader := CopyHeader(header)
	log.Debug("header", "header", header.Hash().String())
	log.Debug("testheader", "testheader", testheader.Hash().String())
	// txs := Mocktxs()
	// block := NewBlock(header, txs)
	// Log.Infoln("blockheader:", block.Header().Hash().String())

	// Log.Infoln("block-string:", block.String())
	// Log.Infoln("header-string:", header.String())
	// Log.Infoln("block-header-string:", block.Header().String())
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
func TestBlockEncoding(t *testing.T) {
	headerEnc := common.FromHex("f8eca0bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff498845506eb07a000000000000000000000000000000000000000000000000000000000000000000c887465737464617461948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a00a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9ef842a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000")
	bodyEnc := common.FromHex("f855f853e880825208808402faf080940095e7baea6a6c7c4c2dfeb977efac326af552d8808474657374808080e9808252080c8402faf080940095e7baea6a6c7c4c2dfe1237efac326af552d806857465737432808080")

	// var block Block
	// if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
	// 	fmt.Println("decode error: ", err)
	// 	t.Fatal("decode error: ", err)
	// }

	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}

	// 构造block，测试block的header和body的信息encode后的结果
	block := MockBlockWithtxs()
	// fmt.Println("len:", len(block.Transactions()))
	// for _, tx := range block.Transactions() {
	// 	fmt.Println("tx.data:", tx.data)
	// }
	// fmt.Println(" block.Transactions()[0].Hash()", block.Transactions()[0].Hash())
	check("len(Transactions)", len(block.Transactions()), 2)
	// check("Transactions[0].Hash", block.Transactions()[0].Hash(), tx1.Hash())

	ourBlockheaderEnc, err := rlp.EncodeToBytes(block.Header())
	fmt.Println("ourBlockheaderEnc:", ourBlockheaderEnc)
	if err != nil {
		// fmt.Println("encode error: ", err)
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourBlockheaderEnc, headerEnc) {
		// fmt.Printf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockheaderEnc, headerEnc)
		t.Errorf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockheaderEnc, headerEnc)
	}
	ourBlockbodyEnc, err := rlp.EncodeToBytes(block.Body())
	fmt.Println("ourBlockbodyEnc:", ourBlockbodyEnc)
	if err != nil {
		fmt.Println("encode error: ", err)
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourBlockbodyEnc, bodyEnc) {
		// fmt.Printf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockbodyEnc, bodyEnc)
		t.Errorf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockbodyEnc, bodyEnc)
	}

	//解码后的数据测试
	var decodeheader *Header
	err = rlp.Decode(bytes.NewReader(ourBlockheaderEnc), &decodeheader)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("header:", decodeheader.Hash())

	var decodebody *Body
	err = rlp.Decode(bytes.NewReader(ourBlockbodyEnc), &decodebody)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("body:", decodebody.Transactions[0].data)

	decodeblock := NewBlockWithHeader(decodeheader).WithBody(decodebody.Transactions)
	Transaction := decodeblock.Transactions()
	for i := 0; i < Transaction.Len(); i++ {
		fmt.Println("交易", i, "内容：", Transaction[i].data)
	}

	// check("Difficulty", decodeblock.Difficulty(), big.NewInt(131072))
	// check("GasLimit", decodeblock.GasLimit(), uint64(3141592))
	// check("GasUsed", decodeblock.GasUsed(), uint64(21000))
	// check("Coinbase", decodeblock.Coinbase(), common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"))
	// check("MixDigest", decodeblock.MixDigest(), common.HexToHash("bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff498"))
	// check("Root", decodeblock.Root(), common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017"))
	// check("Hash", decodeblock.Hash(), common.HexToHash("0a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9e"))
	// check("Nonce", decodeblock.Nonce(), uint64(0xa13a5a8c8f2bb1c4))
	// check("Time", decodeblock.Time(), big.NewInt(1426516743))
	// check("Size", decodeblock.Size(), common.StorageSize(len(blockEnc)))

}

package rawdb

import (
	"fmt"
	"math/big"
	"testing"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/crypto/sha3"
	"xchain-go/ethdb"
	"xchain-go/rlp"
)

func TestHeaderStorage(t *testing.T) {
	Log := log(prefix)

	db := ethdb.NewMemDatabase()
	dposContext := &basic.DposContextProto{}
	//构造两个hash
	candidate := []byte{1, 2}
	epoch := []byte{1, 3}
	dposContext.CandidateHash = basic.MockHash(candidate)
	dposContext.EpochHash = basic.MockHash(epoch)
	Log.Infoln("dposContext.CandidateHash", dposContext.CandidateHash.String())
	Log.Infoln("dposContext.EpochHash", dposContext.EpochHash.String())

	header := &basic.Header{Extradata: []byte("test Block"), Number: big.NewInt(12), Timestamp: big.NewInt(1535706356), Validator: common.BytesToAddress([]byte{1, 12}), ParentHash: common.BytesToHash([]byte{1, 12}), DposContext: dposContext}
	if entry := ReadHeader(db, header.Hash(), header.Number.Uint64()); entry != nil {
		t.Fatalf("此header找到了返回结果：%v", entry)
	}
	WriteHeader(db, header)
	//从db中查找已写入的header
	//先查找header对应的结构体
	if entry := ReadHeader(db, header.Hash(), header.Number.Uint64()); entry == nil {
		t.Fatalf("没有找到已存储的header")
	} else if entry.Hash() != header.Hash() {
		t.Fatalf("检索到的header信息不匹配: have %v, want %v", entry, header)
	}
	//再查找header的rlp
	if entry := ReadHeaderRLP(db, header.Hash(), header.Number.Uint64()); entry == nil {
		t.Fatalf("没有找到已存储的headerrlp")
	} else {
		var decodedata *basic.Header
		if err := rlp.DecodeBytes(entry, &decodedata); err != nil {
			t.Fatalf("解析headerrlp失败")
		} else if decodedata.Hash() != header.Hash() {
			t.Fatalf("检索到的headerrlp信息不匹配: have %v, want %v", decodedata.Hash().String(), header.Hash().String())
		}
		Log.Infoln("解析到的DposContext，DposContext,CandidateHash:", decodedata.DposContext.CandidateHash.String(), "EpochHash", decodedata.DposContext.EpochHash.String())
	}
	DeleteHeader(db, header.Hash(), header.Number.Uint64())
	if entry := ReadHeader(db, header.Hash(), header.Number.Uint64()); entry != nil {
		t.Fatalf("删除的header没有成功，仍返回了：%v", entry)
	}
}

// 测试body的存取
func TestBodyStorage(t *testing.T) {
	db := ethdb.NewMemDatabase()
	// Create a test body to move around the database and make sure it's really new
	//新建一个测试的body，去进行数据的存取。
	//首先确保库中没有这个hash
	//构造hash的时候，临时直接使用body的内容去计算hash，作为存储用的key中的hash；key中的number，用12作为number

	body := &basic.Body{Bodydata: []byte("test Block")}

	//计算hash的方法
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, body)
	hash := common.BytesToHash(hw.Sum(nil))

	number := uint64(12)
	//检查库中，确保不存在此条记录
	if entry := ReadBody(db, hash, number); entry != nil {
		t.Fatalf("Non existent body returned: %v", entry)
	}
	//开始写入和校验body的存取
	WriteBody(db, hash, number, body)
	//读取和校验body，然后根据读取到的body的rlp的值进行校验
	if entry := ReadBody(db, hash, number); entry == nil {
		t.Fatalf("Stored body not found")
	}
	//todo:解析出来的body进行校验
	if entry := ReadBodyRLP(db, hash, number); entry == nil {
		t.Fatalf("没有找到存储的body的rlp值")
	} else {
		//若查到了rlp的值，rlp解码，计算解出来的body的hash与原hash是否一致
		var decodebody *basic.Body
		err := rlp.DecodeBytes(entry, &decodebody)
		if err != nil {
			t.Error("bodyrlp解码失败")
		}
		hasher := sha3.NewKeccak256()
		rlp.Encode(hasher, decodebody)
		hashdecode := common.BytesToHash(hasher.Sum(nil))

		if hash != hashdecode {
			t.Fatalf("存储的rlp值与实际的不匹配: have %v, want %v", hashdecode, hash)
		}
		// fmt.Printf("存储的rlp值与实际值: have %v, want %v", hashdecode, hash)

	}
	// Delete the header and verify the execution
	DeleteBody(db, hash, number)
	if entry := ReadBody(db, hash, number); entry != nil {
		t.Fatalf("Deleted header returned: %v", entry)
	}
}
func TestBlockStorage(t *testing.T) {
	db := ethdb.NewMemDatabase()
	//构造一个block，并确保库中没有该block的相关信息
	// //构造两个hash
	// dposContext := &basic.DposContextProto{}
	// candidate := []byte{1, 2}
	// epoch := []byte{1, 3}
	// dposContext.CandidateHash = basic.MockHash(candidate)
	// dposContext.EpochHash = basic.MockHash(epoch)
	// Log.Infoln("dposContext.CandidateHash", dposContext.CandidateHash.String())
	// Log.Infoln("dposContext.EpochHash", dposContext.EpochHash.String())
	header := &basic.Header{Extradata: []byte("test Block"), Number: big.NewInt(12), Timestamp: big.NewInt(1535706356), Validator: common.BytesToAddress([]byte{1, 12}), ParentHash: common.BytesToHash([]byte{1, 12}), DposContext: dposContext}
	body := &basic.Body{Bodydata: []byte("test Block")}
	block := basic.NewBlock(header, body)
	Log.Infoln("block.Header().Hash()", block.Header().Hash().String())
	Log.Infoln("header.Hash()", header.Hash().String())
	if entry := ReadBlock(db, block.Header().Hash(), block.NumberU64()); entry != nil {
		t.Errorf("此block找到了返回的block结果,%v", entry)
	}
	if entry := ReadHeader(db, block.Header().Hash(), block.NumberU64()); entry != nil {
		t.Errorf("此block找到了返回的header结果,%v", entry)
	}
	if entry := ReadBody(db, block.Header().Hash(), block.NumberU64()); entry != nil {
		t.Errorf("此block找到了返回的body结果,%v", entry)
	}
	Log.Infoln("该测试block在数据库中不存在，可以开始测试存储校验")
	//写入并且校验库中的block信息
	WriteBlock(db, block)
	//校验
	if entry := ReadBlock(db, block.Header().Hash(), block.NumberU64()); entry == nil {
		t.Fatalf("没找到存储的block")
	} else if entry.Header().Hash() != block.Header().Hash() {
		Log.Infoln("entry", entry)
		t.Fatalf("查找到的block与实际不匹配have %v, want %v", entry.Header().Hash().String(), block.Header().Hash().String())
	}
	if entry := ReadHeader(db, block.Header().Hash(), block.NumberU64()); entry == nil {
		t.Fatalf("没找到存储的block")
	} else if entry.Hash() != block.Header().Hash() {
		t.Fatalf("查找到的block的header与实际不匹配have %v, want %v", entry.Hash().String(), block.Header().Hash().String())
	}
	if entry := ReadBody(db, block.Header().Hash(), block.NumberU64()); entry == nil {
		t.Fatalf("没找到存储的blockhave %v, want %v", entry, block.Body())
	}
	//todo:读取到的body的校验

	//删除测试的block
	DeleteBlock(db, block.Header().Hash(), block.NumberU64())
	if entry := ReadBlock(db, block.Header().Hash(), block.NumberU64()); entry != nil {
		t.Fatalf("Deleted block returned: %v", entry)
	}
	if entry := ReadHeader(db, block.Header().Hash(), block.NumberU64()); entry != nil {
		t.Fatalf("Deleted header returned: %v", entry)
	}
	if entry := ReadBody(db, block.Header().Hash(), block.NumberU64()); entry != nil {
		t.Fatalf("Deleted body returned: %v", entry)
	}
	Log.Infoln("删除成功")
}

func hash1(header *basic.Header) string {
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, header)
	hash1 := common.BytesToHash(hasher.Sum(nil))
	return hash1.String()
}
func TestHash(t *testing.T) {
	header := &basic.Header{Extradata: []byte("test Block"), Number: big.NewInt(12), Timestamp: big.NewInt(1535706356), Validator: common.BytesToAddress([]byte{1, 12}), ParentHash: common.BytesToHash([]byte{1, 12})}
	Log.Infoln("hash1:", hash1(header))

	hash2 := header.Hash()
	fmt.Printf("hash2:%v", hash2.String())

}

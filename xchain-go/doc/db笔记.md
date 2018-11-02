
# db存储使用说明文档
基础是使用了leveldb
常用的方法见代码示例
```
//goleveldb代码示例
package main

import (
	"fmt"
	"git/tfd/xchain-go/core/basic"
	"git/tfd/xchain-go/ethdb"
)

const dbFile = "test.db"
func main() {
	//创建并打开数据库
	db, err := basic.OpenDatabase(dbFile, 512, 512)
	if err != nil {
		fmt.Println("failed to open database,the err info ：%v ", err)
	}
	fmt.Println("打开数据库成功")
	defer db.Close() //关闭数据库

	//写入5条数据
	db.Put([]byte("key1"), []byte("value1"))
	db.Put([]byte("key2"), []byte("value2"))
	db.Put([]byte("key3"), []byte("value3"))
	db.Put([]byte("key4"), []byte("value4"))
	db.Put([]byte("key5"), []byte("value5"))

	//循环遍历数据
	fmt.Println("循环遍历数据")
	iter := db.(*ethdb.LDBDatabase).NewIterator()
	for iter.Next() {
		fmt.Printf("key:%s, value:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()

	//读取某条数据
	data, _ := db.Get([]byte("key2"))
	fmt.Printf("读取单条数据key2:%s\n", data)

	// //批量写入数据
	// batch := new(leveldb.Batch)
	// batch.Put([]byte("key6"), []byte(strconv.Itoa(10000)))
	// batch.Put([]byte("key7"), []byte(strconv.Itoa(20000)))
	// batch.Delete([]byte("key4"))
	// db.(*ethdb.LDBDatabase).Write(batch, nil)

	// //查找数据
	// key := "key7"
	// iter = db.(*ethdb.LDBDatabase).NewIterator()
	// for ok := iter.Seek([]byte(key)); ok; ok = iter.Next() {
	// 	fmt.Printf("查找数据:%s, value:%s\n", iter.Key(), iter.Value())
	// }
	// iter.Release()

	//按key范围遍历数据
	fmt.Println("按key范围遍历数据")
	iter = db.(*ethdb.LDBDatabase).NewIterator()
	for iter.Next() {
		fmt.Printf("key:%s, value:%s\n", iter.Key(), iter.Value())
	}
	iter.Release()
}
```
# block数据存储db说明
package路径:
`core/rawdb`
**1. 存取block的方法**

- 存储block
> func WriteBlock(db DatabaseWriter, block *basic.Block) 

 参数：block、db

- 读取block
> func ReadBlock(db DatabaseReader, hash common.Hash, number uint64) *basic.Block 

参数：db、blockhash、blocknumber
 
**2. 存储的key-value**
 |key|value|拼接key的方法或key关键字|说明|
 |:--|:--|:--|:--|
 |[]byte("LastBlock")|hash|headBlockKey|跟踪最新的已知完整块的哈希值|
 |"H"+hash|encodeBlockNumber(number)|headerNumberKey|通过hash快速检索到number|
 |"h"+number+"n"|hash|headerHashKey|通过number快速检索到hash|
 |"h"+number+hash |rlp(header)|headerKey|存储header的rlp编码值|
 |"b"+number+hash |rlp(body)|blockBodyKey|存储body的rlp编码值|
 **3. 存储方法对应的功能**
 |方法名|功能|对应的key方法|说明|
 |:--|:--|:--|:--|
 |WriteCanonicalHash|根据number存储hash|headerHashKey|通过number快速检索到hash|
 |WriteHeadBlockHash|存储最新的区块的hash|headBlockKey|跟踪最新的已知完整块的哈希值|
 |WriteHeader|1.根据hash存储number|1.headerNumberKey|通过hash快速检索到number|
 |WriteHeader|2.存储header的rlp编码值|2.headerKey|存储header的rlp编码值|
 |WriteBody|存储body的rlp编码值|blockBodyKey|存储body的rlp编码值
 |WriteBlock|调用WriteHeader和WriteBody来存储block|
 **4. 读取方法对应的功能**
 |方法名|功能|对应的key方法|说明|
 |:--|:--|:--|:--|
 |ReadCanonicalHash|根据number读取hash|headerHashKey|通过number快速检索到hash|
 |ReadHeadBlockHash|读取最新的区块的hash|headBlockKey|跟踪最新的已知完整块的哈希值|
 |ReadHeaderNumber|根据hash读取number|headerNumberKey|通过hash快速检索到number|
 |ReadHeader|读取header的rlp编码值，并解析|headerKey|读取header并解析|
 |ReadBody|读取body的rlp编码值，并解析|blockBodyKey|存储body并解析|
 |ReadBlock|调用ReadHeader和ReadBody，拼接成block|
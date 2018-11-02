//Package rawdb包含了一组低级数据库.
//用于定义block结构存储db时的，key的关键字、数据项前缀、后缀，及拼接key的函数实现。
package rawdb

import (
	"encoding/binary"
	"xchain-go/common"
)

// The fields below define the low level database schema prefixing.
var (
	// databaseVerisionKey tracks the current database version.
	// databaseVerisionKey跟踪当前数据库版本
	databaseVerisionKey = []byte("DatabaseVersion")

	// headHeaderKey tracks the latest know header's hash.
	// headHeaderKey跟踪最新的已知的header的哈希值
	headHeaderKey = []byte("LastHeader")

	// headBlockKey tracks the latest know full block's hash.
	// headBlockKey跟踪最新的已知完整块的哈希值
	headBlockKey = []byte("LastBlock")

	// headFastBlockKey tracks the latest known incomplete block's hash duirng fast sync.
	// headFastBlockKey跟踪在快速同步期间的最新的已知不完整块的哈希值
	headFastBlockKey = []byte("LastFast")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`, used for indexes).
	// 数据项前缀（使用单个字节以避免混合数据类型，避免使用`i`，用于索引）。

	headerPrefix    = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	numSuffix       = []byte("n") // headerPrefix + num (uint64 big endian) + numSuffix -> hash
	blockHashPrefix = []byte("H") // blockHashPrefix + hash -> num (uint64 big endian)

	blockBodyPrefix = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body

	configPrefix = []byte("ethereum-config-") // config prefix for the db

)

// encodeBlockNumber encodes a block number as big endian uint64
// encodeBlockNumber将blocknumber转为大端数
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// headerKey = headerPrefix + num (uint64 big endian) + hash
// headerKey对应存储header
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(headerPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

// headerHashKey = headerPrefix + num (uint64 big endian) + numSuffix
// headerHashKey 对应存储headerhash，可以通过number快速检索到hash
func headerHashKey(number uint64) []byte {
	return append(append(headerPrefix, encodeBlockNumber(number)...), numSuffix...)
}

// headerNumberKey = blockHashPrefix + hash
// headerNumberKey对应存储number，可以通过hash快速检索到number
func headerNumberKey(hash common.Hash) []byte {
	return append(blockHashPrefix, hash.Bytes()...)
}

// blockBodyKey = blockBodyPrefix + num (uint64 big endian) + hash
// blockBodyKey对应存储body的值
func blockBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(blockBodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

// configKey = configPrefix + hash
func configKey(hash common.Hash) []byte {
	return append(configPrefix, hash.Bytes()...)
}

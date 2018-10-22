// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rawdb

import (
	mylog "mylog2"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/rlp"
)

var prefix = "rawdb"

func log(prefix string) *mylog.SimpleLogger {
	common.InitLog(prefix)
	return common.Logger.NewSessionLogger()
}

// WriteCanonicalHash stores the hash assigned to a canonical block number.
// WriteCanonicalHash存储分配给规范块编号的哈希
func WriteCanonicalHash(db DatabaseWriter, hash common.Hash, number uint64) {
	Log := log(prefix)

	if err := db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		Log.Infoln("Failed to store number to hash mapping", "err", err)
	}
}

// WriteHeadHeaderHash stores the hash of the current canonical head header.
// 当前规范头标题的哈希值
func WriteHeadHeaderHash(db DatabaseWriter, hash common.Hash) {
	Log := log(prefix)

	if err := db.Put(headHeaderKey, hash.Bytes()); err != nil {
		Log.Infoln("Failed to store last header's hash", "err", err)
	}
}

// WriteHeadBlockHash stores the head block's hash.
// header中的blockhash
func WriteHeadBlockHash(db DatabaseWriter, hash common.Hash) {
	Log := log(prefix)

	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
		Log.Infoln("Failed to store last block's hash", "err", err)
	}
}

// WriteHeadFastBlockHash stores the hash of the current fast-sync head block.
func WriteHeadFastBlockHash(db DatabaseWriter, hash common.Hash) {
	Log := log(prefix)

	if err := db.Put(headFastBlockKey, hash.Bytes()); err != nil {
		Log.Infoln("Failed to store last fast block's hash", "err", err)
	}
}

// WriteHeader stores a block header into the database and also stores the hash-
// to-number mapping.
// key:blockhash,value:number
// key:number+hash,value:rlp(header)
func WriteHeader(db DatabaseWriter, header *basic.Header) {
	Log := log(prefix)

	// Write the hash -> number mapping
	var (
		hash    = header.Hash()
		number  = header.Number.Uint64()
		encoded = encodeBlockNumber(number)
	)
	Log.Infoln("headerhash:\n", hash.String())
	key := headerNumberKey(hash)
	if err := db.Put(key, encoded); err != nil {
		Log.Infoln("Failed to store hash to number mapping", "err", err)
	}
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	Log.Infoln("headerrlp:\n", data)
	if err != nil {
		Log.Infoln("Failed to RLP encode header", "err", err)
	}
	key = headerKey(number, hash)
	if err := db.Put(key, data); err != nil {
		Log.Infoln("Failed to store header", "err", err)
	}
	Log.Infoln("存储header成功")
}

// WriteBodyRLP stores an RLP encoded block body into the database.
func WriteBodyRLP(db DatabaseWriter, hash common.Hash, number uint64, rlp rlp.RawValue) {
	Log := log(prefix)

	if err := db.Put(blockBodyKey(number, hash), rlp); err != nil {
		Log.Infoln("Failed to store block body", "err", err)
	}
}

// WriteBody storea a block body into the database.
func WriteBody(db DatabaseWriter, hash common.Hash, number uint64, body *basic.Body) {
	Log := log(prefix)

	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		Log.Infoln("Failed to RLP encode body", "err", err)
	}
	WriteBodyRLP(db, hash, number, data)
}

// WriteBlock serializes a block into the database, header and body separately.
func WriteBlock(db DatabaseWriter, block *basic.Block) {
	WriteBody(db, block.Header().Hash(), block.NumberU64(), block.Body())
	WriteHeader(db, block.Header())
}

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
	"bytes"
	"encoding/binary"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/rlp"

	log "github.com/inconshreveable/log15"
)

// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
// ReadCanonicalHash检索分配给规范块编号的哈希。--通过number检索hash值
func ReadCanonicalHash(db DatabaseReader, number uint64) common.Hash {
	data, _ := db.Get(headerHashKey(number))
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// ReadHeaderNumber returns the header number assigned to a hash.
// ReadHeaderNumber返回分配给hash的块号。--根据块hash读取块号
func ReadHeaderNumber(db DatabaseReader, hash common.Hash) *uint64 {
	data, _ := db.Get(headerNumberKey(hash))
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

// ReadHeadHeaderHash retrieves the hash of the current canonical head header.
// ReadHeadHeaderHash返回当前块的块hash
func ReadHeadHeaderHash(db DatabaseReader) common.Hash {
	data, _ := db.Get(headHeaderKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// ReadHeadBlockHash retrieves the hash of the current canonical head block.
// ReadHeadBlockHash
func ReadHeadBlockHash(db DatabaseReader) common.Hash {
	data, _ := db.Get(headBlockKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// ReadHeadFastBlockHash retrieves the hash of the current fast-sync head block.
func ReadHeadFastBlockHash(db DatabaseReader) common.Hash {
	data, _ := db.Get(headFastBlockKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// ReadBlock retrieves an entire block corresponding to the hash, assembling it
// back from the stored header and body. If either the header or body could not
// be retrieved nil is returned.
//
// Note, due to concurrent download of header and block body the header and thus
// canonical hash can be stored in the database but the body data not (yet).
// ReadBlock 是将header和body读取之后，拼接成block
func ReadBlock(db DatabaseReader, hash common.Hash, number uint64) *basic.Block {

	header := ReadHeader(db, hash, number)
	if header == nil {
		log.Info("获取block的header时，header为空")
		return nil
	}
	body := ReadBody(db, hash, number)
	if body == nil {
		log.Info("获取block的body时，body为空")
		return nil
	}
	return basic.NewBlockWithHeader(header).WithBody(body.Transactions)
}

// ReadHeader retrieves the block header corresponding to the hash.
// ReadHeader 根据hash和number读取出header.返回rlp(header)在decode后的Header对象
func ReadHeader(db DatabaseReader, hash common.Hash, number uint64) *basic.Header {

	data := ReadHeaderRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	// Log.Infoln("读取到的header\n", data)
	// header := new(basic.Header)
	var header basic.Header

	// var dposcontext *basic.DposContextProto
	if err := rlp.DecodeBytes(data, &header); err != nil {
		log.Info("Invalid block header RLP", "hash", hash, "err", err)
		return nil
	}
	// Log.Infoln("header", header.Hash().String())
	return &header
}

// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
// ReadHeaderRLP返回存储的rlp(header)
func ReadHeaderRLP(db DatabaseReader, hash common.Hash, number uint64) rlp.RawValue {
	data, _ := db.Get(headerKey(number, hash))
	return data
}

// ReadBodyRLP retrieves the block body (transactions and uncles) in RLP encoding.
// ReadBodyRLP返回存储的rlp(body)
func ReadBodyRLP(db DatabaseReader, hash common.Hash, number uint64) rlp.RawValue {
	data, _ := db.Get(blockBodyKey(number, hash))
	return data
}

// ReadBody retrieves the block body corresponding to the hash.
// ReadHeader返回rlp(body)在decode后的Body对象
func ReadBody(db DatabaseReader, hash common.Hash, number uint64) *basic.Body {

	data := ReadBodyRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	body := new(basic.Body)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		log.Info("Invalid block body RLP", "hash", hash, "err", err)
		return nil
	}
	return body
}

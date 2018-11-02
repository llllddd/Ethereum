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
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/rlp"

	log "github.com/inconshreveable/log15"
)

// WriteCanonicalHash 传入hash和number，根据number存储hash
func WriteCanonicalHash(db DatabaseWriter, hash common.Hash, number uint64) error {

	if err := db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		log.Error("Failed to store number to hash mapping", "err", err)
		return err
	}
	return nil
}

// WriteHeadHeaderHash 传入hash，在当前db中根据最新区块"LastHeader"存储hash
func WriteHeadHeaderHash(db DatabaseWriter, hash common.Hash) error {

	if err := db.Put(headHeaderKey, hash.Bytes()); err != nil {
		log.Error("Failed to store last header's hash", "err", err)
		return err
	}
	return nil

}

// WriteHeadBlockHash 传入hash，在当前db中根据最新区块"LastBlock"存储hash
func WriteHeadBlockHash(db DatabaseWriter, hash common.Hash) error {

	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
		log.Error("Failed to store last block's hash", "err", err)
		return err
	}
	return nil
}

// WriteHeadFastBlockHash stores the hash of the current fast-sync head block.
func WriteHeadFastBlockHash(db DatabaseWriter, hash common.Hash) error {

	if err := db.Put(headFastBlockKey, hash.Bytes()); err != nil {
		log.Error("Failed to store last fast block's hash", "err", err)
		return err
	}
	return nil
}

// WriteHeader stores a block header into the database and also stores the hash-
// to-number mapping.
// key:blockhash,value:number
// key:number+hash,value:rlp(header)
func WriteHeader(db DatabaseWriter, header *basic.Header) error {

	// Write the hash -> number mapping
	var (
		hash    = header.Hash()
		number  = header.Number.Uint64()
		encoded = encodeBlockNumber(number)
	)
	log.Debug("headerhash", "headerhash", hash.String())
	key := headerNumberKey(hash)
	if err := db.Put(key, encoded); err != nil {
		log.Error("Failed to store hash to number mapping", "err", err)
		return err
	}
	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	log.Debug("headerrlp:\n", data)
	if err != nil {
		log.Error("Failed to RLP encode header", "err", err)
		return err
	}
	key = headerKey(number, hash)
	if err := db.Put(key, data); err != nil {
		log.Error("Failed to store header", "err", err)
		return err
	}
	log.Debug("存储header成功")
	return nil
}

// WriteBodyRLP stores an RLP encoded block body into the database.
func WriteBodyRLP(db DatabaseWriter, hash common.Hash, number uint64, rlp rlp.RawValue) error {

	if err := db.Put(blockBodyKey(number, hash), rlp); err != nil {
		log.Error("Failed to store block body", "err", err)
		return err
	}
	return nil
}

// WriteBody storea a block body into the database.
func WriteBody(db DatabaseWriter, hash common.Hash, number uint64, body *basic.Body) error {

	data, err := rlp.EncodeToBytes(body)
	log.Debug("blockBodyKey", "blockBodyKey", blockBodyKey(number, hash))
	if err != nil {
		log.Error("Failed to RLP encode body", "err", err)
		return err
	}
	WriteBodyRLP(db, hash, number, data)
	log.Debug("blockBodyKey", "blockBodyKey", blockBodyKey(number, hash))
	return nil
}

// WriteBlock serializes a block into the database, header and body separately.
func WriteBlock(db DatabaseWriter, block *basic.Block) error {
	if err := WriteBody(db, block.Header().Hash(), block.NumberU64(), block.Body()); err != nil {
		return err
	}
	if err := WriteHeader(db, block.Header()); err != nil {
		return err
	}
	return nil
}

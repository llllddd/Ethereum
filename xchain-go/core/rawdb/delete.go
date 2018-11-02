//实现删除db中的存储的块信息功能
package rawdb

import (
	"xchain-go/common"

	log "github.com/inconshreveable/log15"
)

// DeleteCanonicalHash removes the number to hash canonical mapping.
func DeleteCanonicalHash(db DatabaseDeleter, number uint64) {
	if err := db.Delete(headerHashKey(number)); err != nil {
		log.Error("Failed to delete number to hash mapping", "err", err)
	}
}

// DeleteHeader removes all block header data associated with a hash.
func DeleteHeader(db DatabaseDeleter, hash common.Hash, number uint64) {

	if err := db.Delete(headerKey(number, hash)); err != nil {
		log.Error("Failed to delete header", "err", err)
	}
	if err := db.Delete(headerNumberKey(hash)); err != nil {
		log.Error("Failed to delete hash to number mapping", "err", err)
	}
}

// DeleteBody removes all block body data associated with a hash.
func DeleteBody(db DatabaseDeleter, hash common.Hash, number uint64) {

	if err := db.Delete(blockBodyKey(number, hash)); err != nil {
		log.Error("Failed to delete block body", "err", err)
	}
}

// DeleteBlock removes all block data associated with a hash.
func DeleteBlock(db DatabaseDeleter, hash common.Hash, number uint64) {
	DeleteHeader(db, hash, number)
	DeleteBody(db, hash, number)
}

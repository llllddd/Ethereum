package main

import (
	"fmt"
	"hash"
	"sdk/sha3"
)

const (
	epochLength = 30000
)

type hasher func(dest []byte, data []byte)

func makeHasher(h hash.Hash) hasher {
	// sha3.state supports Read to get the sum, use it to avoid the overhead of Sum.
	// Read alters the state but we reset the hash before every operation.
	type readerHash interface {
		hash.Hash
		Read([]byte) (int, error)
	}
	rh, ok := h.(readerHash)
	if !ok {
		panic("can't find Read method on hash")
	}
	outputLen := rh.Size()
	fmt.Println("哈希值长度: ", outputLen)
	return func(dest []byte, data []byte) {
		rh.Reset()
		rh.Write(data)
		rh.Read(dest[:outputLen])
	}
}

func CalHash(block uint64) []byte {
	seed := make([]byte, 32)
	if block < epochLength {
		return seed
	}
	hasher := makeHasher(sha3.NewKeccak256())
	for i := 0; i < int(block/epochLength); i++ {
		hasher(seed, seed)
		fmt.Printf("第--%d--次: %x ", i, seed)
	}
	return seed
}

func CalHash1(block uint64) []byte {
	seed := make([]byte, 32)
	hasher := sha3.NewKeccak256()
	for i := 0; i < int(block/epochLength); i++ {
		hasher.Reset()
		hasher.Write(seed)
		seed = hasher.Sum(nil)
		fmt.Printf("第--%d--次: %x ", i, seed)
	}
	return seed
}

func main() {
	block := uint64(2*30000 + 1)
	result := CalHash1(block)
	fmt.Printf("普通哈希种子为: %x\n\n", result)
	result1 := CalHash(block)
	fmt.Printf("种子为: %x\n", result1)
}

package main

import (
	"fmt"
	"hash"
	"reflect"
	"sdk/sha3"

	"github.com/ethereum/go-ethereum/common/bitutil"

	"encoding/binary"
	"unsafe"
)

const (
	epochLength = 30000
	hashBytes   = 64
	hashWords   = 16
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

func GenerateCache(dest []uint32, seed []byte) {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&dest))
	header.Len *= 4
	header.Cap *= 4
	cache := *(*[]byte)(unsafe.Pointer(&header))

	//计算需要填充的64字节的哈希值的个数
	size := uint64(len(cache))
	rows := int(size) / hashBytes

	//第一个哈希是由种子计算得到，之后的哈希是由前一个哈希值计算得到。
	keccak256 := makeHasher(sha3.NewKeccak256())
	keccak256(cache, seed)
	//每隔64字节填充一个新的哈希值。
	for offset := uint64(hashBytes); offset < size; offset += hashBytes {
		keccak256(cache[offset:], cache[offset-hashBytes:offset])
	}
	//使用RandMemoHash算法计算最终的cache
	temp := make([]byte, hashBytes)

	for i := 0; i < 3; i++ {
		for j := 0; j < rows; j++ {
			var (
				srcOff = ((j - 1 + rows) % rows) * hashBytes
				dstOff = j * hashBytes
				xorOff = (binary.LittleEndian.Uint32(cache[dstOff:]) % uint32(rows)) * hashBytes
			)
			bitutil.XORBytes(temp, cache[srcOff:srcOff+hashBytes], cache[xorOff:xorOff+hashBytes])
			keccak256(cache[dstOff:], temp)
		}
	}
}

func GenerateDataset(cache []uint32, index uint32, keccak512 hasher) []byte {
	rows := uint32(len(cache) / 16) //64字节1长度
	mix := make([]byte, hashBytes)

	binary.LittleEndian.PutUint32(mix, cache[(index%rows)*hashWords]^index)
	fmt.Println("1...: ", mix)
	for i := 1; i < hashWords; i++ {
		binary.LittleEndian.PutUint32(mix[i*4:], cache[(index%rows)*hashWords+uint32(i)])
	}

	fmt.Println("2...: ", mix)
	keccak512(mix, mix)
	return mix
}

func main() {
	block := uint64(2*30000 + 1)
	result := CalHash1(block)
	fmt.Printf("普通哈希种子为: %x\n\n", result)
	result1 := CalHash(block)
	fmt.Printf("种子为: %x\n", result1)

	cache := make([]uint32, 16)
	GenerateCache(cache, result)
	fmt.Printf("cache: \n ", cache)
	keccak512 := makeHasher(sha3.NewKeccak512())

	dataitem := GenerateDataset(cache, 0, keccak512)
	fmt.Println("data: ", dataitem)
}

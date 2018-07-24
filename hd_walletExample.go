package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"sdk/sha3"

	"github.com/go-ethereum-hdwallet"
	bip39 "github.com/tyler-smith/go-bip39"
)

/*
func Mnemonic(entropy []byte)(string,error){
	//计算随机数的比特长要满足在128~256比特间且满足能被32整除
	entropyLength := len(entropy)*8
	if(entropyLength%32) !=0 || entropyLength >256 || entropyLength <128{
		return nil,err
	}
	//在随机数之后要加上校验和以验证完整性，校验和的长度为字节数除以4，即比特数除以32。
	checkSumLength:= entropyLength/32
	//助记码是将随机数分组，每组11比特，分组的数量也是助记码中单词的数量。
	mnemonicCodeLength := entropyLength/11
	//使用SHA256算法计算随机数的hash值，并将HASH值的前8比特加到随机数之后作为校验和。
	entropy = addCheckSum(entropy)








}
*/
func addChecksum1(data []byte) []byte {
	// Get first byte of sha256
	hasher := sha256.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	firstChecksumByte := hash[0]

	// len() is in bytes so we divide by 4
	checksumBitLength := uint(len(data) / 4)

	// For each bit of check sum we want we shift the data one the left
	// and then set the (new) right most bit equal to checksum bit at that index
	// staring from the left
	dataBigInt := new(big.Int).SetBytes(data)
	BigOne := big.NewInt(1)
	BigTwo := big.NewInt(2)
	for i := uint(0); i < checksumBitLength; i++ {
		// Bitshift 1 left
		dataBigInt.Mul(dataBigInt, BigTwo)

		// Set rightmost bit if leftmost checksum bit is set
		if uint8(firstChecksumByte&(1<<(7-i))) > 0 {
			dataBigInt.Or(dataBigInt, BigOne)
		}
	}
	fmt.Println("SHA256Hash为:", hash)
	return dataBigInt.Bytes()
}

func addCheckSum(entropy []byte) []byte {
	hasher := sha3.NewKeccak256()
	hasher.Write(entropy)
	hash := hasher.Sum(nil)
	//校验和的比特数
	checkSumLength := len(entropy) / 4

	entropy = append(entropy, hash[:checkSumLength]...)
	return entropy
}

func main() {
	//256bitseed
	_, seed := hdwallet.NewSeed()
	fmt.Println("随机种子:", seed)
	//生成助记码，种子为128~256比特
	entropy, _ := bip39.NewEntropy(256)
	//将生成的随机数转换为助记码
	code, _ := bip39.NewMnemonic(entropy)
	fmt.Println("种子", entropy, "对应的助记码为", code)
	code1 := addChecksum1(entropy)
	fmt.Println(code1)
}

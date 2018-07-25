package main

import (
	"fmt"
	"math/big"
	"strings"
	"sdk/sha3"
	"encoding/binary"
	"github.com/go-ethereum-hdwallet"
	bip39 "github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)



func Mnemonic(entropy []byte)(string,error){
	//计算随机数的比特长要满足在128~256比特间且满足能被32整除
	entropyLength := len(entropy)*8
	if(entropyLength%32) !=0 || entropyLength >256 || entropyLength <128{
		fmt.Println("ERROR")
	}
	//在随机数之后要加上校验和以验证完整性，校验和的长度为字节数除以4，即比特数除以32。
	checkSumLength:= entropyLength/32
	//助记码是将随机数分组，每组11比特，分组的数量也是助记码中单词的数量。
	mnemonicCodeLength := (entropyLength+checkSumLength)/11
	//使用SHA256算法计算随机数的hash值，并将HASH值的前8比特加到随机数之后作为校验和。
	entropy = addCheckSum(entropy)
	//将添加了校验码的随机数分组,每组11比特
	entropyInt := new(big.Int).SetBytes(entropy)

	words := make([]string,mnemonicCodeLength)

	word := big.NewInt(0)
	last11bitmask := big.NewInt(2047)
	
        wordList := wordlists.ChineseSimplified
	for i:= mnemonicCodeLength - 1;i>=0;i--{
		word.And(entropyInt, last11bitmask)
		entropyInt.Div(entropyInt, big.NewInt(2048))

		wordBytes := padByteSlice(word.Bytes(),2)

		words[i]=wordList[binary.BigEndian.Uint16(wordBytes)]
	}

	return strings.Join(words," "),nil
}

func padByteSlice(slice []byte,length int)[]byte{
	offset := length - len(slice)
	if offset <= 0{
		return slice
	}
	newSlice := make([]byte,length)
	copy(newSlice[offset:],slice)
	return newSlice
}


func addCheckSum(entropy []byte) []byte {
	hasher := sha3.NewKeccak256()
	hasher.Write(entropy)
	hash := hasher.Sum(nil)
	//校验和的比特数
	checkSumLength := len(entropy) / 4
	bytelength := checkSumLength/8
	entropy = append(entropy, hash[:bytelength]...)
	return entropy
}



func main() {
	//256bitseed
	_, seed := hdwallet.NewSeed()
	fmt.Println("随机种子:", seed)
	//生成助记码，种子为128~256比特
	entropy, _ := bip39.NewEntropy(256)
	//将生成的随机数转换为助记码
	code, _ := Mnemonic(entropy)
	fmt.Println("种子", entropy)
	fmt.Println("对应的助记码为", code)
}

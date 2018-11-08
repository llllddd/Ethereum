package main

import (
	"fmt"
	"xchain-go/crypto"
)

func main() {
	//	_, sk1 := Dilithium_weak.GenerateSk()

	pk2, sk2, _ := crypto.GenerateKey()
	fmt.Println(pk2)
	pk1, _ := crypto.GeneratePubKey(sk2)
	fmt.Println(pk1)

	add := crypto.PubkeyToAddress(pk2)
	fmt.Println(add, len(add))

	var msg1 = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0xb, 0xc, 0xd, 0xe, 0xf, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	//		var msg2 = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0xb, 0xc, 0xd, 0xe, 0xf, 1}
	//	sm1 := Dilithium_weak.Sign(msg1, sk1)
	sm1, _ := crypto.Sign(msg1, sk2)
	//	sm2[1500] = 0
	fmt.Println("-----------------------------\n", sm1)

	t := crypto.VerifySignature(pk2, msg1, sm1)
	fmt.Println("验证结果", t)

}

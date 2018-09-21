package main

import (
	"bytes"
	"crypto/ecdsa"
	//	"crypto/elliptic"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"io"
	"math/big"
	"math/rand"
	"time"
)

func genPrivatekey(rand io.Reader) (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func getPubkey(sk *ecdsa.PrivateKey) *ecdsa.PublicKey {
	return &sk.PublicKey
}

func Encrypt(M_x, M_y *big.Int, pk *ecdsa.PublicKey) (*big.Int, *big.Int, *big.Int, *big.Int) {
	curve := secp256k1.S256()

	rand.Seed(time.Now().Unix())
	rr := rand.New(rand.NewSource(time.Now().Unix()))
	r_temp := new(big.Int).Rand(rr, curve.N)

	r := r_temp.Bytes()

	//	r := []byte{111, 111}

	A_x, A_y := curve.ScalarBaseMult(r)

	Temp1, Temp2 := curve.ScalarMult(pk.X, pk.Y, r)

	//	fmt.Println("r*Pk", Temp1, Temp2)
	//	B_x, B_y := curve.Add(M_x, M_y, Temp1, Temp2)
	B_x := new(big.Int).Add(M_x, Temp1)
	B_y := new(big.Int).Add(M_y, Temp2)
	return A_x, A_y, B_x, B_y
}

func Decrypt(C_1, C_2, C_3, C_4 *big.Int, sk *ecdsa.PrivateKey) (*big.Int, *big.Int) {
	curve := secp256k1.S256()
	sk_ := sk.D.Bytes()

	temp_1, temp_2 := curve.ScalarMult(C_1, C_2, sk_)
	//	fmt.Println("ddddd", temp_1, temp_2)
	negtemp_1 := new(big.Int).Neg(temp_1)
	negtemp_2 := new(big.Int).Neg(temp_2)
	//	fmt.Println("Oncurve:", curve.IsOnCurve(negtemp_1, negtemp_2))
	//	fmt.Println("ccccc", negtemp_1, negtemp_2)
	M_x := new(big.Int).Add(C_3, negtemp_1)
	M_y := new(big.Int).Add(C_4, negtemp_2)
	return M_x, M_y
}

func main() {
	seed := []byte{110, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111, 111}

	bytesBuffer := bytes.NewBuffer(seed)

	sk, _ := genPrivatekey(bytesBuffer)

	pk := getPubkey(sk)

	fmt.Println("公钥为: ", pk)
	fmt.Println("明文为公钥为了简便")
	curve := secp256k1.S256()
	
	M_1, M_2 := pk.X, pk.Y
	TT := curve.IsOnCurve(M_1, M_2)
	fmt.Println("OnCurve:", TT)

	C_1, C_2, C_3, C_4 := Encrypt(M_1, M_2, pk)

	fmt.Println("密文:", C_1, C_2, C_3, C_4)

	P_1, P_2 := Decrypt(C_1, C_2, C_3, C_4, sk)

	fmt.Println("解密后的明文:", P_1, P_2)

}

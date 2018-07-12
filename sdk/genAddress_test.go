package genaddandsign

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"sdk/secp256k1"
	"testing"
)

var hash1 = []byte{102, 240, 242, 127, 60, 155, 194, 113, 15, 92, 233, 22, 129, 152, 226, 226, 18, 166, 147, 214, 245, 183, 139, 218, 102, 202, 4, 255, 115, 59, 130, 254}

var x = "85778012444621206970083398062621418315245209882652708636297984772786400062929"

var y = "73736342751682697214599007260010111457386809005892716907395380321104741225476"

var address = [20]byte{129, 152, 226, 226, 18, 166, 147, 214, 245, 183, 139, 218, 102, 202, 4, 255, 115, 59, 130, 254}

const TestCount = 1000

//定义测试的函数判断S256曲线生成私钥并判断产生的公钥是否在曲线上
func TestGenAccountSek(t *testing.T) {
	for i := 0; i < TestCount; i++ {
		priv, err := genAccountSek(rand.Reader)
		if err != nil {
			t.Errorf("error: %s", err)
			return
		}
		if !secp256k1.S256().IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
			t.Errorf("public key invalid: %s", err)
		}
	}
}

func TestGetBytePkFromSk(t *testing.T) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkeybytes := getBytePkFromSk(key)
	X, Y := elliptic.Unmarshal(secp256k1.S256(), pubkeybytes)
	if !secp256k1.S256().IsOnCurve(X, Y) {
		t.Errorf("publickey can't transform: %s ", err)
	}
}

func TestGenHashFromPub(t *testing.T) {
	var pkX, pkY big.Int
	pukX, _ := pkX.SetString(x, 10)
	pukY, _ := pkY.SetString(y, 10)
	pub := ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     pukX,
		Y:     pukY,
	}
	temphash := genHashFromPub(&pub)
	if !bytes.Equal(hash1, temphash) {
		t.Errorf("can't get valid hash")
	}
}

func TestGenAddressFromPub(t *testing.T) {
	var pkX, pkY big.Int
	pukX, _ := pkX.SetString(x, 10)
	pukY, _ := pkY.SetString(y, 10)
	pub := ecdsa.PublicKey{
		Curve: secp256k1.S256(),
		X:     pukX,
		Y:     pukY,
	}
	tempaddr := genAddressFromPub(&pub)
	var sl1 []byte = tempaddr[:]
	var sl2 []byte = address[:]
	if !bytes.Equal(sl1, sl2) {
		t.Errorf("can't get valid address: want %d,have %d", sl2, sl1)
	}
}

func TestSignBySecp256k1AndRecover(t *testing.T) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	seckey := key.D.Bytes()
	pubkey1 := elliptic.Marshal(secp256k1.S256(), key.PublicKey.X, key.PublicKey.Y)
	c := 32
	msg := make([]byte, c)
	rand.Read(msg)
	sig, err := signBySecp256k1(msg, seckey)
	if err != nil {
		t.Errorf("signature error: %s", err)
	}
	pubkey2, err := RecoverTxsignPubkey(msg, sig)
	if err != nil {
		t.Errorf("recover error: %s", err)
	}
	if !bytes.Equal(pubkey1, pubkey2) {
		t.Errorf("pubkey mismatch: want: %x have: %x", pubkey1, pubkey2)
	}
}

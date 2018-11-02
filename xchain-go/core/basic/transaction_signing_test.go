package basic

import (
	"math/big"
	"testing"

	"xchain-go/common"
	"xchain-go/crypto"
	"xchain-go/rlp"
)

func TestEIP155Signing(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	signer := NewEIP155Signer(big.NewInt(18))
	tx, err := SignTx(NewTransaction(new(big.Int), addr, new(big.Int), new(big.Int), nil), signer, key)

	if err != nil {
		t.Fatal(err)
	}

	from, err := Sender(signer, tx)
	if err != nil {
		t.Fatal(err)
	}

	if from != addr {
		t.Errorf("expected from and address to be equal. Got %x want %x", from, addr)
	}
	// tx.data.Hash = *tx.Hash()
	// txrlp, err := rlp.EncodeToBytes(tx)
	// if err != nil {
	// 	fmt.Println("err:", err)
	// }
	// fmt.Println("txrlp:", common.ToHex(txrlp))
	// fmt.Printf("expected from and address to be equal. Got %v want %v", from.String(), addr.String())
	// var txdecode Transaction
	// err = rlp.DecodeBytes(txrlp, &txdecode)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("decode:", txdecode.Hash())
	// fmt.Println("beforeencode,", tx.Hash())
}

func TestEIP155SigningVitalik(t *testing.T) {
	// Test vectors come from http://vitalik.ca/files/eip155_testvec.txt
	for i, test := range []struct {
		txRlp, addr string
	}{
		{
			"f86080825208808094e6cd5c336104c27a5fd9c90a891aed5b5097c1e4808048a014345d6c320bc570c2cf500ce061e12992bf2ae571a580a4cef9cb6e547d8f07a07cd5b6ae9068b0986ffb864ce86d82cc222c94eba8496c8d09683ebef4613d97",
			"0xe6CD5C336104c27A5FD9C90A891aeD5B5097c1e4",
		},
		{
			"f860808252088080945961b93240041ff20f4ff2483f54ce2dfa972a05808048a098a966e9f96fb3e32433228cdeb97401fd7b85d01105918bc08541edbc7374d3a0731979d25b69f2e6b6685c75764ddd7256f732baf71e34e0fb4bd67b58ee9bb6",
			"0x5961B93240041ff20f4fF2483f54Ce2dfa972a05",
		},
		{
			"f860808252088080944173db5657f47ba7777c2acddc83dec26ac5f497808048a0bae8ca64041554915ae21e5980f68c3202448d89679d7d3e222425d0f94c5c19a0031a509185bb9319800df6fdedce4dbcd9a4d78308cf3b8e11d0788dc2dc60a0",
			"0x4173db5657f47bA7777C2AcddC83DeC26ac5F497",
		},
		{
			"f86080825208808094ee8a6815d96056901a2bd2153566bdaa6761105e808048a05e15b5ecb7a423afdae7ae9d0189e5c2d585ea048e8870f2c3cfbb767282c269a073b52e9f26f65bfd246d83aa6f1834b31665179fd6c18422e212db14721f93df",
			"0xee8a6815D96056901A2bD2153566BdAa6761105E",
		},
		{
			"f8608082520880809416ef19f89e942175551976ef434e9c50b6fc37dc808047a09b24cef8206e6513644fb564187698986e1207ac4f5b6871eab5e14d52222e4ea02c78fd4eeddf9276a87ffdf7ed88d9cedf73480b5022be7f0a43d9b3f209f6ce",
			"0x16EF19F89e942175551976ef434e9C50b6FC37DC",
		},
		{
			"f860808252088080942fa3c330adaebed608a424dcde684e7b8e03053a808048a0126ffab44f98757480cee655275d7af3397cfe00936e26d815c184ef5f323149a05f68f737f85e0faa2365b4628693336e7fa213436359b7f42ef27930f4962a16",
			"0x2fA3C330aDAEBED608A424dcDe684e7b8E03053A",
		},
		// {"f866068504a817c80683023e3894353535353535353535353535353535353535353581d88025a06455bf8ea6e7463a1046a0b52804526e119b4bf5136279614e0b1e8e296a4e2fa06455bf8ea6e7463a1046a0b52804526e119b4bf5136279614e0b1e8e296a4e2d", "0xf1f571dc362a0e5b2696b8e775f8491d3e50de35"},
		// {"f867078504a817c807830290409435353535353535353535353535353535353535358201578025a052f1a9b320cab38e5da8a8f97989383aab0a49165fc91c737310e4f7e9821021a052f1a9b320cab38e5da8a8f97989383aab0a49165fc91c737310e4f7e9821021", "0xd37922162ab7cea97c97a87551ed02c9a38b7332"},
		// {"f867088504a817c8088302e2489435353535353535353535353535353535353535358202008025a064b1702d9298fee62dfeccc57d322a463ad55ca201256d01f62b45b2e1c21c12a064b1702d9298fee62dfeccc57d322a463ad55ca201256d01f62b45b2e1c21c10", "0x9bddad43f934d313c2b79ca28a432dd2b7281029"},
		// {"f867098504a817c809830334509435353535353535353535353535353535353535358202d98025a052f8f61201b2b11a78d6e866abc9c3db2ae8631fa656bfe5cb53668255367afba052f8f61201b2b11a78d6e866abc9c3db2ae8631fa656bfe5cb53668255367afb", "0x3c24d7329e92f84f08556ceb6df1cdb0104ca49f"},
	} {
		signer := NewEIP155Signer(big.NewInt(18))

		var tx *Transaction
		err := rlp.DecodeBytes(common.Hex2Bytes(test.txRlp), &tx)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		from, err := Sender(signer, tx)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		addr := common.HexToAddress(test.addr)
		if from != addr {
			t.Errorf("%d: expected %x got %x", i, addr, from)
		}

	}
}

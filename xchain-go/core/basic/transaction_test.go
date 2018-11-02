package basic

import (
	//	"bytes"
	//	"crypto/ecdsa"
	//	"encoding/json"
	"math/big"
	"testing"

	"xchain-go/common"
	//	"xchain-go/crypto"
	"xchain-go/rlp"
)

type tests struct {
	name      string
	Timestamp *big.Int
	GasLimit  *big.Int
	Address   common.Address
	Amount    *big.Int
	Data      []byte
	want      string
}

func TestHash(t *testing.T) {
	tests1 := []tests{
		{"transaction中各参数格式正确", big.NewInt(0), big.NewInt(5000000), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), []byte("test"), "0x62c15a9d8ca7fe3dbb9ef1e77dfaa6f7a0942c021bef0f20ed47bc63267a5119"},
		{"transation中data为空", big.NewInt(0), big.NewInt(5000000), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(0), []byte{}, "0xa7ecf5cbe68f1a672e92c3b9c9e19e971ad1472a5938e01a6075ae88c6b12e97"},
		{"transaction中amount数值很大", big.NewInt(0), big.NewInt(5000000), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(999999999999999999), []byte{}, "0x423db747fa5ff83f13ef9d83df05e04d5d7e6cd4500eb3f99891bf0a23c4dcc3"},
	}

	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			tx := NewTransaction(tt.Timestamp, tt.Address, tt.Amount, tt.GasLimit, tt.Data)
			got := tx.Hash()
			if got.String() != tt.want {
				t.Error("测试结果:", got.String(), "!=预期结果want:", tt.want)
			}
		})
	}
}

func TestSize(t *testing.T) {
	tests2 := []tests{
		{"空交易size", big.NewInt(0), big.NewInt(5000000), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), nil, "37.00 B"},
		{"空交易size携带数据的交易size", big.NewInt(0), big.NewInt(5000000), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), []byte{1, 2, 3}, "40.00 B"},
		{"不携带数据但其他信息完整的交易size", big.NewInt(127), big.NewInt(5000000), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(127), nil, "37.00 B"},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			tx := NewTransaction(tt.Timestamp, tt.Address, tt.Amount, tt.GasLimit, tt.Data)
			data, _ := rlp.EncodeToBytes(tx)
			size := uint64(len(data))
			tx.size.Store(common.StorageSize(rlp.ListSize(size)))
			got := tx.Size()
			if got.String() != tt.want {
				t.Error("测试结果:", got.String(), "!=预期结果want:", tt.want)
			}
		})
	}
}

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
	Address   common.Address
	Amount    *big.Int
	Data      []byte
	want      string
}

func TestHash(t *testing.T) {
	tests1 := []tests{
		{"transaction中各参数格式正确", big.NewInt(0), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), []byte("test"), "0xf151e0ad3b2d367c2aa46dc9954f3069a4dbe547551232d34d4fd81151e3abb7"},
		{"transation中data为空", big.NewInt(0), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(0), []byte{}, "0xec072cd930213c95ddeed95f743f8ec47374974d85fd55982ed6248720158895"},
		{"transaction中amount数值很大", big.NewInt(0), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(999999999999999999), []byte{}, "0xaebdaae0121c599dd65ff40d4681eff0c3b866d48b9897973920faf5adb868ad"},
	}

	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			tx := NewTransaction(tt.Timestamp, tt.Address, tt.Amount, tt.Data)
			got := tx.Hash()
			if got.String() != tt.want {
				t.Error("测试结果:", got.String(), "!=预期结果want:", tt.want)
			}
		})
	}
}

func TestSize(t *testing.T) {
	tests2 := []tests{
		{"空交易size", big.NewInt(0), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), nil, "31.00 B"},
		{"空交易size携带数据的交易size", big.NewInt(0), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), []byte{1, 2, 3}, "34.00 B"},
		{"不携带数据但其他信息完整的交易size", big.NewInt(127), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(127), nil, "31.00 B"},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			tx := NewTransaction(tt.Timestamp, tt.Address, tt.Amount, tt.Data)
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

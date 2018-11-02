package core

import (
	"math/big"
	"testing"
	"xchain-go/common"
	"xchain-go/core/basic"
)

var (
	tx_0 = basic.NewTransaction(big.NewInt(1537926486), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), big.NewInt(5000000), []byte{})

	tx_1 = basic.NewTransaction(big.NewInt(1537926486), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), big.NewInt(5000000), []byte{})

	tx_2 = basic.NewTransaction(big.NewInt(1537926488), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), big.NewInt(5000000), []byte{})

	tx_3 = basic.NewTransaction(big.NewInt(1537926489), common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d8"), big.NewInt(0), big.NewInt(5000000), []byte{})
)

func TestAdd(t *testing.T) {

	tests := []struct {
		name      string
		Timestamp *big.Int
		tx        *basic.Transaction
	}{
		{"Tx1", big.NewInt(1537926486), tx_0},
		{"Tx2", big.NewInt(1537926486), tx_1},
		{"Tx3", big.NewInt(1537926488), tx_2},
		{"Tx4", big.NewInt(1537926489), tx_3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := newtxList()
			ls.Add(tt.tx)
			for ti, tr := range ls.items {
				if tr.Timestamp().Cmp(tt.Timestamp) != 0 || ti.Cmp(tt.Timestamp) != 0 {
					t.Error("测试结果:", tr.Timestamp(), "!=预期结果want:", tt.Timestamp)
				}
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	ls1 := newtxList()
	ls1.Add(tx_0)
	ls1.Add(tx_1)
	ls1.Add(tx_2)
	ls1.Add(tx_3)
	want1 := []*big.Int{big.NewInt(1537926486), big.NewInt(1537926486), big.NewInt(1537926488), big.NewInt(1537926489)}
	ls2 := newtxList()
	ls2.Add(tx_1)
	ls2.Add(tx_2)
	ls2.Add(tx_3)
	want2 := []*big.Int{big.NewInt(1537926486), big.NewInt(1537926488), big.NewInt(1537926489)}

	tests := []struct {
		name string
		list *txList
		want []*big.Int
	}{
		{"有相同时间戳的交易排序", ls1, want1},
		{"时间戳各不相同的交易排序", ls2, want2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txs := tt.list.Flatten()
			for i, tx := range txs {
				if tx.Timestamp().Cmp(tt.want[i]) != 0 {
					t.Error("测试结果:", tx.Timestamp(), "!=预期结果want:", tt.want[i])
				}
			}
		})
	}

}

func TestRemove(t *testing.T) {
	ls1 := newtxList()
	ls1.Add(tx_0)
	ls1.Add(tx_1)
	ls1.Add(tx_2)
	ls1.Add(tx_3)

	ls2 := newtxList()
	ls2.Add(tx_1)
	ls2.Add(tx_2)
	ls2.Add(tx_3)
	want := big.NewInt(1537926486)
	tests := []struct {
		name string
		list *txList
		want *big.Int
	}{
		{"在有相同时间戳的交易列表删除指定时间戳交易", ls1, want},
		{"在时间戳各不相同的交易排列表删除指定时间戳交易", ls2, want},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.list.Remove(big.NewInt(1537926486))
			txs := tt.list.Ready()
			for _, tx := range txs {
				if tx.Timestamp().Cmp(tt.want) == 0 {
					t.Error("测试结果:", tx.Timestamp(), "!=预期结果want:", tt.want)
				}
			}
		})
	}

}

func TestForward(t *testing.T) {
	ls1 := newtxList()
	ls1.Add(tx_0)
	ls1.Add(tx_1)
	ls1.Add(tx_2)
	ls1.Add(tx_3)
	ls2 := newtxList()
	ls2.Add(tx_1)
	ls2.Add(tx_2)
	ls2.Add(tx_3)
	want := big.NewInt(1537926489)
	tests := []struct {
		name string
		list *txList
		want *big.Int
	}{
		{"在有相同时间戳的交易列表删除小于指定时间戳交易", ls1, want},
		{"在时间戳各不相同的交易排列表删除小于指定时间戳交易", ls2, want},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txs := tt.list.Forward(big.NewInt(1537926489))
			for _, tx := range txs {
				if tx.Timestamp().Cmp(tt.want) != -1 {
					t.Error("测试结果:", tx.Timestamp(), "!=预期结果want:", tt.want)
				}
			}
		})
	}

}

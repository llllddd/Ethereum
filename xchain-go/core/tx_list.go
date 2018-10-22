package core

import (
	"container/heap"
	"math/big"
	"sort"
	"xchain-go/core/basic"
)

type txHeap []*big.Int

//实现sort接口中的三个方法
func (h txHeap) Len() int           { return len(h) }
func (h txHeap) Less(i, j int) bool { return h[i].Cmp(h[j]) < 0 }
func (h txHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

//这两个函数实现了heap中定义的两个方法，这样就定义了一个堆
func (h *txHeap) Push(x interface{}) {
	*h = append(*h, x.(*big.Int))
}

//从heap中删除元素
func (h *txHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

//使用堆来维护交易的时间戳
type txList struct {
	items map[*big.Int]*basic.Transaction
	index *txHeap
}

//初始化txList
func newtxList() *txList {
	return &txList{
		items: make(map[*big.Int]*basic.Transaction),
		index: new(txHeap),
	}
}

//Flatten 将交易按照时间戳排序(升序),快速排序后返回交易列表.
func (m *txList) Flatten() basic.Transactions {
	txs := make(basic.Transactions, 0, len(m.items))
	for _, tx := range m.items {
		txs = append(txs, tx)
	}
	sort.Sort(basic.TransactionsByTimestamp(txs))
	return txs
}

//将小于给定范围的交易删除并返回，堆排序查找
func (m *txList) Forward(threshold *big.Int) basic.Transactions {
	var removed basic.Transactions
	for m.index.Len() > 0 && (*m.index)[0].Cmp(threshold) < 0 {
		timestamp := heap.Pop(m.index).(*big.Int)
		removed = append(removed, m.items[timestamp])
		delete(m.items, timestamp)
	}
	return removed
}

//移除指定时间戳的交易
func (m *txList) Remove(timestamp *big.Int) bool {
	length := m.index.Len()
	for i := length - 1; i >= 0; i-- { //由于有相同时间戳的交易所以需要全部遍历可能耗时
		if (*m.index)[i].Cmp(timestamp) == 0 {
			heap.Remove(m.index, i)
		}
	}

	for t, _ := range m.items {
		if t.Cmp(timestamp) == 0 {
			delete(m.items, t)
		}
	}
	return true
}

//向txlist中添加交易
func (m *txList) Add(tx *basic.Transaction) bool {
	timestamp := tx.Timestamp()
	if m.items[timestamp] == nil {
		heap.Push(m.index, timestamp)
	}
	m.items[timestamp] = tx
	return true
}

//拿出txlist中的所有交易
func (m *txList) Ready() basic.Transactions {
	var ready basic.Transactions
	for _, tx := range m.items {
		ready = append(ready, tx)
	}
	return ready
}

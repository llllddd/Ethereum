package main

import (
	"sync"
	"xchain-go/common"
	"xchain-go/ethdb"
)

const (
	blockInterval    = int64(10)    //出块间隔
	epochInterval    = int64(86400) //一轮出块周期
	maxValidatorSize = 21           //验证节点个数
)

type Dpos struct {
	db ethdb.Database
	mu sync.RWMutex
}

// 节点类型
type Node struct {
	Name  common.Address //节点名称
	Votes int            // 被选举的票数
}

// func (d *Dpos) CheckValidator(lastBlock Block, now int64) error {
// 	// TODO : 校验Deadline?

// }

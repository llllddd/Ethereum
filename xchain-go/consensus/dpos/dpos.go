package dpos

import (
	"bytes"
	"errors"
	"xchain-go/common"
	"xchain-go/consensus"
	"xchain-go/core/basic"
	"xchain-go/ethdb"
)

const (
	extraVanity      = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal        = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
	timeOfFirstBlock = int64(0)
	epochInterval    = int64(86400)             //一轮周期间隔
	blockInterval    = int64(10)                //区块间隔
	maxValidatorSize = 4                        //最大验证节点个数
	safeSize         = maxValidatorSize*2/3 + 1 //最少验证节点个数
)

type Dpos struct {
	db     ethdb.Database // Database to store and retrieve snapshot checkpoints
	signer common.Address //出块节点

}

func PrevSlot(now int64) int64 {
	return int64((now-1)/blockInterval) * blockInterval
}

func NextSlot(now int64) int64 {
	return int64((now+blockInterval-1)/blockInterval) * blockInterval
}

// 检查当前的 validator 是否为当前节点，如果是的话则通过 CreateNewWork 来创建一个新的打块任务
func (d *Dpos) CheckValidator(lastBlock *basic.Block, now int64) error {

	// 查询当前是否为出块时间
	if err := d.checkDeadline(lastBlock, now); err != nil {
		return err
	}

	// 找出当前区块验证者
	validator, err := lookupValidator(d.db, now)
	if err != nil {
		return err
	}

	if (validator == common.Address{}) || bytes.Compare(validator.Bytes(), d.signer.Bytes()) != 0 {
		return errors.New("invalid block validator")
	}
	return nil
}

// 判断当前是否为出块周期
func (d *Dpos) checkDeadline(lastBlock *basic.Block, now int64) error {

	prevSlot := PrevSlot(now)
	nextSlot := NextSlot(now)

	if lastBlock.Time().Int64() >= nextSlot {
		return errors.New("mint the future block")
	}
	// 若lastBlock的时间和prevSlot相等，说明上一个区块时间已经完成，可以出下一个块，或者下一个区块开始时间和现在时间差距不到1，则说明可以出块
	if lastBlock.Time().Int64() == prevSlot || nextSlot-now <= 1 {
		return nil
	}
	return errors.New("wait for last block arrived")
}

// 准备出块数据
func (d *Dpos) Prepare(chain consensus.ChainReader, header *basic.Header) error {

	number := header.Number.Uint64()
	// TODO:增加Extra长度限制
	if len(header.Extradata) < extraVanity {
		header.Extradata = append(header.Extradata, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extradata))...)
	}
	header.Extradata = header.Extradata[:extraVanity]
	header.Extradata = append(header.Extradata, make([]byte, extraSeal)...)
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return errors.New("unknown ancestor")
	}

	// 将Validator设置为当前节点的signer
	header.Validator = d.signer
	return nil
}

package dpos

import (
	"encoding/binary"
	"errors"

	"math/big"
	"math/rand"
	"sort"
	"xchain-go/common"
	"xchain-go/core/basic"
	"xchain-go/crypto"
	"xchain-go/ethdb"
	// "xchain-go/core/basic"
)

// // 验证者详细信息，包含出块数，对应投票人信息，域名
// type Validator struct {
// 	ValidatorAddr common.Address      `json:"address"`
// 	MintCnt       int64               `json:"mintCnt"` //出块数
// 	Delegators    `json:"delegators"` //投票人地址及对应投票数
// 	// DNS DNS `json:"dns"`//域名
// }

// //一个周期内，超级节点（验证者）信息，key:周期编号，value:验证者列表 Jsonstring(Validator)
// type Epoch map[int64]map[common.Address]string

type EpochContext struct {
	TimeStamp int64 `json:"timestamp"`
	// DposContext basic.DposContext `json:"dposContext"`
	statedb *basic.StateDB `json:"stateDB"`
}

// 查询本次出块账号
// 1. 确认当前是否为出块周期
// 2. 从validators中查询当前出块的账号
func lookupValidator(db ethdb.Database, now int64) (validator common.Address, err error) {
	Log := common.Logger.NewSessionLogger()
	validator = common.Address{}
	offset := now % epochInterval
	epochNum := now / epochInterval
	// 若offset%blockInterval != 0 ,则现在不是出块的时间，所以直接返回空地址和错误
	if offset%blockInterval != 0 {
		Log.Errorf("invalid time to mint the block")
		return common.Address{}, errors.New("invalid time to mint the block")
	}
	// 否则，offset/blockInterval 查询目前是这一个周期中的第几个区块
	offset /= blockInterval

	// 从 epochTrie 中查询 validators
	validators := basic.GetValidators(db, epochNum)

	validatorSize := len(validators)
	if validatorSize == 0 {
		Log.Errorf("failed to lookup validator")
		return common.Address{}, errors.New("failed to lookup validator")
	}
	// 每个节点出一个块，并且依照顺序出块，查询此轮验证账户地址
	offset %= int64(validatorSize)
	return validators[offset], nil
}

func tryElect(db ethdb.Database, genesis, parent *basic.Header, timestamp int64) error {

	Log := common.Logger.NewSessionLogger()
	// 查看当前是否为选举周期
	genesisEpoch := genesis.Timestamp.Int64() / epochInterval
	prevEpoch := parent.Timestamp.Int64() / epochInterval
	currentEpoch := timestamp / epochInterval
	Log.Infoln("初始周期Number: %v ,上一个周期Number: %v ,本轮周期Number: %v\n", genesisEpoch, prevEpoch, currentEpoch)
	// TODO： 为什么要有这一步？
	prevEpochIsGenesis := prevEpoch == genesisEpoch
	if prevEpochIsGenesis && prevEpoch < currentEpoch {
		prevEpoch = currentEpoch - 1
	}

	// 根据当前块和上一块的时间计算当前块和上一块是否属于同一个周期
	// 如果是同一个周期，意味者当前块不是周期的第一个块，不需要触发选举
	// 如果不是同一个周期，说明当前块是该周期的第一块，需要触发选举
	for i := prevEpoch; i < currentEpoch; i++ {
		// if prevEpoch is not genesis, kickout not active candidate
		// 如果前一个周期不是创世周期，触发踢出候选人规则
		// 踢出规则：看上一周期是否存在候选人出块少于特定阈值（50%）如果存在，则踢出
		Log.Debugln("prevEpochIsGenesis:", prevEpochIsGenesis)
		if !prevEpochIsGenesis {
			// 若上一轮出块数不符合要求，则踢出该验证者
			if err := kickoutValidator(db, currentEpoch, timestamp); err != nil {
				return err
			}
		}

		oldCandidates := basic.ReadCandidates(db)
		// 若候选者数量不足，则报错
		if int64(len(oldCandidates)) < safeSize {
			Log.Errorf("too few candidates")
			return errors.New("too few candidates")
		}

		candidates := sortableAddresses{}
		for _, candidate := range oldCandidates {
			candidates = append(candidates, &sortableAddress{candidate.CandiAddr, candidate.VotedNum})
		}

		// 对候选者进行排序
		sort.Sort(candidates)
		// 取前N个，设为验证者
		if len(candidates) > maxValidatorSize {
			candidates = candidates[:maxValidatorSize]
		}

		// 打乱验证者顺序
		// shuffle candidates
		// 打乱验证人列表，由于使用seed是由父区块的hash以及当前周期编号组合能够，
		// 所以每个节点计算出来的验证人列表也会一致
		seed := int64(binary.LittleEndian.Uint32(crypto.Keccak512(parent.Hash().Bytes()))) + i
		r := rand.New(rand.NewSource(seed))
		for i := len(candidates) - 1; i > 0; i-- {
			j := int(r.Int31n(int32(i + 1)))
			Log.DDDebugf(">>>shuffle candidates i = %v,j = %v\n", i, j)
			candidates[i], candidates[j] = candidates[j], candidates[i]
		}
		sortedValidators := make([]common.Address, 0)
		for _, candidate := range candidates {
			sortedValidators = append(sortedValidators, candidate.address)
		}

		// 将更新后的验证者列表存入到epochTrie树中,其中 i = prevEpoch
		basic.SetValidators(db, sortedValidators, i+1)
		Log.Infoln("Come to new epoch", "prevEpoch", i, "nextEpoch", i+1)
	}

	return nil

}

// 剔除验证者
func kickoutValidator(db ethdb.Database, epochNum int64, timestamp int64) error {
	Log := common.Logger.NewSessionLogger()
	// 读取当前周期
	// epochNum := timestamp / epochInterval
	// 取上一轮验证者
	epochNum = epochNum - 1
	validators := basic.ReadValidators(db, epochNum)
	Log.Debugf("kickoutValidator epochNum:%v \n", epochNum)
	if len(validators) == 0 {
		Log.Errorf("no validator could be kickout")
		return errors.New("no validator could be kickout")
	}

	epochDuration := epochInterval
	// 在第一个周期的情况下,需要做特殊处理
	if timestamp-timeOfFirstBlock < epochInterval {
		epochDuration = timestamp - timeOfFirstBlock
	}

	needKickoutValidators := sortableAddresses{}

	// 获取需要踢出的验证者列表
	for _, validator := range validators {
		// 获取validator的出块数量
		cnt := validator.MintCnt
		// 若出块数量小于最大出块数量的一半，则认为出块数量不达标
		if cnt < epochDuration/blockInterval/maxValidatorSize/2 {
			needKickoutValidators = append(needKickoutValidators, &sortableAddress{validator.ValidatorAddr, big.NewInt(cnt)})
		}

	}

	needKickoutValidatorCnt := len(needKickoutValidators)
	// 不需要踢出验证者
	if needKickoutValidatorCnt <= 0 {
		return nil
	}
	// 排序
	sorts := sort.Reverse(needKickoutValidators)
	sort.Sort(sorts)

	candidates := basic.ReadCandidates(db)
	candidateCount := int64(len(candidates))
	for i, validator := range needKickoutValidators {
		// ensure candidate count greater than or equal to safeSize
		if candidateCount <= safeSize {
			Log.Infoln("No more candidate can be kickout", "candidateCount", candidateCount, "needKickoutCount", len(needKickoutValidators)-i)
			return nil
		}

		if err := basic.KickoutCandidate(db, validator.address); err != nil {
			return err
		}
		// if kickout success, candidateCount minus 1
		candidateCount--
		Log.Infoln("Kickout candidate", "candidate", validator.address.String(), "mintCnt", validator.weight.String())
	}
	return nil

}

type sortableAddress struct {
	address common.Address
	weight  *big.Int
}

type sortableAddresses []*sortableAddress

func (p sortableAddresses) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p sortableAddresses) Len() int      { return len(p) }
func (p sortableAddresses) Less(i, j int) bool {
	if p[i].weight.Cmp(p[j].weight) < 0 {
		return false
	} else if p[i].weight.Cmp(p[j].weight) > 0 {
		return true
	} else {
		return p[i].address.String() < p[j].address.String()
	}
}

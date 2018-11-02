package basic

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"
	"xchain-go/common"
	"xchain-go/config"
	"xchain-go/ethdb"
	"xchain-go/rlp"
)

const maxCandidates = 100
const lockTime = 86400

var (
	epochPrefix      = []byte("epoch") //epoch 列表 for the db
	candidatesPrefix = []byte("candidates")
	accountPrefix    = []byte("address")
)

type DelegatorInfo struct {
	Votes    *big.Int `json:"votes"`     //投票数
	VoteTime *big.Int `json:"timestamp"` //投票时间
}
type Delegators map[common.Address]DelegatorInfo

// 候选人详细信息，包含获得的总票数，及给其投票的投票人信息
type Candidate struct {
	CandiAddr  common.Address      `json:"address"`  //候选人之地
	VotedNum   *big.Int            `json:"votedNum"` //总票数
	Delegators `json:"delegators"` //投票人地址及对应投票数
	// DNS DNS  `json:"dns"`//域名
}

// 候选人列表
type Candidates map[common.Address]Candidate

// TODO: 修改StateDB结构
type StateDB struct {
	db   ethdb.Database
	lock sync.Mutex
}

// 验证者详细信息，包含出块数，对应投票人信息，域名
type Validator struct {
	ValidatorAddr common.Address      `json:"address"`
	MintCnt       int64               `json:"mintCnt"` //出块数
	Delegators    `json:"delegators"` //投票人地址及对应投票数
	// DNS DNS `json:"dns"`//域名
}

//一个周期内，超级节点（验证者）信息，key:周期编号，value:验证者列表 Jsonstring(Validator)
type Epoch map[int64]map[common.Address]Validator

type DposContext struct {
	// // 记录每个周期的验证人列表及其详情（包含验证者地址，出块数量，投票人，DNS）
	// Epoch Epoch
	// // 记录当前候选人列表，及其详情（包含候选人地址，总票数，投票人，DNS）
	// Candidates Candidates
	db ethdb.Database
	// TODO:删除Account,此处仅用作测试
	Account Account
}

// TODO:修改Account结构
type Account map[common.Address]*big.Int

type DposContextProto struct {
	EpochHash     common.Hash `json:"epochRoot"        gencodec:"required"`
	CandidateHash common.Hash `json:"candidateRoot"    gencodec:"required"`
}

func NewDposContextFromProto(db ethdb.Database, ctxProto *DposContextProto) (*DposContext, error) {

	return &DposContext{
		db: db,
		//qiqi-todo:补充返回的account
		// Account:map[ctxProto.CandidateHash]
	}, nil
}

// MockDposProto 返回构造的DposContextProto，目前构造的两个参数都为空hash
func MockDposProto() *DposContextProto {
	return &DposContextProto{
		EpochHash:     common.Hash{},
		CandidateHash: common.Hash{},

		// EpochHash:     d.epochTrie.Hash(),
		// DelegateHash:  d.delegateTrie.Hash(),
		// CandidateHash: d.candidateTrie.Hash(),
		// VoteHash:      d.voteTrie.Hash(),
		// MintCntHash:   d.mintCntTrie.Hash(),
	}
}

// 添加候选人
func BecomeCandidate(db ethdb.Database, candidateAddr common.Address) error {
	Log := common.Logger.NewSessionLogger()
	// TODO: 1.进行地址校验
	// TODO: 2.判断候选人是否符合条件
	candidates := ReadCandidates(db)
	if len(candidates) == 0 {
		newCandidate := Candidate{CandiAddr: candidateAddr, VotedNum: big.NewInt(0), Delegators: make(Delegators)}
		candidates = make(map[common.Address]Candidate)
		candidates[candidateAddr] = newCandidate

		// 写入数据库
		WriteCandidates(db, candidates)
		return nil
	}
	candidateString, _ := json.Marshal(candidates)
	Log.Debugln("当前的候选人列表:", string(candidateString))
	// fmt.Println("ReadCandidates candidates:", candidates)
	// 候选人数量是否达到上限N
	if len(candidates) >= maxCandidates {
		return errors.New(fmt.Sprintf("候选人数已达上限：%v", maxCandidates))
	}
	// 判断该候选人是否存在
	if _, ok := candidates[candidateAddr]; ok {
		return errors.New(fmt.Sprintf("%v 已存在", candidateAddr.String()))
	}

	// 初始化候选人信息
	newCandidate := Candidate{CandiAddr: candidateAddr, VotedNum: big.NewInt(0), Delegators: make(Delegators)}
	Log.Debugln("newCandidate:", newCandidate)

	candidates[candidateAddr] = newCandidate

	// 写入数据库
	WriteCandidates(db, candidates)

	return nil
}

// 添加投票人
func (candidates Candidates) AddDelegator(delegatorAddr, candidateAddr common.Address, amount *big.Int) error {
	Log := common.Logger.NewSessionLogger()
	Log.Infoln(">> AddDelegator STARTS!")
	if _, ok := candidates[candidateAddr]; !ok {
		return errors.New("候选人不存在！")
	}

	candidate := candidates[candidateAddr]
	delegateInfo := DelegatorInfo{Votes: amount, VoteTime: big.NewInt(time.Now().Unix())}
	// 因为amount是指针，如果这里不重新赋值，在改变delegateInfo.Votes的值的时候，amount的值也会被改变。
	amountValue := big.NewInt(amount.Int64())
	if _, ok := candidate.Delegators[delegatorAddr]; ok {
		Log.Debugln("原来包含的投票数：", delegateInfo.Votes)
		delegateInfo.Votes.Add(delegateInfo.Votes, candidate.Delegators[delegatorAddr].Votes)
		Log.Debugln("增加后的投票数", delegateInfo.Votes)
	}
	candidate.Delegators[delegatorAddr] = delegateInfo
	candidate.VotedNum.Add(candidate.VotedNum, amountValue)

	// 打印candidate信息
	can, _ := json.Marshal(candidate)
	Log.Infoln("AddDelegator 完成后,候选人列表:", string(can))

	Log.Infoln("AddDelegator ENDS!")
	return nil

}

// 踢出投票人
func KickoutCandidate(db ethdb.Database, candidateAddr common.Address) error {
	Log := common.Logger.NewSessionLogger()
	Log.Infoln("KickoutCandidate STARTS!")
	// TODO：增加删除候选人的条件
	// TODO: 1.进行地址校验
	// TODO: 2.判断候选人是否符合条件
	candidates := ReadCandidates(db)
	candidateString, _ := json.Marshal(candidates)
	Log.Debugln("候选人列表读取结果:", string(candidateString))

	// 在候选人列表中，删除对应的候选人
	if _, ok := candidates[candidateAddr]; ok {
		delete(candidates, candidateAddr)
		Log.Infoln("删除成功！")
		// 更新数据库中的Candidates
		WriteCandidates(db, candidates)
		return nil
	}

	Log.Infoln("KickoutCandidate ENDS!")
	return errors.New("候选人不存在！")

}

// 发起投票
func Delegate(db ethdb.Database, delegatorAddr, candidateAddr common.Address, amount *big.Int) error {
	Log := common.Logger.NewSessionLogger()
	Log.Infoln(">>> Delegate STARTS!")
	candidates := ReadCandidates(db)
	account := ReadAccounts(db)
	// 判断候选人是否存在于候选人列表中
	if _, ok := candidates[candidateAddr]; !ok {
		return errors.New(fmt.Sprintf("所投候选人 %v 不存在，投票失败！", candidateAddr.String()))
	}

	if _, ok := account[delegatorAddr]; !ok {
		return errors.New(fmt.Sprintf("投票人账号 %v 不存在，投票失败！", delegatorAddr.String()))
	}
	// 查询delegatorAddr余额是否大于票数
	if account.GetBalance(delegatorAddr).Cmp(amount) == -1 {
		return errors.New(fmt.Sprintf("投票人 %v 余额不足，投票失败！", delegatorAddr))
	}

	// 进行投票，更新余额，更新Candidates列表
	// 更新余额
	account.SubBalance(delegatorAddr, amount)

	// 更新Candidate里面的投票人信息,总票数
	err := candidates.AddDelegator(delegatorAddr, candidateAddr, amount)
	if err != nil {
		return err
	}

	WriteCandidates(db, candidates)
	WriteAccounts(db, account)
	Log.Infoln(">>> Delegate ENDS!")
	return nil

}

// 取消投票
func UnDelegate(db ethdb.Database, delegatorAddr, candidateAddr common.Address, amount *big.Int) error {
	Log := common.Logger.NewSessionLogger()
	Log.Infoln(">>> UnDelegate STARTS!")
	candidates := ReadCandidates(db)
	// 读取配置文件，若配置文件中"lock"字段为false,则不判断锁定时间
	err := config.ParseConfig("../../config/conf.json")
	if err != nil {
		Log.Errorln("读取配置文件失败！err:", err)
	}

	if _, ok := candidates[candidateAddr]; !ok {
		return errors.New("候选人不存在！")
	}

	candidate := candidates[candidateAddr]
	// 查看投票人是否存在
	if _, ok := candidate.Delegators[delegatorAddr]; !ok {
		return errors.New(fmt.Sprintf("候选人 %v 中不存在投票人 %v", candidateAddr.String(), delegatorAddr.String()))
	}

	deleVoteTime := candidate.Delegators[delegatorAddr].VoteTime
	timestamp := big.NewInt(time.Now().Unix())
	timeInterval := new(big.Int)
	timeInterval.Sub(deleVoteTime, timestamp)

	// 时间间隔 < 锁定时间，返回错误。
	if config.Cfg.Lock == true && timeInterval.Cmp(big.NewInt(lockTime)) != 1 {
		Log.Infof("票数已被锁定，上一轮投票时间：%v ,当前时间：%v ,最短锁定时间：%v\n", deleVoteTime, timestamp, lockTime)
		return errors.New("票数已被锁定!")
	}

	// 取消投票票数 > 已投票数，返回错误
	if candidate.Delegators[delegatorAddr].Votes.Cmp(amount) == -1 {
		Log.Infoln("已投票数 %v 小于期望取消投票的票数 %v\n", candidate.Delegators[delegatorAddr].Votes, amount)
		return errors.New("已投票数小于期望取消投票的票数")
	}

	//  取消投票数 < 已投票数，相应投票人总票数 - 取消投票数
	if candidate.Delegators[delegatorAddr].Votes.Cmp(amount) == 1 {
		Log.Infoln("取消投票数 < 已投票数")
		candidate.Delegators[delegatorAddr].Votes.Sub(candidate.Delegators[delegatorAddr].Votes, amount)
	}
	// 取消投票数 = 已投票数 ，删除相应投票人
	if candidate.Delegators[delegatorAddr].Votes.Cmp(amount) == 0 {
		Log.Infoln("取消投票数 =  已投票数")
		delete(candidate.Delegators, delegatorAddr)
	}

	//总票数中删除取消投票数
	candidate.VotedNum = candidate.VotedNum.Sub(candidate.VotedNum, amount)
	// 更新数据库中候选人列表
	WriteCandidates(db, candidates)
	Log.Infoln(">>> UnDelegate ENDS!")
	return nil

}

// TODO:修改GetBalance方法
func (acc Account) GetBalance(addr common.Address) *big.Int {
	return acc[addr]
}

// TODO:修改GetBalance方法
func (acc Account) SubBalance(addr common.Address, amount *big.Int) {
	acc[addr] = acc[addr].Sub(acc[addr], amount)
}

// 从 epochTrie 中查询 验证者地址列表
func GetValidators(db ethdb.Database, epochNum int64) []common.Address {
	var validatorAddrs []common.Address
	var validators = make(map[common.Address]Validator)

	validators = ReadValidators(db, epochNum)

	for _, validator := range validators {
		validatorAddrs = append(validatorAddrs, []common.Address{validator.ValidatorAddr}...)
	}

	return validatorAddrs
}

// 更新验证者列表
func SetValidators(db ethdb.Database, validators []common.Address, epochNum int64) error {

	if len(validators) == 0 {
		return errors.New("验证者列表不能为空!")
	}
	validatorStructs := make(map[common.Address]Validator)
	candidates := ReadCandidates(db)
	for _, validatorAddr := range validators {

		for _, validator := range candidates {
			if validator.CandiAddr == validatorAddr {
				// validatorStructs = append(validatorStructs, Validator{ValidatorAddr: validator.CandiAddr, MintCnt: int64(0), Delegators: validator.Delegators})
				validatorStructs[validatorAddr] = Validator{ValidatorAddr: validator.CandiAddr, MintCnt: int64(0), Delegators: validator.Delegators}
			} else {
				continue
			}
		}

	}
	// 写入validators数据
	validatorsString, _ := json.Marshal(validatorStructs)
	common.Logger.Debugf("SetValidators 存入验证者，epochNum: %v, validatorStructs: %v\n", epochNum, string(validatorsString))
	WriteValidators(db, epochNum, validatorStructs)
	return nil
}

// 从DB读取candidates
func ReadCandidates(db ethdb.Database) (candidates Candidates) {
	Log := common.Logger.NewSessionLogger()
	Log.Infoln("ReadCandidates STARTS!")
	data, _ := db.Get(candidatesKey())
	// fmt.Println("ReadCandidates data:", data)
	if len(data) == 0 {
		return nil
	}
	var decodedata []byte
	if err := rlp.DecodeBytes(data, &decodedata); err != nil {
		Log.Errorf("candidates 解码失败! err:%v\n", err)
		return nil
	}

	if err := json.Unmarshal(decodedata, &candidates); err != nil {
		Log.Errorf("candidates Unmarshal 失败! err:%v\n", err)
	}
	Log.Infoln("ReadCandidates ENDS!")
	return candidates

}

// 从DB读取验证者
func ReadValidators(db ethdb.Database, epochNumber int64) (validators map[common.Address]Validator) {

	data, _ := db.Get(epochKey(epochNumber))
	if len(data) == 0 {
		return nil
	}
	var decodedata []byte
	if err := rlp.DecodeBytes(data, &decodedata); err != nil {
		common.Logger.Errorf("validators 解码失败！data:%v,err:%v\n", data, err)
		return nil
	}
	if err := json.Unmarshal(decodedata, &validators); err != nil {
		common.Logger.Printf("validators Unmarshal 失败! err:%v\n", err)
	}
	return validators
}

// 从DB读取candidates
func ReadAccounts(db ethdb.Database) (account Account) {

	data, _ := db.Get(accountKey())
	// fmt.Println("ReadCandidates data:", data)
	if len(data) == 0 {
		return nil
	}
	var decodedata []byte
	if err := rlp.DecodeBytes(data, &decodedata); err != nil {
		common.Logger.Errorf("account 解码失败! err:%v\n", err)
		return nil
	}

	if err := json.Unmarshal(decodedata, &account); err != nil {
		common.Logger.Errorf("account Unmarshal 失败! err:%v\n", err)
	}

	return account

}

// epochKey = epochPrefix + epochNum
func epochKey(epochNum int64) []byte {
	return append(epochPrefix, []byte(strconv.FormatInt(epochNum, 10))...)
}

// candidatesKey = candidatesPrefix
func candidatesKey() []byte {
	return []byte(candidatesPrefix)
}

// candidatesKey = candidatesPrefix
func accountKey() []byte {
	return []byte(accountPrefix)
}

// 更新数据库中的Candidates列表
func WriteCandidates(db ethdb.Database, candidates Candidates) error {
	Log := common.Logger.NewSessionLogger()
	can, err := json.Marshal(candidates)
	if err != nil {
		Log.Errorf("candidates Json格式转换失败, err:%v\n", err)
		return err
	}
	Log.Debugln("存入候选人列表 :", string(can))
	data, err := rlp.EncodeToBytes(can)
	if err != nil {
		Log.Errorf("Failed to RLP encode candidates, err:%v\n", err)
		return err

	}
	if err = db.Put(candidatesKey(), data); err != nil {
		Log.Errorf("Failed to store candidates body, err:%v\n", err)
		return err
	}
	return nil
}

// 更新数据库中的Candidates列表
func WriteValidators(db ethdb.Database, epochNum int64, validators map[common.Address]Validator) error {
	Log := common.Logger.NewSessionLogger()
	vals, err := json.Marshal(validators)
	if err != nil {
		Log.Errorf("validators Json格式转换失败, err:%v\n", err)
		return err
	}
	Log.Debugln("validators 存入候选人列表 :", string(vals))
	data, err := rlp.EncodeToBytes(vals)
	if err != nil {
		Log.Errorf("Failed to RLP encode validators, err:%v\n", err)
		return err

	}
	if err = db.Put(epochKey(epochNum), data); err != nil {
		Log.Errorf("Failed to store validators body, err:%v\n", err)
		return err
	}

	return nil

}

// 更新数据库中的Account列表
func WriteAccounts(db ethdb.Database, account Account) error {

	acc, err := json.Marshal(account)
	if err != nil {
		common.Logger.Errorf("account Json格式转换失败, err:%v", err)
		return err
	}
	common.Logger.Debugln("account 存入候选人列表 :", string(acc))
	data, err := rlp.EncodeToBytes(acc)
	if err != nil {
		common.Logger.Errorf("Failed to RLP encode candidates, err:%v", err)
		return err

	}
	if err = db.Put(accountKey(), data); err != nil {
		common.Logger.Errorf("Failed to store account body, err:%v", err)
		return err
	}
	return nil
}

// func PrintLog() {
// 	Log := common.Logger.NewSessionLogger()
// 	common.Logger.Debugln("PrintLog hello world!")
// 	common.Logger.Errorln("PrintLog hello world!")
// 	common.Logger.Infoln("PrintLog hello world!")

// 	Log.Debugln("PrintLog hello world!")
// 	Log.Errorln("PrintLog hello world!")
// 	Log.Infoln("PrintLog hello world!")
// }

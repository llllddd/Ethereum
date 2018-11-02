package dpos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"
	"xchain-go/common"
	"xchain-go/core/basic"
)

const testDbFile = "xchain_test.db"

func InitEpochAndDB() {
	common.InitLog("TestEpoch")
	Log := common.Logger.NewSessionLogger()
	if err := os.RemoveAll(testDbFile); err != nil {
		Log.Errorln("删除DB文件失败，err:", err)
	}
	Log.Debugln("删除DB文件成功！")

	db, err := basic.OpenDatabase(testDbFile, 512, 512)
	if err != nil {
		Log.Errorln("failed to open database,the err info : ", err)
	}
	Log.Debugln("打开数据库成功")
	defer db.Close() //关闭数据库

	addr := [20]byte{0, 1, 2}
	addr1 := [20]byte{1, 1, 1}
	addr2 := [20]byte{2, 2, 2}
	addr3 := [20]byte{3, 3, 3}
	addr4 := [20]byte{4, 4, 4}
	delegator1 := [20]byte{6, 7, 8}
	delegator2 := [20]byte{6, 6, 6}
	acc := make(basic.Account)
	acc[delegator1] = big.NewInt(100)
	acc[delegator2] = big.NewInt(100)

	candidele := make(basic.Delegators)
	newCandidate := basic.Candidate{CandiAddr: addr,
		VotedNum:   big.NewInt(0),
		Delegators: candidele,
	}

	cand := make(map[common.Address]basic.Candidate)
	cand[addr] = newCandidate

	err = basic.BecomeCandidate(db, addr)
	if err != nil {
		Log.Errorln(err)
	}
	err = basic.BecomeCandidate(db, addr1)
	if err != nil {
		Log.Errorln(err)
	}
	err = basic.BecomeCandidate(db, addr2)
	if err != nil {
		Log.Errorln(err)
	}
	err = basic.BecomeCandidate(db, addr3)
	if err != nil {
		Log.Errorln(err)
	}
	err = basic.BecomeCandidate(db, addr4)
	if err != nil {
		Log.Errorln(err)
	}
	readCandidates := basic.ReadCandidates(db)
	candidateString, _ := json.Marshal(readCandidates)
	Log.Debugln("candidates 读取结果:", string(candidateString))

	var validators []common.Address
	validators = append(validators, addr)
	validators = append(validators, addr1)
	validators = append(validators, addr2)
	validators = append(validators, addr3)
	// validators = append(validators, addr4)

	basic.SetValidators(db, validators, 2)
	// WriteAccounts(db, acc)
	readValidators := basic.ReadValidators(db, 2)

	readValidatorsString, _ := json.Marshal(readValidators)
	Log.Debugln("readValidatorsString 读取结果:", string(readValidatorsString))
}

// 测试每轮踢出不符合条件的验证者
func TestKickoutValidators(t *testing.T) {
	InitEpochAndDB()
	Log := common.Logger.NewSessionLogger()
	db, err := basic.OpenDatabase(testDbFile, 512, 512)
	if err != nil {
		Log.Errorln("failed to open database,the err info : ", err)
	}
	Log.Debugln("打开数据库成功")
	// tests := []struct {
	// 	name string
	// 	args *big.Int
	// 	want *big.Int
	// }{
	// 	{"出块数不足踢出", big.NewInt(11), big.NewInt(10)},
	// 	{"踢出数量太多", big.NewInt(9), big.NewInt(1)},
	// 	{"已投票数与取消票数相等", big.NewInt(1), big.NewInt(0)},
	// }
	candidatesbefore := basic.ReadCandidates(db)
	candidatesbeforeString, _ := json.Marshal(candidatesbefore)
	Log.Debugln("candidatesbefore 读取结果:", string(candidatesbeforeString))

	kickoutValidator(db, 1, epochInterval*2)
	candidates := basic.ReadCandidates(db)
	candiString, _ := json.Marshal(candidates)
	Log.Debugln("candiString 读取结果:", string(candiString))
}

func TestTryElect(t *testing.T) {
	InitEpochAndDB()
	Log := common.Logger.NewSessionLogger()

	db, err := basic.OpenDatabase(testDbFile, 512, 512)
	if err != nil {
		Log.Errorln("failed to open database,the err info : ", err)
	}
	Log.Debugln("打开数据库成功")

	genesis := &basic.Header{
		Timestamp: big.NewInt(0),
	}

	parent := &basic.Header{
		Timestamp: big.NewInt(epochInterval*2 + 1),
	}
	candidates := basic.ReadCandidates(db)
	candidatesString, _ := json.Marshal(candidates)
	Log.Infoln("befor tryElect,候选人列表::", string(candidatesString))

	validators2 := basic.ReadValidators(db, 2)
	validators3 := basic.ReadValidators(db, 3)
	validators2S, _ := json.Marshal(validators2)
	validators3S, _ := json.Marshal(validators3)

	Log.Infoln("befor tryElect,validators2:%v\n", string(validators2S))
	Log.Infoln("befor tryElect,validators3:%v\n", string(validators3S))

	err = tryElect(db, genesis, parent, epochInterval*3+1)
	if err != nil {
		fmt.Println("TestTryElect err:", err)
	}

	acandidates := basic.ReadCandidates(db)
	acandidatesString, _ := json.Marshal(acandidates)
	Log.Infoln("after tryElect,candidates:", string(acandidatesString))

	avalidators2 := basic.ReadValidators(db, 2)
	avalidators3 := basic.ReadValidators(db, 3)
	avalidators2S, _ := json.Marshal(avalidators2)
	avalidators3S, _ := json.Marshal(avalidators3)

	Log.Infoln("after tryElect,validators2:%v\n", string(avalidators2S))
	Log.Infoln("after tryElect,validators3:%v\n", string(avalidators3S))
}

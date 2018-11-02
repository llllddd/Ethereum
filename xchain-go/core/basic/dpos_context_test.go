package basic

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"
	"xchain-go/common"
)

// 初始化数据，在执行其他Test之前，需要先执行Init
// 初始化完成后，DB中包含候选人列表： {"0x0001020000000000000000000000000000000000":{"address":"0x0001020000000000000000000000000000000000","votedNum":0,"delegators":{}}}
// DposContext 包含DB,以及Account:{"0x0606060000000000000000000000000000000000":100,"0x0607080000000000000000000000000000000000":100}
const testDbFile = "xchain_test.db"

func InitDposAndDB() {
	common.InitLog("TEST")
	Log := common.Logger.NewSessionLogger()
	if err := os.RemoveAll(testDbFile); err != nil {
		Log.Errorln("删除DB文件失败，err:", err)
	}
	Log.Infoln("删除DB文件成功！")

	db, err := OpenDatabase(testDbFile, 512, 512)
	if err != nil {
		Log.Errorln("failed to open database,the err info : ", err)
	}
	Log.Infoln("打开数据库成功")
	defer db.Close() //关闭数据库

	addr := [20]byte{0, 1, 2}
	addr1 := [20]byte{1, 1, 1}
	delegator1 := [20]byte{6, 7, 8}
	delegator2 := [20]byte{6, 6, 6}
	acc := make(Account)
	acc[delegator1] = big.NewInt(100)
	acc[delegator2] = big.NewInt(100)

	candidele := make(Delegators)
	newCandidate := Candidate{CandiAddr: addr,
		VotedNum:   big.NewInt(0),
		Delegators: candidele,
	}

	cand := make(map[common.Address]Candidate)
	cand[addr] = newCandidate

	err = BecomeCandidate(db, addr)
	if err != nil {
		fmt.Println(err)
	}
	err = BecomeCandidate(db, addr1)
	if err != nil {
		fmt.Println(err)
	}

	readCandidates := ReadCandidates(db)
	candidateString, _ := json.Marshal(readCandidates)
	Log.DDDebugln("candidates 读取结果:", string(candidateString))

	WriteAccounts(db, acc)
	accounts := ReadAccounts(db)

	accountsString, _ := json.Marshal(accounts)
	Log.DDDebugln("account 读取结果:", string(accountsString))
}

// 添加候选人测试
func TestAddCandidate(t *testing.T) {
	InitDposAndDB()
	db, err := OpenDatabase(testDbFile, 512, 512)
	if err != nil {
		fmt.Println("failed to open database,the err info : ", err)
	}
	fmt.Println("打开数据库成功")
	defer db.Close() //关闭数据库

	tests := []struct {
		name string
		args [20]byte
		want int
	}{
		{"候选人存在", [20]byte{0, 1, 2}, 2},
		{"候选人不存在", [20]byte{4, 5, 6}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BecomeCandidate(db, tt.args)
			if err != nil {
				fmt.Println(err)
			}

			readCandidates := ReadCandidates(db)
			candidateString, _ := json.Marshal(readCandidates)
			fmt.Println("candidates 读取结果:", string(candidateString))
			if got := len(readCandidates); got != tt.want {
				t.Errorf("got = %v,want = %v", got, tt.want)
			}
		})

	}

}

func TestKickoutCandidate(t *testing.T) {
	InitDposAndDB()
	db, err := OpenDatabase(testDbFile, 512, 512)
	if err != nil {
		fmt.Println("failed to open database,the err info : ", err)
	}
	fmt.Println("打开数据库成功")
	defer db.Close() //关闭数据库
	addr := [20]byte{0, 1, 2}
	addrNoExist := [20]byte{7, 7, 7}

	tests := []struct {
		name string
		args [20]byte
		want int
	}{
		{"候选人存在", addr, 1},
		{"候选人不存在", addrNoExist, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := KickoutCandidate(db, tt.args); err != nil {
				fmt.Println("err:", err)
			}

			readCandidates := ReadCandidates(db)
			candidateString, _ := json.Marshal(readCandidates)
			fmt.Println("candidates 读取结果:", string(candidateString))
			if got := len(readCandidates); got != tt.want {
				t.Errorf("failed")
			}
		})
	}

}

func TestDelegate(t *testing.T) {
	InitDposAndDB()
	db, err := OpenDatabase(testDbFile, 512, 512)
	if err != nil {
		fmt.Println("failed to open database,the err info : ", err)
	}
	fmt.Println("打开数据库成功")
	defer db.Close() //关闭数据库

	delegatorExist := [20]byte{6, 7, 8}
	delegatorUNExist := [20]byte{7, 7, 7}
	candidatorExist := [20]byte{0, 1, 2}
	candidatorUNExist := [20]byte{1, 1, 2}

	tests := []struct {
		name       string
		candidator [20]byte
		delegator  [20]byte
		want       int
	}{
		{"候选人不存在", candidatorUNExist, delegatorExist, -1},
		{"投票人不存在", candidatorExist, delegatorUNExist, -1},
		{"投票人已存在", candidatorExist, delegatorExist, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Delegate(db, tt.delegator, tt.candidator, big.NewInt(10))
			if err != nil {
				fmt.Println("Delegate err:", err)
			}
			candidates := ReadCandidates(db)
			candidateString, _ := json.Marshal(candidates)
			fmt.Printf("%s 情况下，candidates 读取结果:%s\n", tt.name, string(candidateString))
			if tt.name == "候选人不存在" {
				if _, ok := candidates[tt.candidator]; ok {
					t.Errorf("候选人实际上存在！")
				}
			} else {
				if got := candidates[tt.candidator].VotedNum.Cmp(big.NewInt(10)); got != tt.want {
					t.Errorf("投票失败！")
				}
			}

		})

	}
}

// 测试取消投票
func TestUnDelegate(t *testing.T) {
	InitDposAndDB()
	db, err := OpenDatabase(testDbFile, 512, 512)
	if err != nil {
		fmt.Println("failed to open database,the err info : ", err)
	}
	fmt.Println("打开数据库成功")
	defer db.Close() //关闭数据库
	delegatorAddr := [20]byte{6, 7, 8}
	candidatorAddr := [20]byte{0, 1, 2}
	err = Delegate(db, delegatorAddr, candidatorAddr, big.NewInt(10))
	if err != nil {
		fmt.Println("Delegate err:", err)
	}

	tests := []struct {
		name string
		args *big.Int
		want *big.Int
	}{
		{"已投票数大于取消票数", big.NewInt(11), big.NewInt(10)},
		{"已投票数小于取消票数", big.NewInt(9), big.NewInt(1)},
		{"已投票数与取消票数相等", big.NewInt(1), big.NewInt(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := UnDelegate(db, delegatorAddr, candidatorAddr, tt.args)
			if err != nil {
				fmt.Println("UnDelegate err: ", err)
			}
			candidates := ReadCandidates(db)
			votedNum := candidates[candidatorAddr].VotedNum
			if got := votedNum; got.Cmp(tt.want) != 0 {
				t.Errorf("取消票数不相等，现在的票数%v, 预期的票数 %v:", got, tt.want)
			}
		})
	}

}

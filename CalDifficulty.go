package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/params"
)

type Header struct {
	Time       *big.Int
	Difficulty *big.Int
	Number     *big.Int
}

var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big9          = big.NewInt(9)
	big10         = big.NewInt(10)
	bigMinus99    = big.NewInt(-99)
	big2999999    = big.NewInt(2999999)
)

func CalcDifficultyByzantium(time uint64, parent *Header) *big.Int {
	// https://github.com/ethereum/EIPs/issues/100.
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
	//        ) + 2^(periodCount - 2)

	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// (2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big9)

	x.Sub(big1, x)

	// max((2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// parent_diff + (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
	y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
	x.Mul(y, x)
	fmt.Println(x)
	x.Add(parent.Difficulty, x)
	fmt.Println(x)
	// minimum difficulty can ever be (before exponential factor)
	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}
	// calculate a fake block number for the ice-age delay:
	//   https://github.com/ethereum/EIPs/pull/669
	//   fake_block_number = min(0, block.number - 3_000_000
	fakeBlockNumber := new(big.Int)
	if parent.Number.Cmp(big2999999) >= 0 {
		fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, big2999999) // Note, parent is 1 less than the actual block number
	}
	// for the exponential factor
	periodCount := fakeBlockNumber
	periodCount.Div(periodCount, expDiffPeriod)
	fmt.Println(periodCount)
	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}
func Timetotimestamp(to string) int64 {
	timelayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timelayout, to, loc)
	sr := theTime.Unix()
	return sr
}

func main() {
	toBeTransformed1 := "2018-01-30 21:41:33" //待转换为世间戳的字符串
	toBeTransformed2 := "2018-01-30 21:41:41"
	sr1 := Timetotimestamp(toBeTransformed1)
	sr2 := Timetotimestamp(toBeTransformed2)
	fmt.Println("父区块时间: ", sr1)
	fmt.Println("带挖掘区块时间戳: ", sr2)
	parentTime := big.NewInt(sr1)
	parentNumber := big.NewInt(5000000)
	fmt.Println("父区块: ", parentNumber)
	parentDifficulty, _ := new(big.Int).SetString("2546613975853490", 10)
	var header *Header = &Header{parentTime, parentDifficulty, parentNumber}

	time := uint64(sr2)
	result := CalcDifficultyByzantium(time, header)
	fmt.Println(result)
}

package core

import "errors"

var (
	// ErrKnownBlock is returned when a block to import is already known locally.
	ErrKnownBlock = errors.New("block already known")

	// ErrGasLimitReached is returned by the gas pool if the amount of gas required
	// by a transaction is higher than what's left in the block.
	ErrGasLimitReached = errors.New("gas limit reached")

	// ErrBlacklistedHash is returned if a block to import is on the blacklist.
	ErrBlacklistedHash = errors.New("blacklisted hash")

	// ErrNonceTooHigh is returned if the nonce of a transaction is higher than the
	// next one expected based on the local chain.
	ErrNonceTooHigh = errors.New("nonce too high")

	ErrTimestampTooLow   = errors.New("交易时间小于最新时间")
	ErrTimeError         = errors.New("时间错误")
	ErrTimestampOutBound = errors.New("交易时间超出打包范围")
	ErrOversizedData     = errors.New("交易尺寸过大")
	ErrNegativeValue     = errors.New("交易金额为负")
	ErrUnderpriced       = errors.New("金额不足太低")
)

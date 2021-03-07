package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type PoolInfo struct {
	LpToken              common.Address `abi:"lpToken"` // Address of LP token contract.
	InitReward           *big.Int       `abi:"initreward"`
	Duration             *big.Int       `abi:"duration"`
	PeriodFinish         *big.Int       `abi:"periodFinish"`
	RewardRate           *big.Int       `abi:"rewardRate"`
	LastUpdateTime       *big.Int       `abi:"lastUpdateTime"`
	RewardPerTokenStored *big.Int       `abi:"rewardPerTokenStored"`
	TotalSupply          *big.Int       `abi:"totalSupply"`
	IsMigrated           bool           `abi:"isMigrated"`
}

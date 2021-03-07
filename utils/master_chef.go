package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tpkeeper/pitaya/types"
	"math/big"
	"fmt"
)

func GetPoolInfo(ethApi string, contractAddr common.Address, index *big.Int) (*types.PoolInfo, error) {
	methodName := "poolInfo"
	var poolInfo types.PoolInfo
	err := GetCall(ethApi, rewardPoolAbiJson, methodName, contractAddr, nil, &poolInfo, index)
	if err != nil {
		return nil, err
	}
	return &poolInfo, nil
}
func GetPoolLength(ethApi string, contractAddr common.Address) (*big.Int, error) {
	methodName := "totalID"
	type Amount struct {
		Amount *big.Int
	}
	amount := Amount{}
	err := GetCall(ethApi, rewardPoolAbiJson, methodName, contractAddr, nil, &amount)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return amount.Amount, nil
}

func GetHbtPerBlock(ethApi string, contractAddr common.Address) (*big.Int, error) {
	methodName := "hbtPerBlock"
	type Amount struct {
		Amount *big.Int
	}
	amount := Amount{}
	err := GetCall(ethApi, rewardPoolAbiJson, methodName, contractAddr, nil, &amount)
	if err != nil {
		return nil, err
	}
	return amount.Amount, nil
}

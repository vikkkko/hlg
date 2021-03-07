package utils

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

func GetTokenBalance(ethApi string, tokenContractAddr, userAddr common.Address, blockNumber *big.Int) (*big.Int, error) {
	packName := "balanceOf"
	type Amount struct {
		Amount *big.Int
	}
	amount := Amount{}

	err := GetCall(ethApi, erc20AbiJson, packName, tokenContractAddr, blockNumber, &amount, userAddr)

	if err != nil {
		return nil, err
	}
	return amount.Amount, nil
}

func GetDecimals(ethApi string, tokenContractAddr common.Address) (uint8, error) {
	packName := "decimals"
	type Amount struct {
		Amount uint8
	}
	amount := Amount{}

	err := GetCall(ethApi, erc20AbiJson, packName, tokenContractAddr, nil, &amount)
	if err != nil {
		return 0, err
	}
	return amount.Amount, nil
}

func GetSymbol(ethApi string, tokenContractAddr common.Address) (string, error) {
	methodName := "symbol"
	var symbol string
	err := GetCall(ethApi, erc20AbiJson, methodName, tokenContractAddr, nil, &symbol)
	if err != nil {
		return "", err
	}
	return symbol, nil
}

func GetTotalSupply(ethApi string, tokenContractAddr common.Address) (*big.Int, error) {
	packName := "totalSupply"
	type Amount struct {
		Amount *big.Int
	}
	amount := Amount{}

	err := GetCall(ethApi, erc20AbiJson, packName, tokenContractAddr, nil, &amount)
	if err != nil {
		return nil, err
	}
	return amount.Amount, nil
}

func GetCall(ethApi, abiJson, methodName string, contractAddr common.Address, blockNumber *big.Int,
	retValue interface{}, args ...interface{}) error {
	erc20Abi, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return err
	}
	packDataBts, err := erc20Abi.Pack(methodName, args...)
	if err != nil {
		return err
	}
	client, err := ethclient.Dial(ethApi)
	if err != nil {
		return err
	}
	ethCallMsg := ethereum.CallMsg{
		From: contractAddr,
		To:   &contractAddr,
		Data: packDataBts}

	resBts, err := client.CallContract(context.Background(), ethCallMsg, blockNumber)
	if err != nil {
		return err
	}
	return erc20Abi.UnpackIntoInterface(retValue, methodName, resBts)
}

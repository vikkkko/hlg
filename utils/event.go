package utils

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/tpkeeper/pitaya/models"
	"math/big"
	"strings"
)

func GetLogs(ethApi string, fromBlock *big.Int, contractAddresses []common.Address, topics [][]common.Hash) ([]types.Log, error) {
	client, err := ethclient.Dial(ethApi)
	if err != nil {
		return nil, err
	}
	query := ethereum.FilterQuery{
		Addresses: contractAddresses,
		Topics:    topics,
		FromBlock: fromBlock,
	}
	return client.FilterLogs(context.Background(), query)
}

//MasterChef
//
//0 event Deposit(address indexed user, uint256 indexed pid, uint256 amount);
//1 event Withdraw(address indexed user, uint256 indexed pid, uint256 amount);
//2 event EmergencyWithdraw(address indexed user, uint256 indexed pid, uint256 amount);
//3 event ProfitLock(address indexed user, uint256 indexed pid, uint256 pt, uint256 times);
//4 event ExtractReward(address indexed user, uint256 indexed pid, uint256 amount);

//HBTLock
//5 event Withdraw(address indexed user,uint256 unlockNumber);

func GetDepositAndWithLogs(ethApi string, fromBlock *big.Int, contractAddrs []common.Address) ([]*models.DepositWithDrawLog, error) {
	masterAbi, err := abi.JSON(strings.NewReader(rewardPoolAbiJson))
	if err != nil {
		return nil, err
	}
	hbtLockAbi, err := abi.JSON(strings.NewReader(hbtLockAbiJson))
	if err != nil {
		return nil, err
	}

	var depositTopic common.Hash
	var withdrawTopic common.Hash
	var emergencyWithdrawTopic common.Hash
	var extractRewardTopic common.Hash
	var profitLockTopic common.Hash
	var hbtLockWithdrawTopic common.Hash
	var playerBookEventTopic common.Hash

	var depositEvent abi.Event
	var withdrawEvent abi.Event
	var emergencyWithdrawEvent abi.Event
	var extractRewardEvent abi.Event
	var profitLockEvent abi.Event
	var hbtLockWithdrawEvent abi.Event
	var playerBookEvent abi.Event

	if event, exist := masterAbi.Events["Deposit"]; exist {
		depositEvent = event
		depositTopic = event.ID
	} else {
		return nil, errors.New("deposit not exist")
	}

	if event, exist := masterAbi.Events["Withdraw"]; exist {
		withdrawEvent = event
		withdrawTopic = event.ID
	} else {
		return nil, errors.New("withdraw not exist")
	}

	if event, exist := masterAbi.Events["EmergencyWithdraw"]; exist {
		emergencyWithdrawEvent = event
		emergencyWithdrawTopic = event.ID
	} else {
		return nil, errors.New("EmergencyWithdraw not exist")
	}

	if event, exist := masterAbi.Events["ExtractReward"]; exist {
		extractRewardEvent = event
		extractRewardTopic = event.ID
	} else {
		return nil, errors.New("ExtractReward not exist")
	}

	if event, exist := masterAbi.Events["ProfitLock"]; exist {
		profitLockEvent = event
		profitLockTopic = event.ID
	} else {
		return nil, errors.New("ProfitLock not exist")
	}

	if event, exist := hbtLockAbi.Events["Withdraw"]; exist {
		hbtLockWithdrawEvent = event
		hbtLockWithdrawTopic = event.ID
	} else {
		return nil, errors.New("hbtLock Withdraw not exist")
	}
	if event, exist := masterAbi.Events["PlayerBookEvent"]; exist {
		playerBookEvent = event
		playerBookEventTopic = event.ID
	} else {
		return nil, errors.New("PlayerBookEvent not exist")
	}

	logs, err := GetLogs(ethApi,
		fromBlock,
		contractAddrs,
		[][]common.Hash{{depositTopic, withdrawTopic, emergencyWithdrawTopic,
			extractRewardTopic, profitLockTopic, hbtLockWithdrawTopic, playerBookEventTopic}})

	if err != nil {
		return nil, err
	}

	depositWithdrawLogs := make([]*models.DepositWithDrawLog, 0)
	for _, log := range logs {
		l := models.DepositWithDrawLog{}

		var amount Amount
		switch hex.EncodeToString(log.Topics[0].Bytes()) {
		case hex.EncodeToString(depositTopic.Bytes()):
			l.Type = 0
			despositOrWithdraw := DepositOrWithdraw{}
			err := masterAbi.UnpackIntoInterface(&amount, "Deposit", log.Data)
			if err != nil {
				return nil, err
			}
			l.AmountBigInt = decimal.NewFromBigInt(amount.Amount, 0).String()
			//直解析前两个，只有前两个是topic，最后一个是data
			err = abi.ParseTopics(&despositOrWithdraw, depositEvent.Inputs[:2], log.Topics[1:])
			if err != nil {
				return nil, err
			}

			l.UserAddr = strings.ToLower(despositOrWithdraw.User.String())
			l.PoolId = despositOrWithdraw.Pid.Uint64()
		case hex.EncodeToString(withdrawTopic.Bytes()):
			l.Type = 1
			err := masterAbi.UnpackIntoInterface(&amount, "Withdraw", log.Data)
			if err != nil {
				return nil, err
			}
			l.AmountBigInt = decimal.NewFromBigInt(amount.Amount, 0).String()

			despositOrWithdraw := DepositOrWithdraw{}

			err = abi.ParseTopics(&despositOrWithdraw, withdrawEvent.Inputs[:2], log.Topics[1:])
			if err != nil {
				return nil, err
			}
			l.UserAddr = strings.ToLower(despositOrWithdraw.User.String())
			l.PoolId = despositOrWithdraw.Pid.Uint64()
		case hex.EncodeToString(emergencyWithdrawTopic.Bytes()):
			l.Type = 2
			err := masterAbi.UnpackIntoInterface(&amount, "EmergencyWithdraw", log.Data)
			if err != nil {
				return nil, err
			}
			l.AmountBigInt = decimal.NewFromBigInt(amount.Amount, 0).String()

			despositOrWithdraw := DepositOrWithdraw{}

			err = abi.ParseTopics(&despositOrWithdraw, emergencyWithdrawEvent.Inputs[:2], log.Topics[1:])
			if err != nil {
				return nil, err
			}
			l.UserAddr = strings.ToLower(despositOrWithdraw.User.String())
			l.PoolId = despositOrWithdraw.Pid.Uint64()
		case hex.EncodeToString(profitLockTopic.Bytes()):
			l.Type = 3
			profitData := ProfitData{}
			err := masterAbi.UnpackIntoInterface(&profitData, "ProfitLock", log.Data)
			if err != nil {
				return nil, err
			}
			l.Times = uint64(decimal.NewFromBigInt(profitData.Times, 0).IntPart())
			l.Pt, _ = decimal.NewFromBigInt(profitData.Pt, -18).Float64()

			despositOrWithdraw := DepositOrWithdraw{}

			err = abi.ParseTopics(&despositOrWithdraw, profitLockEvent.Inputs[:2], log.Topics[1:3])
			if err != nil {
				return nil, err
			}
			l.UserAddr = strings.ToLower(despositOrWithdraw.User.String())
			l.PoolId = despositOrWithdraw.Pid.Uint64()
		case hex.EncodeToString(extractRewardTopic.Bytes()):
			l.Type = 4
			err := masterAbi.UnpackIntoInterface(&amount, "ExtractReward", log.Data)
			if err != nil {
				return nil, err
			}
			l.AmountBigInt = decimal.NewFromBigInt(amount.Amount, 0).String()

			despositOrWithdraw := DepositOrWithdraw{}

			err = abi.ParseTopics(&despositOrWithdraw, extractRewardEvent.Inputs[:2], log.Topics[1:])
			if err != nil {
				return nil, err
			}
			l.UserAddr = strings.ToLower(despositOrWithdraw.User.String())
			l.PoolId = despositOrWithdraw.Pid.Uint64()
		case hex.EncodeToString(hbtLockWithdrawTopic.Bytes()):
			l.Type = 5
			unlockNumber := UnlockNumber{}
			err := hbtLockAbi.UnpackIntoInterface(&unlockNumber, "Withdraw", log.Data)
			if err != nil {
				return nil, err
			}
			l.UnlockNumber, _ = decimal.NewFromBigInt(unlockNumber.UnlockNumber, -18).Float64()

			despositOrWithdraw := DepositOrWithdraw{}

			err = abi.ParseTopics(&despositOrWithdraw, hbtLockWithdrawEvent.Inputs[:1], log.Topics[1:])
			if err != nil {
				return nil, err
			}
			l.UserAddr = strings.ToLower(despositOrWithdraw.User.String())
		case hex.EncodeToString(playerBookEventTopic.Bytes()):
			l.Type = 6
			err := masterAbi.UnpackIntoInterface(&amount, "PlayerBookEvent", log.Data)
			if err != nil {
				return nil, err
			}

			l.AmountBigInt = decimal.NewFromBigInt(amount.Amount, 0).String()

			//过滤为0的log
			if amount.Amount.Int64() == 0 {
				continue
			}
			despositOrWithdraw := DepositOrWithdraw{}

			err = abi.ParseTopics(&despositOrWithdraw, playerBookEvent.Inputs[:2], log.Topics[1:])
			if err != nil {
				return nil, err
			}
			l.UserAddr = strings.ToLower(despositOrWithdraw.User.String())
		default:
			return nil, errors.New("no case match")
		}

		l.TransactionHash = strings.ToLower(log.TxHash.String())
		l.BlockNumber = log.BlockNumber
		client, err := ethclient.Dial(ethApi)
		if err != nil {
			return nil, err
		}
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
		if err != nil {
			return nil, err
		}
		l.BlockTime = block.Time()
		l.LogIndex = log.Index
		depositWithdrawLogs = append(depositWithdrawLogs, &l)
	}

	return depositWithdrawLogs, nil

}

type DepositOrWithdraw struct {
	Amount   *big.Int       `abi:"amount"`
	Pid      *big.Int       `abi:"pid"`
	User     common.Address `abi:"user"`
	FromUser common.Address `abi:"fromUser"`
}
type Amount struct {
	Amount *big.Int `abi:"amount"`
}
type UnlockNumber struct {
	UnlockNumber *big.Int
}

type ProfitData struct {
	Pt    *big.Int
	Times *big.Int
}

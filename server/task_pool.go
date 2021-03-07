package server

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tpkeeper/pitaya/utils"
	"gorm.io/gorm"
	"math/big"
	"strings"
	"fmt"
)

func (svr *Server) ExeTaskPool() error {
	fmt.Println(masterChiefContract)
	length, err := utils.GetPoolLength(svr.ethApi, masterChiefContract)
	if err != nil {
		return err
	}
	for i := int64(0); i < length.Int64(); i++ {
		//获取合约中第i个池子
		poolInContract, err := utils.GetPoolInfo(svr.ethApi, masterChiefContract, big.NewInt(i))
		if err != nil {
			logrus.Errorf("utils.GetPoolInfo err %s", err.Error())
			continue
		}
		logrus.Infof("pool id %d pool addr %v", i, poolInContract.LpToken.String())

		//获取该池子在数据库中的记录
		pool, err := svr.dao.GetPoolByPoolId(uint64(i))
		if err != nil && err != gorm.ErrRecordNotFound {
			logrus.Errorf("dao.GetPoolByPoolAddr %s", err.Error())
			continue
		}

		if len(pool.Symbol) == 0 {
			//获取symbol
			symbol, err := utils.GetSymbol(svr.ethApi, poolInContract.LpToken)
			if err != nil {
				logrus.Errorf("utils.GetSymbol err %s", err.Error())
				continue
			}
			pool.Symbol = symbol
		}

		//获取抵押数量
		//depositAmount, err := utils.GetTokenBalance(svr.ethApi, poolInContract.LpToken, masterChiefContract, nil)
		//if err != nil {
		//	logrus.Errorf("utils.GetTokenBalance err %s", err.Error())
		//	continue
		//}
		//logrus.Infof("depositAmount %s", depositAmount.String())

		depositAmount := poolInContract.TotalSupply

		if pool.Decimals == 0 {
			//获取精度
			decimals, err := utils.GetDecimals(svr.ethApi, poolInContract.LpToken)
			if err != nil {
				logrus.Errorf("utils.GetDecimals err %s", err.Error())
				continue
			}
			logrus.Infof("decimals %d", decimals)
			pool.Decimals = decimals
		}

		//计算真实的抵押数量
		depositFloat, _ := decimal.NewFromBigInt(depositAmount, 0-int32(pool.Decimals)).Float64()
		pool.DepositAmount = depositFloat
		logrus.Infof("deposit float %f", depositFloat)
		//获取该池子抵押代币价格
		//if "0x32fd949e1953b21b7a8232ef4259cd708b4e0847" == strings.ToLower(poolInContract.LpToken.String()) {
		//	price, err := utils.GetPitayaPrice()
		//	if err != nil {
		//		logrus.Errorf("utils.GetPitayaPrice err %s", err.Error())
		//	} else {
		//		pool.Price = price
		//	}
		//} else {

		//todo 映射关系从配置文件中读取
		if name, exist := GetPriceFromCoinw[strings.ToLower(poolInContract.LpToken.String())]; exist {
			priceFloat, err := utils.GetPriceFromCoinw(name)
			if err != nil {
				logrus.Errorf("utils.GetPriceFromCoinw err %s,name %s", err.Error(), name)
			} else {
				pool.Price = priceFloat
			}
		} else {

			price, err := utils.GetLpPrice(poolInContract.LpToken.String())
			if err != nil {
				price, err = utils.GetTokenPrice(poolInContract.LpToken.String())
				if err != nil {
					logrus.Errorf("utils.GetTokenPrice err %s addr: %s",
						err.Error(), poolInContract.LpToken.String())
				} else {
					pool.Price = utils.StrToFloat(price)
				}
			} else {
				pool.Price = utils.StrToFloat(price)
			}
		}
		//}

		//totalAlloc, err := svr.dao.GetTotalAllocPoint()
		//if err == nil {
		pitaya, err := svr.dao.GetPitaya()
		if err == nil {
			//更新该池子每个块产量
			//if totalAlloc != 0 {
			//	pool.RewardPerTokenStored, _ = decimal.NewFromFloat(pitaya.RewardPerTokenStored).
			//		Mul(decimal.NewFromInt(int64(pool.AllocPoint))).
			//		Div(decimal.NewFromInt(int64(totalAlloc))).Float64()
			//}
			durationDeci:= decimal.NewFromBigInt(poolInContract.Duration, 0)
			duration, _ :=durationDeci.Float64()
			initRewardDeci := decimal.NewFromBigInt(poolInContract.InitReward, 0)
			//更新该池子 apy
			if pool.DepositAmount != 0 && pool.Price != 0 && duration > 0 {
				rate, _ := initRewardDeci.Div(durationDeci).Float64()

				pool.Apy = rate * 365 * 24 * 60 * 60 * pitaya.Price * 100 / pool.Price * pool.DepositAmount

			}
		}
		//}

		//更新池子信息
		pool.PoolId = uint64(i)
		pool.Address = poolInContract.LpToken.String()
		pool.RewardPerTokenStored, _ = decimal.NewFromBigInt(poolInContract.RewardPerTokenStored, 0-18).Float64()

		err = svr.dao.UpdateOrInsertPool(pool)
		if err != nil {
			logrus.Errorf("dao.UpdateOrInsertPool err %s", err.Error())
			continue
		}
	}
	return nil
}

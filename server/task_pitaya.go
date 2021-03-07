package server

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tpkeeper/pitaya/utils"
	"gorm.io/gorm"
)

func (svr *Server) ExeTaskPitaya() error {

	//decimals, err := utils.GetDecimals(svr.ethApi, pitayaTokenContract)
	//if err != nil {
	//	return err
	//}

	decimals:=18
	hbt, err := svr.dao.GetPitaya()
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	//更新 total supply
	totalSupply, err := utils.GetTotalSupply(svr.ethApi, pitayaTokenContract)
	if err != nil {
		logrus.Errorf("utils.GetTotalSupply err %s", err.Error())
	} else {
		totalSupplyFloat, _ := decimal.NewFromBigInt(totalSupply, 0-int32(decimals)).Float64()
		hbt.TotalSupply = totalSupplyFloat
	}

	//更新锁定中的hbt数量
	//lockTotal, err := utils.GetTotalSupply(svr.ethApi, hbtLockTokenContract)
	//lockTotal, err := utils.GetTokenBalance(svr.ethApi, pitayaTokenContract, hbtLockTokenContract, nil)
	//if err != nil {
	//	logrus.Error("utils.GetTotalSupply %s", err.Error())
	//} else {
	//	lockTotalFloat, _ := decimal.NewFromBigInt(lockTotal, 0-int32(decimals)).Float64()
	//	hbt.LockAmount = lockTotalFloat
	//}

	//更新price
	//price, err := utils.GetTokenPrice(pitayaTokenContract.String())
	//if err != nil {
	//	logrus.Errorf("utils.GetTokenPrice err %s coins %s", err.Error(), "hbt")
	//} else {
	//	hbt.Price = utils.StrToFloat(price)
	//}
	//如果获取的为0 更换获取方式
	//if hbt.Price == 0 {
	priceFloat, err := utils.GetPitayaPrice()
	if err != nil {
		logrus.Errorf("utils.GetPitayaPrice err %s", err)
	} else {
		hbt.Price = priceFloat
	}
	//}

	//更新mfiPerBlock
	//mfiPerBlock, err := utils.GetHbtPerBlock(svr.ethApi, masterChiefContract)
	//if err != nil {
	//	logrus.Errorf("utils.GetHbtPerBlock err %s", err.Error())
	//} else {
	//	mfiPerBlockFloat, _ := decimal.NewFromBigInt(mfiPerBlock, 0-int32(decimals)).Float64()
	//	hbt.PitayaPerBlock = mfiPerBlockFloat
	//}

	//更新apy
	//pools, err := svr.dao.GetAllPools()
	//if err == nil {
	//	var total float64
	//	for _, p := range pools {
	//		total += p.Price * p.DepositAmount
	//	}
	//	if total != 0 {
	//		hbt.Apy = hbt.PitayaPerBlock * 365 * 6750 * hbt.Price * 100 / total
	//	}
	//}

	return svr.dao.UpdateOrInsertHbt(hbt)
}

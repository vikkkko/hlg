package server

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

type PoolInfo struct {
	Index    uint64  `json:"index"`
	Addr     string  `json:"addr"`
	Apy      float64 `json:"apy"`
	Weight   uint64  `json:"weight"`
	Name     string  `json:"name"`
	LogoUrl1 string  `json:"logo_url_1"`
	LogoUrl2 string  `json:"logo_url_2"`
	//Price         float64 `json:"price"`
	//PitayaPrice   float64 `PitayaPrice`
	//DepositAmount float64 `json:"deposit_amount"`
	//DepositValue  float64 `json:"deposit_value"`
	//HbtPerDay     float64 `json:"hbt_per_day"`
}

func (svr *Server) routerGetAllPools(ctx *gin.Context) {
	pools, err := svr.dao.GetAllPools()
	if err != nil {
		Err(ctx, err.Error())
		return
	}

	rspPools := make([]PoolInfo, 0)
	for _, p := range pools {
		pool := PoolInfo{}
		pool.Apy = p.Apy
		pool.Addr = p.Address
		pool.Index = p.PoolId
		//pool.DepositAmount = p.DepositAmount
		//pool.DepositValue = p.Price * p.DepositAmount
		//pool.HbtPerDay = p.RewardPerTokenStored * 5760
		//pool.Price = p.Price
		rspPools = append(rspPools, pool)
	}
	Ok(ctx, rspPools)
}

func (svr *Server) routerGetPoolByAddr(ctx *gin.Context) {
	addr := ctx.Query("addr")
	if !common.IsHexAddress(addr) {
		Err(ctx, "not a eth address")
		return
	}
	p, err := svr.dao.GetPoolByPoolAddr(strings.ToLower(addr))
	if err != nil {
		Err(ctx, err.Error())
		return
	}
	//hbt, err := svr.dao.GetPitaya()
	//if err != nil {
	//	Err(ctx, err.Error())
	//	return
	//}
	poolInfo, err := svr.dao.GetPoolInfoByPoolAddr(strings.ToLower(addr))
	if err != nil {
		logrus.Errorf("svr.dao.GetPoolInfoByPoolAddr err %s", err)
		Err(ctx, err.Error())
		return
	}

	pool := PoolInfo{}
	pool.Apy = p.Apy
	pool.Addr = p.Address
	pool.Index = p.PoolId
	pool.Name = poolInfo.Name
	pool.Weight = poolInfo.Weight
	pool.LogoUrl1 = poolInfo.LogoUrl1
	pool.LogoUrl2 = poolInfo.LogoUrl2
	//pool.DepositAmount = p.DepositAmount
	//pool.DepositValue = p.Price * p.DepositAmount
	//pool.HbtPerDay = p.RewardPerTokenStored * 5760
	//pool.Price = p.Price
	//pool.PitayaPrice = hbt.Price

	Ok(ctx, pool)
}

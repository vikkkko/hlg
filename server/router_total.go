package server

import "github.com/gin-gonic/gin"

type Total struct {
	TotalDepositValue float64 `json:"total_deposit_value"`
	TotalProfit       float64 `json:"total_profit"` //累计收益
}

func (svr *Server) routerGetTotal(ctx *gin.Context) {
	pools, err := svr.dao.GetAllPools()
	if err != nil {
		Err(ctx, err.Error())
		return
	}

	var total Total
	for _, p := range pools {
		total.TotalDepositValue += p.Price * p.DepositAmount
	}

	hbt, err := svr.dao.GetPitaya()
	if err != nil {
		Err(ctx, err.Error())
		return
	}

	total.TotalProfit = hbt.Price * hbt.TotalSupply
	Ok(ctx, total)
}

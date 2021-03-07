package server

import "github.com/gin-gonic/gin"

type HbtInfo struct {
	Price       float64 `json:"price"`
	MarketCap   float64 `json:"market_cap"`
	TotalSupply float64 `json:"total_supply"`
	LockAmount  float64 `json:"lock_amount"`
	LockValue   float64 `json:"lock_value"`
	HbtPerBlock float64 `json:"hbt_per_block"`
}

func (svr *Server) routerGetHbtInfo(ctx *gin.Context) {
	hbt, err := svr.dao.GetPitaya()
	if err != nil {
		Err(ctx, err.Error())
		return
	}

	hbtInfo := HbtInfo{}
	hbtInfo.Price = hbt.Price
	hbtInfo.MarketCap = hbt.Price * hbt.TotalSupply
	hbtInfo.TotalSupply = hbt.TotalSupply
	//hbtInfo.LockAmount = hbt.LockAmount
	//hbtInfo.LockValue = hbt.LockAmount * hbt.Price
	hbtInfo.HbtPerBlock = hbt.PitayaPerBlock

	Ok(ctx, hbtInfo)
}

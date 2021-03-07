package server

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tpkeeper/pitaya/models"
	"gorm.io/gorm"
	"strings"
)

type ReqPoolInfo struct {
	TokenAddress string `json:"token_address"`
	Weight       uint64 `json:"weight"`
	Name         string `json:"name"`
	LogoUrl1     string `json:"logo_url_1"`
	LogoUrl2     string `json:"logo_url_2"`
	Toukuang	 uint64 `json:"toukuang"`
}

func (svr *Server) checkLogoUrl(url string) error {
	//目前只验证logo图片url对应的文件是否存在本地服务器
	pos := strings.Index(url, "/static/file/")
	if pos < 0 {
		return fmt.Errorf("%s is not safe", url)
	}
	if !common.FileExist(fmt.Sprintf(".%s", url[pos:])) {
		return fmt.Errorf("%s is not exists", url)
	}
	return nil
}

//上币提案提交
func (svr *Server) handlePostPoolInfo(c *gin.Context) {

	req := ReqPoolInfo{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Err(c, "invalid param")
		logrus.Printf("params bind error, %v, params:%+v", err, req)
		return
	}
	if !common.IsHexAddress(req.TokenAddress) {
		Err(c, "ethAddress is not right")
		return
	}

	coinInfo := models.PoolInfo{
		TokenAddress: req.TokenAddress,
	}

	txDao := svr.dao

	poolInfo, err := txDao.GetPoolInfoByPoolAddr(coinInfo.TokenAddress)
	if err != nil && err != gorm.ErrRecordNotFound {
		Err(c, err.Error())
		logrus.Errorf("txDao.GetPoolInfoByPoolAddr err %s", err.Error())
		return
	}

	poolInfo.TokenAddress = coinInfo.TokenAddress
	poolInfo.Name = req.Name
	poolInfo.Weight = req.Weight
	poolInfo.LogoUrl1 = req.LogoUrl1
	poolInfo.LogoUrl2 = req.LogoUrl2
	poolInfo.Toukuang = req.Toukuang

	err = txDao.UpdateOrInsertPoolInfo(poolInfo)
	if err != nil {
		Err(c, "coin info post error")
		params, _ := json.Marshal(coinInfo)
		logrus.Errorf("txDao.UpdateOrInsertCoinInfo err %s, params:%s", err.Error(), string(params))
		return
	}

	Ok(c)
}

type ReqPoolWeight struct {
	TokenAddress string `json:"token_address"`
	Weight       uint64 `json:"weight"`
	Toukuang	 uint64 `json:"toukuang"`
}

//上币提案提交
func (svr *Server) handleUpdatePoolInfoWeight(c *gin.Context) {

	req := ReqPoolWeight{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Err(c, "invalid param")
		logrus.Printf("params bind error, %v, params:%+v", err, req)
		return
	}
	if !common.IsHexAddress(req.TokenAddress) {
		Err(c, "ethAddress is not right")
		return
	}

	coinInfo := models.PoolInfo{
		TokenAddress: req.TokenAddress,
	}

	txDao := svr.dao

	poolInfo, err := txDao.GetPoolInfoByPoolAddr(coinInfo.TokenAddress)
	if err != nil {
		Err(c, err.Error())
		logrus.Errorf("txDao.GetPoolInfoByPoolAddr err %s", err.Error())
		return
	}

	poolInfo.TokenAddress = coinInfo.TokenAddress
	poolInfo.Weight = req.Weight
	poolInfo.Toukuang = req.Toukuang

	err = txDao.UpdateOrInsertPoolInfo(poolInfo)
	if err != nil {
		Err(c, "coin info post error")
		params, _ := json.Marshal(coinInfo)
		logrus.Errorf("txDao.UpdateOrInsertCoinInfo err %s, params:%s", err.Error(), string(params))
		return
	}

	Ok(c)
}

type RspPoolInfo struct {
	Index    uint64 `json:"index"`
	Address  string `json:"addr"`
	Weight   uint64 `json:"weight"`
	Name     string `json:"name"`
	LogoUrl1 string `json:"logo_url_1"`
	LogoUrl2 string `json:"logo_url_2"`
	Toukuang uint64 `json:"toukuang"`
}

func (svr *Server) handleGetPoolInfo(c *gin.Context) {
	pools, err := svr.dao.GetAllPools()
	if err != nil {
		Err(c, err.Error())
		return
	}

	rspPools := make([]RspPoolInfo, 0)
	for i, pool := range pools {
		poolInfo, err := svr.dao.GetPoolInfoByPoolAddr(pool.Address)
		if err != nil {
			logrus.Errorf("svr.dao.GetPoolInfoByPoolAddr err %s", err)
			continue
		}
		rspPool := RspPoolInfo{
			Index:    uint64(i),
			Address:  pool.Address,
			Weight:   poolInfo.Weight,
			Name:     poolInfo.Name,
			LogoUrl1: poolInfo.LogoUrl1,
			LogoUrl2: poolInfo.LogoUrl2,
			Toukuang: poolInfo.Toukuang,
		}

		rspPools = append(rspPools, rspPool)

	}
	Ok(c, rspPools)

}

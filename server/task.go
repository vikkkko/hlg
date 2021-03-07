package server

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/tpkeeper/pitaya/config"
	"time"
)

var masterChiefContract = common.Address{}
var pitayaTokenContract = common.Address{}

var GetPriceFromCoinw = map[string]string{
	"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": "eth",
	"0xdac17f958d2ee523a2206206994597c13d831ec7": "usdt",
	"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48": "usdc",
	"0x6b175474e89094c44da98b954eedeac495271d0f": "dai",
	"0x0db02e277d4364136623e9ade2e2e5f06202ad11": "hpy",
	"0x97cc1e26582763bf31f3434c552dfa7f748816ce": "hc",
	"0xc348f56e479d94bc117305df01ca6c2ebdfcc76e": "dash",
	"0x378052d5d6b316f12f55d4622b7e9cd31ba4cf0a": "qtum",
	"0xc97265e4a0b19c173a5cf9f394602a149654bdce": "fil",
	"0xb35c540e4412363858c8a49ec0aa8c9bc76c8cca": "bch",
	"0x90e4f6d98e53a974c06cf047e6fcbb2f2ce482ef": "btc",
	"0x957339c0586ba22d472ef6f579749ee9439bf85d": "iost",
	"0x2689a1d35ad5d656c1fb9468dd007ead6c3fde6c": "nuls",
}

func (svr *Server) Task() {
	masterChiefContract = common.HexToAddress(config.GetMasterChiefContract())
	pitayaTokenContract = common.HexToAddress(config.GetPitayaTokenContract())
	ticker := time.NewTicker(time.Duration(svr.tickerInterval) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			logrus.Println("Task start **********")
			err := svr.ExeTaskPitaya()
			if err != nil {
				logrus.Errorf("Task err: %s", err.Error())
			}
			logrus.Println("Task End **********")

			logrus.Println("TaskPool start **********")
			err = svr.ExeTaskPool()
			if err != nil {
				logrus.Errorf("ExePool err: %s", err.Error())
			}
			logrus.Println("TaskPool End **********")

		}

	}
}

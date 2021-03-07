package utils

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tpkeeper/pitaya/types"
	"io/ioutil"
	"net/http"
	"strconv"
)

var hpyPriceUrl = "http://47.88.158.38//hpymarket/api/all_coin_info?coins=%s"

func GetPriceFromCoinw(coins string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf(hpyPriceUrl, coins))
	if err != nil {
		return 0, fmt.Errorf("get hpyPriceUrl err %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("get price Url statue %d", resp.StatusCode)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ioutil.ReadAll err %s", err.Error())
	}
	rspPrice := types.RspPrice{}
	json.Unmarshal(bts, &rspPrice)

	if rspPrice.Status != "200" || len(rspPrice.Data) == 0 {
		return 0, fmt.Errorf("get hpyPriceUrl status %s,data len %d", rspPrice.Status, len(rspPrice.Data))
	}
	var priceFloat float64
	switch rspPrice.Data[0].Price.(type) {
	case string:
		priceDeci, err := decimal.NewFromString(rspPrice.Data[0].Price.(string))
		if err != nil {
			return 0, err
		}
		priceFloat, _ = priceDeci.Float64()
	case float64:
		priceFloat = rspPrice.Data[0].Price.(float64)
	}

	return priceFloat, nil
}

func StrToFloat(str string) float64 {
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return v
}

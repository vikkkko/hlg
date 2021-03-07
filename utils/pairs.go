package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var transactionGraphqlUniswapV2 = "https://api.thegraph.com/subgraphs/name/ianlapham/uniswapv2"
var transactionGraphqlMainNetMiniV2 = "https://api.thegraph.com/subgraphs/name/noberk/chapter3"
var transactionGraphqlTestNetMiniV2 = "https://api.thegraph.com/subgraphs/name/noberk/chapter4"

var transactionGraphqlApi = transactionGraphqlUniswapV2
//var transactionGraphqlApi = transactionGraphqlMainNetMiniV2

func GetLpPrice(id string) (string, error) {
	id = strings.ToLower(id)
	type requestPairs struct {
		Data struct {
			Pair struct {
				ID          string `json:"id"`
				ReserveUSD  string `json:"reserveUSD"`
				TotalSupply string `json:"totalSupply"`
			} `json:"pair"`
		} `json:"data"`
	}
	queryStr := `
            { 
                pair(id:"%s") {
   				 	id
    				totalSupply
    				reserveUSD
  				} 
            }
        `
	jsonData := map[string]string{
		"query": fmt.Sprintf(queryStr, id),
	}
	jsonValue, _ := json.Marshal(jsonData)
	fmt.Println(string(jsonValue))
	request, err := http.NewRequest("POST", transactionGraphqlApi, bytes.NewBuffer(jsonValue))
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("status %d", response.StatusCode))
	}

	repsdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var pausesRespone requestPairs
	err = json.Unmarshal(repsdata, &pausesRespone)
	if err != nil {
		return "", err
	}

	reserveUsdDeci, err := decimal.NewFromString(pausesRespone.Data.Pair.ReserveUSD)
	if err != nil {
		return "", err
	}
	totalSupplyDeci, err := decimal.NewFromString(pausesRespone.Data.Pair.TotalSupply)
	if err != nil {
		return "", err
	}

	priceDeci := reserveUsdDeci.Div(totalSupplyDeci)

	return priceDeci.StringFixed(6), nil
}

func GetTokenPrice(id string) (string, error) {
	return "1.000000", nil
	id = strings.ToLower(id)
	type requestPairs struct {
		Data struct {
			TokenDayDatas []struct {
				ID       string `json:"id"`
				PriceUSD string `json:"priceUSD"`
			} `json:"tokenDayDatas"`
		} `json:"data"`
	}

	queryStr := `
            { 
                tokenDayDatas(first:1,orderBy: date, orderDirection: desc,where:{token:"%s"}) {
   				 	id
    				priceUSD
  				} 
            }
        `
	jsonData := map[string]string{
		"query": fmt.Sprintf(queryStr, id),
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", transactionGraphqlApi, bytes.NewBuffer(jsonValue))
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("status %d", response.StatusCode))
	}

	repsdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var pausesRespone requestPairs
	err = json.Unmarshal(repsdata, &pausesRespone)
	if err != nil {
		return "", err
	}

	if len(pausesRespone.Data.TokenDayDatas) == 0 {
		return "", errors.New("tokenDayDatas len is 0")
	}
	priceUsdDeci, err := decimal.NewFromString(pausesRespone.Data.TokenDayDatas[0].PriceUSD)
	if err != nil {
		return "", err
	}

	return priceUsdDeci.StringFixed(6), nil
}

var queryPriceStr = `{	pair(id:"0x13cf64aacb033bb6a8cecd6d1e9eaf8ab4250022"){
							id
							token1Price
						}
					}`
type requestTransactions struct {
	Data struct {
		Pair struct {
			Id          string `json:"id"`
			Token0Price string `json:"token1Price"`
		} `json:"pair"`
	} `json:"data"`
}

func GetPitayaPrice() (float64, error) {
	jsonData := map[string]string{
		"query": queryPriceStr,
	}

	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", transactionGraphqlApi, bytes.NewBuffer(jsonValue))
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, errors.New(fmt.Sprintf("status:%d", response.StatusCode))
	}

	repsdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	var transactionsRespone requestTransactions
	err = json.Unmarshal(repsdata, &transactionsRespone)
	if err != nil {
		return 0, err
	}

	return StrToFloat(transactionsRespone.Data.Pair.Token0Price), nil

}
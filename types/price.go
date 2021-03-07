package types

type RspPrice struct {
	Status string         `json:"status"`
	Data   []RspPriceData `json:"data"`
	Msg    string         `json:"msg"`
}

type RspPriceData struct {
	ID        string      `json:"id"`
	Coin      string      `json:"coin"`
	Cn        string      `json:"cn"`
	En        string      `json:"en"`
	Price     interface{} `json:"price"`
	T         string      `json:"t"`
	Volume24  string      `json:"volume24"`
	Logo      string      `json:"logo"`
	Circulate string      `json:"circulate"`
	Supply    string      `json:"supply"`
	Percent   float64     `json:"percent"`
	Max24     string      `json:"max24"`
	Min24     string      `json:"min24"`
}

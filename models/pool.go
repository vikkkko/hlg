package models

type Pool struct {
	BaseModel
	PoolId               uint64  `gorm:"index"` //在合约中的索引
	Address              string  //抵押代币 erc合约地址
	Decimals             uint8   //抵押代币精度
	Symbol               string  //符号
	Price                float64 //价格
	DepositAmount        float64 //抵押数量
	AllocPoint           uint64  //权重
	RewardPerTokenStored float64
	Apy                  float64
}

type PoolInfo struct {
	BaseModel
	TokenAddress string `json:"token_address"`
	Weight       uint64 `json:"weight"`
	Name         string `json:"name"`
	LogoUrl1     string `json:"logo_url_1"`
	LogoUrl2     string `json:"logo_url_2"`
	Toukuang	 uint64	`json:"toukuang"`
}

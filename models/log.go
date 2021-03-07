package models

type DepositWithDrawLog struct {
	BaseModel       `json:"-"`
	TransactionHash string  `gorm:"uniqueIndex:index_tx_type" json:"transaction_hash"`
	LogIndex        uint    `gorm:"uniqueIndex:index_tx_type" json:"-"`
	BlockNumber     uint64  `gorm:"index" json:"block_number"`
	BlockTime       uint64  `json:"block_time"`
	UserAddr        string  `json:"user_addr"`
	PoolId          uint64  `json:"pool_id"`
	Amount          float64 `json:"amount"`
	AmountBigInt    string  `json:"-"`
	Symbol          string  `json:"symbol"`
	Pt              float64 `json:"pt"`                                    //type 为3时候有意义
	Times           uint64  `json:"times"`                                 //type 为3时候有意义
	UnlockNumber    float64 `json:"unlock_number"`                         //type 为5时有意义
	Type            uint8   `gorm:"uniqueIndex:index_tx_type" json:"type"` //todo 还有其他类型 0 deposit 1 withraw 2 EmergencyWithdraw 3 ProfitLock 4 ExtractReward
}

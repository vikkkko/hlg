package models

type Pitaya struct {
	BaseModel
	TotalSupply    float64
	Price          float64
	PitayaPerBlock float64
	Apy            float64
	LockAmount     float64 //锁定数量
}

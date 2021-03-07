package dao

import "github.com/tpkeeper/pitaya/models"

func (d *Dao) UpdateOrInsertPool(pool *models.Pool) error {
	return d.db.Save(pool).Error
}

func (d *Dao) GetAllPools() (pools []*models.Pool, err error) {
	err = d.db.Find(&pools).Error
	return
}

func (d *Dao) GetPoolByPoolAddr(addr string) (pool *models.Pool, err error) {
	pool = &models.Pool{}
	err = d.db.Take(pool, "address = ?", addr).Error
	return
}

func (d *Dao) GetPoolByPoolId(id uint64) (pool *models.Pool, err error) {
	pool = &models.Pool{}
	err = d.db.Take(pool, "pool_id = ?", id).Error
	return
}

func (d *Dao) GetTotalAllocPoint() (amountRet uint64, err error) {
	type Amount struct {
		Amount uint64
	}
	amount := Amount{}
	err = d.db.Raw("select sum(alloc_point) as amount from pools").Scan(&amount).Error
	amountRet = amount.Amount
	return
}

func (d *Dao) UpdateOrInsertPoolInfo(pool *models.PoolInfo) error {
	return d.db.Save(pool).Error
}

func (d *Dao) GetPoolInfoByPoolAddr(addr string) (pool *models.PoolInfo, err error) {
	pool = &models.PoolInfo{}
	err = d.db.Take(pool, "token_address = ?", addr).Error
	return
}

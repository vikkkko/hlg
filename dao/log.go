package dao

import (
	"github.com/tpkeeper/pitaya/models"
	"strings"
)

func (d *Dao) UpdateOrInsertDepositWithDrawLog(log *models.DepositWithDrawLog) error {
	return d.db.Save(log).Error
}

func (d *Dao) GetDepositWithDrawLogByTxHashLogIndex(txHash string, logIndex uint) (log *models.DepositWithDrawLog, err error) {
	txHash = strings.ToLower(txHash)
	log = &models.DepositWithDrawLog{}
	err = d.db.Take(log, "transaction_hash = ? and log_index = ?", txHash, logIndex).Error
	return
}
func (d *Dao) GetDepositWithDrawLogsLimit(limit int) (logs []*models.DepositWithDrawLog, err error) {
	err = d.db.Order("block_number desc").Limit(limit).Find(&logs).Error
	return
}

func (d *Dao) GetDepositWithDrawLogsByUserLimit(page int, pageSize int, user string) (logs []*models.DepositWithDrawLog, err error) {
	if page <= 0 {
		page = 1
	}
	err = d.db.Order("block_number desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs, "user_addr = ?", user).Error
	return
}

func (d *Dao) GetDepositWithDrawLogsLimit2(limit int, id uint64) (logs []*models.DepositWithDrawLog, err error) {
	err = d.db.Order("block_number desc").Limit(limit).Find(&logs, "pool_id = ?", id).Error
	return
}

func (d *Dao) GetLogsTopBlockNumber() (amount uint64, err error) {
	type BlockNumber struct {
		BlockNumber int64
	}
	block := BlockNumber{}
	err = d.db.Raw("select block_number from deposit_with_draw_logs order by block_number desc limit 1").
		Scan(&block).Error
	amount = uint64(block.BlockNumber)
	return amount, err
}

func (d *Dao) GetLogsTotalNumberByUser(user string) (number int64, err error) {
	//type Amount struct{
	//	Amount int64
	//}
	//number = Amount{}
	err = d.db.Raw("select count(*) as number from deposit_with_draw_logs where user_addr = ?", user).Scan(&number).Error
	return
}

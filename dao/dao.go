package dao

import (
	"errors"
	"gorm.io/gorm"
)

type Dao struct {
	db   *gorm.DB
	isTx bool
}

func New(db *gorm.DB) *Dao {
	return &Dao{
		db: db,
	}
}

//事务开始，调用完此方法必须调用 rollback 或 commitTransaction
func (d *Dao) NewTxDao() *Dao {
	txDao := New(d.db.Begin())
	txDao.isTx = true
	return txDao
}

//仅供事务型dao调用
func (d *Dao) Rollback() error {
	if d.isTx {
		return d.db.Rollback().Error
	}
	return errors.New("not tx")
}

//仅供事务型dao调用
func (d *Dao) CommitTransaction() error {
	if d.isTx {
		return d.db.Commit().Error
	}
	return errors.New("not tx")
}

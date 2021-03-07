package models

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tpkeeper/pitaya/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	maxOpenConn = 10 //connect pool size
	maxIdleConn = 30
)

type Config struct {
	Host, Port, User, Pass, DBName string
	Mode                           uint8
}

//don't use soft delete
type BaseModel struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDB(cfg *Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DBName)

	logLevel := logger.Error
	if cfg.Mode == config.Debug {
		logLevel = logger.Info
	}

	db, err = gorm.Open(
		mysql.New(
			mysql.Config{
				DSN:                       dsn,   // data source name
				DefaultStringSize:         256,   // default size for string fields
				DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
				DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
				DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
				SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
			}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logLevel)})

	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDb.SetMaxIdleConns(maxIdleConn)
	sqlDb.SetMaxOpenConns(maxOpenConn)

	db.AutoMigrate(Pool{}, Pitaya{}, DepositWithDrawLog{}, PoolInfo{})
	logrus.Debug("new DB success")
	return
}

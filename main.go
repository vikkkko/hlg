package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tpkeeper/pitaya/config"
	"github.com/tpkeeper/pitaya/dao"
	"github.com/tpkeeper/pitaya/log"
	"github.com/tpkeeper/pitaya/models"
	"github.com/tpkeeper/pitaya/server"
	"os"
	"runtime"
	"runtime/debug"
)

func _main() error {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("loadConfig err: %s", err.Error())
		return err
	}
	log.InitLogFile(cfg.LogFilePath)
	logrus.Infof("config info: logFilePath: %s, db: %v, mode: %v ethApi:%s, MasterChiefContract: %v, PitayaTokenContract: %v",
		cfg.LogFilePath, cfg.Db, cfg.Mode, cfg.EthApi, cfg.MasterChiefContract, cfg.PitayaTokenContract)

	//init db
	db, err := models.NewDB(&models.Config{
		Host:   cfg.Db.Host,
		Port:   cfg.Db.Port,
		User:   cfg.Db.User,
		Pass:   cfg.Db.Pwd,
		DBName: cfg.Db.Name,
		Mode:   cfg.Mode})
	if err != nil {
		logrus.Errorf("db err: %s", err.Error())
		return err
	}
	logrus.Infof("db connect success")

	//interrupt signal
	ctx := server.ShutdownListener()

	defer func() {
		sqlDb, err := db.DB()
		if err != nil {
			logrus.Errorf("db.DB() err: %s", err.Error())
			return
		}
		logrus.Infof("shutting down the db ...")
		sqlDb.Close()
	}()

	//server
	server, err := server.NewServer(cfg, dao.New(db))
	if err != nil {
		logrus.Errorf("new server err: %s", err.Error())
		return err
	}
	server.Start()
	defer func() {
		logrus.Infof("shutting down server ...")
		server.Stop()
	}()

	<-ctx.Done()
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(40)
	err := _main()
	if err != nil {
		os.Exit(1)
	}
}

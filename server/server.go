package server

import (
	"github.com/sirupsen/logrus"
	"github.com/tpkeeper/pitaya/config"
	"github.com/tpkeeper/pitaya/dao"
	"net/http"
	"time"
)

type Server struct {
	listenAddr     string
	mode           uint8
	ethApi         string
	httpServer     *http.Server
	tickerInterval int64

	dao *dao.Dao
}

func NewServer(cfg *config.Config, dao *dao.Dao) (*Server, error) {
	s := &Server{
		listenAddr:     cfg.ListenAddr,
		mode:           cfg.Mode,
		tickerInterval: cfg.Ticker,
		ethApi:         cfg.EthApi,
		dao:            dao,
	}

	handler := s.InitRouters()

	s.httpServer = &http.Server{
		Addr:         s.listenAddr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return s, nil
}

func (svr *Server) Start() {
	go func() {
		logrus.Infof("Gin server start on %s", svr.listenAddr)
		err := svr.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed { //排除 调用close shutdown 调用导致的err
			logrus.Errorf("Gin server start err: %s", err.Error())
			shutdownRequestChannel <- struct{}{} //shutdown server
			return
		}
		logrus.Infof("Gin server done on %s", svr.listenAddr)
	}()

	go svr.Task()
}

func (svr *Server) Stop() {
	if svr.httpServer != nil {
		err := svr.httpServer.Close()
		if err != nil {
			logrus.Errorf("Problem shutdown Gin server :%s", err.Error())
		}
	}

}

package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var (
	Release = uint8(0)
	Debug   = uint8(1)
	cfg     = Config{}
)

type Config struct {
	Ticker              int64 // block duration
	EthApi              string
	ListenAddr          string
	Mode                uint8 //release debug test
	MasterChiefContract string
	PitayaTokenContract string
	Db                  Db
	LogFilePath         string
}

type Db struct {
	Host string
	Port string
	User string
	Pwd  string
	Name string
}

func Load() (*Config, error) {
	configFilePath := flag.String("C", "conf.toml", "Config file path")
	flag.Parse()

	if err := loadSysConfig(*configFilePath, &cfg); err != nil {
		return nil, err
	}
	if cfg.LogFilePath == "" {
		cfg.LogFilePath = "./log_data"
	}
	return &cfg, nil
}

func loadSysConfig(path string, config *Config) error {
	_, err := os.Open(path)
	if err != nil {
		return err
	}
	if _, err := toml.DecodeFile(path, config); err != nil {
		return err
	}
	fmt.Println("load sysConfig success")
	return nil
}

func GetMasterChiefContract() string {
	return cfg.MasterChiefContract
}

func GetPitayaTokenContract() string {
	return cfg.PitayaTokenContract
}

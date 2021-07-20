package config

import (
	"github.com/BurntSushi/toml"
	"github.com/feitianlove/golib/common/logger"
)

type Config struct {
	WebLog       *logger.LogConf
	WebAccessLog *logger.LogConf
	MysqlLog     *logger.LogConf
	CtrlLog      *logger.LogConf
	MysqlConf    *MysqlConf
	Worker       *Worker
	Master       *Master
}

type MysqlConf struct {
	User     string
	Passwd   string
	Host     string
	Port     int64
	Database string
}

type Worker struct {
	ListenPort int64
	Domain     string
	StaticDir  string
}

type Master struct {
	ListenPort int64
	Domain     string
	Token      string
}

func InitConfig() (*Config, error) {
	var config = defaultConfig()
	_, err := toml.DecodeFile("./etc/worker.conf", config)
	if err != nil {
		return nil, err
	}
	return config, err
}

func defaultConfig() *Config {
	return &Config{
		WebLog: &logger.LogConf{
			LogLevel:      "info",
			LogPath:       "/Users/fenghui/goCode/worker/log/web.log",
			LogReserveDay: 1,
			ReportCaller:  true,
		},
		WebAccessLog: &logger.LogConf{
			LogLevel:      "info",
			LogPath:       "/Users/fenghui/goCode/worker/log/web_access.log",
			LogReserveDay: 1,
			ReportCaller:  true,
		},
		MysqlLog: &logger.LogConf{
			LogLevel:      "info",
			LogPath:       "/Users/fenghui/goCode/worker/log/mysql.log",
			LogReserveDay: 1,
			ReportCaller:  true,
		},
		CtrlLog: &logger.LogConf{
			LogLevel:      "info",
			LogPath:       "/Users/fenghui/goCode/worker/log/ctrl.log",
			LogReserveDay: 1,
			ReportCaller:  true,
		},
		MysqlConf: &MysqlConf{
			User:     "",
			Passwd:   "",
			Host:     "",
			Port:     0,
			Database: "",
		},
	}
}

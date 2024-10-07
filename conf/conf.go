package conf

import (
	"github.com/lpphub/golib/env"
)

var (
	RConf resourceConf
)

type resourceConf struct {
	Mysql mysqlConf
}

type mysqlConf struct {
	Addr     string `yaml:"addr"`
	DataBase string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func InitConf() {
	env.LoadConf("conf/conf.yaml", &RConf)
}

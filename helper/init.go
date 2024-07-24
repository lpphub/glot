package helper

import (
	"github.com/lpphub/golib/zlog"
	"glot/conf"
)

func preInit() {
	// 加载配置文件
	conf.InitConf()
	// 日志配置
	zlog.InitLog(zlog.WithBufSwitch(false))
}

func InitResource() {
	preInit()

	initDb()
	initCache()
}

func Clear() {
	zlog.Close()
}

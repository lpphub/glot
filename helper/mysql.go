package helper

import (
	"fmt"
	"github.com/lpphub/golib/zlog"
	"glot/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDb() {
	var (
		err    error
		dbConf = conf.RConf.Mysql

		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=True&loc=Asia%%2FShanghai",
			dbConf.User, dbConf.Password, dbConf.Addr, dbConf.DataBase)
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: zlog.NewGormLogger(dbConf.DataBase, dbConf.Addr), // 重写日志
	})
	if err != nil {
		panic("init db error: " + err.Error())
	}
}

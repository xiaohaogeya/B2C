package models

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/prometheus/common/log"
)

var (
	DB  *gorm.DB
	err error
)

func init() {
	mysqlAdmin, _ := beego.AppConfig.String("MysqlUser")
	mysqlPwd, _ := beego.AppConfig.String("MysqlPwd")
	mysqlDB, _ := beego.AppConfig.String("MysqlDB")
	DB, err =
		gorm.Open("mysql", mysqlAdmin+":"+mysqlPwd+"@/"+mysqlDB+"?charset=utf8"+
			"&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
		log.Error("连接MySql数据库失败")
	} else {
		log.Info("连接MySql数据库成功")
	}
}

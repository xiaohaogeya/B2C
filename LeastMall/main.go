package main

import (
	"LeastMall/common"
	"LeastMall/models"
	_ "LeastMall/routers"
	"encoding/gob"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/beego/beego/v2/server/web/session/redis"
)

func main() {
	// 添加方法到map，用于前端HTML代码调用
	_ = beego.AddFuncMap("timestampToDate", common.TimestampToDate)
	_ = beego.AddFuncMap("formatImage", common.FormatImage)
	_ = beego.AddFuncMap("mul", common.Mul)
	_ = beego.AddFuncMap("formatAttribute", common.FormatAttribute)

	models.DB.LogMode(true)
	_ = beego.AddFuncMap("setting", models.GetSettingByColumn)

	// 后台配置允许跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"127.0.0.1"}, // 允许访问所有源
		AllowMethods: []string{ // 可选参数
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{ // 允许的Header种类
			"Origin",
			"Authorization",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type"},
		AllowCredentials: true, //是否允许cookie
	}))

	//注册模型
	gob.Register(models.Administrator{})

	//配置redis用于存储session
	beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	redisConn, _ := beego.AppConfig.String("redisConn")
	beego.BConfig.WebConfig.Session.SessionProviderConfig = redisConn
	beego.Run()
}

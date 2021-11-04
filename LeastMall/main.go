package main

import (
	"LeastMall/common"
	"LeastMall/models"
	_ "LeastMall/routers"
	"github.com/beego/beego/v2/adapter/plugins/cors"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// 添加方法到map，用于前端HTML代码调用
	_ = beego.AddFuncMap("timestampToDate", common.TimestampToDate)
	_ = beego.AddFuncMap("formatImage", common.FormatImage)
	_ = beego.AddFuncMap("mul", common.Mul)
	_ = beego.AddFuncMap("formatAttribute", common.FormatAttribute)

	//models.DB.LogMode(true)
	_ = beego.AddFuncMap("setting", models.GetSettingByColumn)

	// 后台配置允许跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"127.0.0.1"},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
	}))

	beego.Run()
}

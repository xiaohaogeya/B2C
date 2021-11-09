package routers

import (
	"LeastMall/common"
	"LeastMall/controllers/backend"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	adminPath, _ := beego.AppConfig.String("AdminPath")
	ns := beego.NewNamespace("/"+adminPath,
		beego.NSBefore(common.BackendAuth),

		// 管理员管理
		beego.NSRouter("/administrator", &backend.AdministratorController{}),
		beego.NSRouter("/administrator/add", &backend.AdministratorController{}, "get:Add"),
		beego.NSRouter("/administrator/edit", &backend.AdministratorController{}, "get:Edit"),
		beego.NSRouter("/administrator/goadd", &backend.AdministratorController{}, "post:GoAdd"),
		beego.NSRouter("/administrator/goedit", &backend.AdministratorController{}, "post:GoEdit"),
		beego.NSRouter("/administrator/delete", &backend.AdministratorController{}, "get:Delete"),

		//权限管理
		beego.NSRouter("/auth", &backend.AuthController{}),
		beego.NSRouter("/auth/add", &backend.AuthController{}, "get:Add"),
		beego.NSRouter("/auth/edit", &backend.AuthController{}, "get:Edit"),
		beego.NSRouter("/auth/goadd", &backend.AuthController{}, "post:GoAdd"),
		beego.NSRouter("/auth/goedit", &backend.AuthController{}, "post:GoEdit"),
		beego.NSRouter("/auth/delete", &backend.AuthController{}, "get:Delete"),
	)
	beego.AddNamespace(ns)
}

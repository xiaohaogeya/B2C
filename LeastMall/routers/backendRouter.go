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

		// 后台管理
		beego.NSRouter("/", &backend.MainController{}),
		beego.NSRouter("/welcome", &backend.MainController{}, "get:Welcome"),
		beego.NSRouter("/main/changestatus", &backend.MainController{}, "get:ChangeStatus"),
		beego.NSRouter("/main/editnum", &backend.MainController{}, "get:EditNum"),

		// 登录登出
		beego.NSRouter("/login", &backend.LoginController{}),
		beego.NSRouter("/login/gologin", &backend.LoginController{}, "post:GoLogin"),
		beego.NSRouter("/login/loginout", &backend.LoginController{}, "get:LoginOut"),

		// 管理员管理
		beego.NSRouter("/administrator", &backend.AdministratorController{}),
		beego.NSRouter("/administrator/add", &backend.AdministratorController{}, "get:Add"),
		beego.NSRouter("/administrator/edit", &backend.AdministratorController{}, "get:Edit"),
		beego.NSRouter("/administrator/goadd", &backend.AdministratorController{}, "post:GoAdd"),
		beego.NSRouter("/administrator/goedit", &backend.AdministratorController{}, "post:GoEdit"),
		beego.NSRouter("/administrator/delete", &backend.AdministratorController{}, "get:Delete"),

		// 权限管理
		beego.NSRouter("/auth", &backend.AuthController{}),
		beego.NSRouter("/auth/add", &backend.AuthController{}, "get:Add"),
		beego.NSRouter("/auth/edit", &backend.AuthController{}, "get:Edit"),
		beego.NSRouter("/auth/goadd", &backend.AuthController{}, "post:GoAdd"),
		beego.NSRouter("/auth/goedit", &backend.AuthController{}, "post:GoEdit"),
		beego.NSRouter("/auth/delete", &backend.AuthController{}, "get:Delete"),

		// 轮播图管理
		beego.NSRouter("/banner", &backend.BannerController{}),
		beego.NSRouter("/banner/add", &backend.BannerController{}, "get:Add"),
		beego.NSRouter("/banner/edit", &backend.BannerController{}, "get:Edit"),
		beego.NSRouter("/banner/goadd", &backend.BannerController{}, "post:GoAdd"),
		beego.NSRouter("/banner/goedit", &backend.BannerController{}, "post:GoEdit"),
		beego.NSRouter("/banner/delete", &backend.BannerController{}, "get:Delete"),

		// 导航管理
		beego.NSRouter("/menu", &backend.MenuController{}),
		beego.NSRouter("/menu/add", &backend.MenuController{}, "get:Add"),
		beego.NSRouter("/menu/edit", &backend.MenuController{}, "get:Edit"),
		beego.NSRouter("/menu/goadd", &backend.MenuController{}, "post:GoAdd"),
		beego.NSRouter("/menu/goedit", &backend.MenuController{}, "post:GoEdit"),
		beego.NSRouter("/menu/delete", &backend.MenuController{}, "get:Delete"),
	)
	beego.AddNamespace(ns)
}

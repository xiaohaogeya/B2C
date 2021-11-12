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

		// 订单管理
		beego.NSRouter("/order", &backend.OrderController{}),
		beego.NSRouter("/order/detail", &backend.OrderController{}, "get:Detail"),
		beego.NSRouter("/order/edit", &backend.OrderController{}, "get:Edit"),
		beego.NSRouter("/order/goEdit", &backend.OrderController{}, "post:GoEdit"),
		beego.NSRouter("/order/delete", &backend.OrderController{}, "get:Delete"),

		// 商品分类管理
		beego.NSRouter("/productCate", &backend.ProductCateController{}),
		beego.NSRouter("/productCate/add", &backend.ProductCateController{}, "get:Add"),
		beego.NSRouter("/productCate/edit", &backend.ProductCateController{}, "get:Edit"),
		beego.NSRouter("/productCate/goadd", &backend.ProductCateController{}, "post:GoAdd"),
		beego.NSRouter("/productCate/goedit", &backend.ProductCateController{}, "post:GoEdit"),
		beego.NSRouter("/productCate/delete", &backend.ProductCateController{}, "get:Delete"),

		// 商品管理
		beego.NSRouter("/product", &backend.ProductController{}),
		beego.NSRouter("/product/add", &backend.ProductController{}, "get:Add"),
		beego.NSRouter("/product/goadd", &backend.ProductController{}, "post:GoAdd"),
		beego.NSRouter("/product/edit", &backend.ProductController{}, "get:Edit"),
		beego.NSRouter("/product/goedit", &backend.ProductController{}, "post:GoEdit"),
		beego.NSRouter("/product/delete", &backend.ProductController{}, "get:Delete"),
		beego.NSRouter("/product/goUpload", &backend.ProductController{}, "post:GoUpload"),
		beego.NSRouter("/product/getProductTypeAttribute", &backend.ProductController{}, "get:GetProductTypeAttribute"),
		beego.NSRouter("/product/changeProductImageColor", &backend.ProductController{}, "get:ChangeProductImageColor"),
		beego.NSRouter("/product/removeProductImage", &backend.ProductController{}, "get:RemoveProductImage"),
	)
	beego.AddNamespace(ns)
}

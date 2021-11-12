package routers

import (
	"LeastMall/common"
	"LeastMall/controllers/frontend"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// 首页
	beego.Router("/", &frontend.IndexController{})

	// 地址 配置中间件判断权限
	beego.InsertFilter("/address/*", beego.BeforeRouter, common.FrontendAuth)
	beego.Router("/address/add", &frontend.AddressController{}, "post:AddAddress")
	beego.Router("/address/getOneAddress", &frontend.AddressController{}, "get:GetOneAddress")
	beego.Router("/address/updateAddress", &frontend.AddressController{}, "post:UpdateAddress")
	beego.Router("/address/changeDefaultAddress", &frontend.AddressController{}, "put:ChangeDefaultAddress")

	// auth
	beego.Router("/auth/sendCode", &frontend.AuthController{}, "get:SendCode")
	beego.Router("/auth/doRegister", &frontend.AuthController{}, "post:GoRegister")
	beego.Router("/auth/validateSmsCode", &frontend.AuthController{}, "get:ValidateSmsCode")
	beego.Router("/auth/login", &frontend.AuthController{}, "get:Login")
	beego.Router("/auth/registerStep1", &frontend.AuthController{}, "get:RegisterStep1")
	beego.Router("/auth/registerStep2", &frontend.AuthController{}, "get:RegisterStep2")
	beego.Router("/auth/registerStep3", &frontend.AuthController{}, "get:RegisterStep3")
	beego.Router("/auth/login", &frontend.AuthController{}, "get:Login")
	beego.Router("/auth/goLogin", &frontend.AuthController{}, "post:GoLogin")
	beego.Router("/auth/loginOut", &frontend.AuthController{}, "get:LoginOut")

	// 购物车
	beego.Router("/cart", &frontend.CartController{})
	beego.Router("/cart/addCart", &frontend.CartController{}, "get:AddCart")
	beego.Router("/cart/incCart", &frontend.CartController{}, "get:IncCart")
	beego.Router("/cart/decCart", &frontend.CartController{}, "get:DecCart")
	beego.Router("/cart/delCart", &frontend.CartController{}, "get:DelCart")
	beego.Router("/cart/changeOneCart", &frontend.CartController{}, "get:ChangeOneCart")
	beego.Router("/cart/changeAllCart", &frontend.CartController{}, "get:ChangeAllCart")

	// 结算 配置中间件判断权限
	beego.InsertFilter("/buy/*", beego.BeforeRouter, common.FrontendAuth)
	beego.Router("/buy/checkout", &frontend.CheckoutController{}, "get:Checkout")
	beego.Router("/buy/doOrder", &frontend.CheckoutController{}, "post:GoOrder")
	beego.Router("/buy/confirm", &frontend.CheckoutController{}, "get:Confirm")
	beego.Router("/buy/orderPayStatus", &frontend.CheckoutController{}, "get:OrderPayStatus")

	// 支付宝支付
	beego.Router("/alipay", &frontend.PayController{}, "get:Alipay")
	beego.Router("/alipayNotify", &frontend.PayController{}, "post:AlipayNotify")
	beego.Router("/alipayReturn", &frontend.PayController{}, "get:AlipayReturn")

	// 产品
	beego.Router("/category_:id([0-9]+).html", &frontend.ProductController{}, "get:CategoryList")
	beego.Router("/item_:id([0-9]+).html", &frontend.ProductController{}, "get:ProductItem")
	beego.Router("/product/getImgList", &frontend.ProductController{}, "get:GetImgList")
	beego.Router("/product/collect", &frontend.ProductController{}, "get:Collect")

}

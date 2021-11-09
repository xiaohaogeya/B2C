package routers

import (
	"LeastMall/common"
	"LeastMall/controllers/frontend"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// 配置中间件判断权限
	beego.InsertFilter("/address/*", beego.BeforeRouter, common.FrontendAuth)
	beego.Router("/address/add", &frontend.AddressController{}, "post:AddAddress")
	beego.Router("/address/getOneAddress", &frontend.AddressController{}, "get:GetOneAddress")
	beego.Router("/address/updateAddress", &frontend.AddressController{}, "post:UpdateAddress")
	beego.Router("/address/changeDefaultAddress", &frontend.AddressController{}, "put:ChangeDefaultAddress")
}

package routers

import (
	"LeastMall/controllers/frontend"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// //配置中间件判断权限
	beego.Router("/address/getOneAddressList/", &frontend.AddressController{}, "get:GetOneAddressList")
}
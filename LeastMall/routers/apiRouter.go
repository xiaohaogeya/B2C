package routers

import (
	"LeastMall/controllers/api"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/", &api.V1Controller{}),
		beego.NSRouter("/menu", &api.V1Controller{}, "get:Menu"),
	)
	beego.AddNamespace(ns)
}

package api

import beego "github.com/beego/beego/v2/server/web"

type V2Controller struct {
	beego.Controller
}

func (c *V2Controller) Get() {
	c.Ctx.WriteString("api v2")
}

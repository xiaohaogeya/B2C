package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
	//s, err := beego.AppConfig.Int("HttpPort")
	//if err != nil {
	//	return
	//}
	//fmt.Println(s)
	//fmt.Println(time.Now().Unix())
	//str := common.TimestampToDate(time.Now().Unix())
	//fmt.Println(str)
}

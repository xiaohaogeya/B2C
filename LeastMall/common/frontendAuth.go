/*
	前台权限认证文件
*/
package common

import (
	"LeastMall/models"
	"github.com/beego/beego/v2/server/web/context"
)

func FrontendAuth(ctx *context.Context)  {
	// 前台用户有没有登陆
	user := models.User{}
	models.Cookie.Get(ctx, "userinfo", &user)
	if user.Id == 0 {
		ctx.Redirect(302, "/auth/login")
	}
}
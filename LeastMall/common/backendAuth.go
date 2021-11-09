/*
	后台权限认证文件
*/
package common

import (
	"LeastMall/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"net/url"
	"strings"
)

// BackendAuth 后台权限判断
func BackendAuth(ctx *context.Context) {
	pathname := ctx.Request.URL.String()
	adminPath, _ := beego.AppConfig.String("AdminPath")
	userinfo, ok := ctx.Input.Session("userinfo").(models.Administrator)
	if !(ok && userinfo.Username != "") {
		if pathname != "/"+adminPath+"/login" && pathname != "/"+adminPath+"/login/gologin" && pathname != "/"+adminPath+"/login/verificode" {
			//ctx.Redirect(302, "/"+adminPath+"/login")
		}
	} else {
		pathname = strings.Replace(pathname, adminPath, "", 1)
		urlPath, _ := url.Parse(pathname)
		if userinfo.IsSuper == 0 && !excludeAuthPath(urlPath.Path) {
			roleId := userinfo.RoleId
			var roleAuth []models.RoleAuth
			models.DB.Where("role_id=?", roleId).Find(&roleAuth)

			roleAuthMap := make(map[int]int)
			for _, v := range roleAuth {
				roleAuthMap[v.AuthId] = v.AuthId
			}

			auth := models.Auth{}
			models.DB.Where("url=?", urlPath.Path).Find(&auth)
			if _, ok := roleAuthMap[auth.Id]; !ok {
				ctx.WriteString("没有权限")
				return
			}
		}
	}
}

// excludeAuthPath 检验路径权限
func excludeAuthPath(urlPath string) bool {
	authPath, _ := beego.AppConfig.String("ExcludeAuthPath")
	excludeAuthPathSlice := strings.Split(authPath, ",")
	for _, v := range excludeAuthPathSlice {
		if v == urlPath {
			return true
		}
	}
	return false
}

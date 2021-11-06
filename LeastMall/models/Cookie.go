package models

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

//定义结构体  缓存结构体 私有
type cookie struct{}

// Set 写入数据的方法
func (c *cookie) Set(ctx *context.Context, key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	secureCookie, _ := beego.AppConfig.String("SecureCookie")
	domain, _ := beego.AppConfig.String("DoMain")
	ctx.SetSecureCookie(secureCookie, key, string(bytes), 3600*24*30, "/", domain, nil, true) // 有效期30天
}

// Remove 删除数据的方法
func (c *cookie) Remove(ctx *context.Context, key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	secureCookie, _ := beego.AppConfig.String("SecureCookie")
	domain, _ := beego.AppConfig.String("DoMain")
	ctx.SetSecureCookie(secureCookie, key, string(bytes), -1, "/", domain, nil, true)
}

// Get 获取数据的方法
func (c *cookie) Get(ctx *context.Context, key string, obj interface{}) bool {
	secureCookie, _ := beego.AppConfig.String("SecureCookie")
	data, ok := ctx.GetSecureCookie(secureCookie, key)
	if !ok {
		return false
	}
	_ = json.Unmarshal([]byte(data), obj)
	return true
}

// Cookie 实例化结构体
var Cookie = &cookie{}

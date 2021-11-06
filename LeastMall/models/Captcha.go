package models

import (
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/server/web/captcha"
)

var Cpt *captcha.Captcha

func init() {
	store := cache.NewMemoryCache()
	Cpt = captcha.NewWithFilter("/captcha/", store)
	Cpt.ChallengeNums = 4
	Cpt.StdWidth = 100
	Cpt.StdHeight = 40
}

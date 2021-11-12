package frontend

import (
	"LeastMall/models"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	//调用初始化功能
	c.BaseInit()

	//获取轮播图 注意获取的时候要写地址
	var banner []models.Banner
	if hasBanner := models.CacheDb.Get("banner", &banner); hasBanner == true {
		c.Data["bannerList"] = banner
	} else {
		models.DB.Where("status=1 AND banner_type=1").Order("sort desc").Find(&banner)
		c.Data["bannerList"] = banner
		models.CacheDb.Set("banner", banner)
	}

	//获取手机商品列表
	var redisPhone []models.Product
	if hasPhone := models.CacheDb.Get("phone", &redisPhone); hasPhone == true {
		c.Data["phoneList"] = redisPhone
	} else {
		phone := models.GetProductByCategory(1, "hot", 8)
		c.Data["phoneList"] = phone
		models.CacheDb.Set("phone", phone)
	}
	//获取电视商品列表
	var redisTv []models.Product
	if hasTv := models.CacheDb.Get("tv", &redisTv); hasTv == true {
		c.Data["tvList"] = redisTv
	} else {
		tv := models.GetProductByCategory(4, "best", 8)
		c.Data["tvList"] = tv
		models.CacheDb.Set("tv", tv)
	}

	c.TplName = "frontend/index/index.html"
}

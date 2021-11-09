package backend

import "LeastMall/models"

type BannerController struct {
	BaseController
}

func (c *BannerController) Get()  {
	var bannerList []models.Banner
	models.DB.Find(&bannerList)

	c.Data["bannerList"] = bannerList
	c.TplName = "backend/banner/index.html"
}

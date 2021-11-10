package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"github.com/prometheus/common/log"
	"os"
	"strconv"
)

type BannerController struct {
	BaseController
}

func (c *BannerController) Get() {
	var bannerList []models.Banner
	models.DB.Find(&bannerList)

	c.Data["bannerList"] = bannerList
	c.TplName = "backend/banner/index.html"
}

func (c *BannerController) Add() {
	c.TplName = "backend/banner/add.html"
}

func (c *BannerController) GoAdd() {
	bannerType, err1 := c.GetInt("banner_type")
	title := c.GetString("title")
	link := c.GetString("link")
	sort, err2 := c.GetInt("sort")
	status, err3 := c.GetInt("status")
	if err1 != nil || err3 != nil {
		c.Error("非法请求", "/banner")
		return
	}
	if err2 != nil {
		c.Error("排序表单里面输入的数据不合法", "/banner/add")
		return
	}
	bannerImgSrc, err4 := c.UploadImg("banner_img")
	if err4 == nil {
		banner := models.Banner{
			Title:      title,
			BannerType: bannerType,
			BannerImg:  bannerImgSrc,
			Link:       link,
			Sort:       sort,
			Status:     status,
			AddTime:    int(common.GetUnix()),
		}
		models.DB.Create(&banner)
		c.Success("增加轮播图成功", "/banner")
	} else {
		c.Error("增加轮播图失败", "/banner/add")
		return
	}
}

func (c *BannerController) Edit() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("非法请求", "/banner")
		return
	}
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	c.Data["banner"] = banner
	c.TplName = "backend/banner/edit.html"
}

func (c *BannerController) GoEdit() {
	id, err := c.GetInt("id")
	bannerType, err1 := c.GetInt("banner_type")
	title := c.GetString("title")
	link := c.GetString("link")
	sort, err2 := c.GetInt("sort")
	status, err3 := c.GetInt("status")
	if err != nil || err1 != nil || err3 != nil {
		c.Error("非法请求", "/banner")
		return
	}
	if err2 != nil {
		c.Error("排序表单里面输入的数据不合法", "/banner/edit?id="+strconv.Itoa(id))
		return
	}
	bannerImgSrc, _ := c.UploadImg("banner_img")
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	banner.Title = title
	banner.BannerType = bannerType
	banner.Link = link
	banner.Sort = sort
	banner.Status = status
	if bannerImgSrc != "" {
		banner.BannerImg = bannerImgSrc
	}
	err5 := models.DB.Save(&banner).Error
	if err5 != nil {
		c.Error("修改轮播图失败", "/banner/edit?id="+strconv.Itoa(id))
		return
	}
	c.Success("修改轮播图成功", "/banner")
}

func (c *BannerController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/banner")
		return
	}
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	address := "/Users/xiaohaoge/go/src/B2C/LeastMall/" + banner.BannerImg // 物理机地址
	test := os.Remove(address)
	if test != nil {
		log.Error(test)
		// 可以向前端提示物理机图片删除失败。。。
		//c.Error("删除物理机上图片错误", "/banner")
		//return
	}
	models.DB.Delete(&banner)
	c.Success("删除轮播图成功", "/banner")
}

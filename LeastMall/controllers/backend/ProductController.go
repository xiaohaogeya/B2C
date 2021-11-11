package backend

import (
	"LeastMall/models"
	"math"
)

type ProductController struct {
	BaseController
}

func (c *ProductController) Get() {
	page, err := c.GetInt("page")
	if err != nil || page == 0 {
		page = 1
	}
	keyword := c.GetString("keyword")
	pageSize := 5

	where := "1=1" // 1=1 作用相当于select * from table；这里为拼接字符串方便
	if len(keyword) > 0 {
		where += " AND title like \"%" + keyword + "%\""
	}
	var productList []models.Product
	models.DB.Where(where).Offset((page - 1) * pageSize).Limit(pageSize).Find(&productList)
	var count int
	models.DB.Where(where).Table("product").Count(&count)
	c.Data["productList"] = productList
	c.Data["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	c.Data["page"] = page
	c.TplName = "backend/product/index.html"
}

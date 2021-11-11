package backend

import (
	"LeastMall/models"
	"math"
	"strconv"
)

type OrderController struct {
	BaseController
}

func (c *OrderController) Get() {
	page, err := c.GetInt("page")
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize := 5
	keyword := c.GetString("keyword")

	var orderList []models.Order
	var count int

	if keyword == "" {
		models.DB.Table("order").Count(&count)
		models.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&orderList)
	} else {
		models.DB.Where("phone=?", keyword).Table("order").Count(&count)
		models.DB.Where("phone=?", keyword).Offset((page - 1) * pageSize).Limit(pageSize).Find(&orderList)
	}

	c.Data["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	c.Data["page"] = page
	c.Data["order"] = orderList
	c.TplName = "backend/order/order.html"
}

func (c *OrderController) Detail() {
	c.Ctx.WriteString("详情页面")
}
func (c *OrderController) Edit() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/order")
		return
	}
	order := models.Order{}
	models.DB.Where("id=?", id).Find(&order)
	c.Data["order"] = order
	c.TplName = "backend/order/edit.html"
}
func (c *OrderController) GoEdit() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/order")
		return
	}
	orderId := c.GetString("order_id")
	allPrice := c.GetString("all_price")
	name := c.GetString("name")
	phone := c.GetString("phone")
	address := c.GetString("address")
	zipcode := c.GetString("zipcode")
	payStatus, _ := c.GetInt("pay_status")
	payType, _ := c.GetInt("pay_type")
	orderStatus, _ := c.GetInt("order_status")
	order := models.Order{}
	models.DB.Where("id=?", id).Find(&order)
	order.OrderId = orderId
	order.AllPrice, _ = strconv.ParseFloat(allPrice, 64)
	order.Name = name
	order.Phone = phone
	order.Address = address
	order.Zipcode = zipcode
	order.PayStatus = payStatus
	order.PayType = payType
	order.OrderStatus = orderStatus
	models.DB.Save(&order)
	c.Success("订单修改成功", "/order")
}
func (c *OrderController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/order")
		return
	}
	order := models.Order{}
	models.DB.Where("id=?", id).Delete(&order)
	c.Success("删除订单记录成功", "/order")
}

package frontend

import (
	"LeastMall/models"
)

type AddressController struct {
	BaseController
}

// AddAddress 添加地址
func (c *AddressController) AddAddress() {
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	name := c.GetString("name")
	phone := c.GetString("phone")
	address := c.GetString("address")
	zipcode := c.GetString("zipcode")

	var addressCount int
	models.DB.Where("uid=?", user.Id).Table("address").Count(&addressCount)

	if addressCount > 10 {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"message": "增加收货地址失败，收货地址数量超过限制",
		}
		_ = c.ServeJSON()
		return
	}
	models.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	addressResult := models.Address{
		Uid:            user.Id,
		Name:           name,
		Address:        address,
		Phone:          phone,
		Zipcode:        zipcode,
		DefaultAddress: 1,
	}
	models.DB.Create(&addressResult)
	var allAddressResult []models.Address
	models.DB.Where("uid=?", user.Id).Find(&allAddressResult)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"result":  allAddressResult,
	}
	_ = c.ServeJSON()
}

// GetOneAddress 获取单个地址
func (c *AddressController) GetOneAddress() {
	addressId, err := c.GetInt("address_id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "传入参数错误",
		}
		_ = c.ServeJSON()
		return
	}

	address := models.Address{}
	models.DB.Where("id=?", addressId).Find(&address)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"result":  address,
	}
	_ = c.ServeJSON()
}

// UpdateAddress 更新地址
func (c *AddressController) UpdateAddress() {
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	addressId, err := c.GetInt("address_id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "传入参数错误",
		}
		_ = c.ServeJSON()
		return
	}

	name := c.GetString("name")
	phone := c.GetString("phone")
	address := c.GetString("address")
	zipcode := c.GetString("zipcode")

	models.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	addressModel := models.Address{}
	models.DB.Where("id=?", addressId).Find(&addressModel)
	addressModel.Name = name
	addressModel.Phone = phone
	addressModel.Address = address
	addressModel.Zipcode = zipcode
	addressModel.DefaultAddress = 1
	models.DB.Save(&addressModel)

	// 查询当前用户的所有收货地址并返回
	var allAddressResult []models.Address
	models.DB.Where("uid=?", user.Id).Order("default_address desc").Find(&allAddressResult)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"result":  allAddressResult,
	}
	_ = c.ServeJSON()
}

// ChangeDefaultAddress 修改默认地址
func (c *AddressController) ChangeDefaultAddress()  {
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	addressId, err := c.GetInt("address_id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "传入参数错误",
		}
		_ = c.ServeJSON()
		return
	}

	models.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	models.DB.Table("address").Where("id=?", addressId).Updates(map[string]interface{}{"default_address": 1})

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"result":  "更新默认收获地址成功",
	}
	_ = c.ServeJSON()
}
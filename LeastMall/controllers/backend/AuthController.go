package backend

import (
	"LeastMall/models"
	"strconv"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) Get() {
	var authList []models.Auth
	models.DB.Preload("AuthItem").Where("module_id=0").Find(&authList)
	c.Data["authList"] = authList
	c.TplName = "backend/auth/index.html"
}

func (c *AuthController) Add() {
	var authList []models.Auth
	models.DB.Where("module_id=0").Find(&authList)
	c.Data["authList"] = authList
	c.TplName = "backend/auth/add.html"
}

func (c *AuthController) GoAdd() {
	moduleName := c.GetString("module_name")
	iType, err1 := c.GetInt("type")
	actionName := c.GetString("action_name")
	url := c.GetString("url")
	moduleId, err2 := c.GetInt("module_id")
	sort, err3 := c.GetInt("sort")
	description := c.GetString("description")
	status, err4 := c.GetInt("status")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Error("传入参数错误", "/auth/add")
		return
	}

	auth := models.Auth{
		ModuleName:  moduleName,
		Type:        iType,
		ActionName:  actionName,
		Url:         url,
		ModuleId:    moduleId,
		Sort:        sort,
		Description: description,
		Status:      status,
	}
	err := models.DB.Create(&auth).Error
	if err != nil {
		c.Error("增加数据失败", "/auth/add")
		return
	}
	c.Success("增加数据成功", "/auth")
}

func (c *AuthController) Edit() {
	id, err1 := c.GetInt("id")
	if err1 != nil {
		c.Error("传入参数错误", "/auth")
		return
	}
	auth := models.Auth{Id: id}
	models.DB.Find(&auth)
	c.Data["auth"] = auth
	var authList []models.Auth
	models.DB.Where("module_id=0").Find(&authList)
	c.Data["authList"] = authList
	c.TplName = "backend/auth/edit.html"
}

func (c *AuthController) GoEdit() {
	moduleName := c.GetString("module_name")
	iType, err1 := c.GetInt("type")
	actionName := c.GetString("action_name")
	url := c.GetString("url")
	moduleId, err2 := c.GetInt("module_id")
	sort, err3 := c.GetInt("sort")
	description := c.GetString("description")
	status, err4 := c.GetInt("status")
	id, err5 := c.GetInt("id")
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		c.Error("传入参数错误", "/auth")
		return
	}
	auth := models.Auth{Id: id}
	models.DB.Find(&auth)
	auth.ModuleName = moduleName
	auth.Type = iType
	auth.ActionName = actionName
	auth.Url = url
	auth.ModuleId = moduleId
	auth.Sort = sort
	auth.Description = description
	auth.Status = status
	err6 := models.DB.Save(&auth).Error
	if err6 != nil {
		c.Error("修改权限失败", "/auth/edit?id="+strconv.Itoa(id))
		return
	}
	c.Success("修改权限成功", "/auth")
}

func (c *AuthController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/role")
		return
	}
	auth := models.Auth{Id: id}
	models.DB.Find(&auth)
	if auth.ModuleId == 0 {
		var auth2 []models.Auth
		models.DB.Where("module_id=?", auth.Id).Find(&auth2)
		if len(auth2) > 0 {
			c.Error("请删除当前顶级模块下面的菜单或操作！", "/auth")
			return
		}
	}
	models.DB.Delete(&auth)
	c.Success("删除成功", "/auth")
}

package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"strconv"
	"strings"
)

type RoleController struct {
	BaseController
}

func (c *RoleController) Get() {
	var roleList []models.Role
	models.DB.Find(&roleList)
	c.Data["roleList"] = roleList
	c.TplName = "backend/role/index.html"
}

func (c *RoleController) Add() {
	c.TplName = "backend/role/add.html"
}

func (c *RoleController) GoAdd() {
	title := strings.Trim(c.GetString("title"), "")
	description := strings.Trim(c.GetString("description"), "")
	if title == "" {
		c.Error("标题不能为空", "/role/add")
		return
	}
	var roleList []models.Role
	models.DB.Where("title=?", title).Find(&roleList)
	if len(roleList) != 0 {
		c.Error("该部门已存在！", "/role/add")
		return
	}
	role := models.Role{}
	role.Title = title
	role.Description = description
	role.Status = 1
	role.AddTime = int(common.GetUnix())
	err := models.DB.Create(&role).Error
	if err != nil {
		c.Error("增加部门失败", "/role/add")
	} else {
		c.Success("增加部门成功", "/role")
	}
}

func (c *RoleController) Edit() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/role")
		return
	}
	role := models.Role{Id: id}
	models.DB.Find(&role)
	c.Data["role"] = role
	c.TplName = "backend/role/edit.html"
}

func (c *RoleController) GoEdit() {
	title := strings.Trim(c.GetString("title"), "")
	description := strings.Trim(c.GetString("description"), "")
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/role")
		return
	}
	role := models.Role{Id: id}
	models.DB.Find(&role)
	role.Title = title
	role.Description = description
	err2 := models.DB.Save(&role).Error
	if err2 != nil {
		c.Error("修改部门失败", "/role/edit?id="+strconv.Itoa(id))
	} else {
		c.Success("修改部门成功", "/role")
	}
}

func (c *RoleController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/role")
		return
	}
	role := models.Role{Id: id}
	var administrator []models.Administrator
	roleAuth := models.RoleAuth{}
	models.DB.Where("role_id=?", id).Delete(&roleAuth)
	models.DB.Preload("Role").Where("role_id=?", id).Find(&administrator)
	if len(administrator) > 0 {
		c.Error("该部门还有未处理的员工，无法删除该部门", "/role")
		return
	}
	models.DB.Delete(&role)
	c.Success("删除部门成功", "/role")
}

func (c *RoleController) Auth() {
	roleId, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/role")
		return
	}
	var auth []models.Auth
	models.DB.Preload("AuthItem").Where("module_id=0").Find(&auth)
	//获取当前部门拥有的权限，并把权限ID放在一个MAP对象里面
	var roleAuth []models.RoleAuth
	models.DB.Where("role_id=?", roleId).Find(&roleAuth)
	roleAuthMap := make(map[int]int)
	for _, v := range roleAuth {
		roleAuthMap[v.AuthId] = v.AuthId
	}
	for i := 0; i < len(auth); i++ {
		if _, ok := roleAuthMap[auth[i].Id]; ok {
			auth[i].Checked = true
		}
		for j := 0; j < len(auth[i].AuthItem); j++ {
			if _, ok := roleAuthMap[auth[i].AuthItem[j].Id]; ok {
				auth[i].AuthItem[j].Checked = true
			}
		}
	}

	c.Data["authList"] = auth
	c.Data["roleId"] = roleId
	c.TplName = "backend/role/auth.html"
}

func (c *RoleController) GoAuth() {
	roleId, err := c.GetInt("id")
	if err != nil {
		c.Error("传入参数错误", "/role")
		return
	}
	authNode := c.GetStrings("auth_node")
	roleAuth := models.RoleAuth{}
	models.DB.Where("role_id=?", roleId).Delete(&roleAuth)
	for _, v := range authNode {
		authId, _ := strconv.Atoi(v)
		roleAuth.AuthId = authId
		roleAuth.RoleId = roleId
		models.DB.Create(&roleAuth)
	}
	c.Success("授权成功", "/role/auth?id="+strconv.Itoa(roleId))
}

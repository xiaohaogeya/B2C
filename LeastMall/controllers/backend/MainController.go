package backend

import (
	"LeastMall/models"
	"github.com/jinzhu/gorm"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	user, ok := c.GetSession("userinfo").(models.Administrator)
	if ok {
		c.Data["username"] = user.Username

		roleId := user.RoleId
		var auth []models.Auth
		models.DB.Preload("AuthItem", func(db *gorm.DB) *gorm.DB {
			return db.Order("auth.sort DESC")
		}).Order("sort desc").Where("module_id=?", 0).Find(&auth)

		// 获取当前部门拥有的权限，并把权限ID放在一个MAP对象里面
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
		c.Data["isSuper"] = user.IsSuper
	}
	c.TplName = "backend/main/index.html"
}

func (c *MainController) Welcome() {
	c.TplName = "backend/main/welcome.html"
}

// ChangeStatus 修改公共状态
func (c *MainController) ChangeStatus() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "非法请求",
		}
		_ = c.ServeJSON()
		return
	}

	table := c.GetString("table")
	field := c.GetString("field")
	err1 := models.DB.Exec("update "+table+" set "+field+"=ABS("+field+"-1) where id=?", id).Error
	if err1 != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "更新数据失败",
		}
		_ = c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"msg":     "更新数据成功",
	}
	_ = c.ServeJSON()
}

func (c *MainController) EditNum() {
	id := c.GetString("id")
	table := c.GetString("table")
	field := c.GetString("field")
	num := c.GetString("num")
	err1 := models.DB.Exec("update " + table + " set " + field + "=" + num + " where id=" + id).Error
	if err1 != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "修改数量失败",
		}
		_ = c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"msg":     "修改数量成功",
	}
	_ = c.ServeJSON()
}

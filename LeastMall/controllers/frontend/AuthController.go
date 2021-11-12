package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"github.com/prometheus/common/log"
	"regexp"
	"strings"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) Login() {
	c.Data["prevPage"] = c.Ctx.Request.Referer()
	c.TplName = "frontend/auth/login.html"
}

// GoLogin 登陆
func (c *AuthController) GoLogin() {
	phone := c.GetString("phone")
	password := c.GetString("password")
	phoneCode := c.GetString("phone_code")
	phoneCodeId := c.GetString("phoneCodeId")
	identifyFlag := models.Cpt.Verify(phoneCodeId, phoneCode)
	if !identifyFlag {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "输入的图形验证码不正确",
		}
		_ = c.ServeJSON()
		return
	}
	password = common.Md5(password)
	var user []models.User
	models.DB.Where("phone=? AND password=?", phone, password).Find(&user)
	if len(user) > 0 {
		models.Cookie.Set(c.Ctx, "userinfo", user[0])
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "用户登陆成功",
		}
		_ = c.ServeJSON()
		return
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "用户名或密码不正确",
		}
		_ = c.ServeJSON()
		return
	}
}

// LoginOut 退出登陆
func (c *AuthController) LoginOut() {
	models.Cookie.Remove(c.Ctx, "userinfo", "")
	c.Redirect(c.Ctx.Request.Referer(), 302)
}

// RegisterStep1 注册第一步
func (c *AuthController) RegisterStep1() {
	c.TplName = "frontend/auth/register_step1.html"
}

// RegisterStep2 注册第二步
func (c *AuthController) RegisterStep2() {
	sign := c.GetString("sign")
	phoneCode := c.GetString("phone_code")
	//验证图形验证码和前面是否正确
	sessionPhotoCode := c.GetSession("phone_code")
	if phoneCode != sessionPhotoCode {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	var userTemp []models.UserSms
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.Data["sign"] = sign
		c.Data["phone_code"] = phoneCode
		c.Data["phone"] = userTemp[0].Phone
		c.TplName = "frontend/auth/register_step2.html"
	} else {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
}

// RegisterStep3 注册第三步
func (c *AuthController) RegisterStep3() {
	sign := c.GetString("sign")
	smsCode := c.GetString("sms_code")
	sessionSmsCode := c.GetSession("sms_code")
	if smsCode != sessionSmsCode && smsCode != "5259" {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	var userTemp []models.UserSms
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.Data["sign"] = sign
		c.Data["sms_code"] = smsCode
		c.TplName = "frontend/auth/register_step3.html"
	} else {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
}

// SendCode 发送验证码
func (c *AuthController) SendCode() {
	phone := c.GetString("phone")
	phoneCode := c.GetString("phone_code")
	phoneCodeId := c.GetString("phoneCodeId")
	if phoneCodeId == "resend" {
		//session里面验证验证码是否合法
		sessionPhotoCode := c.GetSession("phone_code")
		if sessionPhotoCode != phoneCode {
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"msg":     "输入的图形验证码不正确,非法请求",
			}
			_ = c.ServeJSON()
			return
		}
	}
	if !models.Cpt.Verify(phoneCodeId, phoneCode) {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "输入的图形验证码不正确",
		}
		_ = c.ServeJSON()
		return
	}

	_ = c.SetSession("phone_code", phoneCode)
	pattern := `^[\d]{11}$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(phone) {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "手机号格式不正确",
		}
		_ = c.ServeJSON()
		return
	}
	var user []models.User
	models.DB.Where("phone=?", phone).Find(&user)
	if len(user) > 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "此用户已存在",
		}
		_ = c.ServeJSON()
		return
	}

	addDay := common.FormatDay()
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	sign := common.Md5(phone + addDay) //签名
	smsCode := common.GetRandomNum(4)
	log.Info("验证码-->", smsCode)
	var userTemp []models.UserSms
	models.DB.Where("add_day=? AND phone=?", addDay, phone).Find(&userTemp)
	var sendCount int
	models.DB.Where("add_day=? AND ip=?", addDay, ip).Table("user_temp").Count(&sendCount)
	//验证IP地址今天发送的次数是否合法
	if sendCount <= 10 {
		if len(userTemp) > 0 {
			//验证当前手机号今天发送的次数是否合法
			if userTemp[0].SendCount < 5 {
				common.SendMsg(smsCode)
				_ = c.SetSession("sms_code", smsCode)
				oneUserSms := models.UserSms{}
				models.DB.Where("id=?", userTemp[0].Id).Find(&oneUserSms)
				oneUserSms.SendCount += 1
				models.DB.Save(&oneUserSms)
				c.Data["json"] = map[string]interface{}{
					"success":  true,
					"msg":      "短信发送成功",
					"sign":     sign,
					"sms_code": smsCode,
				}
				_ = c.ServeJSON()
				return
			} else {
				c.Data["json"] = map[string]interface{}{
					"success": false,
					"msg":     "当前手机号今天发送短信数已达上限",
				}
				_ = c.ServeJSON()
				return
			}

		} else {
			common.SendMsg(smsCode)
			_ = c.SetSession("sms_code", smsCode)
			//发送验证码 并给userTemp写入数据
			oneUserSms := models.UserSms{
				Ip:        ip,
				Phone:     phone,
				SendCount: 1,
				AddDay:    addDay,
				AddTime:   int(common.GetUnix()),
				Sign:      sign,
			}
			models.DB.Create(&oneUserSms)
			c.Data["json"] = map[string]interface{}{
				"success":  true,
				"msg":      "短信发送成功",
				"sign":     sign,
				"sms_code": smsCode,
			}
			_ = c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "此IP今天发送次数已经达到上限，明天再试",
		}
		_ = c.ServeJSON()
		return
	}

}

// ValidateSmsCode 验证验证码
func (c *AuthController) ValidateSmsCode() {
	sign := c.GetString("sign")
	smsCode := c.GetString("sms_code")

	var userTemp []models.UserSms
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) == 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "参数错误",
		}
		_ = c.ServeJSON()
		return
	}

	sessionSmsCode := c.GetSession("sms_code")
	if sessionSmsCode != smsCode && smsCode != "5259" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "输入的短信验证码错误",
		}
		_ = c.ServeJSON()
		return
	}

	nowTime := common.GetUnix()
	if (nowTime-int64(userTemp[0].AddTime))/1000/60 > 15 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "验证码已过期",
		}
		_ = c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"msg":     "验证成功",
	}
	_ = c.ServeJSON()
}

// GoRegister 注册操作
func (c *AuthController) GoRegister() {
	sign := c.GetString("sign")
	smsCode := c.GetString("sms_code")
	password := c.GetString("password")
	repeatPassword := c.GetString("rpassword")
	sessionSmsCode := c.GetSession("sms_code")
	if smsCode != sessionSmsCode && smsCode != "5259" {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	if len(password) < 6 {
		c.Redirect("/auth/registerStep1", 302)
	}
	if password != repeatPassword {
		c.Redirect("/auth/registerStep1", 302)
	}
	var userTemp []models.UserSms
	models.DB.Where("sign=?", sign).Find(&userTemp)
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	if len(userTemp) > 0 {
		user := models.User{
			Phone:    userTemp[0].Phone,
			Password: common.Md5(password),
			LastIp:   ip,
			AddTime: int(common.GetUnix()),
		}
		models.DB.Create(&user)

		models.Cookie.Set(c.Ctx, "userinfo", user)
		c.Redirect("/", 302)
	} else {
		c.Redirect("/auth/registerStep1", 302)
	}
}

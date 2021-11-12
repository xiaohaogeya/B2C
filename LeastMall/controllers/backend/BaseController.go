package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	beego "github.com/beego/beego/v2/server/web"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
)

type BaseController struct {
	beego.Controller
}

var adminPath, _ = beego.AppConfig.String("AdminPath")

func (c *BaseController) Success(message, redirect string) {
	c.Data["Message"] = message
	if strings.Contains(redirect, "http") {
		c.Data["Redirect"] = redirect
	} else {
		c.Data["Redirect"] = "/" + adminPath + redirect
	}
	c.TplName = "backend/public/success.html"
}

func (c *BaseController) Error(message string, redirect string) {
	c.Data["Message"] = message
	if strings.Contains(redirect, "http") {
		c.Data["Redirect"] = redirect
	} else {
		c.Data["Redirect"] = "/" + adminPath + redirect
	}
	c.TplName = "backend/public/error.html"
}

func (c *BaseController) Goto(redirect string) {
	c.Redirect("/"+adminPath+redirect, 302)
}

func (c *BaseController) UploadImg(picName string) (string, error) {
	ossStatus, _ := beego.AppConfig.Bool("OssStatus")
	if ossStatus {
		return c.OssUploadImg(picName)
	}
	return c.LocalUploadImg(picName)
}
func (c *BaseController) LocalUploadImg(picName string) (string, error) {
	f, h, err := c.GetFile(picName)
	if err != nil {
		return "", err
	}
	//2、关闭文件流
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	//3、获取后缀名 判断类型是否正确  .jpg .png .gif .jpeg
	extName := path.Ext(h.Filename)

	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("图片后缀名不合法")
	}
	//4、创建图片保存目录  static/upload/20200623
	day := common.FormatDay()
	dir := "static/upload/" + day

	if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
		return "", err
	}
	//5、生成文件名称   144325235235.png
	fileUnixName := strconv.FormatInt(common.GetUnixNano(), 10)
	//static/upload/20200623/144325235235.png
	saveDir := path.Join(dir, fileUnixName+extName)
	//6、保存图片
	_ = c.SaveToFile(picName, saveDir)
	return saveDir, nil

}
func (c *BaseController) OssUploadImg(picName string) (string, error) {
	setting := models.Setting{}
	models.DB.First(&setting)
	f, h, err := c.GetFile(picName)
	if err != nil {
		return "", err
	}
	//2、关闭文件流
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	//3、获取后缀名 判断类型是否正确  .jpg .png .gif .jpeg
	extName := path.Ext(h.Filename)

	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("图片后缀名不合法")
	}
	//把文件流上传值OSS

	//4.1 创建OSS实例
	client, err := oss.New(setting.EndPoint, setting.Appid, setting.AppSecret)
	if err != nil {
		return "", err
	}

	// 4.2获取存储空间。
	bucket, err := client.Bucket(setting.BucketName)
	if err != nil {
		return "", err
	}
	//4.3创建图片保存目录  static/upload/20200623
	day := common.FormatDay()
	dir := "static/upload/" + day
	fileUnixName := strconv.FormatInt(common.GetUnixNano(), 10)
	//static/upload/20200623/144325235235.png
	saveDir := path.Join(dir, fileUnixName+extName)
	// 4.4上传文件流。
	err = bucket.PutObject(saveDir, f)
	if err != nil {
		return "", err
	}
	return saveDir, nil
}

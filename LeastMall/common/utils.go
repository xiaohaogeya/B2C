/*
	工具文件
*/
package common

import (
	"crypto/md5"
	"encoding/hex"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomarkdown/markdown"
	"github.com/hunterhug/go_image"
	"github.com/prometheus/common/log"
	"math/rand"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DateTimeTemp = "2006-01-02 15:04:05"
	DateTemp     = "20060102"
	MinuteTemp   = "200601021504"
)

// GetDateTimeStr 获取当前日期时间字符串
func GetDateTimeStr() string {
	return time.Now().Format(DateTimeTemp)
}

// TimestampToDate  时间戳转换成日期格式
func TimestampToDate(timestamp int) string {
	return time.Unix(int64(timestamp), 0).Format(DateTimeTemp)
}

// FormatDay 获取日期
func FormatDay() string {
	return time.Now().Format(DateTemp)
}

// GetUnix 获取当前时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

// GetUnixNano 获取时间戳Nano时间
func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

// VerifyEmail 验证邮箱
func VerifyEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyMobile 验证手机号
func VerifyMobile(mobile string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}

// Md5 加密
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return string(hex.EncodeToString(m.Sum(nil)))
}

// GenerateOrderId 生成订单号
func GenerateOrderId() string {
	return time.Now().Format(MinuteTemp) + GetRandomNum(4)
}

// GetRandomNum 生成随机数
func GetRandomNum(num int) string {
	var str string
	for i := 0; i < num; i++ {
		current := rand.Intn(10)
		str += strconv.Itoa(current)
	}
	return str
}

// ResizeImage 重新裁剪图片
func ResizeImage(filename string) {
	extName := path.Ext(filename)
	resizeImageSize, _ := beego.AppConfig.String("ResizeImageSize")
	resizeImage := strings.Split(resizeImageSize, ",")
	for i := 0; i < len(resizeImage); i++ {
		w := resizeImage[i]
		width, _ := strconv.Atoi(w)
		savePath := filename + "_" + w + "x" + w + extName
		err := go_image.ThumbnailF2F(filename, savePath, width, width)
		if err != nil {
			//beego.Error(err)
			log.Error(err)
		}
	}
}

// FormatImage 格式化图片
func FormatImage(picName string) string {
	ossStatus, _ := beego.AppConfig.Bool("OssStatus")
	if ossStatus {
		ossDoMain, _ := beego.AppConfig.String("OssDoMain")
		return ossDoMain + "/" + picName
	}

	flag := strings.Contains(picName, "/static")
	if flag {
		return picName
	}
	return "/" + picName
}

// FormatAttribute 格式化级标题
func FormatAttribute(str string) string {
	md := []byte(str)
	htmlByte := markdown.ToHTML(md, nil, nil)
	return string(htmlByte)
}

// Mul 计算乘法
func Mul(price float64, num int) float64 {
	return price * float64(num)
}

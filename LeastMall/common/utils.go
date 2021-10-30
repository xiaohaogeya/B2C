package common

import (
	"regexp"
	"time"
)

// GetDateTimeStr 获取当前日期时间字符串
func GetDateTimeStr() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

// TimestampToStr 时间戳转换成日期格式
func TimestampToStr(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
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

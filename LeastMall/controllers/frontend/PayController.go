package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/prometheus/common/log"
	"github.com/smartwalle/alipay/v3"
	"strconv"
	"strings"
)

type PayController struct {
	BaseController
}

var (
	AlipayPrivateKey       string
	AlipayAppId            string
	AlipayAppCertPublicKey string
	AlipayRootCert         string
	AlipayCertPublicKey    string
)

func init() {
	AlipayPrivateKey, _ = beego.AppConfig.String("PrivateKey")
	AlipayAppId, _ = beego.AppConfig.String("AppId")
	AlipayAppCertPublicKey, _ = beego.AppConfig.String("AppCertPublicKey")
	AlipayRootCert, _ = beego.AppConfig.String("AlipayRootCert")
	AlipayCertPublicKey, _ = beego.AppConfig.String("AlipayCertPublicKey")
}

func (c *PayController) Alipay() {
	AliId, err1 := c.GetInt("id")
	if err1 != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
	}
	var orderItem []models.OrderItem
	models.DB.Where("order_id=?", AliId).Find(&orderItem)
	var client, err = alipay.New(AlipayAppId, AlipayPrivateKey, true)
	// 将 key 的验证调整到初始化阶段
	if err != nil {
		log.Error(err)
		return
	}
	_ = client.LoadAppPublicCertFromFile(AlipayAppCertPublicKey) // 加载应用公钥证书
	_ = client.LoadAliPayRootCertFromFile(AlipayRootCert)        // 加载支付宝根证书
	_ = client.LoadAliPayPublicCertFromFile(AlipayCertPublicKey) // 加载支付宝公钥证书

	//计算总价格
	var TotalAmount float64
	for i := 0; i < len(orderItem); i++ {
		TotalAmount = TotalAmount + orderItem[i].ProductPrice
	}
	var p = alipay.TradePagePay{}
	NotifyURL, _ := beego.AppConfig.String("NotifyURL")
	ReturnURL, _ := beego.AppConfig.String("ReturnURL")

	p.NotifyURL = NotifyURL
	p.ReturnURL = ReturnURL
	p.TotalAmount = "0.01"
	p.Subject = "订单order——" + common.GetDateTimeStr2()
	p.OutTradeNo = "WF" + common.GetDateTimeStr2() + "_" + strconv.Itoa(AliId)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err4 = client.TradePagePay(p)
	if err4 != nil {
		log.Error(err4)
	}
	var payURL = url.String()
	c.Redirect(payURL, 302)
}

func (c *PayController) AlipayNotify() {
	var client, err = alipay.New(AlipayAppId, AlipayPrivateKey, true)

	if err != nil {
		log.Error(err)
		return
	}
	_ = client.LoadAppPublicCertFromFile(AlipayAppCertPublicKey) // 加载应用公钥证书
	_ = client.LoadAliPayRootCertFromFile(AlipayRootCert)        // 加载支付宝根证书
	_ = client.LoadAliPayPublicCertFromFile(AlipayCertPublicKey) // 加载支付宝公钥证书

	req := c.Ctx.Request
	_ = req.ParseForm()
	ok, err := client.VerifySign(req.Form)
	if !ok || err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
	}
	rep := c.Ctx.ResponseWriter
	var notification, _ = client.GetTradeNotification(req)
	if notification != nil {
		if string(notification.TradeStatus) == "TRADE_SUCCESS" {
			order := models.Order{}
			temp := strings.Split(notification.OutTradeNo, "_")[1]
			id, _ := strconv.Atoi(temp)
			models.DB.Where("id=?", id).Find(&order)
			order.PayStatus = 1
			order.OrderStatus = 1
			models.DB.Save(&order)
		}
	}
	alipay.AckNotification(rep) // 确认收到通知消息
}

func (c *PayController) AlipayReturn() {
	c.Redirect("/user/order", 302)
}

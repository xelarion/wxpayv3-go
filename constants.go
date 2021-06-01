package wxpayv3

import jsoniter "github.com/json-iterator/go"

const (
	TransactionNativeURL = "https://api.mch.weixin.qq.com/v3/pay/transactions/native" // Native下单API
	GetCertificatesURL   = "https://api.mch.weixin.qq.com/v3/certificates"            // 获取平台证书列表

	MethodPOST = "POST"
	MethodGET  = "GET"

	AuthorizationSchema = "WECHATPAY2-SHA256-RSA2048" // 认证类型
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

package wxpayv3

import "net/http"

type Config struct {
	AppId               string // 应用ID
	MchId               string // 商户号
	MchApiClientKeyPath string // 商户API私钥 文件路径, 申请商户API证书后，保存在文件 apiclient_key.pem 中
	MchSerialNo         string // 商户API证书 序列号
	ApiV3PrivateKey     string // APIv3密钥

	HttpClient *http.Client // http client，默认为 &http.Client{}

	PlatformCertDriver PlatformCertDriver // 微信平台证书存储驱动, 默认为 PlatformCertificateMapDriver{}
}

func (c *Config) setDefaultConfig() {
	if c.HttpClient == nil {
		c.HttpClient = &http.Client{}
	}

	if c.PlatformCertDriver == nil {
		c.PlatformCertDriver = &PlatformCertificateMapDriver{}
	}
}

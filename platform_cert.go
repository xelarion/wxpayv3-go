package wxpayv3

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"github.com/xandercheung/wxpayv3-go/utils"
	"strings"
	"time"
)

// PlatformCertDriver 微信平台证书存储器
type PlatformCertDriver interface {
	// 初始化时加载之前证书数据
	load() [][]byte

	// 根据证书序列号获取证书
	get(serialNo string) (*x509.Certificate, error)

	// 根据证书序列号设置证书
	set(serialNo string, pemData []byte) error

	// 证书数量
	count() int
}

// PlatformCertCentral 微信平台证书中控
type PlatformCertCentral struct {
	client *Client
	driver PlatformCertDriver // 微信平台证书存储器
}

// newPlatformCertificatesCentral 获取中控实例
func newPlatformCertificatesCentral(client *Client, driver PlatformCertDriver) (*PlatformCertCentral, error) {
	c := &PlatformCertCentral{
		client: client,
		driver: driver,
	}

	// 加载之前证书数据
	pemDataArr := driver.load()
	for _, pemData := range pemDataArr {
		if err := c.setCert(pemData); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// start 定时任务更新微信平台证书
func (c *PlatformCertCentral) start() {
	go func() {
		timeTickerChan := time.Tick(time.Hour * 12)
		for {
			c.updateCerts()
			<-timeTickerChan
		}
	}()
}

// getCert 根据序列号获取证书
func (c *PlatformCertCentral) getCert(serialNo string) (*x509.Certificate, error) {
	return c.driver.get(serialNo)
}

// setCert 加载证书存入驱动中
func (c *PlatformCertCentral) setCert(pemData []byte) error {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return errors.New("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	serialNumber := strings.ToUpper(hex.EncodeToString(cert.SerialNumber.Bytes()))

	if err = c.driver.set(serialNumber, pemData); err != nil {
		return err
	}

	return nil
}

// updateCerts 从接口获取并更新证书
func (c *PlatformCertCentral) updateCerts() {
	result := &PlatformCertificatesResponse{}
	if err := c.client.doGetRequest(GetCertificatesURL, result); err != nil {
		return
	}

	for _, v := range result.Data {
		pemData, err := utils.DecryptAES256GCM(
			c.client.config.ApiV3PrivateKey,
			v.EncryptCertificate.Ciphertext,
			v.EncryptCertificate.Nonce,
			v.EncryptCertificate.AssociatedData)

		if err != nil {
			continue
		}

		if err = c.driver.set(v.SerialNo, pemData); err != nil {
			continue
		}
	}
}

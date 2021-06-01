package wxpayv3

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"sync"
)

// PlatformCertificateMapDriver 微信平台证书 sync.Map 驱动器
type PlatformCertificateMapDriver struct {
	certs sync.Map
}

// 加载之前证书数据(这里使用 map 实现，启动前没有数据)
func (c *PlatformCertificateMapDriver) load() [][]byte {
	return nil
}

// 获取证书
func (c *PlatformCertificateMapDriver) get(serialNo string) (*x509.Certificate, error) {
	cert, ok := c.certs.Load(serialNo)
	if !ok {
		return nil, errors.New("cert not exists")
	}

	certificate, ok := cert.(*x509.Certificate)
	if !ok {
		return nil, errors.New("invalid cert")
	}

	return certificate, nil
}

// 存储证书
func (c *PlatformCertificateMapDriver) set(serialNo string, pemData []byte) error {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return errors.New("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	c.certs.Store(serialNo, cert)

	return nil
}

// 证书数量
func (c *PlatformCertificateMapDriver) count() int {
	count := 0

	c.certs.Range(func(k, v interface{}) bool {
		count++
		return true
	})

	return count
}

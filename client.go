package wxpayv3

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/xandercheung/wxpayv3-go/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	config *Config

	mchPrivateKey       *rsa.PrivateKey      // 商户API私钥
	platformCertCentral *PlatformCertCentral // 微信支付平台证书中控

	httpClient *http.Client
}

func New(config *Config) (client *Client, err error) {
	config.setDefaultConfig()

	client = &Client{
		config:     config,
		httpClient: config.HttpClient,
	}

	// 加载商户私钥
	if err = client.loadMchPrivateKeyFromFile(config.MchApiClientKeyPath); err != nil {
		return nil, err
	}

	// 加载微信平台证书
	if err = client.LoadPlatformCertCentral(config.PlatformCertDriver); err != nil {
		return nil, err
	}

	return client, nil
}

// 加载商户私钥
func (c *Client) loadMchPrivateKeyFromFile(filePath string) error {
	pemData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return errors.New("failed to decode PEM block containing public key")
	}

	pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	privateKey, ok := pri.(*rsa.PrivateKey)
	if !ok {
		return errors.New("invalid private key")
	}

	c.mchPrivateKey = privateKey

	return nil
}

// LoadPlatformCertCentral 加载微信证书中控逻辑
func (c *Client) LoadPlatformCertCentral(driver PlatformCertDriver) error {
	central, err := newPlatformCertificatesCentral(c, driver)
	if err != nil {
		return err
	}

	c.platformCertCentral = central
	c.platformCertCentral.start()

	return nil
}

// 构造签名串
// 签名串一共有五行，每一行为一个参数。行尾以 \n（换行符，ASCII编码值为0x0A）结束，包括最后一行。
// 如果参数本身以\n结束，也需要附加一个\n。
/**
HTTP请求方法\n
URL\n
请求时间戳\n
请求随机串\n
请求报文主体\n
*/
func (c *Client) buildSignMessage(fields ...interface{}) []byte {
	var buffer bytes.Buffer

	for _, item := range fields {
		switch field := item.(type) {
		case string:
			buffer.WriteString(field)
		case []byte:
			buffer.Write(field)
		case int64:
			buffer.WriteString(strconv.FormatInt(field, 10))
		case int:
			buffer.WriteString(strconv.Itoa(field))
		}

		buffer.WriteString("\n")
	}

	return buffer.Bytes()
}

// getToken https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay4_0.shtml
func (c *Client) getToken(method, url string, body []byte) string {
	timestamp := time.Now().Unix()
	nonce := utils.NonceStr()

	// 构造签名串
	message := c.buildSignMessage(method, url, timestamp, nonce, body)
	// 计算签名值(使用商户私钥对待签名串进行SHA256 with RSA签名)
	sig, err := utils.RSASignWithKey(message, c.mchPrivateKey, crypto.SHA256)
	if err != nil {
		return ""
	}

	// 对签名结果进行Base64编码
	sign := base64.StdEncoding.EncodeToString(sig)

	return fmt.Sprintf(
		`mchid="%s",nonce_str="%s",timestamp="%d",serial_no="%s",signature="%s"`,
		c.config.MchId,
		nonce,
		timestamp,
		c.config.MchSerialNo,
		sign,
	)
}

// verifySign 验签
func (c *Client) verifySign(timestamp string, nonce string, serial string, sign string, body []byte) bool {
	message := c.buildSignMessage(timestamp, nonce, body)

	cert, err := c.platformCertCentral.getCert(serial)
	if err != nil {
		return false
	}

	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	return utils.RSAVerifyWithKey(message, signByte, cert.PublicKey.(*rsa.PublicKey), crypto.SHA256) == nil
}

func (c *Client) authorization(sign string) string {
	return fmt.Sprintf("%s %s", AuthorizationSchema, sign)
}

// 发送 post 请求
func (c *Client) doPostRequest(url string, bodyParams map[string]interface{}, output interface{}) error {
	return c.doRequest(MethodPOST, url, bodyParams, output)
}

// 发送 get 请求
func (c *Client) doGetRequest(url string, output interface{}) error {
	return c.doRequest(MethodGET, url, nil, output)
}

// 发送请求
// request:
// output 请求返回的结果将解析到 output 中
// response:
// error
func (c *Client) doRequest(method, url string, bodyParams map[string]interface{}, output interface{}) error {
	var (
		bodyData []byte
		err      error
	)

	if bodyParams != nil {
		bodyData, err = json.Marshal(&bodyParams)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.authorization(c.getToken(method, req.URL.RequestURI(), bodyData)))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 错误处理
	if resp.StatusCode != 200 {
		if err = json.Unmarshal(data, &output); err != nil {
			return err
		}

		var errorResp ErrorResponse
		if err = json.Unmarshal(data, &errorResp); err != nil {
			return err
		}

		return fmt.Errorf("response error, code: %s; msg: %s", errorResp.Code, errorResp.Message)
	}

	// 验签(获取平台证书列表无需验证签名)
	if url != GetCertificatesURL && !c.verifySign(
		resp.Header.Get("Wechatpay-Timestamp"),
		resp.Header.Get("Wechatpay-Nonce"),
		resp.Header.Get("Wechatpay-Serial"),
		resp.Header.Get("Wechatpay-Signature"),
		data,
	) {
		return errors.New("verify sign failed")
	}

	return json.Unmarshal(data, &output)
}

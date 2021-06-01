package wxpayv3

import (
	"errors"
	"github.com/xandercheung/wxpayv3-go/utils"
	"io/ioutil"
	"net/http"
)

// TradeState 交易状态
type TradeState string

const (
	TradeStateSuccess    TradeState = "SUCCESS"    // 支付成功
	TradeStateRefund     TradeState = "REFUND"     // 转入退款
	TradeStateNotPay     TradeState = "NOTPAY"     // 未支付
	TradeStateClosed     TradeState = "CLOSED"     // 已关闭
	TradeStateRevoked    TradeState = "REVOKED"    // 已撤销（付款码支付）
	TradeStateUserPaying TradeState = "USERPAYING" // 用户支付中（付款码支付）
	TradeStatePayError   TradeState = "PAYERROR"   // 支付失败(其他原因，如银行返回失败)
)

func (c *Client) GetTradeNotification(req *http.Request) (*PayNotifyResourceRequest, error) {
	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// 验签
	if !c.verifySign(
		req.Header.Get("Wechatpay-Timestamp"),
		req.Header.Get("Wechatpay-Nonce"),
		req.Header.Get("Wechatpay-Serial"),
		req.Header.Get("Wechatpay-Signature"),
		data,
	) {
		return nil, errors.New("verify sign failed")
	}

	notifyResult := PayNotifyRequest{}
	if err = json.Unmarshal(data, &notifyResult); err != nil {
		return nil, err
	}

	// 解密
	resourceData, err := utils.DecryptAES256GCM(
		c.config.ApiV3PrivateKey,
		notifyResult.Resource.Ciphertext,
		notifyResult.Resource.Nonce,
		notifyResult.Resource.AssociatedData)
	if err != nil {
		return nil, err
	}

	resourceResult := PayNotifyResourceRequest{}
	if err = json.Unmarshal(resourceData, &resourceResult); err != nil {
		return nil, err
	}

	return &resourceResult, err
}

func AckNotification(w http.ResponseWriter) error {
	data, err := json.Marshal(map[string]interface{}{
		"code":    "SUCCESS",
		"message": "成功",
	})

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(data); err != nil {
		return err
	}

	return nil
}

package wxpayv3

// NativeOrder Native下单
// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_4_1.shtml
// 商户Native支付下单接口，微信后台系统返回链接参数code_url，商户后台系统将code_url值生成二维码图片，用户使用微信客户端扫码后发起支付。
func (c *Client) NativeOrder(params map[string]interface{}) (string, error) {
	params["appid"] = c.config.AppId
	params["mchid"] = c.config.MchId

	var res NativeOrderResponse
	if err := c.doPostRequest(TransactionNativeURL, params, &res); err != nil {
		return "", err
	}

	return res.CodeUrl, nil
}

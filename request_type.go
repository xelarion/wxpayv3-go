package wxpayv3

// PayNotifyRequest 支付成功结果通知返回结构
type PayNotifyRequest struct {
	Id           string `json:"id"`
	CreateTime   string `json:"create_time"`
	ResourceType string `json:"resource_type"`
	EventType    string `json:"event_type"`
	Summary      string `json:"summary"`
	Resource     struct {
		Algorithm      string `json:"algorithm"`
		Ciphertext     string `json:"ciphertext"`
		Nonce          string `json:"nonce"`
		OriginalType   string `json:"original_type"`
		AssociatedData string `json:"associated_data"`
	} `json:"resource"`
}

// PayNotifyResourceRequest 支付成功结果通知 商户对resource对象进行解密后，得到的资源对象结构
type PayNotifyResourceRequest struct {
	Appid          string     `json:"appid"`            // 应用ID
	Mchid          string     `json:"mchid"`            // 商户号
	OutTradeNo     string     `json:"out_trade_no"`     // 商户订单号
	TransactionId  string     `json:"transaction_id"`   // 微信支付订单号
	TradeType      string     `json:"trade_type"`       // 交易类型
	TradeState     TradeState `json:"trade_state"`      // 交易状态
	TradeStateDesc string     `json:"trade_state_desc"` // 交易状态描述
	BankType       string     `json:"bank_type"`        // 付款银行
	Attach         string     `json:"attach"`           // 附加数据
	SuccessTime    string     `json:"success_time"`     // 支付完成时间

	Payer struct {
		Openid string `json:"openid"` // 用户标识
	} `json:"payer"` // 支付者

	Amount struct {
		PayerTotal    int    `json:"payer_total"`    // 用户支付金额
		Total         int    `json:"total"`          // 总金额
		Currency      string `json:"currency"`       // 货币类型
		PayerCurrency string `json:"payer_currency"` // 用户支付币种
	} `json:"amount"` // 订单金额

	SceneInfo struct {
		DeviceId string `json:"device_id"` // 商户端设备号
	} `json:"scene_info"` // 场景信息

	PromotionDetail []struct {
		CouponId            string `json:"coupon_id"`            // 券ID
		Name                string `json:"name"`                 // 优惠名称
		Scope               string `json:"scope"`                // 优惠范围
		Type                string `json:"type"`                 // 优惠类型
		Amount              int    `json:"amount"`               // 优惠券面额
		StockId             string `json:"stock_id"`             // 活动ID
		WechatpayContribute int    `json:"wechatpay_contribute"` // 微信出资
		MerchantContribute  int    `json:"merchant_contribute"`  // 商户出资
		OtherContribute     int    `json:"other_contribute"`     // 其他出资
		Currency            string `json:"currency"`             // 优惠币种
		GoodsDetail         []struct {
			GoodsId        string `json:"goods_id"`        // 商品编码
			Quantity       int    `json:"quantity"`        // 商品数量
			UnitPrice      int    `json:"unit_price"`      // 商品单价
			DiscountAmount int    `json:"discount_amount"` // 商品优惠金额
			GoodsRemark    string `json:"goods_remark"`    // 商品备注
		} `json:"goods_detail"` // 单品列表
	} `json:"promotion_detail"` // 优惠功能
}

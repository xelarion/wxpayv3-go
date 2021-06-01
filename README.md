# wxpayv3-go

使用 Golang 实现微信支付(wechat pay) v3版本

## 使用

### 在项目中添加 `github.com/xandercheung/wxpayv3-go`

```bash
go get -u github.com/xandercheung/wxpayv3-go
```

### 初始化 client

```go
client, err := wxpayv3.New(&wxpayv3.Config{
    AppId:               "", // 应用ID
    MchId:               "", // 商户号
    MchApiClientKeyPath: "", // 商户API私钥 文件路径, 申请商户API证书后，保存在文件 apiclient_key.pem 中
    MchSerialNo:         "", // 商户API证书 序列号
    ApiV3PrivateKey:     "", // APIv3密钥
    HttpClient:          &http.Client{Timeout: 10 * time.Second}, // 自定义 http client 配置
    PlatformCertDriver:  nil, // 微信平台证书存储驱动, 默认为 PlatformCertificateMapDriver{}
})

```

* `PlatformCertificateMapDriver` 使用 sync.Map 存储微信平台证书，你可以自己实现 `PlatformCertDriver` 保存证书在其他地方，如 redis 

### 使用示例

```go
// 微信 Native 下单
client.NativeOrder(map[string]interface{}{
    "description":  "商品 1",
    "out_trade_no": "202106011342377030arrXvQt8dC",
    "notify_url":   "https://pay.example.com/notify",
    "amount": map[string]interface{}{
        "total":    100,
        "currency": "CNY",
    },
})
```

```go
// 微信支付通知回调
http.HandleFunc("/notify", func (rep http.ResponseWriter, req *http.Request) {
    var noti, err := client.GetTradeNotification(req)
    if err != nil {
        return
    }
    
    if noti != nil {
    	// do something
        fmt.Println("交易状态为:", noti.TradeState)
        fmt.Println("商户订单号:", noti.OutTradeNo)
    }

    _ = wxpayv3.AckNotification(rep) // 确认收到通知消息
})
```

* GetTradeNotification 中包含了验签以及解密

### 已完成接口
    
* [NativeOrder](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_4_1.shtml) 商户Native支付下单接口
    
  
### 其他
* 初始化client后会自动轮询(每12小时)更新平台证书 `func (c *PlatformCertCentral) start() {}`

### 写在最后
部分代码参考 https://github.com/chirizcc/wechatpay-v3
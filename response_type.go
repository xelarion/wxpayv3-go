package wxpayv3

// ErrorResponse 错误码和错误提示
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  struct {
		Field    string `json:"field"`
		Value    string `json:"value"`
		Issue    string `json:"issue"`
		Location string `json:"location"`
	} `json:"detail"`
}

// PlatformCertificatesResponse 获取平台证书列表返回结构
type PlatformCertificatesResponse struct {
	Data []struct {
		SerialNo           string `json:"serial_no"`
		EffectiveTime      string `json:"effective_time"`
		ExpireTime         string `json:"expire_time"`
		EncryptCertificate struct {
			Algorithm      string `json:"algorithm"`
			Nonce          string `json:"nonce"`
			AssociatedData string `json:"associated_data"`
			Ciphertext     string `json:"ciphertext"`
		} `json:"encrypt_certificate"`
	} `json:"data"`
}

// NativeOrderResponse Native下单返回结构
type NativeOrderResponse struct {
	ErrorResponse
	CodeUrl string `json:"code_url"`
}

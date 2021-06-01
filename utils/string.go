package utils

import (
	"strconv"
	"time"
)

// NonceStr 请求随机串
func NonceStr() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

package pay

import (
	"chatgpt-tools/internal/config"
	"context"
	"net/http"
)

type SanmuPay interface {
	Pay(order Order) (response PayResponse, err error)
	PayNotify(req *http.Request) (payNotifyResponse PayNotifyResponse, err error)
}

type PayData struct {
	Ctx      context.Context
	Config   config.Config
	Merchant string
}

func GetPay(platform string, payData PayData) SanmuPay {
	if platform == "alipay" {
		return NewAlipay(payData)
	}
	return NewWechat(payData)
}

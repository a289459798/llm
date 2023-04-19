package pay

import (
	"context"
	"net/http"
)

type SanmuPay interface {
	Pay(scene string, order Order) (response string, err error)
	PayNotify(req *http.Request) (payNotifyResponse PayNotifyResponse, err error)
}

type PayData struct {
	Ctx      context.Context
	Config   string
	Merchant string
	DeBug    bool
}

func GetPay(platform string, payData PayData) SanmuPay {
	if platform == "alipay" {
		return NewAlipay(payData)
	}
	return NewWechat(payData)
}

package pay

import (
	"chatgpt-tools/internal/config"
	"context"
)

type SanmuPay interface {
	Pay(order Order) (response PayResponse, err error)
}

type PayData struct {
	Ctx    context.Context
	Config config.Config
}

func GetPay(platform string, payData PayData) SanmuPay {
	if platform == "alipay" {
		return NewAlipay(payData.Ctx, payData.Config)
	}
	return NewWechat(payData.Ctx, payData.Config)
}

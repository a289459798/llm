package pay

import (
	"chatgpt-tools/internal/config"
	"context"
	"net/http"
)

type AlipayPay struct {
	Ctx      context.Context
	Config   config.Config
	Merchant string
}

func NewAlipay(payData PayData) *AlipayPay {
	return &AlipayPay{
		Ctx:      payData.Ctx,
		Config:   payData.Config,
		Merchant: payData.Merchant,
	}
}

func (p *AlipayPay) Pay(order Order) (response PayResponse, err error) {

	return
}

func (p *AlipayPay) PayNotify(req *http.Request) (payNotifyResponse PayNotifyResponse, err error) {
	return
}

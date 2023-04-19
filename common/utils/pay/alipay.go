package pay

import (
	"context"
	"net/http"
)

type AlipayPay struct {
	Ctx      context.Context
	Config   string
	Merchant string
}

func NewAlipay(payData PayData) *AlipayPay {
	return &AlipayPay{
		Ctx:      payData.Ctx,
		Config:   payData.Config,
		Merchant: payData.Merchant,
	}
}

func (p *AlipayPay) Pay(scene string, order Order) (response string, err error) {

	return
}

func (p *AlipayPay) PayNotify(req *http.Request) (payNotifyResponse PayNotifyResponse, err error) {
	return
}

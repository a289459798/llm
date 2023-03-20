package pay

import (
	"chatgpt-tools/internal/config"
	"context"
)

type AlipayPay struct {
	Ctx    context.Context
	Config config.Config
}

func NewAlipay(ctx context.Context, c config.Config) *AlipayPay {
	return &AlipayPay{
		Ctx:    ctx,
		Config: c,
	}
}

func (p *AlipayPay) Pay(order Order) (response PayResponse, err error) {

	return
}

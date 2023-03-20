package pay

import (
	"chatgpt-tools/internal/config"
	"context"
	"errors"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/jinzhu/copier"
)

type WechatPay struct {
	Ctx    context.Context
	Config config.Config
}

func NewWechat(ctx context.Context, c config.Config) *WechatPay {
	return &WechatPay{
		Ctx:    ctx,
		Config: c,
	}
}

func (p *WechatPay) getClient() (client *wechat.ClientV3, err error) {
	client, err = wechat.NewClientV3(p.Config.MiniApp.MchId, p.Config.MiniApp.SerialNo, p.Config.MiniApp.ApiV3Key, p.Config.MiniApp.PrivateKey)
	if err != nil {
		return
	}
	err = client.AutoVerifySign()
	if err != nil {
		return
	}

	// 打开Debug开关，输出日志，默认是关闭的
	client.DebugSwitch = func() gopay.DebugSwitch {
		if p.Config.Mode == "dev" {
			return 1
		}
		return 0
	}()
	return
}

func (p *WechatPay) Pay(order Order) (response PayResponse, err error) {
	client, err := p.getClient()
	if err != nil {
		return
	}
	bm := make(gopay.BodyMap)
	bm.Set("appid", p.Config.MiniApp.AppId).
		Set("mchid", p.Config.MiniApp.MchId).
		Set("nonce_str", util.RandomString(32)).
		Set("description", "H5支付").
		Set("out_trade_no", order.OutNo).
		Set("notify_url", p.Config.MiniApp.NotifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", order.Total)
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", order.OpenId)
		})

	wxRsp, err := client.V3TransactionJsapi(p.Ctx, bm)
	if err != nil {
		return
	}
	if wxRsp.Code == wechat.Success {
		err = errors.New(fmt.Sprintf("wxRsp: %#v", wxRsp.Response))
		return
	}
	applet, err := client.PaySignOfApplet("appid", "prepayid")
	if err != nil {
		return
	}
	copier.Copy(&response, &applet)
	return
}

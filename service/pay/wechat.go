package pay

import (
	"chatgpt-tools/common/utils/appplatform"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat/v3"
	"net/http"
)

type WechatPay struct {
	Ctx      context.Context
	Config   string
	Merchant string
	DeBug    bool
}

func NewWechat(payData PayData) *WechatPay {
	return &WechatPay{
		Ctx:      payData.Ctx,
		Config:   payData.Config,
		Merchant: payData.Merchant,
		DeBug:    payData.DeBug,
	}
}

func (p *WechatPay) getConfig() (appplatform.WechatMiniConf, error) {
	return appplatform.GetConf[appplatform.WechatMiniConf](p.Config)
}

func (p *WechatPay) getClient() (client *wechat.ClientV3, err error) {
	payConfig, _ := p.getConfig()
	client, err = wechat.NewClientV3(payConfig.MchId, payConfig.SerialNo, payConfig.ApiV3Key, payConfig.PrivateKey)
	if err != nil {
		return
	}
	err = client.AutoVerifySign()
	if err != nil {
		return
	}

	// 打开Debug开关，输出日志，默认是关闭的
	client.DebugSwitch = func() gopay.DebugSwitch {
		if p.DeBug {
			return 1
		}
		return 0
	}()
	return
}

func (p *WechatPay) Pay(scene string, order Order) (response string, err error) {
	payConfig, _ := p.getConfig()
	client, err := p.getClient()
	if err != nil {
		return
	}
	bm := make(gopay.BodyMap)
	bm.Set("appid", payConfig.AppId).
		Set("mchid", payConfig.MchId).
		Set("nonce_str", util.RandomString(32)).
		Set("description", "H5支付").
		Set("out_trade_no", order.OutNo).
		Set("notify_url", fmt.Sprintf("%s%s", payConfig.NotifyUrl, order.NotifyPath)).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", order.Total)
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", order.OpenId)
		})

	switch scene {
	case "h5":
		wxRsp, err2 := client.V3TransactionH5(p.Ctx, bm)
		if err2 != nil {
			err = err2
			return
		}
		if wxRsp.Code == wechat.Success {
			err = errors.New(fmt.Sprintf("wxRsp: %#v", wxRsp.Response))
			return
		}
		response = wxRsp.Response.H5Url
		break
	case "jsapi":
		wxRsp, err2 := client.V3TransactionJsapi(p.Ctx, bm)
		if err2 != nil {
			err = err2
			return
		}
		if wxRsp.Code == wechat.Success {
			err = errors.New(fmt.Sprintf("wxRsp: %#v", wxRsp.Response))
			return
		}
		applet, err2 := client.PaySignOfApplet("appid", "prepayid")
		if err2 != nil {
			err = err2
			return
		}
		data, err2 := json.Marshal(applet)
		if err2 != nil {
			err = err2
			return
		}
		response = string(data)
		break
	case "native":
		wxRsp, err2 := client.V3TransactionNative(p.Ctx, bm)
		if err2 != nil {
			err = err2
			return
		}
		if wxRsp.Code == wechat.Success {
			err = errors.New(fmt.Sprintf("wxRsp: %#v", wxRsp.Response))
			return
		}
		response = wxRsp.Response.CodeUrl
		break
	}

	return
}

func (p *WechatPay) PayNotify(req *http.Request) (payNotifyResponse PayNotifyResponse, err error) {
	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		return
	}
	client, err := p.getClient()
	// 获取微信平台证书
	certMap := client.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		return
	}
	payConfig, _ := p.getConfig()
	result, err := notifyReq.DecryptCipherText(payConfig.ApiV3Key)
	if err != nil {
		return
	}
	payNotifyResponse = PayNotifyResponse{
		OutTradeNo:    result.OutTradeNo,
		TransactionId: result.TransactionId,
		Attach:        result.Attach,
		SuccessTime:   result.SuccessTime,
		Amount:        PayAmoount{Total: result.Amount.Total},
	}
	return
}

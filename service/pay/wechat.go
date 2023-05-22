package pay

import (
	"chatgpt-tools/common/utils/appplatform"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"math"
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
	payConfig, err := p.getConfig()
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
		Set("description", order.Body).
		Set("out_trade_no", order.OutNo).
		Set("notify_url", order.NotifyPath).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", float32(math.Round(float64(order.Total*100))))
		})

	switch scene {
	case "h5":
		wxRsp, err2 := client.V3TransactionH5(p.Ctx, bm)
		if err2 != nil {
			err = err2
			return
		}
		if wxRsp.Code != wechat.Success {
			err = errors.New(wxRsp.Error)
			return
		}
		response = wxRsp.Response.H5Url
		break
	case "jsapi":
		bm.SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", order.OpenId)
		})
		wxRsp, err2 := client.V3TransactionJsapi(p.Ctx, bm)
		if err2 != nil {
			err = err2
			return
		}
		if wxRsp.Code != wechat.Success {
			err = errors.New(wxRsp.Error)
			return
		}
		applet, err2 := client.PaySignOfApplet(payConfig.AppId, wxRsp.Response.PrepayId)
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
		if wxRsp.Code != wechat.Success {
			err = errors.New(wxRsp.Error)
			return
		}
		response = wxRsp.Response.CodeUrl
		break
	case "app":
		wxRsp, err2 := client.V3TransactionApp(p.Ctx, bm)
		if err2 != nil {
			err = err2
			return
		}
		if wxRsp.Code != wechat.Success {
			err = errors.New(wxRsp.Error)
			return
		}
		applet, err2 := client.PaySignOfApp(payConfig.AppId, wxRsp.Response.PrepayId)
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
	}

	return
}

func (p *WechatPay) PayNotify(req *http.Request) (payNotifyResponse PayNotifyResponse, err error) {
	notifyReq, err := wechat.V3ParseNotify(req)
	if err != nil {
		return
	}
	fmt.Println(notifyReq.SignInfo.SignBody)
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
		Amount:        PayAmoount{Total: float32(result.Amount.Total) / 100},
	}
	return
}

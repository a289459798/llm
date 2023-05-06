package pay

import (
	"chatgpt-tools/common/utils/appplatform"
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"net/http"
)

type AlipayPay struct {
	Ctx      context.Context
	Config   string
	Merchant string
	DeBug    bool
}

func NewAlipay(payData PayData) *AlipayPay {
	return &AlipayPay{
		Ctx:      payData.Ctx,
		Config:   payData.Config,
		Merchant: payData.Merchant,
		DeBug:    payData.DeBug,
	}
}

func (p *AlipayPay) getConfig() (appplatform.AlipayConf, error) {
	return appplatform.GetConf[appplatform.AlipayConf](p.Config)
}

func (p *AlipayPay) getClient() (client *alipay.Client, err error) {
	payConfig, err := p.getConfig()
	client, err = alipay.NewClient(payConfig.AppId, payConfig.PrivateKey, false)
	if err != nil {
		return
	}
	client.DebugSwitch = func() gopay.DebugSwitch {
		if p.DeBug {
			return gopay.DebugOn
		}
		return gopay.DebugOff
	}()
	return
}

func (p *AlipayPay) Pay(scene string, order Order) (response string, err error) {
	payConfig, err := p.getConfig()
	client, err := p.getClient()
	if err != nil {
		return
	}

	// 设置支付宝请求 公共参数
	//    注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).  // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2). // 设置签名类型，不设置默认 RSA2
							SetNotifyUrl(order.NotifyPath)

	// 自动同步验签（只支持证书模式）
	// 传入 alipayCertPublicKey_RSA2.crt 内容
	client.AutoVerifySign([]byte(payConfig.AlipayPublicKey))

	// 证书内容
	err = client.SetCertSnByContent([]byte(payConfig.AppPublicKey), []byte(payConfig.RootKey), []byte(payConfig.AlipayPublicKey))
	if err != nil {
		return
	}

	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("subject", order.Body).
		Set("out_trade_no", order.OutNo).
		Set("total_amount", order.Total)

	response, err = client.TradeAppPay(p.Ctx, bm)
	if err != nil {
		return
	}

	return
}

func (p *AlipayPay) PayNotify(req *http.Request) (payNotifyResponse PayNotifyResponse, err error) {

	// 解析异步通知的参数
	notifyReq, err := alipay.ParseNotifyToBodyMap(req)
	if err != nil {
		return
	}

	payConfig, err := p.getConfig()
	// 支付宝异步通知验签（公钥证书模式）
	_, err = alipay.VerifySignWithCert(payConfig.AlipayPublicKey, notifyReq)
	if err != nil {
		return
	}

	fmt.Println(notifyReq)

	payNotifyResponse = PayNotifyResponse{
		OutTradeNo:    notifyReq.GetString("out_trade_no"),
		TransactionId: notifyReq.GetString("trade_no"),
		SuccessTime:   notifyReq.GetString("notify_time"),
		Amount:        PayAmoount{Total: notifyReq.GetInterface("total_amount").(float32)},
	}

	return
}

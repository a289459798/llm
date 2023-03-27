package appplatform

import (
	"context"
	"encoding/base64"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
)

type WechatMini struct {
	Ctx  context.Context
	Conf WechatMiniConf
}

func NewWechatMini(data AppData, conf WechatMiniConf) *WechatMini {
	return &WechatMini{
		Ctx:  data.Ctx,
		Conf: conf,
	}
}

func (a *WechatMini) GetSession(req SessionReq) (session Session, err error) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     a.Conf.AppId,
		AppSecret: a.Conf.AppSecret,
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)
	auth := mini.GetAuth()
	sess, err := auth.Code2Session(req.Code)
	if err != nil {
		return
	}
	session.OpenID = sess.OpenID
	session.UnionID = sess.UnionID
	return
}

func (a *WechatMini) GetQRcode(req QrcodeReq) (qrc string, err error) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     a.Conf.AppId,
		AppSecret: a.Conf.AppSecret,
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)
	qr, err := mini.GetQRCode().GetWXACodeUnlimit(qrcode.QRCoder{
		Width: 100,
		Path:  req.Path,
		Scene: req.Scene,
		EnvVersion: func() string {
			return "release"
		}(),
	})
	if err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(qr), nil
}

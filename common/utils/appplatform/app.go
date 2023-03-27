package appplatform

import (
	"context"
	"encoding/json"
)

type SanmuApp interface {
	GetSession(req SessionReq) (session Session, err error)
	GetQRcode(req QrcodeReq) (qrc string, err error)
}

type AppData struct {
	Ctx  context.Context
	Conf string
}

func GetApp(platform string, payData AppData) (SanmuApp, error) {
	conf, err := GetConf[WechatMiniConf](payData.Conf)
	if err != nil {
		return nil, err
	}
	return NewWechatMini(payData, conf), nil
}

func GetConf[T ConfData](conf string) (T, error) {
	var err error
	var val T
	if err != nil {
		return val, err
	}
	if err := json.Unmarshal([]byte(conf), &val); err != nil {
		return val, err
	}
	return val, nil
}

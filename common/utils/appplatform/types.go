package appplatform

type SessionReq struct {
	Code string
}
type QrcodeReq struct {
	Path  string
	Scene string
}

type Session struct {
	OpenID  string
	UnionID string
}

type ConfData interface {
	WechatMiniConf
}

type WechatMiniConf struct {
	AppId      string
	AppSecret  string
	MchId      string
	SerialNo   string
	ApiV3Key   string
	PrivateKey string
	NotifyUrl  string
}

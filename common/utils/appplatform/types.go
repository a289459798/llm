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
	WechatMiniConf | WechatOfficialConf | AlipayConf
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

type WechatOfficialConf struct {
	AppId          string
	AppSecret      string
	Token          string
	EncodingAESKey string
}

type AlipayConf struct {
	AppId           string
	PrivateKey      string
	AppPublicKey    string
	AlipayPublicKey string
	RootKey         string
}

package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type WeChatPayConf struct {
	AppId      string
	MchId      string
	SerialNo   string
	ApiV3Key   string
	PrivateKey string
	NotifyUrl  string
}

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	OpenAIKey string
	Auth      struct {
		AccessSecret string
		AccessExpire int64
	}
	MiniApp struct {
		AppId     string
		AppSecret string
	}
	WeChatPayConf struct {
		Default WeChatPayConf
	}
	OfficialAccount struct {
		AppId     string
		AppSecret string
		Token     string
	}
	Qiniu struct {
		Domain string
		Bucket string
		Ak     string
		SK     string
	}
}

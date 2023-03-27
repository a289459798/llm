package config

import (
	"github.com/zeromicro/go-zero/rest"
)

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

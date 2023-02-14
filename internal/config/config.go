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
}

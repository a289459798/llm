package wechat

import (
	"context"
	"errors"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateLogic {
	return &ValidateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidateLogic) Validate(req types.WechatValidateRequest, w http.ResponseWriter, r *http.Request) (resp string, err error) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	var cfg = &offConfig.Config{
		AppID:          "wx774aefe5b682fe6f",
		AppSecret:      "164c9a3e9090983ef35668293df86045",
		Token:          "2CA4jdhM",
		EncodingAESKey: "9mV4fqzpiaZ53UzhPZuzd4wsCw42v8neGRksLpVyUXy",
		Cache:          memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	// 传入request和responseWriter
	validate := officialAccount.GetServer(r, w).Validate()

	if !validate {
		return "", errors.New("error")
	}
	return req.Echostr, nil
}

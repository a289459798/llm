package wechat

import (
	"chatgpt-tools/common/utils/appplatform"
	"chatgpt-tools/model"
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
	appInfo := model.App{AppKey: req.AppKey}.Info(l.svcCtx.Db)
	if appInfo.ID == 0 {
		return "", errors.New("App-Key 错误")
	}
	c, _ := appplatform.GetConf[appplatform.WechatOfficialConf](appInfo.Conf)
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	var cfg = &offConfig.Config{
		AppID:          c.AppId,
		AppSecret:      c.AppSecret,
		Token:          c.Token,
		EncodingAESKey: c.EncodingAESKey,
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

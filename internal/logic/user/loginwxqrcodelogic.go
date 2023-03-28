package user

import (
	"chatgpt-tools/common/utils/appplatform"
	"chatgpt-tools/model"
	"context"
	"errors"
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginWXQrcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginWXQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginWXQrcodeLogic {
	return &LoginWXQrcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginWXQrcodeLogic) LoginWXQrcode(r *http.Request) (resp *types.LoginWXQrcodeResponse, err error) {
	appKey := r.Header.Get("App-Key")
	if appKey == "" {
		return nil, errors.New("App-Key 错误")
	}
	appInfo := model.App{AppKey: appKey}.Info(l.svcCtx.Db)
	if appInfo.ID == 0 {
		return nil, errors.New("App-Key 错误")
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
	tq := &basic.Request{}
	tq.ExpireSeconds = 604800
	tq.ActionName = "QR_STR_SCENE"
	tq.ActionInfo.Scene.SceneStr = "login"
	ticket, err := officialAccount.GetBasic().GetQRTicket(tq)
	if err != nil {
		return nil, err
	}
	return &types.LoginWXQrcodeResponse{
		Qrcode: fmt.Sprintf(basic.ShowQRCode(ticket)),
	}, nil
}

package user

import (
	"chatgpt-tools/common/utils/appplatform"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest, r *http.Request) (resp *types.InfoResponse, err error) {
	appKey := r.Header.Get("App-Key")
	if appKey == "" {
		return nil, errors.New("App-Key 错误")
	}
	appInfo := model.App{AppKey: appKey}.Info(l.svcCtx.Db)
	if appInfo.ID == 0 {
		return nil, errors.New("App-Key 错误")
	}

	app, err := appplatform.GetApp(appInfo.Platform, appplatform.AppData{
		Ctx:  l.ctx,
		Conf: appInfo.Conf,
	})
	if err != nil {
		return nil, err
	}
	session, err := app.GetSession(appplatform.SessionReq{
		Code: req.Code,
	})
	aiUser, tokenString, err := model.AIUser{}.Login(l.svcCtx.Db, model.UserLogin{
		OpenID:       session.OpenID,
		UnionID:      session.UnionID,
		Channel:      req.Channel,
		AppKey:       appKey,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
	})

	if err != nil {
		return nil, err
	}

	return &types.InfoResponse{
		Token:  tokenString,
		Uid:    aiUser.Uid,
		OpenId: aiUser.OpenId,
	}, nil
}

package user

import (
	"chatgpt-tools/model"
	"context"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginWithAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginWithAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginWithAppLogic {
	return &LoginWithAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginWithAppLogic) LoginWithApp(req *types.LoginAppRequest) (resp *types.InfoResponse, err error) {
	aiUser, tokenString, err := model.AIUser{}.Login(l.svcCtx.Db, model.UserLogin{
		OpenID:       req.OpenId,
		UnionID:      req.UnionID,
		Channel:      req.Channel,
		AppKey:       "app",
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

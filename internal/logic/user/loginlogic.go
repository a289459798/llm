package user

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.InfoResponse, err error) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     l.svcCtx.Config.MiniApp.AppId,
		AppSecret: l.svcCtx.Config.MiniApp.AppSecret,
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)
	auth := mini.GetAuth()
	session, err := auth.Code2Session(req.Code)
	if err != nil {
		return nil, err
	}
	aiUser := &model.AIUser{}
	l.svcCtx.Db.Where("open_id = ?", session.OpenID).First(aiUser)
	if aiUser.ID == 0 {
		tx := l.svcCtx.Db.Begin()
		// 创建用户
		user := &model.User{}
		tx.Create(user)
		aiUser.OpenId = session.OpenID
		aiUser.UnionId = session.UnionID
		aiUser.Platform = req.Platform
		aiUser.Channel = req.Channel
		aiUser.Uid = user.ID
		err = tx.Create(&aiUser).Error
		if err != nil {
			tx.Rollback()
			return nil, errors.New("错误")
		}

		tx.Commit()
	}

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Unix() + l.svcCtx.Config.Auth.AccessExpire
	claims["iat"] = time.Now().Unix()
	claims["uid"] = aiUser.ID
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))

	if err != nil {
		return nil, err
	}

	return &types.InfoResponse{
		Token:  tokenString,
		Uid:    aiUser.ID,
		OpenId: aiUser.OpenId,
	}, nil
}

package user

import (
	"chatgpt-tools/common/utils/appplatform"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
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
	aiUser := &model.AIUser{}
	l.svcCtx.Db.Where("open_id = ?", session.OpenID).First(aiUser)
	if aiUser.Uid == 0 {
		// 判断UnionID是否存在
		l.svcCtx.Db.Where("union_id = ?", session.UnionID).First(aiUser)
		if aiUser.Uid == 0 {
			tx := l.svcCtx.Db.Begin()
			// 创建用户
			user := &model.User{}
			tx.Create(user)
			aiUser.OpenId = session.OpenID
			aiUser.UnionId = session.UnionID
			aiUser.AppKey = appKey
			aiUser.Channel = req.Channel
			aiUser.Uid = user.ID
			err = tx.Create(&aiUser).Error
			if err != nil {
				tx.Rollback()
				return nil, errors.New("错误")
			}

			tx.Commit()
		} else {
			newUser := model.AIUser{}
			newUser.OpenId = session.OpenID
			newUser.UnionId = session.UnionID
			newUser.AppKey = appKey
			newUser.Channel = req.Channel
			newUser.Uid = aiUser.Uid
			err = l.svcCtx.Db.Create(&aiUser).Error
			if err != nil {
				return nil, errors.New("错误")
			}
		}
	} else if aiUser.UnionId == "" {
		aiUser.UnionId = session.UnionID
		l.svcCtx.Db.Save(aiUser)
	}

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Unix() + l.svcCtx.Config.Auth.AccessExpire
	claims["iat"] = time.Now().Unix()
	claims["uid"] = aiUser.Uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))

	if err != nil {
		return nil, err
	}

	return &types.InfoResponse{
		Token:  tokenString,
		Uid:    aiUser.Uid,
		OpenId: aiUser.OpenId,
	}, nil
}

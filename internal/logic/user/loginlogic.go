package user

import (
	"chatgpt-tools/model"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

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
		AppID:     "xxx",
		AppSecret: "xxx",
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)
	auth := mini.GetAuth()
	session, err := auth.Code2Session(req.Code)
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	l.svcCtx.Db.Where("open_id = ?", session.OpenID).First(user)
	if user == nil {
		user.OpenId = session.OpenID
		user.UnionId = session.UnionID
		user.Amount = 20
		l.svcCtx.Db.Create(&user)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Unix() + 86400*365,
		"nbf":    time.Now().Unix(),
		"uid":    user.ID,
		"openid": user.OpenId,
	})

	tokenString, err := token.SignedString([]byte(l.svcCtx.Config.JwtSecret))
	if err != nil {
		return nil, err
	}

	return &types.InfoResponse{
		Amount: user.Amount,
		Token:  tokenString,
	}, nil
}

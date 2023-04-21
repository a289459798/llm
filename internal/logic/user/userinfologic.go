package user

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.InfoRequest) (resp *types.InfoResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := &model.AIUser{}
	l.svcCtx.Db.Where("uid = ?", uid).Preload("Vip").Preload("Vip.Vip").Preload("Distributor").First(user)

	if user.Uid == 0 || user.UnionId == "" {
		return nil, errors.New(string(http.StatusUnauthorized))
	}
	amount := model.NewAccount(l.svcCtx.Db).GetAccount(uint32(uid), time.Now())
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Unix() + l.svcCtx.Config.Auth.AccessExpire
	claims["iat"] = time.Now().Unix()
	claims["uid"] = user.Uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))

	return &types.InfoResponse{
		Amount:    amount.Amount,
		Uid:       uint32(uid),
		OpenId:    user.OpenId,
		Vip:       user.IsVip(),
		Code:      fmt.Sprintf("%b", uid),
		Token:     tokenString,
		Group:     user.IsJoinGroup(l.svcCtx.Db),
		VipName:   user.Vip.Vip.Name,
		VipExpiry: user.Vip.VipExpiry.Format("2006-01-02"),
		IsPartner: user.Distributor.ID > 0 && user.Distributor.Status,
	}, nil
}

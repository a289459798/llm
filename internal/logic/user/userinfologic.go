package user

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
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
	amount := model.NewAccount(l.svcCtx.Db).GetAccount(uint32(uid), time.Now())
	return &types.InfoResponse{
		Amount: amount.ChatAmount - amount.ChatUse,
		Uid:    uint32(uid),
	}, nil
}

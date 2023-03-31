package notify

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnreadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnreadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnreadLogic {
	return &UnreadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnreadLogic) Unread() (resp *types.UserNotifyUnreadResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	var count int64
	l.svcCtx.Db.Model(&model.AINotify{}).Where("uid = ?", uid).Where("status = 0").Count(&count)

	return &types.UserNotifyUnreadResponse{
		Status: count > 0,
	}, nil
}

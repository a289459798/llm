package common

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidChatLogic {
	return &ValidChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidChatLogic) ValidChat(req *types.ValidRequest) (resp *types.ValidResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	amount := model.NewAccount(l.svcCtx.Db).GetAccount(uint32(uid), time.Now())
	if amount.ChatAmount <= amount.ChatUse {
		return nil, errors.New("次数已用完")
	}
	return &types.ValidResponse{
		Data: string(amount.ChatAmount - amount.ChatUse),
	}, nil
}

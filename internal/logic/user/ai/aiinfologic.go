package ai

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiInfoLogic {
	return &AiInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiInfoLogic) AiInfo(req *types.InfoRequest) (resp *types.AIInfoResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	ai := &model.AI{}
	l.svcCtx.Db.Where("uid = ?", uid).Preload("Role").Find(&ai)
	if ai.ID == 0 {
		return nil, nil
	}
	return &types.AIInfoResponse{
		Name:    ai.Name,
		Photo:   ai.Image,
		Welcome: ai.Role.Welcome,
		Call:    ai.Call,
	}, nil
}

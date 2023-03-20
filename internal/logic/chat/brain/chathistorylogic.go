package brain

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHistoryLogic {
	return &ChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatHistoryLogic) ChatHistory(req types.ChatHistoryRequest) (resp *types.ChatHistoryResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	records := []model.Record{}
	l.svcCtx.Db.Where("uid = ?", uid).
		Where("chat_id = ?", req.ChatId).
		Order("id desc").
		Find(&records)
	history := []types.ChatHistory{}
	model := "gpt3.5"
	for _, m := range records {
		if m.Model != "" {
			model = m.Model
		}
		history = append([]types.ChatHistory{
			{
				Q: m.Content,
				A: m.Result,
			},
		}, history...)
	}

	return &types.ChatHistoryResponse{ChatId: req.ChatId, History: history, Model: model}, nil
}

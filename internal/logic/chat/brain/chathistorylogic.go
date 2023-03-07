package brain

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

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

func (l *ChatHistoryLogic) ChatHistory() (resp *types.ChatHistoryResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	record := &model.Record{}
	today := time.Now().Format("2006-01-02")
	l.svcCtx.Db.Where("uid = ?", uid).
		Order("id desc").
		Where("type = ?", "chat/chat").
		Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").
		Select("chat_id").Find(&record)

	chatId := ""
	history := []types.ChatHistory{}
	if record.ChatId != "" {
		chatId = record.ChatId
		records := []model.Record{}
		l.svcCtx.Db.Where("uid = ?", uid).
			Where("chat_id = ?", chatId).
			Order("id desc").
			Limit(3).
			Find(&records)
		for _, m := range records {
			history = append([]types.ChatHistory{
				{
					Q: m.Content,
					A: m.Result,
				},
			}, history...)
		}
	}

	return &types.ChatHistoryResponse{ChatId: chatId, History: history}, nil
}

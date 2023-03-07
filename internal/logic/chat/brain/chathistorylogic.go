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
		Order("id desc").Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").
		Select("chat_id").Find(&record)

	chatId := ""
	if record.ChatId != "" {
		chatId = record.ChatId
	}

	return &types.ChatHistoryResponse{ChatId: chatId}, nil
}

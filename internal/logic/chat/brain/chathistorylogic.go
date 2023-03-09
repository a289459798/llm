package brain

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"fmt"
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
	today := time.Now().Format("2006-01-02")
	totalRecord := &struct {
		ChatId string
		Count  int
	}{}
	l.svcCtx.Db.Model(&model.Record{}).Where("uid = ?", uid).
		Order("id desc").
		Where("type = ?", "chat/chat").
		Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").
		Group("chat_id").
		Select("chat_id, count(*) as count").Find(&totalRecord)
	fmt.Println(totalRecord)

	chatId := ""
	history := []types.ChatHistory{}
	if totalRecord.Count >= 5 {
		chatId = totalRecord.ChatId
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

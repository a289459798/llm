package history

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"math"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatListLogic {
	return &ChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatListLogic) ChatList(req types.PageRequest) (resp *types.ChatHistoryListResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	offset := req.Offset
	limit := req.Limit
	tx := l.svcCtx.Db.Model(&model.Record{}).
		Where("uid = ?", uid).
		Where("type = ?", "chat/chat").
		Where("is_delete = 0").
		Where("title != ?", "").
		Group("chat_id")
	var total int64
	tx.Count(&total)
	maxPage := int(math.Ceil(float64(total) / float64(limit)))
	if maxPage < offset {
		return &types.ChatHistoryListResponse{}, nil
	}
	records := []model.Record{}
	tx.Order("id desc").Offset((offset - 1) * limit).Limit(limit).Select("min(id), title, chat_id, created_at").Find(&records)
	data := []types.ChatHistoryData{}
	for _, record := range records {
		data = append(data, types.ChatHistoryData{
			Q:      record.Title,
			ChatId: record.ChatId,
			Time:   record.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	return &types.ChatHistoryListResponse{
		Pagination: types.Pagination{
			Total:  maxPage,
			Limit:  limit,
			Offset: offset,
		},
		Data: data,
	}, nil
}

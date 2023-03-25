package history

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
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
	isVip := model.AIUser{ID: uint32(uid)}.Find(l.svcCtx.Db).Vip.IsVip()
	if !isVip {
		return nil, errors.New("非VIP")
	}

	offset := req.Offset
	limit := req.Limit
	tx := l.svcCtx.Db.Model(&model.Record{}).
		Where("uid = ?", uid).
		Where("type = ?", "chat/chat").
		Group("chat_id")
	var total int64
	tx.Count(&total)
	maxPage := int(math.Ceil(float64(total) / float64(limit)))
	if maxPage < offset {
		return &types.ChatHistoryListResponse{}, nil
	}
	records := []model.Record{}
	tx.Order("id desc").Offset((offset - 1) * limit).Limit(limit).Select("min(id), chat_id, content,created_at").Find(&records)
	data := []types.ChatHistoryData{}
	for _, record := range records {
		data = append(data, types.ChatHistoryData{
			Q:      record.Content,
			ChatId: record.ChatId,
			Time:   record.CreatedAt.Format("01-02 15:04"),
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

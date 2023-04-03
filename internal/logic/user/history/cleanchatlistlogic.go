package history

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CleanChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCleanChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanChatListLogic {
	return &CleanChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CleanChatListLogic) CleanChatList() (resp *types.ChatHistoryListResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	err = l.svcCtx.Db.Model(&model.Record{}).Where("uid = ?", uid).Where("type = ?", "chat/chat").Update("is_delete", 1).Error

	if err != nil {
		return nil, err
	}
	return &types.ChatHistoryListResponse{}, nil
}

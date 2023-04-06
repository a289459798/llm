package brain

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatDeleteLogic {
	return &ChatDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatDeleteLogic) ChatDelete(req *types.ChatHistoryRequest) (resp *types.ChatHistoryResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	l.svcCtx.Db.Model(&model.Record{}).
		Where("uid = ?", uid).
		Where("chat_id = ?", req.ChatId).
		Update("is_delete", 1)
	return
}

package common

import (
	"chatgpt-tools/model"
	"context"
	"errors"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageLogic {
	return &MessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageLogic) Message() (resp *types.MessageResponse, err error) {
	message := &model.Message{}
	l.svcCtx.Db.Where("status = 1").Find(&message)
	if message.ID == 0 {
		return nil, errors.New("数据为空")
	}

	return &types.MessageResponse{
		Id:      message.ID,
		Title:   message.Title,
		Content: message.Content,
		Link:    message.Link,
	}, nil
}

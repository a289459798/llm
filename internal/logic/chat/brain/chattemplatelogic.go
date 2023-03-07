package brain

import (
	"chatgpt-tools/model"
	"context"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatTemplateLogic {
	return &ChatTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatTemplateLogic) ChatTemplate() (resp *types.ChatTemplateResponse, err error) {

	template := []model.ChatTemplate{}
	l.svcCtx.Db.Where("is_del = 0").Order("sort desc").Select("id, title").Find(&template)

	list := []types.ChatTemplate{}
	for _, chatTemplate := range template {
		list = append(list, types.ChatTemplate{
			TemplateId: chatTemplate.ID,
			Message:    chatTemplate.Title,
		})
	}

	return &types.ChatTemplateResponse{List: list}, nil
}

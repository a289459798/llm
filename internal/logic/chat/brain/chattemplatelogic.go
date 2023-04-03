package brain

import (
	"chatgpt-tools/model"
	"context"
	"fmt"

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

func (l *ChatTemplateLogic) ChatTemplate(req types.ChatTemplateRequest) (resp *types.ChatTemplateResponse, err error) {

	template := []model.ChatTemplate{}
	query := "type is null"
	if req.Type != "" {
		query = fmt.Sprintf("type='%s'", req.Type)
	}
	l.svcCtx.Db.Where(query).Where("is_del = 0").Select("id, title").Order("rand()").Limit(6).Find(&template)

	list := []types.ChatTemplate{}
	for _, chatTemplate := range template {
		list = append(list, types.ChatTemplate{
			TemplateId: chatTemplate.ID,
			Message:    chatTemplate.Title,
		})
	}

	return &types.ChatTemplateResponse{List: list}, nil
}

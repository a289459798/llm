package ai

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiEditLogic {
	return &AiEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiEditLogic) AiEdit(req *types.AIEditRequest) (resp *types.AIInfoResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	if req.RoleId > 0 {
		template := &model.ChatTemplate{}
		l.svcCtx.Db.Where("id = ?", req.RoleId).Find(&template)
		if template.ID == 0 {
			return nil, errors.New("角色不存在")
		}
	}
	ai := &model.AI{}
	l.svcCtx.Db.Where("uid = ?", uid).Find(&ai)
	ai.Uid = uint32(uid)
	ai.Name = req.Name
	ai.Image = req.Photo
	ai.Call = req.Call
	ai.RoleId = req.RoleId
	if ai.ID == 0 {
		l.svcCtx.Db.Create(ai)
	} else {
		l.svcCtx.Db.Save(ai)
	}

	return
}

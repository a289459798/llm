package common

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidTextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidTextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidTextLogic {
	return &ValidTextLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidTextLogic) ValidText(req *types.ValidRequest) (resp *types.ValidResponse, err error) {
	valid := utils.Filter(req.Content)
	if valid != "" {
		return nil, errors.New(req.Content)
	}
	return
}

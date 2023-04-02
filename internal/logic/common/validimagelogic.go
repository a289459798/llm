package common

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type ValidImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidImageLogic {
	return &ValidImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidImageLogic) ValidImage(req *types.ValidImageRequest) (resp *types.ValidResponse, err error) {
	_, err = utils.ImageCheck(l.ctx, l.svcCtx.Config.Qiniu.Ak, l.svcCtx.Config.Qiniu.SK, req.Bucket, req.Key)
	if err != nil {
		return nil, err
	}
	return
}

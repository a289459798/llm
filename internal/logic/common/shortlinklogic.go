package common

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShortLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkLogic {
	return &ShortLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShortLinkLogic) ShortLink(req *types.ShortLinkRequest) (resp *types.QrCodeResponse, err error) {
	return &types.QrCodeResponse{Data: ""}, nil
}

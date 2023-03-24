package common

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadTokenLogic {
	return &UploadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadTokenLogic) UploadToken() (resp *types.UploadTokenResponse, err error) {
	putPolicy := storage.PutPolicy{
		Scope: l.svcCtx.Config.Qiniu.Bucket,
	}
	mac := qbox.NewMac(l.svcCtx.Config.Qiniu.Ak, l.svcCtx.Config.Qiniu.SK)
	upToken := putPolicy.UploadToken(mac)

	return &types.UploadTokenResponse{Token: upToken}, nil
}

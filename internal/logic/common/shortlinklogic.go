package common

import (
	"context"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/config"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

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
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     l.svcCtx.Config.MiniApp.AppId,
		AppSecret: l.svcCtx.Config.MiniApp.AppSecret,
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)
	link, err := mini.GetShortLink().GenerateShortLinkTemp(req.Page, req.Title)
	if err != nil {
		return nil, err
	}

	return &types.QrCodeResponse{Data: link}, nil
}

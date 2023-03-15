package common

import (
	"context"
	"encoding/base64"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QrcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrcodeLogic {
	return &QrcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QrcodeLogic) Qrcode(req types.QrCodeRequest) (resp *types.QrCodeResponse, err error) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     l.svcCtx.Config.MiniApp.AppId,
		AppSecret: l.svcCtx.Config.MiniApp.AppSecret,
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)
	qr, err := mini.GetQRCode().GetWXACodeUnlimit(qrcode.QRCoder{
		Width: 100,
		Path:  req.Path,
		Scene: req.Scene,
		EnvVersion: func() string {
			if l.svcCtx.Config.Mode == "dev" {
				return "trial"
			}
			return "release"
		}(),
	})
	if err != nil {
		return nil, err
	}
	return &types.QrCodeResponse{
		Data: base64.StdEncoding.EncodeToString(qr),
	}, nil
}

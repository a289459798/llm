package common

import (
	"chatgpt-tools/common/utils/appplatform"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"errors"
	"net/http"

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

func (l *QrcodeLogic) Qrcode(req types.QrCodeRequest, r *http.Request) (resp *types.QrCodeResponse, err error) {
	appKey := r.Header.Get("App-Key")
	if appKey == "" {
		return nil, errors.New("App-Key 错误")
	}
	appInfo := model.App{AppKey: appKey}.Info(l.svcCtx.Db)
	if appInfo.ID == 0 {
		return nil, errors.New("App-Key 错误")
	}

	app, err := appplatform.GetApp(appInfo.Platform, appplatform.AppData{
		Ctx:  l.ctx,
		Conf: appInfo.Conf,
	})
	if err != nil {
		return nil, err
	}

	qr, err := app.GetQRcode(appplatform.QrcodeReq{
		Path:  req.Path,
		Scene: req.Scene,
	})
	if err != nil {
		return nil, err
	}
	return &types.QrCodeResponse{
		Data: qr,
	}, nil
}

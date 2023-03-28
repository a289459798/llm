package user

import (
	"chatgpt-tools/model"
	"context"
	"errors"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginCheckLogic {
	return &LoginCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginCheckLogic) LoginCheck(req *types.LoginCheckRequest) (resp *types.InfoResponse, err error) {
	data := &model.ScanScene{}
	l.svcCtx.Db.Where("scene = ?", req.SceneStr).First(&data)
	if data.ID == 0 {
		return nil, errors.New(string(http.StatusNotFound))
	}

	return &types.InfoResponse{
		Token: data.Data,
	}, nil
}
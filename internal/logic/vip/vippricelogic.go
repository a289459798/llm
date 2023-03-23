package vip

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type VipPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipPriceLogic {
	return &VipPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipPriceLogic) VipPrice() (resp *types.VipPriceResponse, err error) {
	vip := []model.Vip{}
	l.svcCtx.Db.Find(&vip)
	res := []types.VipDataResponse{}
	copier.Copy(&res, &vip)
	return &types.VipPriceResponse{
		Data: res,
	}, nil
}

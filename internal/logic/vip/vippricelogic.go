package vip

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
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
	setting := &model.Setting{Name: "vip"}
	vipsSetting, err := setting.Find(l.svcCtx.Db)
	if err != nil {
		return nil, err
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.User{ID: uint32(uid)}.Find(l.svcCtx.Db)
	if user.VipExpiry.Unix() > 0 {
		vipsSetting["first"] = vipsSetting["original"]
	}
	return &types.VipPriceResponse{
		Original: int(vipsSetting["original"].(float64)),
		Price:    int(vipsSetting["first"].(float64)),
	}, nil
}

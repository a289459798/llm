package vip

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

type VipGiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipGiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipGiveLogic {
	return &VipGiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipGiveLogic) VipGive() (resp *types.VipGiveResponse, err error) {
	return &types.VipGiveResponse{
		Day:    0,
		Expiry: "",
	}, nil
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db)
	day := 0
	expiry := ""
	if user.Uid > 0 && user.Vip.ID == 0 {
		// 赠送
		day = 1
		err = user.SetVip(l.svcCtx.Db, &model.VipCode{
			VipId: 1,
			Day:   1,
			Vip: model.Vip{
				Amount: 38,
				Day:    1,
			},
		})
		if err != nil {
			return nil, err
		}
	}
	return &types.VipGiveResponse{
		Day:    day,
		Expiry: expiry,
	}, nil
}

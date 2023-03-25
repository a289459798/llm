package vip

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type VipCxchangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipCxchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipCxchangeLogic {
	return &VipCxchangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipCxchangeLogic) VipCxchange(req *types.VipCxchangeRequest) (resp *types.VipCxchangeResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	vipCode := &model.VipCode{}
	l.svcCtx.Db.Where("uid = ?", uid).Where("code = ?", req.Code).Where("status != 1").Preload("Vip").First(&vipCode)
	if vipCode.ID == 0 {
		return nil, errors.New("兑换码错误")
	}

	// 判断会员等级
	user := model.AIUser{ID: uint32(uid)}.Find(l.svcCtx.Db)
	if user.IsVip() && user.Vip.VipId != vipCode.VipId {
		return nil, errors.New("当前会员与购买会员不符，请联系客服")
	}

	tx := l.svcCtx.Db.Begin()
	tx.SavePoint("start")
	// 修改兑换码状态
	vipCode.Status = true
	err = tx.Save(&vipCode).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = model.AIUser{ID: uint32(uid)}.Find(l.svcCtx.Db).SetVip(tx, vipCode)
	if err != nil {
		tx.RollbackTo("start")
		return nil, err
	}
	tx.Commit()

	return &types.VipCxchangeResponse{}, nil
}

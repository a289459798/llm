package vip

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

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
	l.svcCtx.Db.Where("uid = ?", uid).Where("code = ?", req.Code).Where("status != 1").First(&vipCode)
	if vipCode.ID == 0 {
		return nil, errors.New("兑换码错误")
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
	err = model.User{ID: uint32(uid)}.Find(l.svcCtx.Db).SetVip(tx)
	if err != nil {
		tx.RollbackTo("start")
		return nil, err
	}
	tx.Commit()

	return &types.VipCxchangeResponse{}, nil
}

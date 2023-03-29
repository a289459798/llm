package vip

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"time"

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
		userVip := &model.AIUserVip{}
		userVip.Uid = user.Uid
		userVip.VipExpiry, err = time.ParseInLocation("2006-01-02 15:04:05", time.Now().AddDate(0, 0, 1).Format("2006-01-02")+" 23:59:59", time.Local)
		if err != nil {
			return nil, err
		}
		day = 1
		expiry = userVip.VipExpiry.Format("2006-01-02")
		userVip.VipId = 1
		err = l.svcCtx.Db.Create(userVip).Error
		if err != nil {
			return nil, err
		}
	}
	return &types.VipGiveResponse{
		Day:    day,
		Expiry: expiry,
	}, nil
}

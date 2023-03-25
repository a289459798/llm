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
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.AIUser{ID: uint32(uid)}.Find(l.svcCtx.Db)
	day := 0
	expiry := ""
	if user.ID > 0 && user.Vip.VipExpiry.Unix() < 0 {
		// 赠送
		user.Vip.VipExpiry, err = time.ParseInLocation("2006-01-02 15:04:05", time.Now().AddDate(0, 0, 1).Format("2006-01-02")+" 23:59:59", time.Local)
		if err != nil {
			return nil, err
		}
		day = 1
		expiry = user.Vip.VipExpiry.Format("2006-01-02")
		user.Vip.VipId = 1
		if user.Vip.Uid > 0 {
			l.svcCtx.Db.Save(user.Vip)
		} else {
			l.svcCtx.Db.Create(user.Vip)
		}
	}
	return &types.VipGiveResponse{
		Day:    day,
		Expiry: expiry,
	}, nil
}

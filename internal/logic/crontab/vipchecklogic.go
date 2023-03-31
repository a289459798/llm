package crontab

import (
	"chatgpt-tools/model"
	"context"
	"fmt"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VipCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipCheckLogic {
	return &VipCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipCheckLogic) VipCheck() (resp *types.CronResponse, err error) {
	vipList := []model.AIUserVip{}
	l.svcCtx.Db.Where("vip_expiry between ? and ?", time.Now().Format("2006-01-02 15:04:05"), time.Now().AddDate(0, 0, 1).Format("2006-01-02 15:04:05")).Find(&vipList)
	for _, vip := range vipList {
		// 发送站内信
		l.svcCtx.Db.Create(&model.AINotify{
			Uid:     vip.Uid,
			Title:   "会员即将过期提醒",
			Content: fmt.Sprintf("您的会员将于%s过期，为了不影响您使用，请及时续费", vip.VipExpiry.Format("2006-01-02 15:04:05")),
			Link:    "",
			Status:  false,
		})

	}
	return
}

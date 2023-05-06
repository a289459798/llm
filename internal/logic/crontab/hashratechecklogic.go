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

type HashRateCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHashRateCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HashRateCheckLogic {
	return &HashRateCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HashRateCheckLogic) HashRateCheck() (resp *types.CronResponse, err error) {
	hashRateList := []model.AIUserHashRate{}
	l.svcCtx.Db.Where("expiry between ? and ?", time.Now().Format("2006-01-02 15:04:05"), time.Now().AddDate(0, 0, 1).Format("2006-01-02 15:04:05")).
		Where("amount > use_amount").Find(&hashRateList)
	for _, rate := range hashRateList {
		// 发送站内信
		l.svcCtx.Db.Create(&model.AINotify{
			Uid:     rate.Uid,
			Title:   "算力即将过期提醒",
			Content: fmt.Sprintf("您有部分算力将于%s过期，请及时使用", rate.Expiry.Format("2006-01-02 15:04:05")),
			Link:    "",
			Status:  false,
		})

	}

	return
}

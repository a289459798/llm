package distributor

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DistributorInfoRequest) (resp *types.DistributorInfoResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	distributor := &model.Distributor{}
	l.svcCtx.Db.Where("uid = ?", uid).Where("status = 1").Preload("Level").First(&distributor)
	if distributor.ID == 0 {
		return nil, errors.New("未开通")
	}

	var next types.DistributorInfoNext

	// 获取下一等级
	nextLevel := &model.DistributorLevel{}
	l.svcCtx.Db.Where("id > ?", distributor.LevelId).Order("id asc").First(&nextLevel)
	if nextLevel.ID > 0 {
		next.Name = nextLevel.Name
		next.Ratio = nextLevel.Ratio
		next.Price = nextLevel.UserPrice
		next.User = nextLevel.UserNumber
	}

	monthDate := []string{fmt.Sprintf("%s-01 00:00:00", time.Now().Format("2006-01")), time.Now().Format("2006-01-02 15:04:04")}
	today := []string{fmt.Sprintf("%s 00:00:00", time.Now().Format("2006-01-02")), time.Now().Format("2006-01-02 15:04:04")}
	statistics := types.DistributorStatisticsResponse{
		UserTotal:  model.DistributorRecord{}.TotalWithDate(l.svcCtx.Db, uint32(uid), nil),
		UserMonth:  model.DistributorRecord{}.TotalWithDate(l.svcCtx.Db, uint32(uid), monthDate),
		UserDay:    model.DistributorRecord{}.TotalWithDate(l.svcCtx.Db, uint32(uid), today),
		PayTotal:   model.DistributorPayRecord{}.TotalPayWithDate(l.svcCtx.Db, uint32(uid), nil),
		MoneyTotal: model.DistributorPayRecord{}.TotalMoneyWithDate(l.svcCtx.Db, uint32(uid), nil),
		MoneyMonth: model.DistributorPayRecord{}.TotalMoneyWithDate(l.svcCtx.Db, uint32(uid), monthDate),
		MoneyDay:   model.DistributorPayRecord{}.TotalMoneyWithDate(l.svcCtx.Db, uint32(uid), today),
	}

	return &types.DistributorInfoResponse{
		Level:      distributor.Level.Name,
		Ratio:      distributor.Ratio,
		Link:       fmt.Sprintf("https://chat.smuai.com/c=%d", uid),
		Money:      distributor.Money,
		Statistics: statistics,
		Next:       next,
	}, nil
}

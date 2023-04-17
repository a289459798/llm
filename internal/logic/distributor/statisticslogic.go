package distributor

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatisticsLogic {
	return &StatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatisticsLogic) Statistics(req *types.DistributorStatisticsRequest) (resp *types.DistributorStatisticsResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	monthDate := []string{fmt.Sprintf("%s-01 00:00:00", time.Now().Format("2006-01")), time.Now().Format("2006-01-02 15:04:04")}
	today := []string{fmt.Sprintf("%s 00:00:00", time.Now().Format("2006-01-02")), time.Now().Format("2006-01-02 15:04:04")}
	return &types.DistributorStatisticsResponse{
		UserTotal:  model.DistributorRecord{}.TotalWithDate(l.svcCtx.Db, uint32(uid), nil),
		UserMonth:  model.DistributorRecord{}.TotalWithDate(l.svcCtx.Db, uint32(uid), monthDate),
		UserDay:    model.DistributorRecord{}.TotalWithDate(l.svcCtx.Db, uint32(uid), today),
		PayTotal:   model.DistributorPayRecord{}.TotalPayWithDate(l.svcCtx.Db, uint32(uid), nil),
		MoneyTotal: model.DistributorPayRecord{}.TotalMoneyWithDate(l.svcCtx.Db, uint32(uid), nil),
		MoneyMonth: model.DistributorPayRecord{}.TotalMoneyWithDate(l.svcCtx.Db, uint32(uid), monthDate),
		MoneyDay:   model.DistributorPayRecord{}.TotalMoneyWithDate(l.svcCtx.Db, uint32(uid), today),
	}, nil
}

package history

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SuanliListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSuanliListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SuanliListLogic {
	return &SuanliListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SuanliListLogic) SuanliList(req *types.PageRequest) (resp *types.SuanliHistoryListResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	offset := req.Offset
	limit := req.Limit
	today := time.Now().Format("2006-01-02")
	tx := l.svcCtx.Db.Model(&model.AccountRecord{}).
		Where("uid = ?", uid).
		Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59")
	var total int64
	tx.Count(&total)
	maxPage := int(math.Ceil(float64(total) / float64(limit)))
	if maxPage < offset {
		return &types.SuanliHistoryListResponse{}, nil
	}
	records := []model.AccountRecord{}
	tx.Order("id desc").Offset((offset - 1) * limit).Limit(limit).Find(&records)
	data := []types.SuanliHistoryData{}
	for _, record := range records {
		data = append(data, types.SuanliHistoryData{
			Amount: int(record.Amount),
			Desc:   record.GetType(),
			Time:   record.CreatedAt.Format("15:04:05"),
			Way:    record.Way,
			Type:   record.Type,
		})
	}

	return &types.SuanliHistoryListResponse{
		Pagination: types.Pagination{
			Total:  maxPage,
			Limit:  limit,
			Offset: offset,
		},
		Data: data,
	}, nil
}

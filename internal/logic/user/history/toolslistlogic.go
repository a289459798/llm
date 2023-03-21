package history

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

type ToolsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewToolsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToolsListLogic {
	return &ToolsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToolsListLogic) ToolsList() (resp *types.ToolsHistoryListResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	records := []model.Record{}
	l.svcCtx.Db.Where("uid = ?", uid).
		Group("type").
		Order("id desc").
		Select("type").
		Find(&records)

	data := []types.ToolsHistoryData{}
	for _, record := range records {
		data = append(data, types.ToolsHistoryData{
			Key: record.Type,
		})
	}

	return &types.ToolsHistoryListResponse{
		Data: data,
	}, nil
}

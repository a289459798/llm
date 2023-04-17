package distributor

import (
	"chatgpt-tools/model"
	"context"
	"github.com/jinzhu/copier"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LevelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LevelLogic {
	return &LevelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LevelLogic) Level(req *types.DistributorLevelRequest) (resp *types.DistributorLevelResponse, err error) {
	var level []model.DistributorLevel
	l.svcCtx.Db.Order("id asc").Find(&level)

	var res []types.DistributorLevel
	copier.Copy(&res, &level)
	return &types.DistributorLevelResponse{
		Data: res,
	}, nil
}

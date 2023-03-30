package hashrate

import (
	"chatgpt-tools/model"
	"context"
	"github.com/jinzhu/copier"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HashRatePriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHashRatePriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HashRatePriceLogic {
	return &HashRatePriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HashRatePriceLogic) HashRatePrice() (resp *types.HashRatePriceResponse, err error) {
	hashRate := []model.AIHashRate{}
	l.svcCtx.Db.Find(&hashRate)
	res := []types.HashRateResponse{}
	copier.Copy(&res, &hashRate)
	return &types.HashRatePriceResponse{
		Data: res,
	}, nil
}

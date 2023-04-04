package hashrate

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
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

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := &model.AIUser{}
	l.svcCtx.Db.Where("uid = ?", uid).Preload("Vip").Preload("Vip.Vip").First(user)

	if user.IsVip() {
		for i, re := range res {
			res[i].VipPrice = re.Price * user.Vip.Vip.Discount / 10
		}
	}

	return &types.HashRatePriceResponse{
		Data: res,
	}, nil
}

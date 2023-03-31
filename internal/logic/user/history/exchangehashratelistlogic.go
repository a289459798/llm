package history

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeHashRateListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeHashRateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeHashRateListLogic {
	return &ExchangeHashRateListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeHashRateListLogic) ExchangeHashRateList() (resp *types.HashRateExchangeListResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	exchange := []model.AIUserHashRate{}
	l.svcCtx.Db.Where("uid = ?", uid).Find(&exchange)
	if len(exchange) == 0 {
		return nil, errors.New(string(http.StatusNotFound))
	}
	data := []types.HashRateExchange{}
	for _, rate := range exchange {
		data = append(data, types.HashRateExchange{
			Date:   rate.CreatedAt.Format("2006-01-02"),
			Amount: rate.Amount,
			Use:    rate.UseAmount,
			Expiry: rate.Expiry.Format("2006-01-02 15:04:05"),
			Status: func() uint8 {
				if rate.Amount <= rate.UseAmount {
					return 1
				} else if time.Now().Unix() < rate.Expiry.Unix() && rate.Expiry.Unix()-time.Now().Unix() < 86400 {
					return 2
				} else if time.Now().Unix()-rate.Expiry.Unix() > 0 {
					return 3
				}
				return 0
			}(),
		})
	}
	return &types.HashRateExchangeListResponse{Data: data}, nil
}

package order

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"

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

func (l *DetailLogic) Detail(req *types.OrderInfoRequest) (resp *types.OrderInfoResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	order := &model.Order{}
	l.svcCtx.Db.Where("id = ?", req.OrderId).Where("uid = ?", uid).First(&order)

	if order.ID == 0 {
		return nil, errors.New("订单不存在")
	}

	return &types.OrderInfoResponse{
		OrderId: req.OrderId,
		Status:  order.Status,
	}, nil
}

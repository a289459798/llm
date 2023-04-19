package order

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service/order"
	"context"
	"encoding/json"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.OrderRequest) (resp *types.OrderResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db)
	if user.Uid == 0 {
		return nil, errors.New("用户不存在")
	}

	res, err := order.GetOrder(req.Type, order.OrderData{DB: l.svcCtx.Db}).Create(order.CreateRequest{
		Uid: uint32(uid),
		Items: []model.OrderItem{
			{ItemId: req.ItemId, Number: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	return &types.OrderResponse{
		OrderId: res.OrderId,
		Money:   res.Money,
	}, nil
}

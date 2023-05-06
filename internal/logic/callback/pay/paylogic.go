package pay

import (
	"chatgpt-tools/model"
	"chatgpt-tools/service/order"
	"chatgpt-tools/service/pay"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayLogic {
	return &PayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayLogic) Pay(req *types.PayRequest, r *http.Request) (resp *types.WechatPayResponse, err error) {
	config, err := model.PaySetting{Merchant: req.Merchant}.FindByMerchant(l.svcCtx.Db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	payModel := pay.GetPay(req.Type, pay.PayData{
		Ctx:      l.ctx,
		Config:   config,
		Merchant: req.Merchant,
	})

	payNotify, err := payModel.PayNotify(r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	outNo := payNotify.OutTradeNo
	var orderInfo []model.Order
	l.svcCtx.Db.Where("out_no = ?", outNo).Where("status = ?", model.PayStatusWaitPayment).Find(&orderInfo)

	if len(orderInfo) == 0 {
		return nil, errors.New("订单不存在")
	}

	err = l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		for _, m := range orderInfo {
			err = order.GetOrder(m.OrderType, order.OrderData{DB: l.svcCtx.Db}).Pay(order.PayRequest{
				OrderId: m.ID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.WechatPayResponse{
		Data: "success",
	}, nil
}

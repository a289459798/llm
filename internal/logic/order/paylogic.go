package order

import (
	"chatgpt-tools/model"
	"chatgpt-tools/service/pay"
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"

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

func (l *PayLogic) Pay(req *types.OrderPayRequest) (resp *types.OrderPayResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db)
	if user.Uid == 0 {
		return nil, errors.New("用户不存在")
	}

	order := &model.Order{}
	l.svcCtx.Db.Where("id = ?", req.OrderId).First(&order)

	if order.ID == 0 {
		return nil, errors.New("订单不存在")
	}

	var payStr string
	err = l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		merchant := "default"
		err = tx.Create(&model.OrderPay{
			OutNo:       order.OutNo,
			PayPrice:    order.PayPrice,
			RefundPrice: 0,
			Status:      model.PayStatusWaitPayment,
			Merchant:    merchant,
		}).Error
		if err != nil {
			return err
		}
		payModel := pay.GetPay(req.Platform, pay.PayData{
			Ctx:      l.ctx,
			Config:   "",
			Merchant: merchant,
		})
		payStr, err = payModel.Pay("h5", pay.Order{
			Body:       "Vip充值",
			OutNo:      order.OutNo,
			Total:      order.PayPrice,
			OpenId:     user.OpenId,
			NotifyPath: "callback/pay/wechat/" + merchant,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.OrderPayResponse{
		Data: payStr,
	}, nil
}

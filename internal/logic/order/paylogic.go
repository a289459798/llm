package order

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service/pay"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"

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

	var merchant = ""
	if req.Platform == "wechat" {
		merchant = "wechat_xm"
	} else {
		merchant = "alipay_xm"
	}

	config, err := model.PaySetting{Merchant: merchant}.FindByMerchant(l.svcCtx.Db)
	if err != nil {
		return nil, err
	}

	var payStr string
	err = l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {

		err = tx.Create(&model.OrderPay{
			OutNo:       order.OutNo,
			PayPrice:    order.PayPrice,
			RefundPrice: 0,
			Status:      model.PayStatusWaitPayment,
			PayType:     req.Platform,
			Merchant:    merchant,
		}).Error
		if err != nil {
			return err
		}
		payModel := pay.GetPay(req.Platform, pay.PayData{
			Ctx:      l.ctx,
			Config:   config,
			Merchant: merchant,
		})
		payStr, err = payModel.Pay(req.Scene, pay.Order{
			Body:  "购买商品",
			OutNo: order.OutNo,
			Total: func() float32 {
				if l.svcCtx.Config.Mode == "dev" {
					return 0.01
				}
				return order.PayPrice
			}(),
			OpenId: user.OpenId,
			NotifyPath: func() string {
				if l.svcCtx.Config.Mode == "dev" {
					return fmt.Sprintf("https://api.smuai.com/testpay/callback/pay/%s/%s", req.Platform, merchant)
				}
				return fmt.Sprintf("https://api.smuai.com/callback/pay/%s/%s", req.Platform, merchant)
			}(),
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

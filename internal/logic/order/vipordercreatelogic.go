package order

import (
	"chatgpt-tools/common/utils"
	pay2 "chatgpt-tools/common/utils/pay"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type VipOrderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipOrderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipOrderCreateLogic {
	return &VipOrderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipOrderCreateLogic) VipOrderCreate(req *types.VipPayRequest) (resp *types.PayResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db)
	if user.Uid == 0 {
		return nil, errors.New("用户不存在")
	}

	vip := &model.Vip{}
	l.svcCtx.Db.Where("id = ?", req.VipId).First(&vip)
	if vip.ID == 0 {
		return nil, errors.New("vip不存在")
	}

	order := &model.Order{
		Uid:       uint32(uid),
		OrderNo:   utils.GenerateOrderNo(),
		OutNo:     fmt.Sprintf("VIP%s", utils.GenerateOrderNo()),
		OrderType: "vip",
		CostPrice: 0,
		SellPrice: vip.Price,
		PayPrice:  vip.Price,
		Status:    model.OrderStatusWaitPayment,
	}

	var payStr []byte
	err = l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&order).Error
		if err != nil {
			return err
		}

		merchant := "default"
		err = tx.Create(&model.OrderItem{
			OrderId:   order.ID,
			ItemId:    vip.ID,
			Name:      "会员",
			Image:     "",
			CostPrice: order.CostPrice,
			SellPrice: order.SellPrice,
			PayPrice:  order.PayPrice,
			Number:    1,
		}).Error
		if err != nil {
			return err
		}
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
		payModel := pay2.GetPay(req.Platform, pay2.PayData{
			Ctx:      l.ctx,
			Config:   "",
			Merchant: merchant,
		})
		payData, err := payModel.Pay(pay2.Order{
			Body:       "Vip充值",
			OutNo:      order.OutNo,
			Total:      order.PayPrice,
			OpenId:     user.OpenId,
			NotifyPath: "callback/pay/wechat/" + merchant,
		})
		if err != nil {
			return err
		}
		payStr, err = json.Marshal(payData)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.PayResponse{
		Data: string(payStr),
	}, nil
}

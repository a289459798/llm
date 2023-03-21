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
	setting := &model.Setting{Name: "vip"}
	vipsSetting, err := setting.Find(l.svcCtx.Db)
	if err != nil {
		return nil, err
	}
	user := model.User{ID: uint32(uid)}.Find(l.svcCtx.Db)
	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}
	firstMonth := model.Order{Uid: uint32(uid)}.FirstMonthVip(l.svcCtx.Db)
	if !firstMonth {
		vipsSetting["first"] = vipsSetting["original"]
	}
	order := &model.Order{
		Uid:       uint32(uid),
		OrderNo:   utils.GenerateOrderNo(),
		OutNo:     fmt.Sprintf("VIP%s", utils.GenerateOrderNo()),
		OrderType: "vip",
		CostPrice: 0,
		SellPrice: float32(vipsSetting["original"].(float64)),
		PayPrice:  float32(vipsSetting["first"].(float64)),
		Status:    model.OrderStatusWaitPayment,
	}

	tx := l.svcCtx.Db.Begin()

	err = tx.Create(&order).Error
	if err != nil {
		return nil, err
	}
	err = tx.Create(&model.OrderItem{
		OrderId:   order.ID,
		ItemId:    0,
		Name:      "会员",
		Image:     "",
		CostPrice: order.CostPrice,
		SellPrice: order.SellPrice,
		PayPrice:  order.PayPrice,
		Number:    1,
	}).Error
	if err != nil {
		return nil, err
	}
	err = tx.Create(&model.OrderPay{
		OutNo:       order.OutNo,
		PayPrice:    order.PayPrice,
		RefundPrice: 0,
		Status:      model.PayStatusWaitPayment,
	}).Error
	if err != nil {
		return nil, err
	}
	payModel := pay2.GetPay(req.Platform, pay2.PayData{
		Ctx:    l.ctx,
		Config: l.svcCtx.Config,
	})
	payData, err := payModel.Pay(pay2.Order{
		Body:   "Vip充值",
		OutNo:  order.OutNo,
		Total:  order.PayPrice,
		OpenId: user.OpenId,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	payStr, err := json.Marshal(payData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &types.PayResponse{
		Data: string(payStr),
	}, nil
}

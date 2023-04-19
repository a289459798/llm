package order

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type VipOrder struct {
	DB *gorm.DB
}

func NewVip(data OrderData) *VipOrder {
	return &VipOrder{}
}

func (vipOrder *VipOrder) Create(orderData CreateRequest) (response CreateResponse, err error) {

	vip := &model.Vip{}
	vipOrder.DB.Where("id = ?", orderData.Items[0].ItemId).First(&vip)
	if vip.ID == 0 {
		err = errors.New("vip不存在")
		return
	}

	order := &model.Order{
		Uid:       orderData.Uid,
		OrderNo:   utils.GenerateOrderNo(),
		OutNo:     fmt.Sprintf("VIP%s", utils.GenerateOrderNo()),
		OrderType: "vip",
		CostPrice: 0,
		SellPrice: vip.Price * float32(orderData.Items[0].Number),
		PayPrice:  vip.Price * float32(orderData.Items[0].Number),
		Status:    model.OrderStatusWaitPayment,
	}

	err = vipOrder.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&order).Error
		if err != nil {
			return err
		}
		err = tx.Create(&model.OrderItem{
			OrderId:   order.ID,
			ItemId:    vip.ID,
			Name:      "会员",
			Image:     "",
			CostPrice: order.CostPrice,
			SellPrice: order.SellPrice,
			PayPrice:  order.PayPrice,
			Number:    orderData.Items[0].Number,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err == nil {
		response.OrderId = string(order.ID)
		response.Money = order.PayPrice
	}

	return
}

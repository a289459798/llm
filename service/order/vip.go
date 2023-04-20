package order

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type VipOrder struct {
	DB *gorm.DB
}

func NewVip(data OrderData) *VipOrder {
	return &VipOrder{
		DB: data.DB,
	}
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
		response.OrderId = strconv.Itoa(int(order.ID))
		response.Money = order.PayPrice
	}

	return
}

func (vipOrder *VipOrder) Pay(orderData PayRequest) error {

	orderInfo := &model.Order{}
	vipOrder.DB.Where("id = ?", orderData.OrderId).Preload("Item").Where("status = ?", model.PayStatusWaitPayment).First(&orderInfo)
	if orderInfo.ID == 0 {
		return errors.New("订单不存在")
	}

	err := vipOrder.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range orderInfo.Item {

			vip := &model.Vip{}
			tx.Where("id = ?", item.ItemId).First(&vip)
			if vip.ID == 0 {
				return errors.New("会员不存在")
			}
			err := model.AIUser{Uid: orderInfo.Uid}.Find(tx).SetVip(tx, &model.VipCode{
				VipId: vip.ID,
				Day:   vip.Day,
				Vip: model.Vip{
					Amount: vip.Amount,
				},
			})
			if err != nil {
				return err
			}
		}
		orderInfo.Status = model.PayStatusPayment
		tx.Save(orderInfo)
		tx.Model(&model.OrderPay{}).
			Where("out_no", orderInfo.OutNo).
			Update("status", model.PayStatusPayment).
			Update("pay_time", time.Now().Format("2006-01-02 15:04:05"))
		return nil
	})

	model.Distributor{}.AddMoney(vipOrder.DB, model.DistributorAdd{
		Uid:   orderInfo.Uid,
		Money: orderInfo.PayPrice,
		Way:   0,
		Type:  "vip",
	})

	return err
}

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

type HashRateOrder struct {
	DB *gorm.DB
}

func NewHashRate(data OrderData) *HashRateOrder {
	return &HashRateOrder{
		DB: data.DB,
	}
}

func (hashRateOrder *HashRateOrder) Create(orderData CreateRequest) (response CreateResponse, err error) {

	hashRate := &model.AIHashRate{}
	hashRateOrder.DB.Where("id = ?", orderData.Items[0].ItemId).First(&hashRate)
	if hashRate.ID == 0 {
		err = errors.New("算力不存在")
		return
	}

	order := &model.Order{
		Uid:       orderData.Uid,
		OrderNo:   utils.GenerateOrderNo(),
		OutNo:     fmt.Sprintf("VIP%s", utils.GenerateOrderNo()),
		OrderType: "hashrate",
		CostPrice: 0,
		SellPrice: hashRate.Price * float32(orderData.Items[0].Number),
		PayPrice:  hashRate.Price * float32(orderData.Items[0].Number),
		Status:    model.OrderStatusWaitPayment,
	}

	err = hashRateOrder.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&order).Error
		if err != nil {
			return err
		}
		err = tx.Create(&model.OrderItem{
			OrderId:   order.ID,
			ItemId:    hashRate.ID,
			Name:      "算力",
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

func (hashRateOrder *HashRateOrder) Pay(orderData PayRequest) error {

	orderInfo := &model.Order{}
	hashRateOrder.DB.Where("id = ?", orderData.OrderId).Preload("Item").Where("status = ?", model.PayStatusWaitPayment).First(&orderInfo)
	if orderInfo.ID == 0 {
		return errors.New("订单不存在")
	}

	err := hashRateOrder.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range orderInfo.Item {

			hashRate := &model.AIHashRate{}
			tx.Where("id = ?", item.ItemId).First(&hashRate)
			if hashRate.ID == 0 {
				return errors.New("算力不存在")
			}
			err := tx.Create(&model.AIUserHashRate{
				Uid:       orderInfo.Uid,
				Amount:    hashRate.Amount,
				UseAmount: 0,
				Expiry:    time.Now().AddDate(0, 0, int(hashRate.Day)),
			}).Error
			if err != nil {
				return errors.New("兑换错误")
			}
			account := model.NewAccount(tx).GetAccount(orderInfo.Uid, time.Now())
			err = tx.Create(&model.AccountRecord{
				Uid:           orderInfo.Uid,
				RecordId:      0,
				Way:           1,
				Type:          "exchange",
				Amount:        hashRate.Amount,
				CurrentAmount: account.Amount,
			}).Error
			if err != nil {
				return err
			}
		}
		orderInfo.Status = model.PayStatusPayment
		tx.Save(orderInfo)
		return nil
	})

	return err
}

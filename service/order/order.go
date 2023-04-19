package order

import (
	"gorm.io/gorm"
)

type SMUOrder interface {
	Create(orderI CreateRequest) (response CreateResponse, err error)
}

type OrderData struct {
	DB *gorm.DB
}

func GetOrder(orderType string, orderData OrderData) SMUOrder {
	if orderType == "vip" {
		return NewVip(orderData)
	}
	return NewHashRate(orderData)
}

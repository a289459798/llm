package order

import (
	"gorm.io/gorm"
)

type SMUOrder interface {
	Create(orderData CreateRequest) (response CreateResponse, err error)
	Pay(orderData PayRequest) error
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

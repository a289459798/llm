package order

import "chatgpt-tools/model"

type CreateRequest struct {
	Uid   uint32
	Items []model.OrderItem
}

type CreateResponse struct {
	OrderId string
	Money   float32
}

type PayRequest struct {
	OrderId uint32
}

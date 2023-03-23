package model

import (
	"time"
)

const (
	PayStatusWaitPayment uint8 = 0
	PayStatusPayment     uint8 = 1
	PayStatusRefund      uint8 = 2
)

type OrderPay struct {
	ID          uint32    `gorm:"primary_key" json:"id"`
	OutNo       string    `json:"out_no" gorm:"type:varchar(100);uniqueIndex:ik_out_no"`
	PayPrice    float32   `json:"pay_price" gorm:"type:decimal(10,2)"`
	RefundPrice float32   `json:"refund_price" gorm:"type:decimal(10,2)"`
	Status      uint8     `json:"status"`
	PayType     string    `json:"pay_type" gorm:"type:varchar(20)"`
	Merchant    string    `json:"merchant" gorm:"type:varchar(20)"`
	PayTime     uint8     `json:"pay_time" gorm:"type:TIMESTAMP"`
	CreatedAt   time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt    time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

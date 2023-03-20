package model

import (
	"time"
)

type Order struct {
	ID           uint32    `gorm:"primary_key" json:"id"`
	Uid          uint32    `json:"uid" gorm:"index:ik_uid"`
	OrderNo      string    `json:"order_n0" gorm:"type:varchar(100);uniqueIndex:uk_order_no"`
	OutNo        string    `json:"out_no" gorm:"type:varchar(100)"`
	OrderType    string    `json:"order_type" gorm:"type:varchar(20)"`
	CostPrice    float32   `json:"cost_price"`
	SellPrice    float32   `json:"sell_price"`
	PayPrice     float32   `json:"pay_price"`
	Status       uint8     `json:"status"`
	PayType      string    `json:"pay_type" gorm:"type:varchar(20)"`
	PayTime      uint8     `json:"pay_time" gorm:"type:TIMESTAMP"`
	CancelTime   uint8     `json:"cancel_time" gorm:"type:TIMESTAMP"`
	CompleteTime uint8     `json:"complete_time" gorm:"type:TIMESTAMP"`
	CreatedAt    time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt     time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	OrderStatusWaitPayment uint8 = 0
	OrderStatusPayment     uint8 = 1
	OrderStatusComplete    uint8 = 2
	OrderStatusCancel      uint8 = 3
)

type Order struct {
	ID           uint32      `gorm:"primary_key" json:"id"`
	Uid          uint32      `json:"uid" gorm:"index:ik_uid"`
	OrderNo      string      `json:"order_no" gorm:"type:varchar(100);uniqueIndex:uk_order_no"`
	OutNo        string      `json:"out_no" gorm:"type:varchar(100);index:ik_out_no"`
	OrderType    string      `json:"order_type" gorm:"type:varchar(20)"`
	CostPrice    float32     `json:"cost_price" gorm:"type:decimal(10,2)"`
	SellPrice    float32     `json:"sell_price" gorm:"type:decimal(10,2)"`
	PayPrice     float32     `json:"pay_price" gorm:"type:decimal(10,2)"`
	Status       uint8       `json:"status"`
	CancelTime   time.Time   `json:"cancel_time" gorm:"type:TIMESTAMP"`
	CompleteTime time.Time   `json:"complete_time" gorm:"type:TIMESTAMP"`
	CreatedAt    time.Time   `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt     time.Time   `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
	Item         []OrderItem `gorm:"foreignKey:order_id"`
}

func (o Order) FirstMonthVip(db *gorm.DB) bool {
	db.Where("uid = ?", o.Uid).
		Where("order_type = ?", "vip").
		Where("status = ?", OrderStatusComplete).
		First(&o)
	fmt.Println(o)
	if o.ID > 0 {
		return false
	}
	return true
}

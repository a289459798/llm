package model

import (
	"time"
)

type OrderItem struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	OrderId   uint32    `json:"order_id" gorm:"index:ik_uid"`
	ItemId    uint32    `json:"item_id" `
	Name      string    `json:"name" gorm:"type:varchar(50)"`
	Image     string    `json:"image" gorm:"type:varchar(100)"`
	CostPrice float32   `json:"cost_price" gorm:"type:decimal(10,2)"`
	SellPrice float32   `json:"sell_price" gorm:"type:decimal(10,2)"`
	PayPrice  float32   `json:"pay_price" gorm:"type:decimal(10,2)"`
	Number    uint32    `json:"number"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

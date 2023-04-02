package model

import (
	"time"
)

type Vip struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Name      string    `json:"name" gorm:"type:varchar(20)"`
	Origin    float32   `json:"origin" gorm:"type:decimal(10,1)"`
	Price     float32   `json:"price" gorm:"type:decimal(10,1)"`
	Amount    uint32    `json:"amount"`
	Discount  float32   `json:"discount" gorm:"type:decimal(10,1)"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

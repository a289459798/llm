package model

import (
	"time"
)

type AIHashRate struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Origin    float32   `json:"origin" gorm:"type:decimal(10,1)"`
	Price     float32   `json:"price" gorm:"type:decimal(10,1)"`
	Amount    uint32    `json:"amount"`
	Day       uint32    `json:"day"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

package model

import "time"

type DistributorLevel struct {
	ID         uint32    `gorm:"primary_key" json:"id"`
	Name       string    `json:"name" gorm:"type:varchar(50)"`
	UserNumber uint32    `json:"user_number" `
	UserPrice  float32   `json:"user_price" gorm:"type:decimal(10,2)"`
	Ratio      float32   `json:"ratio" gorm:"type:decimal(10,2)"`
	CreatedAt  time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt   time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

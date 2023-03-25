package model

import (
	"time"
)

type User struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Phone     string    `json:"phone" gorm:"type:varchar(11);index:ik_phone"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

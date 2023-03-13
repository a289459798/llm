package model

import (
	"time"
)

type Contraband struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `json:"uid"`
	Content   string    `gorm:"type:varchar(255)" json:"content"`
	Error     string    `gorm:"type:varchar(50)" json:"error"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

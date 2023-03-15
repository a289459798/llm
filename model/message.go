package model

import (
	"time"
)

type Message struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Title     string    `json:"title" gorm:"type:varchar(50)" `
	Content   string    `json:"content" gorm:"type:varchar(255)" `
	Link      string    `json:"link" gorm:"type:varchar(255)" `
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

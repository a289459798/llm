package model

import (
	"time"
)

type ChatTemplate struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Title     string    `json:"title" gorm:"type:varchar(50)" `
	Question  string    `json:"question" gorm:"type:text" `
	Answer    string    `json:"answer" gorm:"type:varchar(255)" `
	Type      string    `json:"type" gorm:"type:varchar(50);index:idx_type" `
	Welcome   string    `json:"welcome" gorm:"type:varchar(100);" `
	Sort      uint8     `json:"sort"`
	IsDel     bool      `json:"is_del"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

package model

import "time"

type Record struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `json:"uid" gorm:"index:idx_uid"`
	Type      string    `gorm:"type:varchar(20)" json:"type"`
	Content   string    `gorm:"type:varchar(255)" json:"content"`
	Result    string    `json:"result"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

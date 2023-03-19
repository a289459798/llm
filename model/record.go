package model

import "time"

type Record struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `json:"uid" gorm:"index:idx_uid"`
	Type      string    `gorm:"type:varchar(20)" json:"type"`
	ChatId    string    `gorm:"type:varchar(50)" json:"chat_id"`
	Content   string    `gorm:"type:varchar(255)" json:"content"`
	Result    string    `json:"result"`
	Model     string    `json:"model" gorm:"type:varchar(20)"`
	Platform  string    `json:"platform" gorm:"type:varchar(20)"`
	CreatedAt time.Time `gorm:"column:created_at;index:idx_created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

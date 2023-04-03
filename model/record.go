package model

import "time"

type Record struct {
	ID          uint32    `gorm:"primary_key" json:"id"`
	Uid         uint32    `json:"uid" gorm:"index:idx_uid"`
	Type        string    `gorm:"type:varchar(20)" json:"type"`
	ChatId      string    `gorm:"type:varchar(50)" json:"chat_id"`
	Title       string    `gorm:"type:varchar(50)" json:"title"`
	Content     string    `json:"content" gorm:"type:text"`
	ShowContent string    `json:"show_content" gorm:"type:text"`
	Result      string    `json:"result" gorm:"type:text"`
	ShowResult  string    `json:"show_result" gorm:"type:text"`
	Model       string    `json:"model" gorm:"type:varchar(20)"`
	Platform    string    `json:"platform" gorm:"type:varchar(20)"`
	IsDelete    bool      `json:"is_delete" gorm:"default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;index:idx_created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt    time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

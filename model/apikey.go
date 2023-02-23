package model

import "time"

type Apikey struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Channel   string    `gorm:"type:varchar(50)" json:"channel"`
	Ori       string    `gorm:"type:varchar(50)" json:"ori"`
	Key       string    `gorm:"type:varchar(100)" json:"key"`
	Secret    string    `gorm:"type:varchar(100)" json:"secret"`
	Token     string    `gorm:"type:varchar(255)" json:"token"`
	Status    bool      `json:"status"`
	Remark    string    `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

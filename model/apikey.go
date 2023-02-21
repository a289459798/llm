package model

import "time"

type Apikey struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Channel   string    `gorm:"type:varchar(50)" json:"channel"`
	Ori       string    `gorm:"type:varchar(50)" json:"ori"`
	Key       string    `gorm:"type:varchar(100)" json:"key"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

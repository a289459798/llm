package model

import (
	"time"
)

type ScanScene struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Scene     string    `json:"scene" gorm:"type:varchar(64);uniqueIndex:uk_uid_date"`
	Type      string    `json:"type" gorm:"type:varchar(20)"`
	Data      string    `json:"data"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

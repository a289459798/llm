package model

import (
	"time"
)

type RequestLog struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	RequestId string    `json:"request_id" gorm:"type:varchar(64);uniqueIndex:uk_request_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

package model

import (
	"time"
)

type AIUserVip struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `gorm:"after:id;uniqueIndex:uk_uid" json:"uid"`
	VipId     uint32    `json:"vip_id"`
	VipExpiry time.Time `json:"vip_expiry" gorm:"type:TIMESTAMP"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
	Vip       Vip       `gorm:"foreignKey:vip_id;"`
}

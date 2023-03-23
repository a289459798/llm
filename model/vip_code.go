package model

import (
	"time"
)

type VipCode struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `json:"uid" gorm:"uniqueIndex:uk_uid_code"`
	Code      string    `gorm:"type:varchar(64);uniqueIndex:uk_uid_code" json:"code"`
	VipId     uint32    `json:"vip_id"`
	Day       uint32    `json:"day"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
	Vip       Vip       `gorm:"foreignKey:vip_id"`
}

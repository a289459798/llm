package model

import "time"

type Account struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `json:"uid" gorm:"index:idx_uid"`
	Type      uint8     `json:"type" gorm:"COMMENT='0 增加 1 减少'"`
	Current   uint32    `json:"current" gorm:"COMMENT='当前余额'"`
	Amount    uint32    `json:"amount" gorm:"COMMENT='变更余额'"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

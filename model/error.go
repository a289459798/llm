package model

import (
	"gorm.io/gorm"
	"time"
)

type Error struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `json:"uid" gorm:"index:idx_uid"`
	Type      string    `gorm:"type:varchar(20)" json:"type"`
	Question  string    `json:"amount" gorm:"type:text"`
	Error     string    `json:"current_amount" gorm:"type:text"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (err Error) Insert(db *gorm.DB) {
	db.Create(&err)
}

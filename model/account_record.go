package model

import (
	"gorm.io/gorm"
	"time"
)

type AccountRecord struct {
	ID            uint32    `gorm:"primary_key" json:"id"`
	Uid           uint32    `json:"uid" gorm:"uniqueIndex:uk_uid_date"`
	RecordId      uint32    `json:"record_id"`
	Way           uint8     `json:"type" gorm:"COMMENT:方式"`
	Type          string    `gorm:"type:varchar(20);uniqueIndex:uk_uid_date" json:"type"`
	Amount        uint32    `json:"amount"`
	CurrentAmount uint32    `json:"current_amount"`
	CreatedAt     time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create;uniqueIndex:uk_uid_date" json:"created_at,omitempty"`
	UpdateAt      time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}
type AccountRecordModel struct {
	DB *gorm.DB
}

func NewAccountRecord(db *gorm.DB) *AccountRecordModel {
	return &AccountRecordModel{
		DB: db,
	}
}

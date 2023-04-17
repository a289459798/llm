package model

import (
	"gorm.io/gorm"
	"time"
)

type DistributorRecord struct {
	ID             uint32    `gorm:"primary_key" json:"id"`
	DistributorUid uint32    `json:"distributor_uid" gorm:"uniqueIndex:uk_uid"`
	Uid            uint32    `json:"uid" gorm:"uniqueIndex:uk_uid"`
	CreatedAt      time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt       time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (dr DistributorRecord) Create(db *gorm.DB) {
	db.Create(&dr)
}

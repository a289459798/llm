package model

import (
	"time"
)

type DistributorMoneyRecord struct {
	ID             uint32    `gorm:"primary_key" json:"id"`
	DistributorUid uint32    `json:"distributor_uid" gorm:"index:uk_distributor_uid"`
	Money          float32   `json:"money" gorm:"type:decimal(10,2)"`
	Way            uint8     `json:"way"`
	Type           string    `json:"type" gorm:"type:varchar(50)"`
	Remark         string    `json:"remark" gorm:"type:varchar(100)"`
	CreatedAt      time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt       time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

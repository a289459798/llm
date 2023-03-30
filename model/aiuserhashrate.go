package model

import (
	"gorm.io/gorm"
	"time"
)

type AIUserHashRate struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `gorm:"after:id;index:ik_uid" json:"uid"`
	Amount    uint32    `json:"amount"`
	UseAmount uint32    `json:"use_amount"`
	Expiry    time.Time `json:"expiry" gorm:"type:TIMESTAMP"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (aiUserHashRate AIUserHashRate) GetTotalAmount(db *gorm.DB) uint32 {
	db.Where("uid = ?", aiUserHashRate.Uid).
		Where("expiry >= ?", time.Now().Format("2006-01-02 15:04:05")).
		Select("sum(amount-use_amount) as amount").
		First(&aiUserHashRate)
	return aiUserHashRate.Amount
}

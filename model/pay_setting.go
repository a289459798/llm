package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type PaySetting struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Merchant  string    `json:"merchant" gorm:"type:varchar(50);uniqueIndex:uk_merchant"`
	Setting   string    `json:"setting"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (ps PaySetting) FindByMerchant(db *gorm.DB) (string, error) {
	db.Where("merchant = ?", ps.Merchant).Find(&ps)
	if ps.ID == 0 {
		return "", errors.New("商户配置不存在")
	}
	return ps.Setting, nil
}

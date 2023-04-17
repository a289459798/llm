package model

import (
	"gorm.io/gorm"
	"time"
)

type DistributorPayRecord struct {
	ID             uint32    `gorm:"primary_key" json:"id"`
	DistributorUid uint32    `json:"distributor_uid" gorm:"index:uk_distributor_uid"`
	Uid            uint32    `json:"uid" `
	Pay            float32   `json:"pay" gorm:"type:decimal(10,2)"`
	Ratio          float32   `json:"ratio" gorm:"type:decimal(10,2)"`
	Money          float32   `json:"money" gorm:"type:decimal(10,2)"`
	CreatedAt      time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt       time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (dpr DistributorPayRecord) TotalPayWithDate(db *gorm.DB, uid uint32, timeRange []string) float32 {
	tx := db.Model(dpr).Where("distributor_uid = ?", uid)
	if timeRange != nil {
		tx.Where("created_at between ? and ?", timeRange[0], timeRange[1])
	}
	tx.Select("sum(pay) as pay").First(&dpr)
	return dpr.Pay
}

func (dpr DistributorPayRecord) TotalMoneyWithDate(db *gorm.DB, uid uint32, timeRange []string) float32 {
	tx := db.Model(dpr).Where("distributor_uid = ?", uid)
	if timeRange != nil {
		tx.Where("created_at between ? and ?", timeRange[0], timeRange[1])
	}
	tx.Select("sum(money) as money").First(&dpr)
	return dpr.Money
}

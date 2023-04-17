package model

import (
	"gorm.io/gorm"
	"time"
)

type DistributorRecord struct {
	ID             uint32    `gorm:"primary_key" json:"id"`
	DistributorUid uint32    `json:"distributor_uid" gorm:"uniqueIndex:uk_distributor_uid"`
	Uid            uint32    `json:"uid" gorm:"uniqueIndex:uk_uid"`
	CreatedAt      time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt       time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (dr DistributorRecord) Create(db *gorm.DB) {
	err := db.Create(&dr).Error
	if err != nil {
		var amount uint32 = 10
		db.Create(AIUserHashRate{
			Uid:       dr.DistributorUid,
			Amount:    amount,
			UseAmount: 0,
			Expiry:    time.Now().AddDate(0, 0, 5),
		})

		account := NewAccount(db).GetAccount(dr.DistributorUid, time.Now())
		db.Create(AccountRecord{
			Uid:           dr.DistributorUid,
			RecordId:      0,
			Way:           1,
			Type:          "invite",
			Amount:        amount,
			CurrentAmount: account.Amount,
		})

		Distributor{Uid: dr.DistributorUid}.CheckUpgrade(db)
	}
}

func (dr DistributorRecord) TotalWithDate(db *gorm.DB, uid uint32, timeRange []string) uint32 {
	var total int64
	tx := db.Model(DistributorRecord{}).Where("distributor_uid = ?", uid)
	if timeRange != nil {
		tx.Where("created_at between ? and ?", timeRange[0], timeRange[1])
	}
	tx.Count(&total)
	return uint32(total)
}

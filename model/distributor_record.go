package model

import (
	"gorm.io/gorm"
	"time"
)

type DistributorRecord struct {
	ID             uint32    `gorm:"primary_key" json:"id"`
	DistributorUid uint32    `json:"distributor_uid" gorm:"index:ik_distributor_uid"`
	Uid            uint32    `json:"uid" gorm:"uniqueIndex:uk_uid"`
	CreatedAt      time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt       time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (dr DistributorRecord) Create(db *gorm.DB) {
	distributor := &Distributor{}
	db.Where("uid = ?", dr.DistributorUid).Where("status = 1").First(&distributor)
	if distributor.ID == 0 {
		return
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&dr).Error
		if err != nil {
			return err
		}
		var amount uint32 = 20
		err = tx.Create(&AIUserHashRate{
			Uid:       dr.DistributorUid,
			Amount:    amount,
			UseAmount: 0,
			Expiry:    time.Now().AddDate(0, 0, 30),
		}).Error
		if err != nil {
			return err
		}
		account := NewAccount(tx).GetAccount(dr.DistributorUid, time.Now())
		err = tx.Create(&AccountRecord{
			Uid:           dr.DistributorUid,
			RecordId:      0,
			Way:           1,
			Type:          "invite",
			Amount:        amount,
			CurrentAmount: account.Amount,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return
	}
	Distributor{Uid: dr.DistributorUid}.CheckUpgrade(db)
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

package model

import (
	"gorm.io/gorm"
	"time"
)

type Distributor struct {
	ID        uint32           `gorm:"primary_key" json:"id"`
	Uid       uint32           `json:"uid" gorm:"uniqueIndex:uk_uid"`
	LevelId   uint32           `json:"level_id"`
	Ratio     float32          `json:"ratio" gorm:"type:decimal(10,2)"`
	Money     float32          `json:"money" gorm:"type:decimal(10,2)"`
	Status    bool             `json:"status"`
	CreatedAt time.Time        `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time        `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
	Level     DistributorLevel `gorm:"foreignKey:level_id"`
}

type DistributorAdd struct {
	Money  float32
	Uid    uint32
	Way    uint8
	Type   string
	Remark string
}

func (d Distributor) AddMoney(db *gorm.DB, record DistributorAdd) error {

	db.Model(DistributorRecord{}).
		Joins("inner join gpt_distributor on gpt_distributor.uid=gpt_distributor_record.distributor_uid").
		Where("gpt_distributor_record.uid = ?", record.Uid).
		Select("gpt_distributor.*").
		First(&d)
	if d.ID == 0 {
		return nil
	}
	tx := db.Begin()
	money := record.Money
	if record.Way == 0 {
		money = record.Money * d.Ratio / 100
		err := tx.Create(&DistributorPayRecord{
			DistributorUid: 0,
			Uid:            record.Uid,
			Pay:            record.Money,
			Ratio:          d.Ratio,
			Money:          money,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err := tx.Create(&DistributorMoneyRecord{
		DistributorUid: d.Uid,
		Money:          money,
		Way:            record.Way,
		Type:           record.Type,
		Remark:         record.Remark,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 修改余额
	d.Money = func() float32 {
		if record.Way == 0 {
			return d.Money + money
		} else {
			return d.Money - money
		}
	}()
	tx.Save(d)
	tx.Commit()
	return nil
}

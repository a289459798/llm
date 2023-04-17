package model

import (
	"fmt"
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

func (d Distributor) CheckUpgrade(db *gorm.DB) {
	if d.LevelId == 0 {
		db.Where("uid = ?", d.Uid).Where("status = 1").First(&d)
	}
	if d.ID == 0 {
		return
	}
	// 下一级要求
	level := DistributorLevel{}
	db.Where("id > ?", d.LevelId).Order("id asc").First(&level)
	if level.ID == 0 {
		return
	}
	// 查看条件
	totalUser := DistributorRecord{}.TotalWithDate(db, d.Uid, nil)
	totalPay := DistributorPayRecord{}.TotalPayWithDate(db, d.Uid, nil)

	if totalUser >= level.UserNumber && totalPay >= level.UserPrice {
		// 升级
		d.LevelId = level.ID
		db.Save(d)
		// 站内信呢
		db.Create(&AINotify{
			Uid:     d.Uid,
			Title:   "推广员等级变更",
			Content: fmt.Sprintf("恭喜您升级为%s", level.Name),
			Link:    "",
			Status:  false,
		})
	}
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
	err := db.Transaction(func(tx *gorm.DB) error {
		money := record.Money
		if record.Way == 0 {
			money = record.Money * d.Ratio / 100
			err := tx.Create(&DistributorPayRecord{
				DistributorUid: d.Uid,
				Uid:            record.Uid,
				Pay:            record.Money,
				Ratio:          d.Ratio,
				Money:          money,
			}).Error
			if err != nil {
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
		return nil
	})

	if err != nil {
		return err
	}
	d.CheckUpgrade(db)
	return err
}

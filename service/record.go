package service

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/model"
	"gorm.io/gorm"
	"time"
)

type Record struct {
	DB *gorm.DB
}

type RecordParams struct {
	Params string `json:"params"`
}

func NewRecord(db *gorm.DB) *Record {
	return &Record{
		DB: db,
	}
}

func (r *Record) Insert(record *model.Record, params *RecordParams) {
	r.DB.Create(record)

	r.DB.Transaction(func(tx *gorm.DB) error {
		// 消耗次数
		chatUse := uint32(utils.GetSuanLi(record.Uid, record.Type, func() string {
			if params != nil {
				return params.Params
			}
			return ""
		}(), tx))
		amount := model.NewAccount(tx).GetAccount(record.Uid, time.Now())
		// 优先处理临时算力
		if (amount.ChatAmount - amount.ChatUse) >= chatUse {
			amount.ChatUse += chatUse
		} else {
			amount.ChatUse = amount.ChatAmount
			remainUse := chatUse - (amount.ChatAmount - amount.ChatUse)
			hashRate := []model.AIUserHashRate{}
			tx.Where("uid = ?", amount.Uid).Where("expiry >= ?", time.Now().Format("2006-01-02 15:04:05")).Where("amount > use_amount").Order("expiry asc, id asc").Find(&hashRate)
			for _, rate := range hashRate {
				tmpUse := int(remainUse) - (int(rate.Amount) - int(rate.UseAmount))
				if tmpUse <= 0 {
					rate.UseAmount += remainUse
					tx.Save(&rate)
					break
				}
				remainUse = uint32(tmpUse)
				rate.UseAmount = rate.Amount
				tx.Save(&rate)
			}
		}

		tx.Save(&amount)

		// 记录
		tx.Create(&model.AccountRecord{
			Uid:           record.Uid,
			RecordId:      record.ID,
			Way:           0,
			Type:          record.Type,
			Amount:        chatUse,
			CurrentAmount: amount.Amount - chatUse,
		})
		return nil
	})
}

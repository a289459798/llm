package service

import (
	"chatgpt-tools/model"
	"gorm.io/gorm"
	"time"
)

type Record struct {
	DB *gorm.DB
}

func NewRecord(db *gorm.DB) *Record {
	return &Record{
		DB: db,
	}
}

func (r *Record) Insert(record *model.Record) {
	r.DB.Create(record)

	r.DB.Transaction(func(tx *gorm.DB) error {
		// 消耗次数
		var chatUse uint32 = 1
		if record.Type == "image/create" {
			chatUse = 3
		}
		amount := model.NewAccount(tx).GetAccount(record.Uid, time.Now())
		amount.ChatUse += chatUse
		tx.Save(&amount)

		// 记录
		tx.Create(&model.AccountRecord{
			Uid:           record.Uid,
			RecordId:      record.ID,
			Way:           0,
			Type:          record.Type,
			Amount:        chatUse,
			CurrentAmount: amount.ChatAmount - amount.ChatUse,
		})
		return nil
	})
}

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

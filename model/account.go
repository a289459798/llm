package model

import (
	"gorm.io/gorm"
	"math"
	"time"
)

type Account struct {
	ID         uint32    `gorm:"primary_key" json:"id"`
	Uid        uint32    `json:"uid" gorm:"index:idx_uid"`
	ChatAmount uint32    `json:"chat_amount" gorm:"COMMENT:总次数"`
	ChatUse    uint32    `json:"chat_use" gorm:"COMMENT:使用次数"`
	Date       time.Time `json:"date" gorm:"type:date;COMMENT='日期'"`
	CreatedAt  time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt   time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

type AccountModel struct {
	DB *gorm.DB
}

func NewAccount(db *gorm.DB) *AccountModel {
	return &AccountModel{
		DB: db,
	}
}

func (a *AccountModel) GetAccount(uid uint32, date time.Time) *Account {

	account := &Account{}
	a.DB.Where("uid = ?", uid).Where("date = ?", date.Format("2006-01-02")).Find(&account)
	if account.ID == 0 {
		a.DB.Transaction(func(tx *gorm.DB) error {
			var amount uint32 = 10
			// 获取连续天数
			firstAccount := &Account{}
			tx.Where("uid = ?", uid).Order("id desc").Find(&firstAccount)
			if firstAccount.ID > 0 {
				var total int64
				tx.Model(&Account{}).Where("uid = ?", uid).Count(&total)
				t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				t2, _ := time.Parse("2006-01-02", firstAccount.Date.Format("2006-01-02"))
				d := t.Sub(t2)
				day := int64(math.Ceil(d.Hours() / 24))
				if day == total {
					amount += uint32(total)
				}
			}
			account.Uid = uid
			account.ChatAmount = amount
			account.ChatUse = 0
			account.Date = date
			tx.Create(&account)

			tx.Create(&AccountRecord{
				Uid:           uid,
				RecordId:      0,
				Way:           1,
				Type:          "open",
				Amount:        amount,
				CurrentAmount: account.ChatAmount - account.ChatUse,
			})

			return nil
		})

		return account
	}
	return account
}

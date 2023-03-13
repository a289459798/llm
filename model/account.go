package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID         uint32    `gorm:"primary_key" json:"id"`
	Uid        uint32    `json:"uid" gorm:"index:idx_uid"`
	ChatAmount uint32    `json:"chat_amount" gorm:"COMMENT:总次数"`
	ChatUse    uint32    `json:"chat_use" gorm:"COMMENT:使用次数"`
	Date       time.Time `json:"date" gorm:"type:date;COMMENT='日期'"`
	LoginCount uint32    `json:"login_count"`
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
			var amount uint32 = 5
			// 获取连续天数
			yesterdayAccount := &Account{}
			tx.Where("uid = ?", uid).Where("date = ?", time.Now().AddDate(0, 0, -1).Format("2006-01-02")).Find(&yesterdayAccount)
			account.LoginCount = 1
			if yesterdayAccount.ID > 0 {
				amount += yesterdayAccount.LoginCount
				account.LoginCount += yesterdayAccount.LoginCount
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

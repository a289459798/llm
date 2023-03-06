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
		account.Uid = uid
		account.ChatAmount = 10
		account.ChatUse = 0
		account.Date = date
		a.DB.Create(&account)
		return account
	}
	return account
}

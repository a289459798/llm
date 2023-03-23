package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	OpenId    string    `gorm:"type:varchar(64);uniqueIndex:uk_openid" json:"open_id"`
	UnionId   string    `gorm:"type:varchar(64)"json:"union_id"`
	Subscribe bool      `json:"subscribe"`
	JoinGroup bool      `json:"join_group" gorm:"default:0"`
	VipId     uint32    `json:"vip_id"`
	VipExpiry time.Time `json:"vip_expiry" gorm:"type:TIMESTAMP"`
	Platform  string    `json:"platform" gorm:"type:varchar(20)"`
	Channel   string    `json:"channel" gorm:"type:varchar(32)"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (user User) Find(db *gorm.DB) User {
	db.Where("id = ?", user.ID).Find(&user)
	return user
}

func (user User) IsVip() bool {
	return time.Now().Unix() < user.VipExpiry.Unix()
}

func (user User) SetVip(db *gorm.DB) error {
	if user.ID == 0 {
		fmt.Println(222)
		return errors.New("用户不存在")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// 设置VIP过期时间
		user.VipExpiry = user.VipExpiry.AddDate(0, 0, 30)
		tx.Save(user)
		// 增加VIP算力
		var vipAmount uint32 = 200
		amount := NewAccount(tx).GetAccount(user.ID, time.Now())
		amount.ChatAmount += vipAmount
		tx.Save(&amount)

		// 插入算力明细
		err := tx.Create(&AccountRecord{
			Uid:           user.ID,
			RecordId:      0,
			Way:           1,
			Type:          "vip",
			Amount:        vipAmount,
			CurrentAmount: amount.ChatAmount - amount.ChatUse,
		}).Error
		if err != nil {
			fmt.Println(1111)
			return err
		}
		return nil
	})
}

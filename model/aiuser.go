package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type AIUser struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `gorm:"after:id;index:ik_uid" json:"uid"`
	OpenId    string    `gorm:"type:varchar(64);uniqueIndex:uk_openid" json:"open_id"`
	UnionId   string    `gorm:"type:varchar(64)"json:"union_id"`
	Subscribe bool      `json:"subscribe"`
	JoinGroup bool      `json:"join_group" gorm:"default:0"`
	VipExpiry time.Time `json:"vip_expiry" gorm:"type:TIMESTAMP"`
	Platform  string    `json:"platform" gorm:"type:varchar(20)"`
	Channel   string    `json:"channel" gorm:"type:varchar(32)"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (user AIUser) Find(db *gorm.DB) AIUser {
	db.Where("id = ?", user.ID).Find(&user)
	return user
}

func (user AIUser) IsVip() bool {
	return time.Now().Unix() < user.VipExpiry.Unix()
}

func (user AIUser) SetVip(db *gorm.DB, vipCode *VipCode) error {
	if user.ID == 0 {
		return errors.New("用户不存在")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// 设置VIP过期时间
		if user.IsVip() {
			user.VipExpiry = user.VipExpiry.AddDate(0, 0, int(vipCode.Day))
		} else {
			user.VipExpiry, _ = time.ParseInLocation("2006-01-02 15:04:05", time.Now().AddDate(0, 0, int(vipCode.Day)).Format("2006-01-02")+" 23:59:59", time.Local)
		}
		tx.Save(user)
		// 增加VIP算力
		var vipAmount = vipCode.Vip.Amount
		if vipAmount > 0 {
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
				return err
			}
		}

		return nil
	})
}

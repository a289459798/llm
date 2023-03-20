package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	OpenId    string    `gorm:"type:varchar(64)" json:"open_id"`
	UnionId   string    `gorm:"type:varchar(64)"json:"union_id"`
	Subscribe bool      `json:"subscribe"`
	JoinGroup bool      `json:"join_group" gorm:"default:0"`
	VipExpiry time.Time `json:"vip_expiry" gorm:"type:date"`
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

package model

import (
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
	Platform  string    `json:"platform" gorm:"type:varchar(20)"`
	Channel   string    `json:"channel" gorm:"type:varchar(32)"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
	Vip       AIUserVip `gorm:"foreignKey:id;references:uid"`
}

func (user AIUser) Find(db *gorm.DB) AIUser {
	db.Where("id = ?", user.ID).Find(&user)
	return user
}

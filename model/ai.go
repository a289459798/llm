package model

import (
	"gorm.io/gorm"
	"time"
)

type AI struct {
	ID        uint32       `gorm:"primary_key" json:"id"`
	Uid       uint32       `gorm:"index:idx_uid" json:"uid"`
	Name      string       `json:"name" gorm:"type:varchar(50)" `
	Image     string       `json:"image" gorm:"type:varchar(255)" `
	Call      string       `json:"call" gorm:"type:varchar(50)" `
	RoleId    uint32       `json:"role_id"`
	Status    bool         `json:"status"`
	CreatedAt time.Time    `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time    `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
	Role      ChatTemplate `gorm:"foreignKey:id;references:role_id"`
}

func (ai AI) Info(db *gorm.DB) AI {
	db.Model(AI{}).Select("gpt_ai.*").
		Joins("inner join gpt_ai_user on gpt_ai_user.uid=gpt_ai.uid").
		Joins("inner join gpt_ai_user_vip on gpt_ai_user.uid=gpt_ai_user_vip.uid").
		Where("gpt_ai.uid = ?", ai.Uid).
		Where("gpt_ai.status = 1").
		Where("gpt_ai_user_vip.vip_expiry >=?", time.Now().Format("2006-01-02 15:04:05")).
		Preload("Role").
		First(&ai)
	return ai
}

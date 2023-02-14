package model

import "time"

type User struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	OpenId    string    `gorm:"type:varchar(64)" json:"open_id"`
	UnionId   string    `gorm:"type:varchar(64)"json:"union_id"`
	Amount    uint32    `json:"amount"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

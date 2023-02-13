package model

import "time"

type User struct {
	ID        uint      `json:"id"`
	OpenId    string    `json:"open_id"`
	UnionId   string    `json:"union_id"`
	Amount    string    `json:"amount"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

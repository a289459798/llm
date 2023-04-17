package model

import "time"

type Distributor struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `json:"uid" gorm:"uniqueIndex:uk_uid"`
	LevelId   uint32    `json:"level_id"`
	Ratio     float32   `json:"ratio" gorm:"type:decimal(10,2)"`
	Money     float32   `json:"money" gorm:"type:decimal(10,2)"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

package model

import "time"

type Pic2Pic struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `json:"uid" gorm:"index:idx_uid"`
	AppKeyId  uint32    `json:"app_key_id"`
	ImageHash string    `gorm:"type:varchar(50);index:idx_image_hash" json:"image_hash"`
	Prompt    string    `gorm:"type:varchar(100)" json:"prompt"`
	TaskId    string    `gorm:"type:varchar(50)" json:"task_id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

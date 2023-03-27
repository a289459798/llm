package model

import (
	"gorm.io/gorm"
	"time"
)

type App struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	AppKey    string    `json:"app_key" gorm:"type:varchar(32);index:ik_app_key"`
	Platform  string    `json:"platform" gorm:"type:varchar(32)"`
	Conf      string    `json:"conf"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (a App) Info(db *gorm.DB) App {
	db.Where("app_key = ?", a.AppKey).First(&a)
	return a
}

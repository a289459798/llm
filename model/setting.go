package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Setting struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Name      string    `json:"name" gorm:"type:varchar(50);uniqueIndex:uk_name"`
	Setting   string    `json:"setting"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (s Setting) Find(db *gorm.DB) (map[string]interface{}, error) {
	db.Where("name = ?", s.Name).Find(&s)
	setting := make(map[string]interface{})
	err := json.Unmarshal([]byte(s.Setting), &setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

package service

import (
	"chatgpt-tools/model"
	"gorm.io/gorm"
)

type Record struct {
	DB *gorm.DB
}

func NewRecord(db *gorm.DB) *Record {
	return &Record{
		DB: db,
	}
}

func (r *Record) Insert(record *model.Record) {
	r.DB.Create(record)
}

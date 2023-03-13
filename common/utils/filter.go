package utils

import (
	"chatgpt-tools/model"
	"fmt"
	"github.com/importcjj/sensitive"
	"gorm.io/gorm"
)

func Filter(str string, db *gorm.DB) string {
	filter := sensitive.New()
	filter.LoadWordDict("data/sensitive_words_lines.txt")
	valid, s := filter.Validate(str)
	if valid {
		return ""
	}
	db.Create(&model.Contraband{
		Content: str,
		Error:   s,
	})
	return fmt.Sprintf("违禁词：%s，请修改后提交", s)
}

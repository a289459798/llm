package utils

import (
	"fmt"
	"github.com/importcjj/sensitive"
)

func Filter(str string) string {
	filter := sensitive.New()
	filter.LoadWordDict("data/sensitive_words_lines.txt")
	valid, s := filter.Validate(str)
	if valid {
		return ""
	}
	return fmt.Sprintf("违禁词：%s，请修改后提交", s)
}

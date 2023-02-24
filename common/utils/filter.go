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
	return fmt.Sprintf("存在一下违禁词：%s\n人类的问题太深奥，三目无法回答此类问题。", s)
}

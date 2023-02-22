package utils

import (
	"github.com/importcjj/sensitive"
)

func Filter(str string) string {
	filter := sensitive.New()
	filter.LoadWordDict("data/sensitive_words_lines.txt")
	valid, _ := filter.Validate(str)
	if valid {
		return ""
	}
	return "人类的问题太深奥，三目无法回答此类问题"
}

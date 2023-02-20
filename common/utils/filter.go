package utils

import (
	"fmt"
	"github.com/importcjj/sensitive"
)

func Filter(str string) string {
	fmt.Println(str)
	filter := sensitive.New()
	filter.LoadWordDict("./data/sensitive_words_lines.txt")
	return filter.Replace(str, ' ')
}

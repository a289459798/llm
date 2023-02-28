package utils

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

func EncodeURL(s string) string {
	var result string
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.P, r) || unicode.Is(unicode.Hangul, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) || unicode.Is(unicode.Latin, r) || unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Cyrillic, r) || unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Arabic, r) || unicode.Is(unicode.Thai, r) {
			result += url.QueryEscape(string(r))
		} else {
			result += string(r)
		}
	}
	return result
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	src = re.ReplaceAllString(src, "。")
	return strings.TrimSpace(src)
}

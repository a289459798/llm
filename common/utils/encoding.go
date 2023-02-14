package utils

import (
	"net/url"
	"unicode"
)

func EncodeURL(s string) string {
	var result string
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			result += url.QueryEscape(string(r))
		} else {
			result += string(r)
		}
	}
	return result
}

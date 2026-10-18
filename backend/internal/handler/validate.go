package handler

import (
	"strings"
	"unicode/utf8"
)

const maxTextLen = 100

func cleanText(s string, max int) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" || utf8.RuneCountInString(s) > max {
		return "", false
	}
	return s, true
}

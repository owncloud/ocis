package conversions

import (
	"strings"
	"unicode/utf8"
)

// StringToSliceString splits a string into a slice string according to separator
func StringToSliceString(src string, sep string) []string {
	var parts []string
	parsed := strings.Split(src, sep)
	for _, v := range parsed {
		parts = append(parts, strings.TrimSpace(v))
	}

	return parts
}

// Reverse reverses a string
func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

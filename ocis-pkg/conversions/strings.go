package conversions

import (
	"strings"
)

// StringToSliceString splits a string into a slice string according to separator
func StringToSliceString(src string, sep string) []string {
	parsed := strings.Split(src, sep)
	parts := make([]string, 0, len(parsed))
	for _, v := range parsed {
		parts = append(parts, strings.TrimSpace(v))
	}

	return parts
}

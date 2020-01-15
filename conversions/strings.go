package conversions

import "strings"

// StringToSliceString splits a string into a slice string according to separator
func StringToSliceString(src string, sep string) []string {
	var apps []string
	parsed := strings.Split(src, sep)
	for _, v := range parsed {
		apps = append(apps, strings.TrimSpace(v))
	}

	return apps
}

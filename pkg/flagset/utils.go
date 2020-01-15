package flagset

import "strings"

// ParseAppsFlag transforms a string of the format "a, b, c" into []string{"a", "b", "c"}
func ParseAppsFlag(src string) []string {
	var apps []string
	parsed := strings.Split(src, ",")
	for _, v := range parsed {
		apps = append(apps, strings.TrimSpace(v))
	}

	return apps
}

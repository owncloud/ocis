package types

import "strings"

func SplitId(id string) (string, string) {
	ids := strings.Split(id, "$")
	if len(ids) != 2 {
		return id, ""
	}
	return ids[0], ids[1]
}

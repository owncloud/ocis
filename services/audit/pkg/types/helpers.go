package types

import "strings"

func SplitId(id string) (string, string) {
	ids := strings.Split(id, "$")
	return ids[0], ids[1]
}

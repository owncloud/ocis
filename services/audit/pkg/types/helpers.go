package types

import "strings"

// SplitStorageIDFromSpaceID splits the storage- and spaceid- from the given string
func SplitStorageIDFromSpaceID(id string) (string, string) {
	ids := strings.Split(id, "$")
	if len(ids) != 2 {
		return id, ""
	}
	return ids[0], ids[1]
}

package filepathx

import (
	"path/filepath"
	"strings"
)

// JailJoin joins any number of path elements into a single path,
// it protects against directory traversal by removing any "../" elements
// and ensuring that the path is always under the jail.
func JailJoin(jail string, elem ...string) string {
	joined := filepath.Join(append([]string{jail}, elem...)...)
	resolved, err := filepath.Abs(joined)
	if err != nil {
		return jail
	}
	jailResolved, err := filepath.Abs(jail)
	if err != nil {
		return jail
	}
	if !strings.HasPrefix(resolved, jailResolved) {
		return jail
	}
	return resolved
}

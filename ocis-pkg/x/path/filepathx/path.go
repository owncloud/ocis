package filepathx

import (
	"path/filepath"
)

// JailJoin joins any number of path elements into a single path,
// it protects against directory traversal by removing any "../" elements
// and ensuring that the path is always under the jail.
func JailJoin(jail string, elem ...string) string {
	return filepath.Join(jail, filepath.Join(append([]string{"/"}, elem...)...))
}

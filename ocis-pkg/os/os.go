package os

import (
	"os"
	"path/filepath"
)

// MustUserConfigDir generates a default config location for a user based on their OS. This location can be used to store
// any artefacts the app needs for its functioning. It is a pure function. Its only side effect is that results vary
// depending on which operative system we're in.
func MustUserConfigDir(prefix, extension string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, prefix, extension)
}

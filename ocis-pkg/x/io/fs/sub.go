package fsx

import (
	"io/fs"
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// MustSub logic is the same as fs.Sub, but it does not return an error,
// but logs and exits the process if the sub filesystem cannot be loaded.
func MustSub(fsys fs.FS, dir string) fs.FS {
	subFS, err := fs.Sub(fsys, dir)
	if err != nil {
		log.Default().Error().Err(err).Str("package", "fsx").Msg("Cannot load subtree fs")
		os.Exit(1)
	}

	return subFS
}

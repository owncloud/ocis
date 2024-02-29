package fsx

import (
	"fmt"
	"io/fs"
	"os"
)

// MustSub logic is the same as fs.Sub, but it does not return an error,
// but logs and exits the process if the sub filesystem cannot be loaded.
func MustSub(fsys fs.FS, dir string) fs.FS {
	subFS, err := fs.Sub(fsys, dir)
	if err != nil {
		fmt.Printf("unable to load subtree fs\n")
		os.Exit(1)
	}

	return subFS
}

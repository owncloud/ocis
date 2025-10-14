package trash

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// _trashGlobPattern is the glob pattern to find all trash items
	_trashGlobPattern = "spaces/*/*/trash/*/*/*/*"
)

// PurgeTrashEmptyPaths purges empty paths in the trash
func PurgeTrashEmptyPaths(p string, dryRun bool) error {
	// we have all trash nodes in all spaces now
	dirs, err := filepath.Glob(filepath.Join(p, _trashGlobPattern))
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		return errors.New("no trash found. Double check storage path")
	}

	for _, d := range dirs {
		if err := removeEmptyFolder(d, dryRun); err != nil {
			return err
		}
	}
	return nil
}

func removeEmptyFolder(path string, dryRun bool) error {
	if dryRun {
		f, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		if len(f) < 1 {
			fmt.Println("would remove", path)
		}
		return nil
	}
	if err := os.Remove(path); err != nil {
		// we do not really care about the error here
		// if the folder is not empty we will get an error,
		// this is our signal to break out of the recursion
		return nil
	}
	nd := filepath.Dir(path)
	if filepath.Base(nd) == "trash" {
		return nil
	}
	return removeEmptyFolder(nd, dryRun)
}

package trash

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

const (
	// _trashGlobPattern is the glob pattern to find all trash items
	_trashGlobPattern = "storage/users/spaces/*/*/trash/*/*/*/*/*"
)

// DataProvider provides data for the trash folders
type DataProvider struct {
	Events chan interface{}

	fsys     fs.FS
	discpath string
}

// NewDataProvider creates a new provider
func NewDataProvider(fsys fs.FS, discpath string) *DataProvider {
	return &DataProvider{
		Events: make(chan interface{}),

		fsys:     fsys,
		discpath: discpath,
	}
}

// ProduceData produces data for the trash folders
func (dp *DataProvider) ProduceData() error {
	// we have all trash nodes in all spaces now
	dirs, err := fs.Glob(dp.fsys, "storage/users/spaces/*/*/trash/*/*/*/*/*")

	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		return errors.New("no trash found. Double check storage path")
	}

	wg := sync.WaitGroup{}

	for _, l := range dirs {
		wg.Add(1)
		go func() {
			linkpath := filepath.Join(dp.discpath, l)
			r, _ := os.Readlink(linkpath)
			p := filepath.Join(l, "..", r)
			if !hasChildren(dp.fsys, p) {
				dp.Events <- TrashDirs{LinkPath: linkpath, NodePath: p}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		dp.quit()
	}()

	return nil
}

func (dp *DataProvider) quit() {
	close(dp.Events)
}

func hasChildren(fsys fs.FS, path string) bool {
	entries, err := fs.ReadDir(fsys, path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return len(entries) > 0
}

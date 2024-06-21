package trash

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// DataProvider provides data for the trash folders
type DataProvider struct {
	Events chan interface{}

	fsys     fs.FS
	discpath string
	lbs      ListBlobstore
}

// NewDataProvider creates a new provider
func NewDataProvider(fsys fs.FS, discpath string, lbs ListBlobstore) *DataProvider {
	return &DataProvider{
		Events: make(chan interface{}),

		fsys:     fsys,
		discpath: discpath,
		lbs:      lbs,
	}
}

// ProduceData produces data for the trash folders
func (dp *DataProvider) ProduceData() error {
	// we have all trash nodes in all spaces now
	// TODO: this globbing does not work as expected, probably wrong number of stars
	dirs, err := fs.Glob(dp.fsys, "storage/users/spaces/*/*/trash/*/*/*/*")

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
			p := filepath.Join(dp.discpath, l, "..", r)
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
	entries, _ := fs.ReadDir(fsys, path)
	return len(entries) > 0
}

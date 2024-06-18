package trash

import (
	"fmt"
	"io/fs"
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
	// we have all nodes in all spaces now
	dirs, err := fs.Glob(dp.fsys, "storage/users/spaces/*/*/trash/*/*/*/*")
	fmt.Printf("dirs: %v\n", dirs)
	if err != nil {
		return err
	}
	for _, d := range dirs {
		fmt.Printf("dir: %v\n", d)
	}
	return nil
}

package trash

import (
	"fmt"
	"os"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
)

// ListBlobstore required to check blob consistency
type ListBlobstore interface {
	List() ([]*node.Node, error)
	Path(node *node.Node) string
}

type TrashDirs struct {
	LinkPath string
	NodePath string
}

// NewTrashDirs creates a new trash dirs object
func NewTrashDirs() *TrashDirs {
	return &TrashDirs{}
}

// PurgeTrashOrphanedPaths purges orphaned paths in the trash
func PurgeTrashOrphanedPaths(storagepath string, lbs ListBlobstore) error {
	fsys := os.DirFS(storagepath)

	dp := NewDataProvider(fsys, storagepath, lbs)
	if err := dp.ProduceData(); err != nil {
		return err
	}

	t := NewTrashDirs()

	t.GatherData(dp.Events)
	return nil
}

// GatherData gathers data from the data provider
func (t *TrashDirs) GatherData(events <-chan interface{}) {
	for ev := range events {
		switch d := ev.(type) {
		case TrashDirs:
			// we have to remove the directory
			if err := removeDir(d.LinkPath); err != nil {
				fmt.Printf("TrashDirs: %v could not be removed: %v\n", d, err)
			}
			fmt.Printf("TrashDirs: %v is empty and needs to be removed\n", d)
		default:
			// we do not have to handle this here
			continue
		}
	}
}

func removeDir(path string) error {
	return os.RemoveAll(path)
}

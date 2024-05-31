package backup

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/shamaton/msgpack/v2"
)

// ListBlobstore required to check blob consistency
type ListBlobstore interface {
	List() ([]*node.Node, error)
	Path(node *node.Node) string
}

// DataProvider provides data for the consistency check
type DataProvider struct {
	Nodes chan NodeData
	Links chan LinkData
	Blobs chan BlobData
	Quit  chan struct{}

	fsys     fs.FS
	discpath string
	lbs      ListBlobstore
}

// NodeData holds data about the nodes
type NodeData struct {
	NodePath        string
	BlobPath        string
	RequiresSymlink bool
	Inconsistencies []Inconsistency
}

// LinkData about the symlinks
type LinkData struct {
	LinkPath string
	NodePath string
}

// BlobData about the blobs in the blobstore
type BlobData struct {
	BlobPath string
}

// NewProvider creates a new DataProvider object
func NewProvider(fsys fs.FS, discpath string, lbs ListBlobstore) *DataProvider {
	return &DataProvider{
		Nodes: make(chan NodeData),
		Links: make(chan LinkData),
		Blobs: make(chan BlobData),
		Quit:  make(chan struct{}),

		fsys:     fsys,
		discpath: discpath,
		lbs:      lbs,
	}
}

// ProduceData produces data for the consistency check
// Spawns 4 go-routines at the moment. If needed, this can be optimized.
func (dp *DataProvider) ProduceData() error {
	dirs, err := fs.Glob(dp.fsys, "spaces/*/*/nodes/*/*/*/*")
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		return errors.New("no backup found. Double check storage path")
	}

	wg := sync.WaitGroup{}

	// crawl spaces
	wg.Add(1)
	go func() {
		for _, d := range dirs {
			entries, err := fs.ReadDir(dp.fsys, d)
			if err != nil {
				fmt.Println("error reading dir", err)
				continue
			}

			if len(entries) == 0 {
				fmt.Println("empty dir", filepath.Join(dp.discpath, d))
				continue
			}

			for _, e := range entries {
				switch {
				case e.IsDir():
					ls, err := fs.ReadDir(dp.fsys, filepath.Join(d, e.Name()))
					if err != nil {
						fmt.Println("error reading dir", err)
						continue
					}
					for _, l := range ls {
						linkpath := filepath.Join(dp.discpath, d, e.Name(), l.Name())

						r, _ := os.Readlink(linkpath)
						nodePath := filepath.Join(dp.discpath, d, e.Name(), r)
						dp.Links <- LinkData{LinkPath: linkpath, NodePath: nodePath}
					}
					fallthrough
				case filepath.Ext(e.Name()) == "" || _versionRegex.MatchString(e.Name()) || _trashRegex.MatchString(e.Name()):
					np := filepath.Join(dp.discpath, d, e.Name())
					var inc []Inconsistency
					if !dp.filesExist(filepath.Join(d, e.Name())) {
						inc = append(inc, InconsistencyFilesMissing)
					}
					bp, i := dp.getBlobPath(filepath.Join(d, e.Name()))
					if i != "" {
						inc = append(inc, i)
					}

					dp.Nodes <- NodeData{NodePath: np, BlobPath: bp, RequiresSymlink: requiresSymlink(np), Inconsistencies: inc}
				}
			}
		}
		wg.Done()
	}()

	// crawl trash
	wg.Add(1)
	go func() {
		linkpaths, err := fs.Glob(dp.fsys, "spaces/*/*/trash/*/*/*/*/*")
		if err != nil {
			fmt.Println("error reading trash", err)
		}
		for _, l := range linkpaths {
			linkpath := filepath.Join(dp.discpath, l)
			r, _ := os.Readlink(linkpath)
			p := filepath.Join(dp.discpath, l, "..", r)
			dp.Links <- LinkData{LinkPath: linkpath, NodePath: p}
		}
		wg.Done()
	}()

	// crawl blobstore
	wg.Add(1)
	go func() {
		bs, err := dp.lbs.List()
		if err != nil {
			fmt.Println("error listing blobs", err)
		}

		for _, bn := range bs {
			dp.Blobs <- BlobData{BlobPath: dp.lbs.Path(bn)}
		}
		wg.Done()
	}()

	// wait for all crawlers to finish
	go func() {
		wg.Wait()
		dp.quit()
	}()

	return nil
}

func (dp *DataProvider) getBlobPath(path string) (string, Inconsistency) {
	b, err := fs.ReadFile(dp.fsys, path+".mpk")
	if err != nil {
		return "", InconsistencyFilesMissing
	}

	m := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &m); err != nil {
		return "", InconsistencyMalformedFile
	}

	// FIXME: how to check if metadata is complete?

	if bid := m["user.ocis.blobid"]; string(bid) != "" {
		spaceID, _ := getIDsFromPath(filepath.Join(dp.discpath, path))
		return dp.lbs.Path(&node.Node{BlobID: string(bid), SpaceID: spaceID}), ""
	}

	return "", ""
}

func (dp *DataProvider) filesExist(path string) bool {
	check := func(p string) bool {
		_, err := fs.Stat(dp.fsys, p)
		return err == nil
	}
	return check(path) && check(path+".mpk")
}

func (dp *DataProvider) quit() {
	dp.Quit <- struct{}{}
	close(dp.Nodes)
	close(dp.Links)
	close(dp.Blobs)
	close(dp.Quit)
}

func requiresSymlink(path string) bool {
	spaceID, nodeID := getIDsFromPath(path)
	if nodeID != "" && spaceID != "" && (spaceID == nodeID || _versionRegex.MatchString(nodeID)) {
		return false
	}

	return true
}

func getIDsFromPath(path string) (string, string) {
	rawIDs := strings.Split(path, "/nodes/")
	if len(rawIDs) != 2 {
		return "", ""
	}

	s := strings.Split(rawIDs[0], "/spaces/")
	if len(s) != 2 {
		return "", ""
	}

	spaceID := strings.Replace(s[1], "/", "", -1)
	nodeID := strings.Replace(rawIDs[1], "/", "", -1)
	return spaceID, nodeID
}

package backup

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
	fsys     fs.FS
	discpath string
	lbs      ListBlobstore
}

// NodeData holds data about the nodes
type NodeData struct {
	NodePath        string
	BlobPath        string
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
		fsys:     fsys,
		discpath: discpath,
		lbs:      lbs,
	}
}

// ProduceData produces data for the consistency check
func (c *DataProvider) ProduceData() (chan NodeData, chan LinkData, chan BlobData, chan struct{}, error) {
	dirs, err := fs.Glob(c.fsys, "spaces/*/*/nodes/*/*/*/*")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if len(dirs) == 0 {
		return nil, nil, nil, nil, errors.New("no backup found. Double check storage path")
	}

	nodes := make(chan NodeData)
	links := make(chan LinkData)
	blobs := make(chan BlobData)
	quit := make(chan struct{})
	wg := sync.WaitGroup{}

	// crawl spaces
	wg.Add(1)
	go func() {
		for _, d := range dirs {
			entries, err := fs.ReadDir(c.fsys, d)
			if err != nil {
				fmt.Println("error reading dir", err)
				continue
			}

			if len(entries) == 0 {
				fmt.Println("empty dir", filepath.Join(c.discpath, d))
				continue
			}

			for _, e := range entries {
				switch {
				case e.IsDir():
					ls, err := fs.ReadDir(c.fsys, filepath.Join(d, e.Name()))
					if err != nil {
						fmt.Println("error reading dir", err)
						continue
					}
					for _, l := range ls {
						linkpath := filepath.Join(c.discpath, d, e.Name(), l.Name())

						r, _ := os.Readlink(linkpath)
						nodePath := filepath.Join(c.discpath, d, e.Name(), r)
						links <- LinkData{LinkPath: linkpath, NodePath: nodePath}
					}
					fallthrough
				case filepath.Ext(e.Name()) == "" || _versionRegex.MatchString(e.Name()) || _trashRegex.MatchString(e.Name()):
					np := filepath.Join(c.discpath, d, e.Name())
					var inc []Inconsistency
					if !c.filesExist(filepath.Join(d, e.Name())) {
						inc = append(inc, InconsistencyFilesMissing)
					}
					bp, i := c.getBlobPath(filepath.Join(d, e.Name()))
					if i != "" {
						inc = append(inc, i)
					}

					nodes <- NodeData{NodePath: np, BlobPath: bp, Inconsistencies: inc}
				}
			}
		}
		wg.Done()
	}()

	// crawl trash
	wg.Add(1)
	go func() {
		linkpaths, err := fs.Glob(c.fsys, "spaces/*/*/trash/*/*/*/*/*")
		if err != nil {
			fmt.Println("error reading trash", err)
		}
		for _, l := range linkpaths {
			linkpath := filepath.Join(c.discpath, l)
			r, _ := os.Readlink(linkpath)
			p := filepath.Join(c.discpath, l, "..", r)
			links <- LinkData{LinkPath: linkpath, NodePath: p}
		}
		wg.Done()
	}()

	// crawl blobstore
	wg.Add(1)
	go func() {
		bs, err := c.lbs.List()
		if err != nil {
			fmt.Println("error listing blobs", err)
		}

		for _, bn := range bs {
			blobs <- BlobData{BlobPath: c.lbs.Path(bn)}
		}
		wg.Done()
	}()

	// wait for all crawlers to finish
	go func() {
		wg.Wait()
		quit <- struct{}{}
		close(nodes)
		close(links)
		close(blobs)
		close(quit)
	}()

	return nodes, links, blobs, quit, nil
}

func (c *DataProvider) getBlobPath(path string) (string, Inconsistency) {
	b, err := fs.ReadFile(c.fsys, path+".mpk")
	if err != nil {
		return "", InconsistencyFilesMissing
	}

	m := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &m); err != nil {
		return "", InconsistencyMalformedFile
	}

	if bid := m["user.ocis.blobid"]; string(bid) != "" {
		spaceID, _ := getIDsFromPath(filepath.Join(c.discpath, path))
		return c.lbs.Path(&node.Node{BlobID: string(bid), SpaceID: spaceID}), ""
	}

	return "", ""
}

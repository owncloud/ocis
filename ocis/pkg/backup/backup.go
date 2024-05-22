// Package backup contains ocis backup functionality.
package backup

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/shamaton/msgpack/v2"
)

// Inconsistency describes the type of incosistency
type Inconsistency string

var (
	InconsistencyBlobMissing     Inconsistency = "blob missing"
	InconsistencyBlobOrphaned    Inconsistency = "blob orphaned"
	InconsistencyNodeMissing     Inconsistency = "node missing"
	InconsistencyMetadataMissing Inconsistency = "metadata missing"
	InconsistencySymlinkMissing  Inconsistency = "symlink missing"
)

// Consistency holds the node and blob data of a space
type Consistency struct {
	Nodes map[string]Inconsistency
	Links map[string]Inconsistency
	Blobs map[string]Inconsistency

	fsys     fs.FS
	discpath string
}

// ListBlobstore required to check blob consistency
type ListBlobstore interface {
	Upload(node *node.Node, source string) error
	Download(node *node.Node) (io.ReadCloser, error)
	Delete(node *node.Node) error
	List() ([]string, error)
}

// New creates a new Consistency object
func New(fsys fs.FS, discpath string) *Consistency {
	return &Consistency{
		Nodes: make(map[string]Inconsistency),
		Links: make(map[string]Inconsistency),
		Blobs: make(map[string]Inconsistency),

		fsys:     fsys,
		discpath: discpath,
	}
}

// CheckSpaceConsistency checks the consistency of a space
func CheckSpaceConsistency(pathToSpace string, lbs ListBlobstore) error {
	fsys := os.DirFS(pathToSpace)

	con := New(fsys, pathToSpace)
	err := con.Initialize()
	if err != nil {
		return err
	}

	for n := range con.Nodes {
		if _, ok := con.Links[n]; !ok {
			// TODO: This is no inconsistency if
			// * this is the spaceroot
			// * this is a trashed node
			con.Nodes[n] = InconsistencySymlinkMissing
			continue
		}
		delete(con.Links, n)
		delete(con.Nodes, n)
	}

	for l := range con.Links {
		con.Links[l] = InconsistencyNodeMissing
	}

	blobs, err := lbs.List()
	if err != nil {
		return err
	}

	deadBlobs := make(map[string]Inconsistency)
	for _, b := range blobs {
		if _, ok := con.Blobs[b]; !ok {
			deadBlobs[b] = InconsistencyBlobOrphaned
			continue
		}
		delete(con.Blobs, b)
		delete(deadBlobs, b)
	}

	for b := range con.Blobs {
		con.Blobs[b] = InconsistencyBlobMissing
	}
	// get blobs from blobstore.
	// initialize blobstore (s3, local)
	// compare con.Blobs with blobstore.GetBlobs()

	for n := range con.Nodes {
		fmt.Println("Inconsistency", n, con.Nodes[n])
	}
	for l := range con.Links {
		fmt.Println("Inconsistency", l, con.Links[l])
	}
	for b := range con.Blobs {
		fmt.Println("Inconsistency", b, con.Blobs[b])
	}
	for b := range deadBlobs {
		fmt.Println("Inconsistency", b, deadBlobs[b])
	}

	return nil
}

func (c *Consistency) Initialize() error {
	dirs, err := fs.Glob(c.fsys, "nodes/*/*/*/*")
	if err != nil {
		return err
	}

	for _, d := range dirs {
		entries, err := fs.ReadDir(c.fsys, d)
		if err != nil {
			return err
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
					p, err := filepath.EvalSymlinks(filepath.Join(c.discpath, d, e.Name(), l.Name()))
					if err != nil {
						fmt.Println("error evaluating symlink", filepath.Join(d, e.Name(), l.Name()), err)
						continue
					}
					c.Links[p] = ""
				}
			case filepath.Ext(e.Name()) == ".mpk":
				inc, err := c.checkNode(filepath.Join(d, e.Name()))
				if err != nil {
					fmt.Println("error checking node", err)
					continue
				}
				c.Nodes[filepath.Join(c.discpath, d, strings.TrimSuffix(e.Name(), ".mpk"))] = inc
			default:
				// fmt.Println("unknown", e.Name())
			}
		}
	}
	return nil
}

func (c *Consistency) checkNode(path string) (Inconsistency, error) {
	b, err := fs.ReadFile(c.fsys, path)
	if err != nil {
		return "", err
	}

	m := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &m); err != nil {
		return "", err
	}

	if bid := m["user.ocis.blobid"]; string(bid) != "" {
		c.Blobs[string(bid)] = ""
	}

	return "", nil
}

func iterate(fsys fs.FS, path string, d *Consistency) ([]string, error) {
	// open symlink -> NodeMissing
	// remove node from data.Nodes
	// check blob -> BlobMissing
	// remove blob from data.Blobs
	// list children (symlinks!)
	// return children (symlinks!)
	return nil, nil
}

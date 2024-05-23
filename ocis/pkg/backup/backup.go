// Package backup contains ocis backup functionality.
package backup

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/shamaton/msgpack/v2"
)

// Inconsistency describes the type of inconsistency
type Inconsistency string

var (
	// InconsistencyBlobMissing is an inconsistency where a blob is missing in the blobstore
	InconsistencyBlobMissing Inconsistency = "blob missing"
	// InconsistencyBlobOrphaned is an inconsistency where a blob in the blobstore has no reference
	InconsistencyBlobOrphaned Inconsistency = "blob orphaned"
	// InconsistencyNodeMissing is an inconsistency where a symlink points to a non-existing node
	InconsistencyNodeMissing Inconsistency = "node missing"
	// InconsistencyMetadataMissing is an inconsistency where a node is missing metadata
	InconsistencyMetadataMissing Inconsistency = "metadata missing"
	// InconsistencySymlinkMissing is an inconsistency where a node is missing a symlink
	InconsistencySymlinkMissing Inconsistency = "symlink missing"
	// InconsistencyFilesMissing is an inconsistency where a node is missing metadata files like .mpk or .mlock
	InconsistencyFilesMissing Inconsistency = "files missing"
	// InconsistencyMalformedFile is an inconsistency where a node has a malformed metadata file
	InconsistencyMalformedFile Inconsistency = "malformed file"

	// regex to determine if a node is trashed or versioned.
	// 9113a718-8285-4b32-9042-f930f1a58ac2.REV.2024-05-22T07:32:53.89969726Z
	_versionRegex = regexp.MustCompile(`\.REV\.[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]+Z$`)
	//   9113a718-8285-4b32-9042-f930f1a58ac2.T.2024-05-23T08:25:20.006571811Z <- this HAS a symlink
	_trashRegex = regexp.MustCompile(`\.T\.[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]+Z$`)
)

// Consistency holds the node and blob data of a space
type Consistency struct {
	Nodes          map[string][]Inconsistency
	Links          map[string][]Inconsistency
	BlobReferences map[string][]Inconsistency
	Blobs          map[string][]Inconsistency

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
		Nodes:          make(map[string][]Inconsistency),
		Links:          make(map[string][]Inconsistency),
		BlobReferences: make(map[string][]Inconsistency),
		Blobs:          make(map[string][]Inconsistency),

		fsys:     fsys,
		discpath: discpath,
	}
}

// CheckSpaceConsistency checks the consistency of a space
func CheckSpaceConsistency(storagepath string, lbs ListBlobstore) error {
	fsys := os.DirFS(storagepath)

	c := New(fsys, storagepath)
	if err := c.Initialize(); err != nil {
		return err
	}

	if err := c.Evaluate(lbs); err != nil {
		return err
	}

	return c.PrintResults()
}

// Initialize initializes the Consistency object
func (c *Consistency) Initialize() error {
	dirs, err := fs.Glob(c.fsys, "spaces/*/*/nodes/*/*/*/*")
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
					p, err := os.Readlink(filepath.Join(c.discpath, d, e.Name(), l.Name()))
					if err != nil {
						fmt.Println("error reading symlink", err)
					}
					p = filepath.Join(c.discpath, d, e.Name(), p)
					c.Links[p] = []Inconsistency{}
				}
				fallthrough
			case filepath.Ext(e.Name()) == "" || _versionRegex.MatchString(e.Name()) || _trashRegex.MatchString(e.Name()):
				if !c.filesExist(filepath.Join(d, e.Name())) {
					dp := filepath.Join(c.discpath, d, e.Name())
					c.Nodes[dp] = append(c.Nodes[dp], InconsistencyFilesMissing)
				}
				inc := c.checkNode(filepath.Join(d, e.Name()+".mpk"))
				dp := filepath.Join(c.discpath, d, e.Name())
				if inc != "" {
					c.Nodes[dp] = append(c.Nodes[dp], inc)
				} else if len(c.Nodes[dp]) == 0 {
					c.Nodes[dp] = []Inconsistency{}
				}
			}
		}
	}

	links, err := fs.Glob(c.fsys, "spaces/*/*/trash/*/*/*/*/*")
	if err != nil {
		return err
	}
	for _, l := range links {
		p, err := os.Readlink(filepath.Join(c.discpath, l))
		if err != nil {
			fmt.Println("error reading symlink", err)
		}
		p = filepath.Join(c.discpath, l, "..", p)
		c.Links[p] = []Inconsistency{}
	}
	return nil
}

// Evaluate evaluates inconsistencies
func (c *Consistency) Evaluate(lbs ListBlobstore) error {
	for n := range c.Nodes {
		if _, ok := c.Links[n]; !ok && c.requiresSymlink(n) {
			c.Nodes[n] = append(c.Nodes[n], InconsistencySymlinkMissing)
			continue
		}

		deleteInconsistency(c.Links, n)
		deleteInconsistency(c.Nodes, n)
	}

	for l := range c.Links {
		c.Links[l] = append(c.Links[l], InconsistencyNodeMissing)
	}

	blobs, err := lbs.List()
	if err != nil {
		return err
	}

	for _, b := range blobs {
		if _, ok := c.BlobReferences[b]; !ok {
			c.Blobs[b] = append(c.Blobs[b], InconsistencyBlobOrphaned)
			continue
		}
		deleteInconsistency(c.BlobReferences, b)
	}

	for b := range c.BlobReferences {
		c.BlobReferences[b] = append(c.BlobReferences[b], InconsistencyBlobMissing)
	}

	return nil
}

// PrintResults prints the results of the evaluation
func (c *Consistency) PrintResults() error {
	if len(c.Nodes) != 0 {
		fmt.Println("🚨 Inconsistent Nodes:")
	}
	for n := range c.Nodes {
		fmt.Printf("\t👉️ %v\tpath: %s\n", c.Nodes[n], n)
	}
	if len(c.Links) != 0 {
		fmt.Println("🚨 Inconsistent Links:")
	}
	for l := range c.Links {
		fmt.Printf("\t👉️ %v\tpath: %s\n", c.Links[l], l)
	}
	if len(c.Blobs) != 0 {
		fmt.Println("🚨 Inconsistent Blobs:")
	}
	for b := range c.Blobs {
		fmt.Printf("\t👉️ %v\tblob: %s\n", c.Blobs[b], b)
	}
	if len(c.BlobReferences) != 0 {
		fmt.Println("🚨 Inconsistent BlobReferences:")
	}
	for b := range c.BlobReferences {
		fmt.Printf("\t👉️ %v\tblob: %s\n", c.BlobReferences[b], b)
	}
	if len(c.Nodes) == 0 && len(c.Links) == 0 && len(c.Blobs) == 0 && len(c.BlobReferences) == 0 {
		fmt.Printf("💚 No inconsistency found. The backup in '%s' seems to be valid.\n", c.discpath)
	}
	return nil

}

func (c *Consistency) checkNode(path string) Inconsistency {
	b, err := fs.ReadFile(c.fsys, path)
	if err != nil {
		return InconsistencyFilesMissing
	}

	m := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &m); err != nil {
		return InconsistencyMalformedFile
	}

	if bid := m["user.ocis.blobid"]; string(bid) != "" {
		c.BlobReferences[string(bid)] = []Inconsistency{}
	}

	return ""
}

func (c *Consistency) requiresSymlink(path string) bool {
	rawIDs := strings.Split(path, "/nodes/")
	if len(rawIDs) != 2 {
		return true
	}

	s := strings.Split(rawIDs[0], "/spaces/")
	if len(s) != 2 {
		return true
	}

	spaceID := strings.Replace(s[1], "/", "", -1)
	nodeID := strings.Replace(rawIDs[1], "/", "", -1)
	if spaceID == nodeID || _versionRegex.MatchString(nodeID) {
		return false
	}

	return true
}

func (c *Consistency) filesExist(path string) bool {
	check := func(p string) bool {
		_, err := fs.Stat(c.fsys, p)
		return err == nil
	}
	return check(path) && check(path+".mpk")
}

func deleteInconsistency(incs map[string][]Inconsistency, path string) {
	if len(incs[path]) == 0 {
		delete(incs, path)
	}
}

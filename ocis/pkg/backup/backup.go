// Package backup contains ocis backup functionality.
package backup

import (
	"errors"
	"fmt"
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

// ListBlobstore required to check blob consistency
type ListBlobstore interface {
	List() ([]*node.Node, error)
	Path(node *node.Node) string
}

// Consistency holds the node and blob data of a space
type Consistency struct {
	// Storing the data like this might take a lot of memory
	// we might need to optimize this if we run into memory issues
	Nodes          map[string][]Inconsistency
	LinkedNodes    map[string][]Inconsistency
	BlobReferences map[string][]Inconsistency
	Blobs          map[string][]Inconsistency

	nodeToLink map[string]string
	blobToNode map[string]string

	fsys     fs.FS
	discpath string
	lbs      ListBlobstore
}

// New creates a new Consistency object
func New(fsys fs.FS, discpath string, lbs ListBlobstore) *Consistency {
	return &Consistency{
		Nodes:          make(map[string][]Inconsistency),
		LinkedNodes:    make(map[string][]Inconsistency),
		BlobReferences: make(map[string][]Inconsistency),
		Blobs:          make(map[string][]Inconsistency),

		nodeToLink: make(map[string]string),
		blobToNode: make(map[string]string),

		fsys:     fsys,
		discpath: discpath,
		lbs:      lbs,
	}
}

// CheckProviderConsistency checks the consistency of a space
func CheckProviderConsistency(storagepath string, lbs ListBlobstore) error {
	fsys := os.DirFS(storagepath)

	c := New(fsys, storagepath, lbs)
	if err := c.Initialize(); err != nil {
		return err
	}

	if err := c.Evaluate(); err != nil {
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

	if len(dirs) == 0 {
		return errors.New("no backup found. Double check storage path")
	}

	for _, d := range dirs {
		entries, err := fs.ReadDir(c.fsys, d)
		if err != nil {
			return err
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
					c.LinkedNodes[nodePath] = []Inconsistency{}
					c.nodeToLink[nodePath] = linkpath
				}
				fallthrough
			case filepath.Ext(e.Name()) == "" || _versionRegex.MatchString(e.Name()) || _trashRegex.MatchString(e.Name()):
				if !c.filesExist(filepath.Join(d, e.Name())) {
					dp := filepath.Join(c.discpath, d, e.Name())
					c.Nodes[dp] = append(c.Nodes[dp], InconsistencyFilesMissing)
				}
				inc := c.checkNode(filepath.Join(d, e.Name()))
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
		linkpath := filepath.Join(c.discpath, l)
		r, _ := os.Readlink(linkpath)
		p := filepath.Join(c.discpath, l, "..", r)
		c.LinkedNodes[p] = []Inconsistency{}
		c.nodeToLink[p] = linkpath
	}
	return nil
}

// Evaluate evaluates inconsistencies
func (c *Consistency) Evaluate() error {
	for n := range c.Nodes {
		if _, ok := c.LinkedNodes[n]; !ok && c.requiresSymlink(n) {
			c.Nodes[n] = append(c.Nodes[n], InconsistencySymlinkMissing)
			continue
		}

		deleteInconsistency(c.LinkedNodes, n)
		deleteInconsistency(c.Nodes, n)
	}

	// LinkedNodes should be empty now
	for l := range c.LinkedNodes {
		c.LinkedNodes[l] = append(c.LinkedNodes[l], InconsistencyNodeMissing)
	}

	blobs, err := c.lbs.List()
	if err != nil {
		return err
	}

	for _, bn := range blobs {
		p := c.lbs.Path(bn)
		if _, ok := c.BlobReferences[p]; !ok {
			c.Blobs[p] = append(c.Blobs[p], InconsistencyBlobOrphaned)
			continue
		}
		deleteInconsistency(c.BlobReferences, p)
	}

	// BlobReferences should be empty now
	for b := range c.BlobReferences {
		c.BlobReferences[b] = append(c.BlobReferences[b], InconsistencyBlobMissing)
	}

	return nil
}

// PrintResults prints the results of the evaluation
func (c *Consistency) PrintResults() error {
	if len(c.Nodes) != 0 {
		fmt.Println("\n🚨 Inconsistent Nodes:")
	}
	for n := range c.Nodes {
		fmt.Printf("\t👉️ %v\tpath: %s\n", c.Nodes[n], n)
	}
	if len(c.LinkedNodes) != 0 {
		fmt.Println("\n🚨 Inconsistent Links:")
	}
	for l := range c.LinkedNodes {
		fmt.Printf("\t👉️ %v\tpath: %s\n\t\t\t\tmissing node:%s\n", c.LinkedNodes[l], c.nodeToLink[l], l)
	}
	if len(c.Blobs) != 0 {
		fmt.Println("\n🚨 Inconsistent Blobs:")
	}
	for b := range c.Blobs {
		fmt.Printf("\t👉️ %v\tblob: %s\n", c.Blobs[b], b)
	}
	if len(c.BlobReferences) != 0 {
		fmt.Println("\n🚨 Inconsistent BlobReferences:")
	}
	for b := range c.BlobReferences {
		fmt.Printf("\t👉️ %v\tblob: %s\n\t\t\t\treferencing node:%s\n", c.BlobReferences[b], b, c.blobToNode[b])
	}
	if len(c.Nodes) == 0 && len(c.LinkedNodes) == 0 && len(c.Blobs) == 0 && len(c.BlobReferences) == 0 {
		fmt.Printf("💚 No inconsistency found. The backup in '%s' seems to be valid.\n", c.discpath)
	}
	return nil

}

func (c *Consistency) checkNode(path string) Inconsistency {
	b, err := fs.ReadFile(c.fsys, path+".mpk")
	if err != nil {
		return InconsistencyFilesMissing
	}

	m := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &m); err != nil {
		return InconsistencyMalformedFile
	}

	if bid := m["user.ocis.blobid"]; string(bid) != "" {
		spaceID, _ := getIDsFromPath(filepath.Join(c.discpath, path))
		p := c.lbs.Path(&node.Node{BlobID: string(bid), SpaceID: spaceID})
		c.BlobReferences[p] = []Inconsistency{}
		c.blobToNode[p] = filepath.Join(c.discpath, path)
	}

	return ""
}

func (c *Consistency) requiresSymlink(path string) bool {
	spaceID, nodeID := getIDsFromPath(path)
	if nodeID != "" && spaceID != "" && (spaceID == nodeID || _versionRegex.MatchString(nodeID)) {
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

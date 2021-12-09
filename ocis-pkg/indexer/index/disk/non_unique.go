package disk

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	idxerrs "github.com/owncloud/ocis/ocis-pkg/indexer/errors"
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
	"github.com/owncloud/ocis/ocis-pkg/indexer/registry"
)

// NonUnique is able to index an document by a key which might contain non-unique values
//
// /tmp/testfiles-395764020/index.disk/PetByColor/
// ├── Brown
// │   └── rebef-123 -> /tmp/testfiles-395764020/pets/rebef-123
// ├── Green
// │    ├── goefe-789 -> /tmp/testfiles-395764020/pets/goefe-789
// │    └── xadaf-189 -> /tmp/testfiles-395764020/pets/xadaf-189
// └── White
//     └── wefwe-456 -> /tmp/testfiles-395764020/pets/wefwe-456
type NonUnique struct {
	caseInsensitive bool
	indexBy         string
	typeName        string
	filesDir        string
	indexBaseDir    string
	indexRootDir    string
}

func init() {
	registry.IndexConstructorRegistry["disk"]["non_unique"] = NewNonUniqueIndexWithOptions
}

// NewNonUniqueIndexWithOptions instantiates a new UniqueIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewNonUniqueIndexWithOptions(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	return &NonUnique{
		caseInsensitive: opts.CaseInsensitive,
		indexBy:         opts.IndexBy,
		typeName:        opts.TypeName,
		filesDir:        opts.FilesDir,
		indexBaseDir:    path.Join(opts.DataDir, "index.disk"),
		indexRootDir:    path.Join(path.Join(opts.DataDir, "index.disk"), strings.Join([]string{"non_unique", opts.TypeName, opts.IndexBy}, ".")),
	}
}

// Init initializes a unique index.
func (idx *NonUnique) Init() error {
	if _, err := os.Stat(idx.filesDir); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.indexRootDir, 0777); err != nil {
		return err
	}

	return nil
}

// Lookup exact lookup by value.
func (idx *NonUnique) Lookup(v string) ([]string, error) {
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	searchPath := path.Join(idx.indexRootDir, v)
	fi, err := ioutil.ReadDir(searchPath)
	if os.IsNotExist(err) {
		return []string{}, &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	if err != nil {
		return []string{}, err
	}

	ids := make([]string, 0, len(fi))
	for _, f := range fi {
		ids = append(ids, f.Name())
	}

	if len(ids) == 0 {
		return []string{}, &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	return ids, nil
}

// Add adds a value to the index, returns the path to the root-document
func (idx *NonUnique) Add(id, v string) (string, error) {
	if v == "" {
		return "", nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	oldName := path.Join(idx.filesDir, id)
	newName := path.Join(idx.indexRootDir, v, id)

	if err := os.MkdirAll(path.Join(idx.indexRootDir, v), 0777); err != nil {
		return "", err
	}

	err := os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	return newName, err

}

// Remove a value v from an index.
func (idx *NonUnique) Remove(id string, v string) error {
	if v == "" {
		return nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	res, err := filepath.Glob(path.Join(idx.indexRootDir, "/*/", id))
	if err != nil {
		return err
	}

	for _, p := range res {
		if err := os.Remove(p); err != nil {
			return err
		}
	}

	// Remove value directory if it is empty
	valueDir := path.Join(idx.indexRootDir, v)
	fi, err := ioutil.ReadDir(valueDir)
	if err != nil {
		return err
	}

	if len(fi) == 0 {
		if err := os.RemoveAll(valueDir); err != nil {
			return err
		}
	}

	return nil
}

// Update index from <oldV> to <newV>.
func (idx *NonUnique) Update(id, oldV, newV string) (err error) {
	if idx.caseInsensitive {
		oldV = strings.ToLower(oldV)
		newV = strings.ToLower(newV)
	}
	oldDir := path.Join(idx.indexRootDir, oldV)
	oldPath := path.Join(oldDir, id)
	newDir := path.Join(idx.indexRootDir, newV)
	newPath := path.Join(newDir, id)

	if _, err = os.Stat(oldPath); os.IsNotExist(err) {
		return &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: oldV}
	}

	if err != nil {
		return
	}

	if err = os.MkdirAll(newDir, 0777); err != nil {
		return
	}

	if err = os.Rename(oldPath, newPath); err != nil {
		return
	}

	di, err := ioutil.ReadDir(oldDir)
	if err != nil {
		return err
	}

	if len(di) == 0 {
		err = os.RemoveAll(oldDir)
		if err != nil {
			return
		}
	}

	return

}

// Search allows for glob search on the index.
func (idx *NonUnique) Search(pattern string) ([]string, error) {
	if idx.caseInsensitive {
		pattern = strings.ToLower(pattern)
	}
	paths, err := filepath.Glob(path.Join(idx.indexRootDir, pattern, "*"))
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: pattern}
	}

	return paths, nil
}

// CaseInsensitive undocumented.
func (idx *NonUnique) CaseInsensitive() bool {
	return idx.caseInsensitive
}

// IndexBy undocumented.
func (idx *NonUnique) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
func (idx *NonUnique) TypeName() string {
	return idx.typeName
}

// FilesDir undocumented.
func (idx *NonUnique) FilesDir() string {
	return idx.filesDir
}

// Delete deletes the index folder from its storage.
func (idx *NonUnique) Delete() error {
	return os.RemoveAll(idx.indexRootDir)
}

package disk

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	idxerrs "github.com/owncloud/ocis/accounts/pkg/indexer/errors"
	"github.com/owncloud/ocis/accounts/pkg/indexer/index"
	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	"github.com/owncloud/ocis/accounts/pkg/indexer/registry"
)

// Unique ensures that only one document of the same type and key-value combination can exist in the index.
//
// Modeled by creating a indexer-folder per entity and key with symlinks which point to respective documents which contain
// the link-filename as value.
//
// Directory Layout
//
// 		/var/data/index.disk/UniqueUserByEmail/
// 		├── jacky@example.com -> /var/data/users/ewf4ofk-555
// 		├── jones@example.com -> /var/data/users/rulan54-777
// 		└── mikey@example.com -> /var/data/users/abcdefg-123
//
// Example user
//
// 		{
//  		"Id": "ewf4ofk-555",
//  		"UserName": "jacky",
//  		"Email": "jacky@example.com"
// 		}
//
type Unique struct {
	caseInsensitive bool
	indexBy         string
	typeName        string
	filesDir        string
	indexBaseDir    string
	indexRootDir    string
}

func init() {
	registry.IndexConstructorRegistry["disk"]["unique"] = NewUniqueIndexWithOptions
}

// NewUniqueIndexWithOptions instantiates a new UniqueIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewUniqueIndexWithOptions(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	return &Unique{
		caseInsensitive: opts.CaseInsensitive,
		indexBy:         opts.IndexBy,
		typeName:        opts.TypeName,
		filesDir:        opts.FilesDir,
		indexBaseDir:    path.Join(opts.DataDir, "index.disk"),
		indexRootDir:    path.Join(path.Join(opts.DataDir, "index.disk"), strings.Join([]string{"unique", opts.TypeName, opts.IndexBy}, ".")),
	}
}

// Init initializes a unique index.
func (idx *Unique) Init() error {
	if _, err := os.Stat(idx.filesDir); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.indexRootDir, 0777); err != nil {
		return err
	}

	return nil
}

// Add adds a value to the index, returns the path to the root-document
func (idx *Unique) Add(id, v string) (string, error) {
	if v == "" {
		return "", nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	oldName := path.Join(idx.filesDir, id)
	newName := path.Join(idx.indexRootDir, v)
	err := os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	return newName, err
}

// Remove a value v from an index.
func (idx *Unique) Remove(id string, v string) (err error) {
	if v == "" {
		return nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	searchPath := path.Join(idx.indexRootDir, v)
	return os.Remove(searchPath)
}

// Lookup exact lookup by value.
func (idx *Unique) Lookup(v string) (resultPath []string, err error) {
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	searchPath := path.Join(idx.indexRootDir, v)
	if err = isValidSymlink(searchPath); err != nil {
		if os.IsNotExist(err) {
			err = &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}
		return
	}

	p, err := os.Readlink(searchPath)
	if err != nil {
		return []string{}, nil
	}

	return []string{p}, err
}

// Update index from <oldV> to <newV>.
func (idx *Unique) Update(id, oldV, newV string) (err error) {
	if idx.caseInsensitive {
		oldV = strings.ToLower(oldV)
		newV = strings.ToLower(newV)
	}
	oldPath := path.Join(idx.indexRootDir, oldV)
	if err = isValidSymlink(oldPath); err != nil {
		if os.IsNotExist(err) {
			return &idxerrs.NotFoundErr{TypeName: idx.TypeName(), Key: idx.IndexBy(), Value: oldV}
		}

		return
	}

	newPath := path.Join(idx.indexRootDir, newV)
	if err = isValidSymlink(newPath); err == nil {
		return &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: newV}
	}

	if os.IsNotExist(err) {
		err = os.Rename(oldPath, newPath)
	}

	return
}

// Search allows for glob search on the index.
func (idx *Unique) Search(pattern string) ([]string, error) {
	if idx.caseInsensitive {
		pattern = strings.ToLower(pattern)
	}
	paths, err := filepath.Glob(path.Join(idx.indexRootDir, pattern))
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: pattern}
	}

	res := make([]string, 0)
	for _, p := range paths {
		if err := isValidSymlink(p); err != nil {
			return nil, err
		}

		src, err := os.Readlink(p)
		if err != nil {
			return nil, err
		}

		res = append(res, src)
	}

	return res, nil
}

// CaseInsensitive undocumented.
func (idx *Unique) CaseInsensitive() bool {
	return idx.caseInsensitive
}

// IndexBy undocumented.
func (idx *Unique) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
func (idx *Unique) TypeName() string {
	return idx.typeName
}

// FilesDir undocumented.
func (idx *Unique) FilesDir() string {
	return idx.filesDir
}

func isValidSymlink(path string) (err error) {
	var symInfo os.FileInfo
	if symInfo, err = os.Lstat(path); err != nil {
		return
	}

	if symInfo.Mode()&os.ModeSymlink == 0 {
		err = fmt.Errorf("%s is not a valid symlink (bug/corruption?)", path)
		return
	}

	return
}

// Delete deletes the index folder from its storage.
func (idx *Unique) Delete() error {
	return os.Remove(idx.indexRootDir)
}

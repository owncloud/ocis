package disk

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	idxerrs "github.com/owncloud/ocis/accounts/pkg/indexer/errors"

	"github.com/owncloud/ocis/accounts/pkg/indexer/index"
	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	"github.com/owncloud/ocis/accounts/pkg/indexer/registry"
)

// Autoincrement are fields for an index of type autoincrement.
type Autoincrement struct {
	indexBy      string
	typeName     string
	filesDir     string
	indexBaseDir string
	indexRootDir string

	bound *option.Bound
}

// - Creating an autoincrement index has to be thread safe.
// - Validation: autoincrement indexes should only work on integers.

func init() {
	registry.IndexConstructorRegistry["disk"]["autoincrement"] = NewAutoincrementIndex
}

// NewAutoincrementIndex instantiates a new AutoincrementIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewAutoincrementIndex(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	// validate the field
	if opts.Entity == nil {
		panic("invalid autoincrement index: configured without entity")
	}

	k, err := getKind(opts.Entity, opts.IndexBy)
	if !isValidKind(k) || err != nil {
		panic("invalid autoincrement index: configured on non-numeric field")
	}

	return &Autoincrement{
		indexBy:      opts.IndexBy,
		typeName:     opts.TypeName,
		filesDir:     opts.FilesDir,
		bound:        opts.Bound,
		indexBaseDir: path.Join(opts.DataDir, "index.disk"),
		indexRootDir: path.Join(path.Join(opts.DataDir, "index.disk"), strings.Join([]string{"autoincrement", opts.TypeName, opts.IndexBy}, ".")),
	}
}

var (
	validKinds = []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	}
)

// Init initializes an autoincrement index.
func (idx Autoincrement) Init() error {
	if _, err := os.Stat(idx.filesDir); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.indexRootDir, 0777); err != nil {
		return err
	}

	return nil
}

// Lookup exact lookup by value.
func (idx Autoincrement) Lookup(v string) ([]string, error) {
	searchPath := path.Join(idx.indexRootDir, v)
	if err := isValidSymlink(searchPath); err != nil {
		if os.IsNotExist(err) {
			err = &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}

		return nil, err
	}

	p, err := os.Readlink(searchPath)
	if err != nil {
		return []string{}, nil
	}

	return []string{p}, err
}

// Add a new value to the index.
func (idx Autoincrement) Add(id, v string) (string, error) {
	nextID, err := idx.next()
	if err != nil {
		return "", err
	}
	oldName := filepath.Join(idx.filesDir, id)
	var newName string
	if v == "" {
		newName = filepath.Join(idx.indexRootDir, strconv.Itoa(nextID))
	} else {
		newName = filepath.Join(idx.indexRootDir, v)
	}
	err = os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	return newName, err
}

// Remove a value v from an index.
func (idx Autoincrement) Remove(id string, v string) error {
	if v == "" {
		return nil
	}
	searchPath := path.Join(idx.indexRootDir, v)
	return os.Remove(searchPath)
}

// Update index from <oldV> to <newV>.
func (idx Autoincrement) Update(id, oldV, newV string) error {
	oldPath := path.Join(idx.indexRootDir, oldV)
	if err := isValidSymlink(oldPath); err != nil {
		if os.IsNotExist(err) {
			return &idxerrs.NotFoundErr{TypeName: idx.TypeName(), Key: idx.IndexBy(), Value: oldV}
		}

		return err
	}

	newPath := path.Join(idx.indexRootDir, newV)
	err := isValidSymlink(newPath)
	if err == nil {
		return &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: newV}
	}

	if os.IsNotExist(err) {
		err = os.Rename(oldPath, newPath)
	}

	return err
}

// Search allows for glob search on the index.
func (idx Autoincrement) Search(pattern string) ([]string, error) {
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

// IndexBy undocumented.
func (idx Autoincrement) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
func (idx Autoincrement) TypeName() string {
	return idx.typeName
}

// FilesDir  undocumented.
func (idx Autoincrement) FilesDir() string {
	return idx.filesDir
}

func isValidKind(k reflect.Kind) bool {
	for _, v := range validKinds {
		if k == v {
			return true
		}
	}
	return false
}

func getKind(i interface{}, field string) (reflect.Kind, error) {
	r := reflect.ValueOf(i)
	return reflect.Indirect(r).FieldByName(field).Kind(), nil
}

func readDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool {
		a, _ := strconv.Atoi(list[i].Name())
		b, _ := strconv.Atoi(list[j].Name())
		return a < b
	})
	return list, nil
}

func (idx Autoincrement) next() (int, error) {
	files, err := readDir(idx.indexRootDir)
	if err != nil {
		return -1, err
	}

	if len(files) == 0 {
		return int(idx.bound.Lower), nil
	}

	latest, err := strconv.Atoi(path.Base(files[len(files)-1].Name()))
	if err != nil {
		return -1, err
	}

	if int64(latest) < idx.bound.Lower {
		return int(idx.bound.Lower), nil
	}

	return latest + 1, nil
}

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

type Autoincrement struct {
	indexBy      string
	typeName     string
	filesDir     string
	indexBaseDir string
	indexRootDir string
	entity       interface{}
}

// - Creating an autoincrement index has to be thread safe.
// - Validation: autoincrement indexes should only work on integers.

func init() {
	registry.IndexConstructorRegistry["disk"]["autoincrement"] = NewAutoincrementIndex
}

// NewAutoincrementIndex instantiates a new UniqueIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewAutoincrementIndex(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	// validate the field
	if opts.Entity == nil {
		// return error: entity needed for field validation
	}

	k, err := getKind(opts.Entity, opts.IndexBy)
	if !isValidKind(k) || err != nil {
		panic("invalid index in non-numeric field")
	}

	return &Autoincrement{
		indexBy:      opts.IndexBy,
		typeName:     opts.TypeName,
		filesDir:     opts.FilesDir,
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
	}
)

func (idx Autoincrement) Init() error {
	if _, err := os.Stat(idx.filesDir); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.indexRootDir, 0777); err != nil {
		return err
	}

	return nil
}

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

func (idx Autoincrement) Add(id, v string) (string, error) {
	oldName := filepath.Join(idx.filesDir, id)
	newName := filepath.Join(idx.indexRootDir, strconv.Itoa(idx.next()))
	err := os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	return newName, err
}

func (idx Autoincrement) Remove(id string, v string) error {
	panic("implement me")
}

func (idx Autoincrement) Update(id, oldV, newV string) error {
	panic("implement me")
}

func (idx Autoincrement) Search(pattern string) ([]string, error) {
	panic("implement me")
}

func (idx Autoincrement) IndexBy() string {
	panic("implement me")
}

func (idx Autoincrement) TypeName() string {
	panic("implement me")
}

func (idx Autoincrement) FilesDir() string {
	panic("implement me")
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
	// TODO reflect.FieldByName panics. Recover from it.
	// further read: https://blog.golang.org/defer-panic-and-recover
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

func (idx Autoincrement) next() int {
	files, err := readDir(idx.indexRootDir)
	if err != nil {
		// hello handle me pls.
	}

	if len(files) == 0 {
		return 0
	}

	latest, err := strconv.Atoi(path.Base(files[len(files)-1].Name())) // would returning a string be a better interface?
	if err != nil {
		// handle me daddy
	}
	return latest + 1
}

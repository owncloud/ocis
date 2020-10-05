package disk

import (
	"errors"
	"fmt"
	idxerrs "github.com/owncloud/ocis/accounts/pkg/indexer/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	indexBy      string
	typeName     string
	filesDir     string
	indexBaseDir string
	indexRootDir string
}

// NewUniqueIndex instantiates a new UniqueIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewUniqueIndex(typeName, indexBy, filesDir, indexBaseDir string) Unique {
	return Unique{
		indexBy:      indexBy,
		typeName:     typeName,
		filesDir:     filesDir,
		indexBaseDir: indexBaseDir,
		indexRootDir: path.Join(indexBaseDir, strings.Join([]string{"unique", typeName, indexBy}, ".")),
	}
}

func (idx Unique) Init() error {
	if _, err := os.Stat(idx.filesDir); err != nil {
		return err
	}

	if err := os.MkdirAll(idx.indexRootDir, 0777); err != nil {
		return err
	}

	return nil
}

func (idx Unique) Add(id, v string) (string, error) {
	oldName := path.Join(idx.filesDir, id)
	newName := path.Join(idx.indexRootDir, v)
	err := os.Symlink(oldName, newName)
	if errors.Is(err, os.ErrExist) {
		return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	return newName, err
}

func (idx Unique) Remove(id string, v string) (err error) {
	searchPath := path.Join(idx.indexRootDir, v)
	if err = isValidSymlink(searchPath); err != nil {
		return
	}

	return os.Remove(searchPath)
}

// unique.github.com.owncloud.ocis.accounts.pkg.indexer.User.UserName
// unique.github.com.owncloud.ocis.accounts.pkg.indexer.User.UserName/UserName
func (idx Unique) Lookup(v string) (resultPath []string, err error) {
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

func (idx Unique) Update(id, oldV, newV string) (err error) {
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

func (idx Unique) Search(pattern string) ([]string, error) {
	paths, err := filepath.Glob(path.Join(idx.indexRootDir, pattern))
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: pattern}
	}

	res := make([]string, 0, 0)
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

func (idx Unique) IndexBy() string {
	return idx.indexBy
}

func (idx Unique) TypeName() string {
	return idx.typeName
}

func (idx Unique) FilesDir() string {
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

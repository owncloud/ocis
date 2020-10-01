// Package index provides symlink-based index for on-disk document-directories.
package index

import (
	"github.com/rs/zerolog"
	"path"
)

// Index is a facade to configure and query over multiple indices.
type Index struct {
	config  *Config
	indices indexMap
}

type Config struct {
	DataDir          string
	IndexRootDirName string
	Log              zerolog.Logger
}

// IndexType can be implemented to create new index-strategies. See Unique for example.
// Each index implementation is bound to one data-column (IndexBy) and a data-type (TypeName)
type IndexType interface {
	Init() error
	Lookup(v string) ([]string, error)
	Add(id, v string) (string, error)
	Remove(id string, v string) error
	Update(id, oldV, newV string) error
	Search(pattern string) ([]string, error)
	IndexBy() string
	TypeName() string
	FilesDir() string
}

func NewIndex(cfg *Config) *Index {
	return &Index{
		config:  cfg,
		indices: indexMap{},
	}
}


func (i Index) AddUniqueIndex(t interface{}, indexBy, pkName, entityDirName string) error {
	typeName := getTypeFQN(t)
	fullDataPath := path.Join(i.config.DataDir, entityDirName)
	indexPath := path.Join(i.config.DataDir, i.config.IndexRootDirName)

	idx := NewUniqueIndex(typeName, indexBy, fullDataPath, indexPath)

	i.indices.addIndex(typeName, pkName, idx)
	return idx.Init()
}

func (i Index) AddNonUniqueIndex(t interface{}, indexBy, pkName, entityDirName string) error {
	typeName := getTypeFQN(t)
	fullDataPath := path.Join(i.config.DataDir, entityDirName)
	indexPath := path.Join(i.config.DataDir, i.config.IndexRootDirName)

	idx := NewNonUniqueIndex(typeName, indexBy, fullDataPath, indexPath)

	i.indices.addIndex(typeName, pkName, idx)
	return idx.Init()
}


// Add a new entry to the index
func (i Index) Add(t interface{}) error {
	typeName := getTypeFQN(t)

	val, ok := i.indices[typeName]
	if ok {

	}



	return nil


}

// Find a entry by type,field and value.
//  // Find a User type by email
//  man.Find("User", "Email", "foo@example.com")
func (i Index) Find(typeName, key, value string) (pk string, err error) {
	var res = []string{}
	if indices, ok := i.indices[typeName][key]; ok {
		for _, idx := range indices {
			if res, err = idx.Lookup(value); IsNotFoundErr(err) {
				continue
			}

			if err != nil {
				return
			}
		}
	}

	if len(res) == 0 {
		return "", err
	}

	return path.Base(res[0]), err
}

func (i Index) Delete(typeName, pk string) error {

	return nil
}

 */

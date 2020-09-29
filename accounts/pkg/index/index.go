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

// Type can be implemented to create new index-strategies. See Unique for example.
// Each index implementation is bound to one data-column (IndexBy) and a data-type (TypeName)
type Type interface {
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

func (man Index) AddUniqueIndex(typeName, indexBy, entityDirName string) error {
	fullDataPath := path.Join(man.config.DataDir, entityDirName)
	indexPath := path.Join(man.config.DataDir, man.config.IndexRootDirName)

	idx := NewUniqueIndex(typeName, indexBy, fullDataPath, indexPath)
	man.indices.addIndex(idx)

	return idx.Init()
}

func (man Index) AddNormalIndex(typeName, indexBy, entityDirName string) error {
	fullDataPath := path.Join(man.config.DataDir, entityDirName)
	indexPath := path.Join(man.config.DataDir, man.config.IndexRootDirName)

	idx := NewNormalIndex(typeName, indexBy, fullDataPath, indexPath)
	man.indices.addIndex(idx)

	return idx.Init()
}

func (man Index) AddIndex(idx Type) error {
	man.indices.addIndex(idx)
	return idx.Init()
}

// Add a new entry to the index
func (man Index) Add(primaryKey string, entity interface{}) error {
	t, err := getType(entity)
	if err != nil {
		return err
	}

	typeName := t.Type().Name()

	if typeIndices, ok := man.indices[typeName]; ok {
		for _, fieldIndices := range typeIndices {
			for k := range fieldIndices {
				curIdx := fieldIndices[k]
				idxBy := curIdx.IndexBy()
				val := valueOf(entity, idxBy)
				_, err := curIdx.Add(primaryKey, val)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Find a entry by type,field and value.
//  // Find a User type by email
//  man.Find("User", "Email", "foo@example.com")
func (man Index) Find(typeName, key, value string) (pk string, err error) {
	var res = []string{}
	if indices, ok := man.indices[typeName][key]; ok {
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

func (man Index) Delete(typeName, pk string) error {

	return nil
}

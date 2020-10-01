// Package indexer provides symlink-based indexer for on-disk document-directories.
package indexer

import (
	"github.com/rs/zerolog"
	"path"
)

// Indexer is a facade to configure and query over multiple indices.
type Indexer struct {
	config  *Config
	indices indexMap
}

type Config struct {
	DataDir          string
	IndexRootDirName string
	Log              zerolog.Logger
}

// IndexType can be implemented to create new indexer-strategies. See Unique for example.
// Each indexer implementation is bound to one data-column (IndexBy) and a data-type (TypeName)
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

func NewIndex(cfg *Config) *Indexer {
	return &Indexer{
		config:  cfg,
		indices: indexMap{},
	}
}

func (i Indexer) AddUniqueIndex(t interface{}, indexBy, pkName, entityDirName string) error {
	typeName := getTypeFQN(t)
	fullDataPath := path.Join(i.config.DataDir, entityDirName)
	indexPath := path.Join(i.config.DataDir, i.config.IndexRootDirName)

	idx := NewUniqueIndex(typeName, indexBy, fullDataPath, indexPath)

	i.indices.addIndex(typeName, pkName, idx)
	return idx.Init()
}

func (i Indexer) AddNonUniqueIndex(t interface{}, indexBy, pkName, entityDirName string) error {
	typeName := getTypeFQN(t)
	fullDataPath := path.Join(i.config.DataDir, entityDirName)
	indexPath := path.Join(i.config.DataDir, i.config.IndexRootDirName)

	idx := NewNonUniqueIndex(typeName, indexBy, fullDataPath, indexPath)

	i.indices.addIndex(typeName, pkName, idx)
	return idx.Init()
}

// Add a new entry to the indexer
func (i Indexer) Add(t interface{}) error {
	typeName := getTypeFQN(t)
	if fields, ok := i.indices[typeName]; ok {
		for _, indices := range fields.indicesByField {
			for _, idx := range indices {
				pkVal := valueOf(t, fields.pKFieldName)
				idxByVal := valueOf(t, idx.IndexBy())
				if _, err := idx.Add(pkVal, idxByVal); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (i Indexer) FindBy(t interface{}, field string, val string) ([]string, error) {
	typeName := getTypeFQN(t)
	resultPaths := make([]string, 0)
	if fields, ok := i.indices[typeName]; ok {
		for _, idx := range fields.indicesByField[field] {
			res, err := idx.Lookup(val)
			if err != nil {
				if IsNotFoundErr(err) {
					continue
				}

				if err != nil {
					return nil, err
				}
			}

			resultPaths = append(resultPaths, res...)

		}
	}

	result := make([]string, 0, len(resultPaths))
	for _, v := range resultPaths {
		result = append(result, path.Base(v))
	}

	return result, nil
}

func (i Indexer) Delete(t interface{}) error {
	typeName := getTypeFQN(t)
	if fields, ok := i.indices[typeName]; ok {
		for _, indices := range fields.indicesByField {
			for _, idx := range indices {
				pkVal := valueOf(t, fields.pKFieldName)
				idxByVal := valueOf(t, idx.IndexBy())
				if err := idx.Remove(pkVal, idxByVal); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (i Indexer) FindByPartial(t interface{}, field string, pattern string) ([]string, error) {
	typeName := getTypeFQN(t)
	resultPaths := make([]string, 0)
	if fields, ok := i.indices[typeName]; ok {
		for _, idx := range fields.indicesByField[field] {
			res, err := idx.Search(pattern)
			if err != nil {
				if IsNotFoundErr(err) {
					continue
				}

				if err != nil {
					return nil, err
				}
			}

			resultPaths = append(resultPaths, res...)

		}
	}

	result := make([]string, 0, len(resultPaths))
	for _, v := range resultPaths {
		result = append(result, path.Base(v))
	}

	return result, nil

}

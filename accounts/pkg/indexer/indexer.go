// Package indexer provides symlink-based indexer for on-disk document-directories.
package indexer

import (
	"fmt"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/indexer/errors"
	"github.com/owncloud/ocis/accounts/pkg/indexer/index"
	_ "github.com/owncloud/ocis/accounts/pkg/indexer/index/cs3"  // to populate index
	_ "github.com/owncloud/ocis/accounts/pkg/indexer/index/disk" // to populate index
	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	"github.com/owncloud/ocis/accounts/pkg/indexer/registry"
)

// Indexer is a facade to configure and query over multiple indices.
type Indexer struct {
	config  *config.Config
	indices typeMap
}

// IdxAddResult represents the result of an Add call on an index
type IdxAddResult struct {
	Field, Value string
}

// CreateIndexer creates a new Indexer.
func CreateIndexer(cfg *config.Config) *Indexer {
	return &Indexer{
		config:  cfg,
		indices: typeMap{},
	}
}

func getRegistryStrategy(cfg *config.Config) string {
	if cfg.Repo.Disk.Path != "" {
		return "disk"
	}

	return "cs3"
}

func (i Indexer) Reset() error {
	for k := range i.indices {
		delete(i.indices, k)
	}

	// TODO: delete indexes from storage (cs3 / disk)

	return nil
}

// AddIndex adds a new index to the indexer receiver.
func (i Indexer) AddIndex(t interface{}, indexBy, pkName, entityDirName, indexType string, bound *option.Bound, caseInsensitive bool) error {
	strategy := getRegistryStrategy(i.config)
	f := registry.IndexConstructorRegistry[strategy][indexType]
	var idx index.Index

	if strategy == "disk" {
		idx = f(
			option.CaseInsensitive(caseInsensitive),
			option.WithEntity(t),
			option.WithBounds(bound),
			option.WithTypeName(getTypeFQN(t)),
			option.WithIndexBy(indexBy),
			option.WithFilesDir(path.Join(i.config.Repo.Disk.Path, entityDirName)),
			option.WithDataDir(i.config.Repo.Disk.Path),
		)
	} else if strategy == "cs3" {
		idx = f(
			option.CaseInsensitive(caseInsensitive),
			option.WithEntity(t),
			option.WithBounds(bound),
			option.WithTypeName(getTypeFQN(t)),
			option.WithIndexBy(indexBy),
			option.WithDataURL(i.config.Repo.CS3.DataURL),
			option.WithDataPrefix(i.config.Repo.CS3.DataPrefix),
			option.WithJWTSecret(i.config.Repo.CS3.JWTSecret),
			option.WithProviderAddr(i.config.Repo.CS3.ProviderAddr),
			option.WithServiceUserUUID(i.config.ServiceUser.UUID),
			option.WithServiceUserName(i.config.ServiceUser.Username),
		)
	}

	i.indices.addIndex(getTypeFQN(t), pkName, idx)
	return idx.Init()
}

// Add a new entry to the indexer
func (i Indexer) Add(t interface{}) ([]IdxAddResult, error) {
	typeName := getTypeFQN(t)
	var results []IdxAddResult
	if fields, ok := i.indices[typeName]; ok {
		for _, indices := range fields.IndicesByField {
			for _, idx := range indices {
				pkVal := valueOf(t, fields.PKFieldName)
				idxByVal := valueOf(t, idx.IndexBy())
				value, err := idx.Add(pkVal, idxByVal)
				if err != nil {
					return []IdxAddResult{}, err
				}
				if value == "" {
					continue
				}
				results = append(results, IdxAddResult{Field: idx.IndexBy(), Value: value})
			}
		}
	}

	return results, nil
}

// FindBy finds a value on an index by field and value.
func (i Indexer) FindBy(t interface{}, field string, val string) ([]string, error) {
	typeName := getTypeFQN(t)
	resultPaths := make([]string, 0)
	if fields, ok := i.indices[typeName]; ok {
		for _, idx := range fields.IndicesByField[strcase.ToCamel(field)] {
			idxVal := val
			res, err := idx.Lookup(idxVal)
			if err != nil {
				if errors.IsNotFoundErr(err) {
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

// Delete deletes all indexed fields of a given type t on the Indexer.
func (i Indexer) Delete(t interface{}) error {
	typeName := getTypeFQN(t)
	if fields, ok := i.indices[typeName]; ok {
		for _, indices := range fields.IndicesByField {
			for _, idx := range indices {
				pkVal := valueOf(t, fields.PKFieldName)
				idxByVal := valueOf(t, idx.IndexBy())
				if err := idx.Remove(pkVal, idxByVal); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// FindByPartial allows for glob search across all indexes.
func (i Indexer) FindByPartial(t interface{}, field string, pattern string) ([]string, error) {
	typeName := getTypeFQN(t)
	resultPaths := make([]string, 0)
	if fields, ok := i.indices[typeName]; ok {
		for _, idx := range fields.IndicesByField[strcase.ToCamel(field)] {
			res, err := idx.Search(pattern)
			if err != nil {
				if errors.IsNotFoundErr(err) {
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

// Update updates all indexes on a value <from> to a value <to>.
func (i Indexer) Update(from, to interface{}) error {
	typeNameFrom := getTypeFQN(from)
	typeNameTo := getTypeFQN(to)
	if typeNameFrom != typeNameTo {
		return fmt.Errorf("update types do not match: from %v to %v", typeNameFrom, typeNameTo)
	}

	if fields, ok := i.indices[typeNameFrom]; ok {
		for fName, indices := range fields.IndicesByField {
			oldV := valueOf(from, fName)
			newV := valueOf(to, fName)
			pkVal := valueOf(from, fields.PKFieldName)
			for _, idx := range indices {
				if oldV == newV {
					continue
				}
				if oldV == "" {
					if _, err := idx.Add(pkVal, newV); err != nil {
						return err
					}
					continue
				}
				if newV == "" {
					if err := idx.Remove(pkVal, oldV); err != nil {
						return err
					}
					continue
				}
				if err := idx.Update(pkVal, oldV, newV); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

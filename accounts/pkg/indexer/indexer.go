// Package indexer provides symlink-based indexer for on-disk document-directories.
package indexer

import (
	"fmt"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/indexer/errors"
	"github.com/owncloud/ocis/accounts/pkg/indexer/index"
	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	"github.com/owncloud/ocis/accounts/pkg/indexer/registry"
	"github.com/rs/zerolog"
	"path"
)

// Indexer is a facade to configure and query over multiple indices.
type Indexer struct {
	newConfig *config.Config
	config    *Config
	indices   typeMap
}

type Config struct {
	DataDir          string
	IndexRootDirName string
	Log              zerolog.Logger
}

func NewIndexer(cfg *Config) *Indexer {
	return &Indexer{
		config:  cfg,
		indices: typeMap{},
	}
}

func CreateIndexer(cfg *config.Config) *Indexer {
	return &Indexer{
		newConfig: cfg,
		indices:   typeMap{},
	}
}

func getRegistryStrategy(cfg *config.Config) string {
	if cfg.Repo.Disk.Path != "" {
		return "disk"
	}

	return "cs3"
}

func (i Indexer) AddUniqueIndex(t interface{}, indexBy, pkName, entityDirName string) error {
	strategy := getRegistryStrategy(i.newConfig)
	f := registry.IndexConstructorRegistry[strategy]["unique"]
	var idx index.Index

	if strategy == "disk" {
		idx = f(
			option.WithTypeName(getTypeFQN(t)),
			option.WithIndexBy(indexBy),
			option.WithFilesDir(path.Join(i.newConfig.Repo.Disk.Path, entityDirName)),
			option.WithDataDir(i.newConfig.Repo.Disk.Path),
		)
	} else if strategy == "cs3" {
		idx = f(
			option.WithTypeName(getTypeFQN(t)),
			option.WithIndexBy(indexBy),
			option.WithFilesDir(path.Join(i.newConfig.Repo.Disk.Path, entityDirName)),
			option.WithDataDir(i.newConfig.Repo.Disk.Path),
			option.WithDataURL(i.newConfig.Repo.CS3.DataURL),
			option.WithDataPrefix(i.newConfig.Repo.CS3.DataPrefix),
			option.WithJWTSecret(i.newConfig.Repo.CS3.JWTSecret),
		)
	}

	i.indices.addIndex(getTypeFQN(t), pkName, idx)
	return idx.Init()
}

func (i Indexer) AddNonUniqueIndex(t interface{}, indexBy, pkName, entityDirName string) error {
	strategy := getRegistryStrategy(i.newConfig)
	f := registry.IndexConstructorRegistry[strategy]["non_unique"]
	idx := f(
		option.WithTypeName(getTypeFQN(t)),
		option.WithIndexBy(indexBy),
		option.WithFilesDir(path.Join(i.config.DataDir, entityDirName)),
		option.WithIndexBaseDir(path.Join(i.config.DataDir, i.config.IndexRootDirName)),
	)

	i.indices.addIndex(getTypeFQN(t), pkName, idx)
	return idx.Init()
}

// Add a new entry to the indexer
func (i Indexer) Add(t interface{}) error {
	typeName := getTypeFQN(t)
	if fields, ok := i.indices[typeName]; ok {
		for _, indices := range fields.IndicesByField {
			for _, idx := range indices {
				pkVal := valueOf(t, fields.PKFieldName)
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
		for _, idx := range fields.IndicesByField[field] {
			res, err := idx.Lookup(val)
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

func (i Indexer) FindByPartial(t interface{}, field string, pattern string) ([]string, error) {
	typeName := getTypeFQN(t)
	resultPaths := make([]string, 0)
	if fields, ok := i.indices[typeName]; ok {
		for _, idx := range fields.IndicesByField[field] {
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
			for _, index := range indices {
				if oldV == newV {
					continue
				}
				if err := index.Update(pkVal, oldV, newV); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

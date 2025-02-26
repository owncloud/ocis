// Copyright 2018-2022 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package index

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	idxerrs "github.com/cs3org/reva/v2/pkg/storage/utils/indexer/errors"
	"github.com/cs3org/reva/v2/pkg/storage/utils/indexer/option"
	metadata "github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
)

// Unique are fields for an index of type unique.
type Unique struct {
	caseInsensitive bool
	indexBy         option.IndexBy
	typeName        string
	filesDir        string
	indexBaseDir    string
	indexRootDir    string

	storage metadata.Storage
}

// NewUniqueIndexWithOptions instantiates a new UniqueIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewUniqueIndexWithOptions(storage metadata.Storage, o ...option.Option) Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	u := &Unique{
		storage:         storage,
		caseInsensitive: opts.CaseInsensitive,
		indexBy:         opts.IndexBy,
		typeName:        opts.TypeName,
		filesDir:        opts.FilesDir,
		indexBaseDir:    path.Join(opts.Prefix, "index."+storage.Backend()),
		indexRootDir:    path.Join(opts.Prefix, "index."+storage.Backend(), strings.Join([]string{"unique", opts.TypeName, opts.IndexBy.String()}, ".")),
	}

	return u
}

// Init initializes a unique index.
func (idx *Unique) Init() error {
	if err := idx.storage.MakeDirIfNotExist(context.Background(), idx.indexBaseDir); err != nil {
		return err
	}

	return idx.storage.MakeDirIfNotExist(context.Background(), idx.indexRootDir)
}

// Lookup exact lookup by value.
func (idx *Unique) Lookup(v string) ([]string, error) {
	return idx.LookupCtx(context.Background(), v)
}

// LookupCtx retieves multiple exact values and allows passing in a context
func (idx *Unique) LookupCtx(ctx context.Context, values ...string) ([]string, error) {
	var allValues map[string]struct{}
	if len(values) != 1 {
		// prefetch all values with one request
		entries, err := idx.storage.ReadDir(context.Background(), idx.indexRootDir)
		if err != nil {
			return nil, err
		}
		// convert known values to set
		allValues = make(map[string]struct{}, len(entries))
		for _, e := range entries {
			allValues[path.Base(e)] = struct{}{}
		}
	}

	// convert requested values to set
	valueSet := make(map[string]struct{}, len(values))
	if idx.caseInsensitive {
		for _, v := range values {
			valueSet[strings.ToLower(v)] = struct{}{}
		}
	} else {
		for _, v := range values {
			valueSet[v] = struct{}{}
		}
	}

	var matches = make([]string, 0)
	for v := range valueSet {
		if _, ok := allValues[v]; ok || len(allValues) == 0 {
			oldname, err := idx.storage.ResolveSymlink(context.Background(), path.Join(idx.indexRootDir, v))
			if err != nil {
				continue
			}
			matches = append(matches, oldname)
		}
	}

	if len(matches) == 0 {
		var v string
		switch len(values) {
		case 0:
			v = "none"
		case 1:
			v = values[0]
		default:
			v = "multiple"
		}
		return nil, &idxerrs.NotFoundErr{TypeName: idx.typeName, IndexBy: idx.indexBy, Value: v}
	}

	return matches, nil
}

// Add adds a value to the index, returns the path to the root-document
func (idx *Unique) Add(id, v string) (string, error) {
	if v == "" {
		return "", nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	target := path.Join(idx.filesDir, id)
	newName := path.Join(idx.indexRootDir, v)
	if err := idx.storage.CreateSymlink(context.Background(), target, newName); err != nil {
		if os.IsExist(err) {
			return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, IndexBy: idx.indexBy, Value: v}
		}

		return "", err
	}

	return newName, nil
}

// Remove a value v from an index.
func (idx *Unique) Remove(_ string, v string) error {
	if v == "" {
		return nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	searchPath := path.Join(idx.indexRootDir, v)
	_, err := idx.storage.ResolveSymlink(context.Background(), searchPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = &idxerrs.NotFoundErr{TypeName: idx.typeName, IndexBy: idx.indexBy, Value: v}
		}

		return err
	}

	deletePath := path.Join(idx.indexRootDir, v)
	return idx.storage.Delete(context.Background(), deletePath)
}

// Update index from <oldV> to <newV>.
func (idx *Unique) Update(id, oldV, newV string) error {
	if idx.caseInsensitive {
		oldV = strings.ToLower(oldV)
		newV = strings.ToLower(newV)
	}

	if err := idx.Remove(id, oldV); err != nil {
		return err
	}

	if _, err := idx.Add(id, newV); err != nil {
		return err
	}

	return nil
}

// Search allows for glob search on the index.
func (idx *Unique) Search(pattern string) ([]string, error) {
	if idx.caseInsensitive {
		pattern = strings.ToLower(pattern)
	}

	paths, err := idx.storage.ReadDir(context.Background(), idx.indexRootDir)
	if err != nil {
		return nil, err
	}

	searchPath := idx.indexRootDir
	matches := make([]string, 0)
	for _, p := range paths {
		if found, err := filepath.Match(pattern, path.Base(p)); found {
			if err != nil {
				return nil, err
			}

			oldPath, err := idx.storage.ResolveSymlink(context.Background(), path.Join(searchPath, path.Base(p)))
			if err != nil {
				return nil, err
			}
			matches = append(matches, oldPath)
		}
	}

	if len(matches) == 0 {
		return nil, &idxerrs.NotFoundErr{TypeName: idx.typeName, IndexBy: idx.indexBy, Value: pattern}
	}

	return matches, nil
}

// CaseInsensitive undocumented.
func (idx *Unique) CaseInsensitive() bool {
	return idx.caseInsensitive
}

// IndexBy undocumented.
func (idx *Unique) IndexBy() option.IndexBy {
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

// Delete deletes the index folder from its storage.
func (idx *Unique) Delete() error {
	return idx.storage.Delete(context.Background(), idx.indexRootDir)
}

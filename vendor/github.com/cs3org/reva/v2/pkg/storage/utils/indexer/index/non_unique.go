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

// NonUnique are fields for an index of type non_unique.
type NonUnique struct {
	caseInsensitive bool
	indexBy         option.IndexBy
	typeName        string
	filesDir        string
	indexBaseDir    string
	indexRootDir    string

	storage metadata.Storage
}

// NewNonUniqueIndexWithOptions instantiates a new NonUniqueIndex instance.
// /tmp/ocis/accounts/index.cs3/Pets/Bro*
// ├── Brown/
// │   └── rebef-123 -> /tmp/testfiles-395764020/pets/rebef-123
// ├── Green/
// │    ├── goefe-789 -> /tmp/testfiles-395764020/pets/goefe-789
// │    └── xadaf-189 -> /tmp/testfiles-395764020/pets/xadaf-189
// └── White/
// |    └── wefwe-456 -> /tmp/testfiles-395764020/pets/wefwe-456
func NewNonUniqueIndexWithOptions(storage metadata.Storage, o ...option.Option) Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	return &NonUnique{
		storage:         storage,
		caseInsensitive: opts.CaseInsensitive,
		indexBy:         opts.IndexBy,
		typeName:        opts.TypeName,
		filesDir:        opts.FilesDir,
		indexBaseDir:    path.Join(opts.Prefix, "index."+storage.Backend()),
		indexRootDir:    path.Join(opts.Prefix, "index."+storage.Backend(), strings.Join([]string{"non_unique", opts.TypeName, opts.IndexBy.String()}, ".")),
	}
}

// Init initializes a non_unique index.
func (idx *NonUnique) Init() error {
	if err := idx.storage.MakeDirIfNotExist(context.Background(), idx.indexBaseDir); err != nil {
		return err
	}

	return idx.storage.MakeDirIfNotExist(context.Background(), idx.indexRootDir)
}

// Lookup exact lookup by value.
func (idx *NonUnique) Lookup(v string) ([]string, error) {
	return idx.LookupCtx(context.Background(), v)
}

// LookupCtx retieves multiple exact values and allows passing in a context
func (idx *NonUnique) LookupCtx(ctx context.Context, values ...string) ([]string, error) {
	// prefetch all values with one request
	entries, err := idx.storage.ReadDir(context.Background(), idx.indexRootDir)
	if err != nil {
		return nil, err
	}
	// convert known values to set
	allValues := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		allValues[path.Base(e)] = struct{}{}
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

	var matches = map[string]struct{}{}
	for v := range valueSet {
		if _, ok := allValues[v]; ok {
			children, err := idx.storage.ReadDir(context.Background(), filepath.Join(idx.indexRootDir, v))
			if err != nil {
				continue
			}
			for _, c := range children {
				matches[path.Base(c)] = struct{}{}
			}
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

	ret := make([]string, 0, len(matches))
	for m := range matches {
		ret = append(ret, m)
	}
	return ret, nil
}

// Add a new value to the index.
func (idx *NonUnique) Add(id, v string) (string, error) {
	if v == "" {
		return "", nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}

	newName := path.Join(idx.indexRootDir, v)
	if err := idx.storage.MakeDirIfNotExist(context.Background(), newName); err != nil {
		return "", err
	}

	if err := idx.storage.CreateSymlink(context.Background(), id, path.Join(newName, id)); err != nil {
		if os.IsExist(err) {
			return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, IndexBy: idx.indexBy, Value: v}
		}

		return "", err
	}

	return newName, nil
}

// Remove a value v from an index.
func (idx *NonUnique) Remove(id string, v string) error {
	if v == "" {
		return nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}

	deletePath := path.Join(idx.indexRootDir, v, id)
	err := idx.storage.Delete(context.Background(), deletePath)
	if err != nil {
		return err
	}

	toStat := path.Join(idx.indexRootDir, v)
	infos, err := idx.storage.ReadDir(context.Background(), toStat)
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		deletePath = path.Join(idx.indexRootDir, v)
		err := idx.storage.Delete(context.Background(), deletePath)
		if err != nil {
			return err
		}
	}

	return nil
}

// Update index from <oldV> to <newV>.
func (idx *NonUnique) Update(id, oldV, newV string) error {
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
func (idx *NonUnique) Search(pattern string) ([]string, error) {
	if idx.caseInsensitive {
		pattern = strings.ToLower(pattern)
	}

	foldersMatched := make([]string, 0)
	matches := make([]string, 0)
	paths, err := idx.storage.ReadDir(context.Background(), idx.indexRootDir)

	if err != nil {
		return nil, err
	}

	for _, p := range paths {
		if found, err := filepath.Match(pattern, path.Base(p)); found {
			if err != nil {
				return nil, err
			}

			foldersMatched = append(foldersMatched, p)
		}
	}

	for i := range foldersMatched {
		paths, _ := idx.storage.ReadDir(context.Background(), foldersMatched[i])

		for _, p := range paths {
			matches = append(matches, path.Base(p))
		}
	}

	if len(matches) == 0 {
		return nil, &idxerrs.NotFoundErr{TypeName: idx.typeName, IndexBy: idx.indexBy, Value: pattern}
	}

	return matches, nil
}

// CaseInsensitive undocumented.
func (idx *NonUnique) CaseInsensitive() bool {
	return idx.caseInsensitive
}

// IndexBy undocumented.
func (idx *NonUnique) IndexBy() option.IndexBy {
	return idx.indexBy
}

// TypeName undocumented.
func (idx *NonUnique) TypeName() string {
	return idx.typeName
}

// FilesDir  undocumented.
func (idx *NonUnique) FilesDir() string {
	return idx.filesDir
}

// Delete deletes the index folder from its storage.
func (idx *NonUnique) Delete() error {
	return idx.storage.Delete(context.Background(), idx.indexRootDir)
}

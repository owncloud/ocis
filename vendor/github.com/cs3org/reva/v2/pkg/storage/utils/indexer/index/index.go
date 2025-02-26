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

	"github.com/cs3org/reva/v2/pkg/storage/utils/indexer/option"
)

// Index can be implemented to create new indexer-strategies. See Unique for example.
// Each indexer implementation is bound to one data-column (IndexBy) and a data-type (TypeName)
type Index interface {
	Init() error
	Lookup(v string) ([]string, error)
	LookupCtx(ctx context.Context, v ...string) ([]string, error)
	Add(id, v string) (string, error)
	Remove(id string, v string) error
	Update(id, oldV, newV string) error
	Search(pattern string) ([]string, error)
	CaseInsensitive() bool
	IndexBy() option.IndexBy
	TypeName() string
	FilesDir() string
	Delete() error // Delete deletes the index folder from its storage.
}

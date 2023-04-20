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

package option

import (
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
)

// Option defines a single option function.
type Option func(o *Options)

// IndexBy defines how the data is being indexed
type IndexBy interface {
	String() string
}

// IndexByField represents the field that's being used to index the data by
type IndexByField string

// String returns a string representation
func (ibf IndexByField) String() string {
	return string(ibf)
}

// IndexByFunc represents a function that's being used to index the data by
type IndexByFunc struct {
	Name string
	Func func(v interface{}) (string, error)
}

// String returns a string representation
func (ibf IndexByFunc) String() string {
	return ibf.Name
}

// Bound represents a lower and upper bound range for an index.
// todo: if we would like to provide an upper bound then we would need to deal with ranges, in which case this is why the
// upper bound attribute is here.
type Bound struct {
	Lower, Upper int64
}

// Options defines the available options for this package.
type Options struct {
	CaseInsensitive bool
	Bound           *Bound

	TypeName string
	IndexBy  IndexBy
	FilesDir string
	Prefix   string

	Storage metadata.Storage
}

// CaseInsensitive sets the CaseInsensitive field.
func CaseInsensitive(val bool) Option {
	return func(o *Options) {
		o.CaseInsensitive = val
	}
}

// WithBounds sets the Bounds field.
func WithBounds(val *Bound) Option {
	return func(o *Options) {
		o.Bound = val
	}
}

// WithTypeName sets the TypeName option.
func WithTypeName(val string) Option {
	return func(o *Options) {
		o.TypeName = val
	}
}

// WithIndexBy sets the option IndexBy.
func WithIndexBy(val IndexBy) Option {
	return func(o *Options) {
		o.IndexBy = val
	}
}

// WithFilesDir sets the option FilesDir.
func WithFilesDir(val string) Option {
	return func(o *Options) {
		o.FilesDir = val
	}
}

// Copyright 2018-2021 CERN
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

package metadata

import (
	"context"
	"errors"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/cs3org/reva/pkg/storage/utils/decomposedfs/metadata")
}

var errUnconfiguredError = errors.New("no metadata backend configured. Bailing out")

// Backend defines the interface for file attribute backends
type Backend interface {
	Name() string

	All(ctx context.Context, path string) (map[string][]byte, error)
	Get(ctx context.Context, path, key string) ([]byte, error)

	GetInt64(ctx context.Context, path, key string) (int64, error)
	List(ctx context.Context, path string) (attribs []string, err error)
	Set(ctx context.Context, path, key string, val []byte) error
	SetMultiple(ctx context.Context, path string, attribs map[string][]byte, acquireLock bool) error
	Remove(ctx context.Context, path, key string) error

	Purge(path string) error
	Rename(oldPath, newPath string) error
	IsMetaFile(path string) bool
	MetadataPath(path string) string

	AllWithLockedSource(ctx context.Context, path string, source io.Reader) (map[string][]byte, error)
}

// NullBackend is the default stub backend, used to enforce the configuration of a proper backend
type NullBackend struct{}

// Name returns the name of the backend
func (NullBackend) Name() string { return "null" }

// All reads all extended attributes for a node
func (NullBackend) All(ctx context.Context, path string) (map[string][]byte, error) {
	return nil, errUnconfiguredError
}

// Get an extended attribute value for the given key
func (NullBackend) Get(ctx context.Context, path, key string) ([]byte, error) {
	return []byte{}, errUnconfiguredError
}

// GetInt64 reads a string as int64 from the xattrs
func (NullBackend) GetInt64(ctx context.Context, path, key string) (int64, error) {
	return 0, errUnconfiguredError
}

// List retrieves a list of names of extended attributes associated with the
// given path in the file system.
func (NullBackend) List(ctx context.Context, path string) ([]string, error) {
	return nil, errUnconfiguredError
}

// Set sets one attribute for the given path
func (NullBackend) Set(ctx context.Context, path string, key string, val []byte) error {
	return errUnconfiguredError
}

// SetMultiple sets a set of attribute for the given path
func (NullBackend) SetMultiple(ctx context.Context, path string, attribs map[string][]byte, acquireLock bool) error {
	return errUnconfiguredError
}

// Remove removes an extended attribute key
func (NullBackend) Remove(ctx context.Context, path string, key string) error {
	return errUnconfiguredError
}

// IsMetaFile returns whether the given path represents a meta file
func (NullBackend) IsMetaFile(path string) bool { return false }

// Purge purges the data of a given path from any cache that might hold it
func (NullBackend) Purge(purges string) error { return errUnconfiguredError }

// Rename moves the data for a given path to a new path
func (NullBackend) Rename(oldPath, newPath string) error { return errUnconfiguredError }

// MetadataPath returns the path of the file holding the metadata for the given path
func (NullBackend) MetadataPath(path string) string { return "" }

// AllWithLockedSource reads all extended attributes from the given reader
// The path argument is used for storing the data in the cache
func (NullBackend) AllWithLockedSource(ctx context.Context, path string, source io.Reader) (map[string][]byte, error) {
	return nil, errUnconfiguredError
}

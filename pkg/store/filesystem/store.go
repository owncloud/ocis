// Package store implements the go-micro store interface
package store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	mstore "github.com/micro/go-micro/v2/store"
)

// DefaultPath assumes UNIX
var DefaultPath string = "/var/tmp/ocis-store"

// Store interacts with the filesystem to manage account information
type Store struct {
	mountPath string
}

// New returns a new file system store manager
// TODO add mountPath as a flag
func New() Store {
	return Store{
		mountPath: DefaultPath,
	}
}

// Init implements the store interface
func (s *Store) Init(...mstore.Options) {}

// List implements the store interface
func (s *Store) List() ([]*mstore.Record, error) {
	return nil, nil
}

// Read implements the store interface
func (s *Store) Read(key string, opts ...mstore.ReadOption) ([]*mstore.Record, error) {
	return nil, nil
}

// Write implements the store interface
func (s *Store) Write(rec *mstore.Record) error {
	path := filepath.Join(s.mountPath, rec.Key)

	if filepath.IsAbs(path) {
		// TODO WARN, storage is relative to the service directory
	}

	if len(rec.Key) < 1 {
		// TODO log error: empty key
		return fmt.Errorf("%v", "key is empty")
	}

	return ioutil.WriteFile(path, rec.Value, 0644)
}

// Delete implements the store interface
func (s *Store) Delete(key string) error {
	return nil
}

// String implements the store interface, and the stringer interface
func (s *Store) String() string {
	return "store"
}

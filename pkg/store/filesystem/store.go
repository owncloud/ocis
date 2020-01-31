// Package store implements the go-micro store interface
package store

import (
	"fmt"
	"io/ioutil"
	"log"

	mstore "github.com/micro/go-micro/v2/store"
)

// Store interacts with the filesystem to manage account information
type Store struct {
	mountPath string
}

// New returns a new file system store manager
// TODO add mountPath as an option argument
func New() Store {
	return Store{
		mountPath: "/var/tmp/ocis-store",
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
	if len(rec.Key) < 1 {
		return fmt.Errorf("%v", "key is empty")
	}

	if err := ioutil.WriteFile(s.mountPath+"/"+rec.Key, rec.Value, 0644); err != nil {
		log.Panic(err)
	}
	return nil
}

// Delete implements the store interface
func (s *Store) Delete(key string) error {
	return nil
}

// String implements the store interface, and the stringer interface
func (s *Store) String() string {
	return "store"
}

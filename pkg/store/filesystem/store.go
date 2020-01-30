// Package store implements the go-micro store interface
package store

import (
	mstore "github.com/micro/go-micro/v2/store"
)

// Store interacts with the filesystem to manage account information
type Store struct{}

// New returns a new file system store manager
func New() Store {
	return Store{}
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
func (s *Store) Write(*mstore.Record) error {
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

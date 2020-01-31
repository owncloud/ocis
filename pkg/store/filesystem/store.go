// Package store implements the go-micro store interface
package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	mstore "github.com/micro/go-micro/v2/store"
	olog "github.com/owncloud/ocis-pkg/log"
)

// Logger is a global logger
var Logger olog.Logger

// Store interacts with the filesystem to manage account information
type Store struct {
	logger    olog.Logger
	mountPath string
}

// New returns a new file system store manager
// TODO add mountPath as a flag. Accept a *config argument
func New() Store {
	// default to the current working directory if not configured
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// TODO log error
		// TODO deal with this
	}
	return Store{
		mountPath: dir,
		logger:    olog.NewLogger(),
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

	if len(rec.Key) < 1 {
		// TODO log error: empty key
		return fmt.Errorf("%v", "key is empty")
	}

	if err := ioutil.WriteFile(path, rec.Value, 0644); err != nil {
		return err
	}

	s.logger.Info().Msgf("%v bytes written to %v", len(rec.Value), path)
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

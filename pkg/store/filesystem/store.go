// Package store implements the go-micro store interface
package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	mstore "github.com/micro/go-micro/v2/store"
	olog "github.com/owncloud/ocis-pkg/log"
)

// StoreName is the default name for the store container
var StoreName string = "ocis-store"

// Store interacts with the filesystem to manage account information
type Store struct {
	mountPath string
	Logger    olog.Logger
}

// New returns a new file system store manager
// TODO add mountPath as a flag. Accept a *config argument
func New() Store {
	s := Store{
		Logger: olog.NewLogger(),
	}

	// default to the current working directory if not configured
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		s.Logger.Err(err).Msg("initializing accounts store")
	}

	dest := filepath.Join(dir, StoreName)
	if _, err := os.Stat(dest); err != nil {
		s.Logger.Info().Msgf("creating container on %v", dest)
		os.Mkdir(dest, 0700)
	}

	s.mountPath = dest
	return s
}

// Init implements the store interface
// TODO it could prepare the destination path, for instance
func (s *Store) Init(...mstore.Options) {}

// List implements the store interface
func (s *Store) List() ([]*mstore.Record, error) {
	return nil, nil
}

// Read implements the store interface
// this implementation only reads by id.
func (s *Store) Read(key string, opts ...mstore.ReadOption) ([]*mstore.Record, error) {
	contents, err := ioutil.ReadFile(path.Join(s.mountPath, key))
	if err != nil {
		s.Logger.Err(err).Msgf("error reading contents of key %v: file not found", key)
		return []*mstore.Record{}, err
	}

	return []*mstore.Record{
		&mstore.Record{
			Key:   key,
			Value: contents,
		},
	}, nil
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

	s.Logger.Info().Msgf("%v bytes written to %v", len(rec.Value), path)
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

// creates the default container for the ocis-store if it doesn't exist
func init() {

}

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

// TODO add mountPath as a flag. Accept a *config argument
func New() *Store {
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
	return &s
}

// Init implements the store interface
func (s Store) Init(...mstore.Option) error {
	return nil
}

// List returns all the identities in the mountPath folder
func (s Store) List() ([]*mstore.Record, error) {
	records := []*mstore.Record{}
	identities, err := ioutil.ReadDir(s.mountPath)
	if err != nil {
		s.Logger.Err(err).Msgf("error reading %v", s.mountPath)
	}

	s.Logger.Info().Msg("listing identities")
	for _, v := range identities {
		records = append(records, &mstore.Record{
			Key: v.Name(),
		})
	}

	return records, nil
}

// Read implements the store interface. This implementation only reads by id.
func (s Store) Read(key string, opts ...mstore.ReadOption) ([]*mstore.Record, error) {
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
func (s Store) Write(rec *mstore.Record) error {
	path := filepath.Join(s.mountPath, rec.Key)

	if len(rec.Key) < 1 {
		s.Logger.Error().Msg("key cannot be empty")
		return fmt.Errorf("%v", "key is empty")
	}

	if err := ioutil.WriteFile(path, rec.Value, 0644); err != nil {
		return err
	}

	s.Logger.Info().Msgf("%v bytes written to %v", len(rec.Value), path)
	return nil
}

// Delete implements the store interface
func (s Store) Delete(key string) error {
	return nil
}

// String implements the store interface, and the stringer interface
func (s Store) String() string {
	return "store"
}

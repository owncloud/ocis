// Package store implements the go-micro store interface
package store

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/owncloud/ocis-accounts/pkg/account"
	"github.com/owncloud/ocis-accounts/pkg/config"
	olog "github.com/owncloud/ocis-pkg/log"
)

var (
	// StoreName is the default name for the accounts store
	StoreName string = "ocis-store"

	// managerName
	managerName = "filesystem"
)

// StoreName is the default name for the store container

// Store interacts with the filesystem to manage account information
type Store struct {
	mountPath string
	Logger    olog.Logger
}

// New creates a new store
func New(cfg *config.Config) account.Manager {
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

// List returns all the identities in the mountPath folder
func (s Store) List() []*account.Record {
	records := []*account.Record{}
	identities, err := ioutil.ReadDir(s.mountPath)
	if err != nil {
		s.Logger.Err(err).Msgf("error reading %v", s.mountPath)
		return records
	}

	s.Logger.Info().Msg("listing identities")
	for _, v := range identities {
		records = append(records, &account.Record{
			Key: v.Name(),
		})
	}

	return records
}

// Read implements the store interface. This implementation only reads by id.
func (s Store) Read(key string) *account.Record {
	contents, err := ioutil.ReadFile(path.Join(s.mountPath, key))
	if err != nil {
		s.Logger.Err(err).Msgf("error reading contents of key %v: file not found", key)
		return &account.Record{}
	}

	return &account.Record{
		Key:   key,
		Value: contents,
	}
}

// Write implements the store interface
func (s Store) Write(rec *account.Record) *account.Record {
	path := filepath.Join(s.mountPath, rec.Key)

	if len(rec.Key) < 1 {
		s.Logger.Error().Msg("key cannot be empty")
		return &account.Record{}
	}

	if err := ioutil.WriteFile(path, rec.Value, 0644); err != nil {
		return &account.Record{}
	}

	s.Logger.Info().Msgf("%v bytes written to %v", len(rec.Value), path)
	return rec
}
func init() {
	account.Registry[managerName] = New
}

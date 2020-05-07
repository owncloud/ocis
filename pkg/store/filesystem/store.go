// Package store implements the go-micro store interface
package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	// gproto "github.com/golang/protobuf/proto"
	"github.com/owncloud/ocis-accounts/pkg/account"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	olog "github.com/owncloud/ocis-pkg/v2/log"
)

var (
	// StoreName is the default name for the accounts store
	StoreName     = "ocis-store"
	managerName   = "filesystem"
	uuidSpace     = "uuid"
	usernameSpace = "username"
	emailSpace    = "email"
	identitySpace = "identity"
	emptyKeyError = "key cannot be empty"
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
		Logger: olog.NewLogger(olog.Name(cfg.Server.Name)),
	}

	dest := filepath.Join(cfg.MountPath, StoreName)
	if _, err := os.Stat(dest); err != nil {
		s.Logger.Info().Msgf("creating container on %v", dest)
		err := os.MkdirAll(dest, 0700)
		if err != nil {
			s.Logger.Err(err).Msgf("providing container on %v", dest)
		}
	}

	s.mountPath = dest
	return &s
}

// List returns all the identities in the mountPath folder
func (s Store) List() ([]*proto.Record, error) {
	records := []*proto.Record{}
	identities, err := ioutil.ReadDir(path.Join(s.mountPath, uuidSpace))
	if err != nil {
		s.Logger.Err(err).Msgf("error reading %v", s.mountPath)
		return nil, err
	}

	s.Logger.Info().Msg("listing identities")
	for _, v := range identities {
		records = append(records, &proto.Record{
			Key: v.Name(),
		})
	}

	return records, nil
}

func (s Store) readSpace(space string, subspace string, key string) (*proto.Record, error) {
	contents, err := ioutil.ReadFile(path.Join(s.mountPath, space, subspace, key))
	if err != nil {
		s.Logger.Err(err).Str("space", space).Str("subspace", subspace).Str("key", key).Msg("error reading record")
		return nil, err
	}

	rec := NewRecord(
		WithUUID(key),
	)

	if err = json.Unmarshal(contents, rec); err != nil {
		s.Logger.Err(err).Msg("error unmarshaling record")
		return nil, err
	}

	return rec, nil
}

// Read implements the store interface. This implementation only reads by id.
func (s Store) Read(uuid string) (*proto.Record, error) {
	return s.readSpace(uuidSpace, "", uuid)
}

// ReadByUsername implements the store interface. This implementation only reads by username.
func (s Store) ReadByUsername(username string) (*proto.Record, error) {
	return s.readSpace(usernameSpace, "", username)
}

// ReadByEmail implements the store interface. This implementation only reads by email.
func (s Store) ReadByEmail(email string) (*proto.Record, error) {
	i := strings.LastIndex(email, "@")
	if i < 0 {
		// no domain part
		return s.readSpace(emailSpace, "local", email)
	}

	return s.readSpace(emailSpace, email[:i], email[i:])
}

// ReadByIdentity implements the store interface. This implementation only reads by iss & sub.
func (s Store) ReadByIdentity(identity *proto.IdHistory) (*proto.Record, error) {
	if identity.Iss == "" {
		s.Logger.Error().Msg("iss cannot be empty")
		return nil, fmt.Errorf(emptyKeyError)
	}
	return s.readSpace(identitySpace, identity.Iss, identity.Sub)
}

// Write implements the store interface
func (s Store) Write(rec *proto.Record) (*proto.Record, error) {
	if len(rec.Key) < 1 {
		s.Logger.Error().Msg("key cannot be empty")
		return nil, fmt.Errorf(emptyKeyError)
	}

	path := filepath.Join(s.mountPath, uuidSpace, rec.Key)

	contents, err := json.Marshal(rec)
	if err != nil {
		s.Logger.Err(err).Msg("record could not be marshaled")
		return nil, err
	}

	if err := ioutil.WriteFile(path, contents, 0644); err != nil {
		return nil, err
	}
	s.Logger.Info().Int("bytes", len(contents)).Str("path", path).Msg("wrote account")

	// write symlinks for other lookups
	// use hardlinks?
	// TODO what if target already exists? use dirs with multiple symlinks to uuid?
	// TODO what if username or email changes? use uuid folder with a file per property? srsly ... we should use a proper storage for this
	if rec.Payload != nil {
		if rec.Payload.Account != nil {
			if rec.Payload.Account.StandardClaims != nil {
				if rec.Payload.Account.StandardClaims.PreferredUsername != "" {
					p := filepath.Join(s.mountPath, usernameSpace, rec.Payload.Account.StandardClaims.PreferredUsername)
					os.MkdirAll(filepath.Dir(p), 0700)
					os.Symlink(path, p)
				} else {
					s.Logger.Warn().Str("uuid", rec.Key).Msg("has no preferred username in standard claims")
				}
				if rec.Payload.Account.StandardClaims.Email != "" {
					p := filepath.Join(s.mountPath, emailSpace, rec.Payload.Account.StandardClaims.Email)
					os.MkdirAll(filepath.Dir(p), 0700)
					os.Symlink(path, p)
				} else {
					s.Logger.Warn().Str("uuid", rec.Key).Msg("has no email in standard claims")
				}
				if rec.Payload.Account.Issuer != "" {
					if rec.Payload.Account.StandardClaims.Sub != "" {
						p := filepath.Join(s.mountPath, identitySpace, rec.Payload.Account.Issuer, rec.Payload.Account.StandardClaims.Sub)
						os.MkdirAll(filepath.Dir(p), 0700)
						os.Symlink(path, p)
					} else {
						s.Logger.Warn().Str("uuid", rec.Key).Msg("has no sub in standard claims")
					}
				} else {
					s.Logger.Warn().Str("uuid", rec.Key).Msg("has no issuer in account")
				}
			} else {
				s.Logger.Warn().Str("uuid", rec.Key).Msg("has no standard claims in account")
			}
		} else {
			s.Logger.Warn().Str("uuid", rec.Key).Msg("has no account in payload")
		}
	} else {
		s.Logger.Warn().Str("uuid", rec.Key).Msg("has no payload")
	}

	return rec, nil
}

func init() {
	account.Registry[managerName] = New
}

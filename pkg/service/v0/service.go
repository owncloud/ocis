package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-pkg/v2/log"
)

// New returns a new instance of Service
func New(opts ...Option) (s *Service, err error) {
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config
	// read all user and group records

	// for now recreate index on every start
	if err = os.RemoveAll(filepath.Join(cfg.Server.AccountsDataPath, "index.bleve")); err != nil {
		return nil, err
	}

	// check if accounts exist
	accountsDir := filepath.Join(cfg.Server.AccountsDataPath, "accounts")
	var fi os.FileInfo
	if fi, err = os.Stat(accountsDir); err != nil {
		if os.IsNotExist(err) {
			// create accounts directory
			if err = os.MkdirAll(accountsDir, 0700); err != nil {
				return nil, err
			}
			// create default accounts
			accounts := []proto.Account{
				{
					Id:            "4c510ada-c86b-4815-8820-42cdf82c3d51",
					PreferredName: "einstein",
					Mail:          "einstein@example.org",
					DisplayName:   "Albert Einstein",
					UidNumber:     20000,
					GidNumber:     30000,
					PasswordProfile: &proto.PasswordProfile{
						Password: "$6$rounds=35210$sa1u5Pmfo4cr23Vw$RJNGElaDB1D3xorWkfTEGm2Ko.o2QL3E0cimKx23MNxVWVFSkUUeRoC7FqC4RzYDNQBD6cKzovTEaDD.8TDkD.",
					},
					AccountEnabled: true,
				},
				{
					Id:            "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c",
					PreferredName: "marie",
					Mail:          "marie@example.org",
					DisplayName:   "Marie Curie",
					UidNumber:     20001,
					GidNumber:     30000,
					PasswordProfile: &proto.PasswordProfile{
						Password: "$6$rounds=81434$sa1u5Pmfo4cr23Vw$W78cyL884GmuvDpxYPvSRBVzEj02T5QhTTcI8Dv4IKvMooDFGv4bwaWMkH9HfJ0wgpEBW7Lp.4Cad0xE/MYSg1",
					},
					AccountEnabled: true,
				},
				{
					Id:            "932b4540-8d16-481e-8ef4-588e4b6b151c",
					PreferredName: "richard",
					Mail:          "richard@example.org",
					DisplayName:   "Richard Feynman",
					UidNumber:     20002,
					GidNumber:     30000,
					PasswordProfile: &proto.PasswordProfile{
						Password: "$6$rounds=5524$sa1u5Pmfo4cr23Vw$58bQVL/JeUlwM0RY21YKAFMvKvwKLLysGllYXox.vwKT5dHMwdzJjCxwTDMnB2o2pwexC8o/iOXyP2zrhALS40",
					},
					AccountEnabled: true,
				},
				// technical users for kopano and reva
				{
					Id:            "820ba2a1-3f54-4538-80a4-2d73007e30bf",
					PreferredName: "konnectd",
					Mail:          "idp@example.org",
					DisplayName:   "Kopano Konnectd",
					UidNumber:     10000,
					GidNumber:     15000,
					PasswordProfile: &proto.PasswordProfile{
						Password: "$6$rounds=9746$sa1u5Pmfo4cr23Vw$2hnwpkTvUkWX0v6mh8Aw1pbzEXa9EUJzmrey4g2W/8arwWCwhteqU//3aWnA3S0d5T21fOKYteoqlsN1IbTcN.",
					},
					AccountEnabled: true,
				},
				{
					Id:            "bc596f3c-c955-4328-80a0-60d018b4ad57",
					PreferredName: "reva",
					Mail:          "storage@example.org",
					DisplayName:   "Reva Inter Operability Platform",
					UidNumber:     10001,
					GidNumber:     15000,
					PasswordProfile: &proto.PasswordProfile{
						Password: "$6$rounds=91087$sa1u5Pmfo4cr23Vw$wPC3BbMTbP/ytlo0p.f99zJifyO70AUCdKIK9hkhwutBKGCirLmZs/MsWAG6xHjVvmnmHN5NoON7FUGv5pPaN.",
					},
					AccountEnabled: true,
				},
			}
			// TODO groups
			for i := range accounts {
				var bytes []byte
				if bytes, err = json.Marshal(&accounts[i]); err != nil {
					logger.Error().Err(err).Interface("account", &accounts[i]).Msg("could not marshal default account")
					return
				}
				path := filepath.Join(accountsDir, accounts[i].Id)
				if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
					accounts[i].PasswordProfile.Password = "***REMOVED***"
					logger.Error().Err(err).Str("path", path).Interface("account", &accounts[i]).Msg("could not persist default account")
					return
				}
			}

		}
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", accountsDir)
	}

	if err = os.MkdirAll(filepath.Join(cfg.Server.AccountsDataPath, "groups"), 0700); err != nil {
		return nil, err
	}

	mapping := bleve.NewIndexMapping()
	// keep all symbols in terms to allow exact maching, eg. emails
	mapping.DefaultAnalyzer = keyword.Name
	// TODO don't bother to store fields as we will load the account from disk
	//groupsFieldMapping := bleve.NewTextFieldMapping()
	//blogMapping.AddFieldMappingsAt("memberOf", nameFieldMapping)
	// TODO index groups and accounts as different types?

	s = &Service{
		id:     cfg.GRPC.Namespace + "." + cfg.Server.Name,
		log:    logger,
		Config: cfg,
	}

	if s.index, err = bleve.New(filepath.Join(cfg.Server.AccountsDataPath, "index.bleve"), mapping); err != nil {
		return
	}

	if err = s.indexAccounts(filepath.Join(cfg.Server.AccountsDataPath, "accounts")); err != nil {
		return nil, err
	}
	if err = s.indexGroups(filepath.Join(cfg.Server.AccountsDataPath, "groups")); err != nil {
		return nil, err
	}

	// TODO watch folders for new records

	return
}

// Service implements the AccountsServiceHandler interface
type Service struct {
	id     string
	log    log.Logger
	Config *config.Config
	index  bleve.Index
}

func cleanupID(id string) (string, error) {
	id = filepath.Clean(id)
	if id == "." || strings.Contains(id, "/") {
		return "", errors.New("invalid id " + id)
	}
	return id, nil
}

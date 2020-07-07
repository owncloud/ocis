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

	accountsDir := filepath.Join(cfg.Server.AccountsDataPath, "accounts")
	{
		// check if accounts exist
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
						MemberOf: []*proto.Group{
							{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"}, // users
							{Id: "6040aa17-9c64-4fef-9bd0-77234d71bad0"}, // sailing-lovers
							{Id: "dd58e5ec-842e-498b-8800-61f2ec6f911f"}, // violin-haters
							{Id: "262982c1-2362-4afa-bfdf-8cbfef64a06e"}, // physics-lovers
						},
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
						MemberOf: []*proto.Group{
							{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"}, // users
							{Id: "7b87fd49-286e-4a5f-bafd-c535d5dd997a"}, // radium-lovers
							{Id: "cedc21aa-4072-4614-8676-fa9165f598ff"}, // polonium-lovers
							{Id: "262982c1-2362-4afa-bfdf-8cbfef64a06e"}, // physics-lovers
						},
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
						MemberOf: []*proto.Group{
							{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"}, // users
							{Id: "a1726108-01f8-4c30-88df-2b1a9d1cba1a"}, // quantum-lovers
							{Id: "167cbee2-0518-455a-bfb2-031fe0621e5d"}, // philosophy-haters
							{Id: "262982c1-2362-4afa-bfdf-8cbfef64a06e"}, // physics-lovers
						},
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
						MemberOf: []*proto.Group{
							{Id: "34f38767-c937-4eb6-b847-1c175829a2a0"}, // sysusers
						},
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
						MemberOf: []*proto.Group{
							{Id: "34f38767-c937-4eb6-b847-1c175829a2a0"}, // sysusers
						},
					},
				}
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
	}

	groupsDir := filepath.Join(cfg.Server.AccountsDataPath, "groups")
	{
		// check if groups exist
		var fi os.FileInfo
		if fi, err = os.Stat(groupsDir); err != nil {
			if os.IsNotExist(err) {
				// create accounts directory
				if err = os.MkdirAll(groupsDir, 0700); err != nil {
					return nil, err
				}
				// create default accounts
				groups := []proto.Group{
					{Id: "34f38767-c937-4eb6-b847-1c175829a2a0", GidNumber: 15000, OnPremisesSamAccountName: "sysusers", DisplayName: "Technical users", Description: "A group for technical users. They should not show up in sharing dialogs.", Members: []*proto.Account{
						{Id: "820ba2a1-3f54-4538-80a4-2d73007e30bf"}, // konnectd
						{Id: "bc596f3c-c955-4328-80a0-60d018b4ad57"}, // reva
					}},
					{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa", GidNumber: 30000, OnPremisesSamAccountName: "users", DisplayName: "Users", Description: "A group every normal user belongs to.", Members: []*proto.Account{
						{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
						{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
						{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
					}},
					{Id: "6040aa17-9c64-4fef-9bd0-77234d71bad0", GidNumber: 30001, OnPremisesSamAccountName: "sailing-lovers", DisplayName: "Sailing lovers", Members: []*proto.Account{
						{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
					}},
					{Id: "dd58e5ec-842e-498b-8800-61f2ec6f911f", GidNumber: 30002, OnPremisesSamAccountName: "violin-haters", DisplayName: "Violin haters", Members: []*proto.Account{
						{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
					}},
					{Id: "7b87fd49-286e-4a5f-bafd-c535d5dd997a", GidNumber: 30003, OnPremisesSamAccountName: "radium-lovers", DisplayName: "Radium lovers", Members: []*proto.Account{
						{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
					}},
					{Id: "cedc21aa-4072-4614-8676-fa9165f598ff", GidNumber: 30004, OnPremisesSamAccountName: "polonium-lovers", DisplayName: "Polonium lovers", Members: []*proto.Account{
						{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
					}},
					{Id: "a1726108-01f8-4c30-88df-2b1a9d1cba1a", GidNumber: 30005, OnPremisesSamAccountName: "quantum-lovers", DisplayName: "Quantum lovers", Members: []*proto.Account{
						{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
					}},
					{Id: "167cbee2-0518-455a-bfb2-031fe0621e5d", GidNumber: 30006, OnPremisesSamAccountName: "philosophy-haters", DisplayName: "Philosophy haters", Members: []*proto.Account{
						{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
					}},
					{Id: "262982c1-2362-4afa-bfdf-8cbfef64a06e", GidNumber: 30007, OnPremisesSamAccountName: "physics-lovers", DisplayName: "Physics lovers", Members: []*proto.Account{
						{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
						{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
						{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
					}},
				}
				for i := range groups {
					var bytes []byte
					if bytes, err = json.Marshal(&groups[i]); err != nil {
						logger.Error().Err(err).Interface("group", &groups[i]).Msg("could not marshal default group")
						return
					}
					path := filepath.Join(groupsDir, groups[i].Id)
					if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
						logger.Error().Err(err).Str("path", path).Interface("group", &groups[i]).Msg("could not persist default group")
						return
					}
				}
			}
		} else if !fi.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", groupsDir)
		}
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

	indexDir := filepath.Join(cfg.Server.AccountsDataPath, "index.bleve")
	// for now recreate index on every start
	if err = os.RemoveAll(indexDir); err != nil {
		return nil, err
	}
	if s.index, err = bleve.New(indexDir, mapping); err != nil {
		return
	}
	if err = s.indexAccounts(accountsDir); err != nil {
		return nil, err
	}
	if err = s.indexGroups(groupsDir); err != nil {
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

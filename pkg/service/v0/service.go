package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/CiscoM31/godata"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/search/query"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-accounts/pkg/provider"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/rs/zerolog/log"
	"github.com/tredoe/osutil/user/crypt"

	// register crypt functions
	_ "github.com/tredoe/osutil/user/crypt/apr1_crypt"
	_ "github.com/tredoe/osutil/user/crypt/md5_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha256_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

// New returns a new instance of Service
// TODO pass in logger as options
func New(cfg *config.Config) (s *Service, err error) {
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
				},
			}
			// TODO groups
			for i := range accounts {
				var bytes []byte
				if bytes, err = json.Marshal(accounts[i]); err != nil {
					log.Error().Err(err).Interface("account", accounts[i]).Msg("could not marshal default account")
					return
				}
				path := filepath.Join(accountsDir, accounts[i].Id)
				if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
					accounts[i].PasswordProfile.Password = "***REMOVED***"
					log.Error().Err(err).Str("path", path).Interface("account", accounts[i]).Msg("could not persist default account")
					return
				}
			}

		}
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", accountsDir)
	}

	mapping := bleve.NewIndexMapping()
	// keep all symbols in terms to allow exact maching, eg. emails
	mapping.DefaultAnalyzer = keyword.Name
	// TODO don't bother to store fields as we will load the account from disk

	s = &Service{
		Config: cfg,
	}

	if s.index, err = bleve.New(filepath.Join(cfg.Server.AccountsDataPath, "index.bleve"), mapping); err != nil {
		return
	}
	var f *os.File
	if f, err = os.Open(filepath.Join(cfg.Server.AccountsDataPath, "accounts")); err != nil {
		log.Error().Err(err).Str("dir", filepath.Join(cfg.Server.AccountsDataPath, "accounts")).Msg("could not open accounts folder")
		return
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Error().Err(err).Str("dir", filepath.Join(cfg.Server.AccountsDataPath, "accounts")).Msg("could not list accounts folder")
		return
	}
	var data []byte
	for _, file := range list {
		path := filepath.Join(cfg.Server.AccountsDataPath, "accounts", file.Name())
		if data, err = ioutil.ReadFile(path); err != nil {
			log.Error().Err(err).Str("path", path).Msg("could not read account")
			continue
		}
		a := proto.Account{}
		if err = json.Unmarshal(data, &a); err != nil {
			log.Error().Err(err).Str("path", path).Msg("could not unmarshal account")
			continue
		}
		log.Debug().Interface("account", a).Msg("found account")
		if err = s.index.Index(a.Id, a); err != nil {
			log.Error().Err(err).Str("path", path).Interface("account", a).Msg("could not index account")
			continue
		}
	}

	// TODO watch folders for new records

	s.Config = cfg

	return
}

// Service implements the AccountsServiceHandler interface
type Service struct {
	Config *config.Config
	index  bleve.Index
}

// an auth request is currently hardcoded and has to match this regex
// login eq \"teddy\" and password eq \"F&1!b90t111!\"
var authQuery = regexp.MustCompile(`^login eq '(.*)' and password eq '(.*)'$`) // TODO how is ' escaped in the password?

// ListAccounts implements the AccountsServiceHandler interface
// the query contains account properties
// TODO id vs onpremiseimmutableid
func (s Service) ListAccounts(ctx context.Context, in *proto.ListAccountsRequest, res *proto.ListAccountsResponse) (err error) {

	var password string

	var query query.Query

	// check if this looks like an auth request
	match := authQuery.FindStringSubmatch(in.Query)
	if len(match) == 3 {
		in.Query = fmt.Sprintf("preferred_name eq '%s'", match[1]) // todo fetch email? make query configurable
		password = match[2]
		if password == "" {

			return fmt.Errorf("password must not be empty")
		}
	}

	if in.Query != "" {
		// parse the query like an odata filter
		var q *godata.GoDataFilterQuery
		if q, err = godata.ParseFilterString(in.Query); err != nil {
			log.Error().Err(err).Msg("could not parse query")
			return
		}

		// convert to bleve query
		query, err = provider.BuildBleveQuery(q)
		if err != nil {
			log.Error().Err(err).Msg("could not build bleve query")
			return
		}
	} else {
		query = bleve.NewMatchAllQuery()
	}

	log.Debug().Interface("query", query).Msg("using query")

	searchRequest := bleve.NewSearchRequest(query)
	var searchResult *bleve.SearchResult
	searchResult, err = s.index.Search(searchRequest)
	if err != nil {
		log.Error().Err(err).Msg("could not execute bleve search")
		return
	}

	log.Debug().Interface("result", searchResult).Msg("result")

	res.Accounts = make([]*proto.Account, 0)

	for _, hit := range searchResult.Hits {
		path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", hit.ID)

		var data []byte
		data, err = ioutil.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("could not read account")
			continue
		}
		a := proto.Account{}
		err = json.Unmarshal(data, &a)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("could not unmarshal account")
			continue
		}
		log.Debug().Interface("account", a).Msg("found account")

		if password != "" {
			if a.PasswordProfile == nil {
				log.Debug().Interface("account", a).Msg("no password profile")
				return fmt.Errorf("invalid password")
			}
			if !s.passwordIsValid(a.PasswordProfile.Password, password) {
				return fmt.Errorf("invalid password")
			}
		}

		res.Accounts = append(res.Accounts, &a)
	}

	return nil
}

func (s Service) passwordIsValid(hash string, pwd string) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Err(fmt.Errorf("%s", r)).Str("hash", hash).Msg("password lib panicked")
		}
	}()

	c := crypt.NewFromHash(hash)
	return c.Verify(hash, []byte(pwd)) == nil
}

func cleanupID(id string) (string, error) {
	id = filepath.Clean(id)
	if id == "." || strings.Contains(id, "/") {
		return "", errors.New("bad request")
	}
	return id, nil
}

// GetAccount implements the AccountsServiceHandler interface
func (s Service) GetAccount(c context.Context, req *proto.GetAccountRequest, res *proto.Account) (err error) {
	var id string
	if id, err = cleanupID(req.Id); err != nil {
		return
	}
	path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", id)

	var data []byte
	data, err = ioutil.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("could not read account")
		// TODO we need error handling ... eg Not Found
		return fmt.Errorf("account not found")
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("could not unmarshal account")
		return fmt.Errorf("internal server error")
	}

	log.Debug().Interface("account", res).Msg("found account")
	return
}

// CreateAccount implements the AccountsServiceHandler interface
func (s Service) CreateAccount(c context.Context, req *proto.CreateAccountRequest, res *proto.Account) (err error) {
	var id string
	if req.Account == nil {
		return fmt.Errorf("account missing")
	}
	if req.Account.Id == "" {
		req.Account.Id = uuid.Must(uuid.NewV4()).String()
	}

	if id, err = cleanupID(req.Account.Id); err != nil {
		return
	}
	path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", id)

	if req.Account.PasswordProfile != nil && req.Account.PasswordProfile.Password != "" {
		// encrypt password
		c := crypt.New(crypt.SHA512)
		if req.Account.PasswordProfile.Password, err = c.Generate([]byte(req.Account.PasswordProfile.Password), nil); err != nil {
			req.Account.PasswordProfile.Password = "***REMOVED***"
			log.Error().Err(err).Str("id", id).Interface("account", req.Account).Msg("could not hash password")
			return
		}
	}

	var bytes []byte
	if bytes, err = json.Marshal(req.Account); err != nil {
		log.Error().Err(err).Interface("account", req.Account).Msg("could not marshal account")
		return
	}
	if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
		req.Account.PasswordProfile.Password = "***REMOVED***"
		log.Error().Err(err).Str("id", id).Str("path", path).Interface("account", req.Account).Msg("could not persist new account")
		return
	}

	if err = s.index.Index(id, req.Account); err != nil {
		req.Account.PasswordProfile.Password = "***REMOVED***"
		log.Error().Err(err).Str("id", id).Str("path", path).Interface("account", req.Account).Msg("could not index new account")
		return
	}

	return
}

// UpdateAccount implements the AccountsServiceHandler interface
func (s Service) UpdateAccount(c context.Context, req *proto.UpdateAccountRequest, res *proto.Account) (err error) {
	return errors.New("not implemented")
}

// DeleteAccount implements the AccountsServiceHandler interface
func (s Service) DeleteAccount(c context.Context, req *proto.DeleteAccountRequest, res *empty.Empty) (err error) {
	var id string
	if id, err = cleanupID(req.Id); err != nil {
		return
	}
	path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", id)

	if err = os.Remove(path); err != nil {
		log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not remove account")
		return
	}

	if err = s.index.Delete(id); err != nil {
		log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not remove account from index")
		return
	}
	return os.Remove(path)
}

// ListGroups implements the AccountsServiceHandler interface
func (s Service) ListGroups(c context.Context, req *proto.ListGroupsRequest, res *proto.ListGroupsResponse) (err error) {
	return errors.New("not implemented")
}

// GetGroup implements the AccountsServiceHandler interface
func (s Service) GetGroup(c context.Context, req *proto.GetGroupRequest, res *proto.Group) (err error) {
	return errors.New("not implemented")
}

// CreateGroup implements the AccountsServiceHandler interface
func (s Service) CreateGroup(c context.Context, req *proto.CreateGroupRequest, res *proto.Group) (err error) {
	return errors.New("not implemented")
}

// UpdateGroup implements the AccountsServiceHandler interface
func (s Service) UpdateGroup(c context.Context, req *proto.UpdateGroupRequest, res *proto.Group) (err error) {
	return errors.New("not implemented")
}

// DeleteGroup implements the AccountsServiceHandler interface
func (s Service) DeleteGroup(c context.Context, req *proto.DeleteGroupRequest, res *empty.Empty) (err error) {
	return errors.New("not implemented")
}

// RegisterSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func RegisterSettingsBundles(l *olog.Logger) {
	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	service := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)

	requests := []settings.SaveSettingsBundleRequest{
		generateSettingsBundleProfileRequest(),
	}

	for i := range requests {
		res, err := service.SaveSettingsBundle(context.Background(), &requests[i])
		if err != nil {
			l.Err(err).
				Msg("Error registering settings bundle")
		} else {
			l.Info().
				Str("bundle key", res.SettingsBundle.Identifier.BundleKey).
				Msg("Successfully registered settings bundle")
		}
	}
}

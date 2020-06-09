package service

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/CiscoM31/godata"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/golang/protobuf/ptypes/empty"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-accounts/pkg/provider"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/rs/zerolog/log"
)

// New returns a new instance of Service
func New(cfg *config.Config) Service {
	// read all user and group records

	// for now recreate index on every start
	os.RemoveAll(filepath.Join(cfg.Server.AccountsDataPath, "index.bleve"))
	os.MkdirAll(filepath.Join(cfg.Server.AccountsDataPath, "accounts"), 0700)

	mapping := bleve.NewIndexMapping()
	// TODO don't bother to store fields as we will load the account from disk
	index, err := bleve.New(filepath.Join(cfg.Server.AccountsDataPath, "index.bleve"), mapping)
	if err != nil {
		panic(err)
	}
	f, err := os.Open(filepath.Join(cfg.Server.AccountsDataPath, "accounts"))
	if err != nil {
		log.Error().Err(err).Str("dir", filepath.Join(cfg.Server.AccountsDataPath, "accounts")).Msg("could not open acconts folder")
		panic(err)
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Error().Err(err).Str("dir", filepath.Join(cfg.Server.AccountsDataPath, "accounts")).Msg("could not list accounts folder")
		panic(err)
	}
	for _, file := range list {
		path := filepath.Join(cfg.Server.AccountsDataPath, "accounts", file.Name())
		data, err := ioutil.ReadFile(path)
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
		index.Index(a.Id, a)
	}

	// TODO watch folders for new records

	s := Service{
		Config: cfg,
		index:  index,
	}

	return s
}

// Service implements the AccountsServiceHandler interface
type Service struct {
	Config *config.Config
	index  bleve.Index
}

// the auth request is currently hardcoded and has to macth this regex
// userName eq \"teddy\" and password eq \"F&1!b90t111!\"
var authQuery = regexp.MustCompile(`^username eq '(.*)' and password eq '(.*)'$`) // TODO how is ' escaped in the password?

// ListAccounts implements the AccountsServiceHandler interface
// the query contains account properties
// TODO id vs onpremiseimmutableid
func (s Service) ListAccounts(ctx context.Context, in *proto.ListAccountsRequest, res *proto.ListAccountsResponse) (err error) {

	var query query.Query

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
	searchResult, err := s.index.Search(searchRequest)
	log.Debug().Interface("result", searchResult).Msg("result")

	res.Accounts = make([]*proto.Account, 0)

	for _, hit := range searchResult.Hits {
		path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", hit.ID)

		data, err := ioutil.ReadFile(path)
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
		res.Accounts = append(res.Accounts, &a)
	}

	return nil
}

// GetAccount implements the AccountsServiceHandler interface
func (s Service) GetAccount(c context.Context, req *proto.GetAccountRequest, res *proto.Account) (err error) {
	return errors.New("not implemented")
}

// CreateAccount implements the AccountsServiceHandler interface
func (s Service) CreateAccount(c context.Context, req *proto.CreateAccountRequest, res *proto.Account) (err error) {
	return errors.New("not implemented")
}

// UpdateAccount implements the AccountsServiceHandler interface
func (s Service) UpdateAccount(c context.Context, req *proto.UpdateAccountRequest, res *proto.Account) (err error) {
	return errors.New("not implemented")
}

// DeleteAccount implements the AccountsServiceHandler interface
func (s Service) DeleteAccount(c context.Context, req *proto.DeleteAccountRequest, res *empty.Empty) (err error) {
	return errors.New("not implemented")
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

package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/CiscoM31/godata"
	"github.com/golang/protobuf/ptypes/empty"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-accounts/pkg/provider"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/rs/zerolog/log"
	"gopkg.in/ldap.v2"
)

// New returns a new instance of Service
func New(cfg *config.Config) Service {
	s := Service{
		Config: cfg,
	}

	return s
}

// Service implements the AccountsServiceHandler interface
type Service struct {
	Config *config.Config
}

// ListAccounts implements the AccountsServiceHandler interface
func (s Service) ListAccounts(ctx context.Context, in *proto.ListAccountsRequest, res *proto.ListAccountsResponse) (err error) {

	log.Debug().Str("query", in.Query).Int32("page-size", in.PageSize).Str("page-token", in.PageToken).Msg("ListAccounts")

	filter := "(&)" // see Absolute True and False Filters in https://tools.ietf.org/html/rfc4526#section-2

	if in.Query != "" {
		// parse the query like an odata filter
		var q *godata.GoDataFilterQuery
		if q, err = godata.ParseFilterString(in.Query); err != nil {
			return err
		}

		// convert to ldap filter
		filter, err = provider.BuildLDAPFilter(q, &s.Config.LDAP.Schema)
		if err != nil {
			return err
		}
	}

	log.Debug().Str("filter", filter).Msg("using filter")

	var l *ldap.Conn
	l, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", s.Config.LDAP.Hostname, s.Config.LDAP.Port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(s.Config.LDAP.BindDN, s.Config.LDAP.BindPassword)
	if err != nil {
		return err
	}

	// TODO combine the parsed query with a query filter from the config, eg. fmt.Sprintf(s.Config.LDAP.UserFilter, clientID)

	// Search for the given clientID
	searchRequest := ldap.NewSearchRequest(
		s.Config.LDAP.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn", s.Config.LDAP.Schema.AccountID, s.Config.LDAP.Schema.Username, s.Config.LDAP.Schema.DisplayName, s.Config.LDAP.Schema.Mail, s.Config.LDAP.Schema.Groups}, // TODO Groups, Identities?
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return err
	}

	log.Debug().Interface("entries", sr.Entries).Msg("entries")

	res.Accounts = make([]*proto.Account, 0)
	for i := range sr.Entries {
		res.Accounts = append(res.Accounts, &proto.Account{
			AccountId: sr.Entries[i].GetAttributeValue(s.Config.LDAP.Schema.AccountID),
			// TODO identities
			Username:    sr.Entries[i].GetAttributeValue(s.Config.LDAP.Schema.Username),
			DisplayName: sr.Entries[i].GetAttributeValue(s.Config.LDAP.Schema.DisplayName),
			Mail:        sr.Entries[i].GetAttributeValue(s.Config.LDAP.Schema.Mail),
			//Groups:      sr.Entries[i].GetAttributeValues(s.Config.LDAP.Schema.Groups),
		})
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

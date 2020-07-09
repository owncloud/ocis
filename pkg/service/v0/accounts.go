package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/CiscoM31/godata"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-accounts/pkg/provider"
	"github.com/tredoe/osutil/user/crypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	// register crypt functions
	_ "github.com/tredoe/osutil/user/crypt/apr1_crypt"
	_ "github.com/tredoe/osutil/user/crypt/md5_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha256_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

func (s Service) indexAccounts(path string) (err error) {

	var f *os.File
	if f, err = os.Open(path); err != nil {
		s.log.Error().Err(err).Str("dir", path).Msg("could not open accounts folder")
		return
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		s.log.Error().Err(err).Str("dir", path).Msg("could not list accounts folder")
		return
	}
	for _, file := range list {
		a := &proto.Account{}
		if err = s.loadAccount(file.Name(), a); err != nil {
			s.log.Error().Err(err).Str("account", file.Name()).Msg("could not load account")
			continue
		}
		s.log.Debug().Interface("account", a).Msg("found account")
		if err = s.index.Index(a.Id, a); err != nil {
			s.log.Error().Err(err).Interface("account", a).Msg("could not index account")
			continue
		}
	}

	return
}

// an auth request is currently hardcoded and has to match this regex
// login eq \"teddy\" and password eq \"F&1!b90t111!\"
var authQuery = regexp.MustCompile(`^login eq '(.*)' and password eq '(.*)'$`) // TODO how is ' escaped in the password?

func (s Service) loadAccount(id string, a *proto.Account) (err error) {
	path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", id)

	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return merrors.NotFound(s.id, "could not read account: %v", err.Error())
	}

	if err = json.Unmarshal(data, a); err != nil {
		return merrors.InternalServerError(s.id, "could not unmarshal account: %v", err.Error())
	}
	return
}

// loggableAccount redacts the password from the account
func loggableAccount(a *proto.Account) *proto.Account {
	if a != nil && a.PasswordProfile != nil {
		a.PasswordProfile.Password = "***REMOVED***"
	}
	return a
}

func (s Service) writeAccount(a *proto.Account) (err error) {

	// leave only the group id
	s.deflateMemberOf(a)

	var bytes []byte
	if bytes, err = json.Marshal(a); err != nil {
		return merrors.InternalServerError(s.id, "could not marshal account: %v", err.Error())
	}

	path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", a.Id)
	if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
		return merrors.InternalServerError(s.id, "could not write account: %v", err.Error())
	}
	return
}

func (s Service) expandMemberOf(a *proto.Account) {
	if a == nil {
		return
	}
	expanded := []*proto.Group{}
	for i := range a.MemberOf {
		g := &proto.Group{}
		// TODO resolve by name, when a create or update is issued they may not have an id? fall back to searching the group id in the index?
		if err := s.loadGroup(a.MemberOf[i].Id, g); err == nil {
			g.Members = nil // always hide members when expanding
			expanded = append(expanded, g)
		} else {
			// log errors but continue execution for now
			s.log.Error().Err(err).Str("id", a.MemberOf[i].Id).Msg("could not load group")
		}
	}
	a.MemberOf = expanded
}

// deflateMemberOf replaces the groups of a user with an instance that only contains the id
func (s Service) deflateMemberOf(a *proto.Account) {
	if a == nil {
		return
	}
	deflated := []*proto.Group{}
	for i := range a.MemberOf {
		if a.MemberOf[i].Id != "" {
			deflated = append(deflated, &proto.Group{Id: a.MemberOf[i].Id})
		} else {
			// TODO fetch and use an id when group only has a name but no id
			s.log.Error().Str("id", a.Id).Interface("group", a.MemberOf[i]).Msg("resolving groups by name is not implemented yet")
		}
	}
	a.MemberOf = deflated
}

func (s Service) passwordIsValid(hash string, pwd string) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Error().Err(fmt.Errorf("%s", r)).Str("hash", hash).Msg("password lib panicked")
		}
	}()

	c := crypt.NewFromHash(hash)
	return c.Verify(hash, []byte(pwd)) == nil
}

// ListAccounts implements the AccountsServiceHandler interface
// the query contains account properties
func (s Service) ListAccounts(ctx context.Context, in *proto.ListAccountsRequest, out *proto.ListAccountsResponse) (err error) {

	var password string

	var query query.Query

	// check if this looks like an auth request
	match := authQuery.FindStringSubmatch(in.Query)
	if len(match) == 3 {
		in.Query = fmt.Sprintf("preferred_name eq '%s'", match[1]) // todo fetch email? make query configurable
		password = match[2]
		if password == "" {
			return merrors.Unauthorized(s.id, "password must not be empty")
		}
	}

	if in.Query != "" {
		// parse the query like an odata filter
		var q *godata.GoDataFilterQuery
		if q, err = godata.ParseFilterString(in.Query); err != nil {
			s.log.Error().Err(err).Msg("could not parse query")
			return merrors.InternalServerError(s.id, "could not parse query: %v", err.Error())
		}

		// convert to bleve query
		query, err = provider.BuildBleveQuery(q)
		if err != nil {
			s.log.Error().Err(err).Msg("could not build bleve query")
			return merrors.InternalServerError(s.id, "could not build bleve query: %v", err.Error())
		}
	} else {
		query = bleve.NewMatchAllQuery()
	}

	s.log.Debug().Interface("query", query).Msg("using query")

	searchRequest := bleve.NewSearchRequest(query)
	var searchResult *bleve.SearchResult
	searchResult, err = s.index.Search(searchRequest)
	if err != nil {
		s.log.Error().Err(err).Msg("could not execute bleve search")
		return merrors.InternalServerError(s.id, "could not execute bleve search: %v", err.Error())
	}

	s.log.Debug().Interface("result", searchResult).Msg("result")

	out.Accounts = make([]*proto.Account, 0)

	for _, hit := range searchResult.Hits {

		a := &proto.Account{}
		if err = s.loadAccount(hit.ID, a); err != nil {
			s.log.Error().Err(err).Str("account", hit.ID).Msg("could not load account, skipping")
			continue
		}
		var currentHash string
		if a.PasswordProfile != nil {
			currentHash = a.PasswordProfile.Password
		}
		s.log.Debug().Interface("account", loggableAccount(a)).Msg("found account")

		if password != "" {
			if a.PasswordProfile == nil {
				s.log.Debug().Interface("account", loggableAccount(a)).Msg("no password profile")
				return merrors.Unauthorized(s.id, "invalid password")
			}
			if !s.passwordIsValid(currentHash, password) {
				return merrors.Unauthorized(s.id, "invalid password")
			}
		}
		// TODO add groups if requested
		// if in.FieldMask ...
		s.expandMemberOf(a)

		// remove password before returning
		a.PasswordProfile.Password = ""

		out.Accounts = append(out.Accounts, a)
	}

	return
}

// GetAccount implements the AccountsServiceHandler interface
func (s Service) GetAccount(c context.Context, in *proto.GetAccountRequest, out *proto.Account) (err error) {
	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	if err = s.loadAccount(id, out); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not load account")
		return
	}
	s.log.Debug().Interface("account", loggableAccount(out)).Msg("found account")

	// TODO add groups if requested
	// if in.FieldMask ...
	s.expandMemberOf(out)

	// remove password
	out.PasswordProfile.Password = ""

	return
}

// CreateAccount implements the AccountsServiceHandler interface
func (s Service) CreateAccount(c context.Context, in *proto.CreateAccountRequest, out *proto.Account) (err error) {
	var id string
	if in.Account == nil {
		return merrors.BadRequest(s.id, "account missing")
	}
	if in.Account.Id == "" {
		in.Account.Id = uuid.Must(uuid.NewV4()).String()
	}

	if id, err = cleanupID(in.Account.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}
	path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", id)

	if in.Account.PasswordProfile != nil && in.Account.PasswordProfile.Password != "" {
		// encrypt password
		c := crypt.New(crypt.SHA512)
		if in.Account.PasswordProfile.Password, err = c.Generate([]byte(in.Account.PasswordProfile.Password), nil); err != nil {
			s.log.Error().Err(err).Str("id", id).Interface("account", loggableAccount(in.Account)).Msg("could not hash password")
			return merrors.InternalServerError(s.id, "could not hash password: %v", err.Error())
		}
	}

	// extract group id
	// TODO groups should be ignored during create, use groups.AddMember? return error?
	if err = s.writeAccount(in.Account); err != nil {
		s.log.Error().Err(err).Interface("account", loggableAccount(in.Account)).Msg("could not persist new account")
		return
	}

	if err = s.index.Index(id, in.Account); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("path", path).Interface("account", loggableAccount(in.Account)).Msg("could not index new account")
		return merrors.InternalServerError(s.id, "could not index new account: %v", err.Error())
	}

	return
}

// UpdateAccount implements the AccountsServiceHandler interface
// read only fields are ignored
// TODO how can we unset specific values? using the update mask
func (s Service) UpdateAccount(c context.Context, in *proto.UpdateAccountRequest, out *proto.Account) (err error) {
	var id string
	if in.Account == nil {
		return merrors.BadRequest(s.id, "account missing")
	}
	if in.Account.Id == "" {
		return merrors.BadRequest(s.id, "account id missing")
	}

	if id, err = cleanupID(in.Account.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}
	path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", id)

	if err = s.loadAccount(id, out); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not load account")
		return
	}
	s.log.Debug().Interface("account", loggableAccount(out)).Msg("found account")

	t := time.Now()
	tsnow := &timestamppb.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}

	// id read-only
	out.AccountEnabled = in.Account.AccountEnabled
	out.IsResourceAccount = in.Account.IsResourceAccount
	// creation-type read only
	out.Identities = in.Account.Identities
	out.DisplayName = in.Account.DisplayName
	out.PreferredName = in.Account.PreferredName
	out.UidNumber = in.Account.UidNumber
	out.GidNumber = in.Account.GidNumber
	out.Mail = in.Account.Mail // read only?
	out.Description = in.Account.Description

	if in.Account.PasswordProfile != nil && in.Account.PasswordProfile.Password != "" {
		// encrypt password
		c := crypt.New(crypt.SHA512)
		if out.PasswordProfile.Password, err = c.Generate([]byte(in.Account.PasswordProfile.Password), nil); err != nil {
			s.log.Error().Err(err).Str("id", id).Interface("account", loggableAccount(in.Account)).Msg("could not hash password")
			return merrors.InternalServerError(s.id, "could not hash password: %v", err.Error())
		}
		out.PasswordProfile.LastPasswordChangeDateTime = tsnow
	}
	// lastPasswordChangeDateTime calculated, see password
	out.PasswordProfile.PasswordPolicies = in.Account.PasswordProfile.PasswordPolicies
	out.PasswordProfile.ForceChangePasswordNextSignIn = in.Account.PasswordProfile.ForceChangePasswordNextSignIn
	out.PasswordProfile.ForceChangePasswordNextSignInWithMfa = in.Account.PasswordProfile.ForceChangePasswordNextSignInWithMfa

	// memberOf read only
	// createdDateTime read only
	// deleteDateTime read only

	out.OnPremisesSyncEnabled = in.Account.OnPremisesSyncEnabled
	// ... TODO on prem for sync

	if out.ExternalUserState != in.Account.ExternalUserState {
		out.ExternalUserState = in.Account.ExternalUserState
		out.ExternalUserStateChangeDateTime = tsnow
	}
	// out.RefreshTokensValidFromDateTime TODO use to invalidate all existing sessions
	// out.SignInSessionsValidFromDateTime TODO use to invalidate all existing sessions

	if err = s.writeAccount(out); err != nil {
		s.log.Error().Err(err).Interface("account", loggableAccount(out)).Msg("could not persist updated account")
		return
	}

	if err = s.index.Index(id, out); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("path", path).Interface("account", loggableAccount(out)).Msg("could not index new account")
		return merrors.InternalServerError(s.id, "could not index updated account: %v", err.Error())
	}

	// remove password
	out.PasswordProfile.Password = ""

	return
}

// DeleteAccount implements the AccountsServiceHandler interface
func (s Service) DeleteAccount(c context.Context, in *proto.DeleteAccountRequest, out *empty.Empty) (err error) {
	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}
	path := filepath.Join(s.Config.Server.AccountsDataPath, "accounts", id)

	a := &proto.Account{}
	if err = s.loadAccount(id, a); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not load account")
		return
	}

	// delete member relationship in groups
	for i := range a.MemberOf {
		err = s.RemoveMember(c, &proto.RemoveMemberRequest{
			GroupId:   a.MemberOf[i].Id,
			AccountId: id,
		}, a.MemberOf[i])
		if err != nil {
			s.log.Error().Err(err).Str("accountid", id).Str("groupid", a.MemberOf[i].Id).Msg("could not remove group membership")
		}
	}

	if err = os.Remove(path); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not remove account")
		return merrors.InternalServerError(s.id, "could not remove account: %v", err.Error())
	}

	if err = s.index.Delete(id); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not remove account from index")
		return merrors.InternalServerError(s.id, "could not remove account from index: %v", err.Error())
	}

	s.log.Info().Str("id", id).Msg("deleted account")
	return
}

package service

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/CiscoM31/godata"
	"github.com/blevesearch/bleve"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/accounts/pkg/provider"
	"github.com/owncloud/ocis/accounts/pkg/storage"
	"github.com/owncloud/ocis/ocis-pkg/roles"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	settings_svc "github.com/owncloud/ocis/settings/pkg/service/v0"
	"github.com/rs/zerolog"
	"github.com/tredoe/osutil/user/crypt"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/timestamppb"

	// register crypt functions
	_ "github.com/tredoe/osutil/user/crypt/apr1_crypt"
	_ "github.com/tredoe/osutil/user/crypt/md5_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha256_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

// accLock mutually exclude readers from writers on account files
var accLock sync.Mutex

func (s Service) indexAccount(id string) error {
	a := &proto.BleveAccount{
		BleveType: "account",
	}
	if err := s.repo.LoadAccount(context.Background(), id, &a.Account); err != nil {
		s.log.Error().Err(err).Str("account", id).Msg("could not load account")
		return err
	}
	s.log.Debug().Interface("account", a).Msg("found account")
	if err := s.index.Index(a.Id, a); err != nil {
		s.log.Error().Err(err).Interface("account", a).Msg("could not index account")
		return err
	}
	return nil
}

// an auth request is currently hardcoded and has to match this regex
// login eq \"teddy\" and password eq \"F&1!b90t111!\"
var authQuery = regexp.MustCompile(`^login eq '(.*)' and password eq '(.*)'$`) // TODO how is ' escaped in the password?

func (s Service) expandMemberOf(a *proto.Account) {
	if a == nil {
		return
	}
	expanded := []*proto.Group{}
	for i := range a.MemberOf {
		g := &proto.Group{}
		// TODO resolve by name, when a create or update is issued they may not have an id? fall back to searching the group id in the index?
		if err := s.repo.LoadGroup(context.Background(), a.MemberOf[i].Id, g); err == nil {
			g.Members = nil // always hide members when expanding
			expanded = append(expanded, g)
		} else {
			// log errors but continue execution for now
			s.log.Error().Err(err).Str("id", a.MemberOf[i].Id).Msg("could not load group")
		}
	}
	a.MemberOf = expanded
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

func (s Service) hasAccountManagementPermissions(ctx context.Context) bool {
	// get roles from context
	roleIDs, ok := roles.ReadRoleIDsFromContext(ctx)
	if !ok {
		/**
		 * FIXME: with this we are skipping permission checks on all requests that are coming in without roleIDs in the
		 * metadata context. This is a huge security impairment, as that's the case not only for grpc requests but also
		 * for unauthenticated http requests and http requests coming in without hitting the ocis-proxy first.
		 */
		return true
	}

	// check if permission is present in roles of the authenticated account
	return s.RoleManager.FindPermissionByID(ctx, roleIDs, AccountManagementPermissionID) != nil
}

// ListAccounts implements the AccountsServiceHandler interface
// the query contains account properties
func (s Service) ListAccounts(ctx context.Context, in *proto.ListAccountsRequest, out *proto.ListAccountsResponse) (err error) {
	if !s.hasAccountManagementPermissions(ctx) {
		return merrors.Forbidden(s.id, "no permission for ListAccounts")
	}

	accLock.Lock()
	defer accLock.Unlock()
	var password string

	// check if this looks like an auth request
	match := authQuery.FindStringSubmatch(in.Query)
	if len(match) == 3 {
		in.Query = fmt.Sprintf("on_premises_sam_account_name eq '%s'", match[1]) // todo fetch email? make query configurable
		password = match[2]
		if password == "" {
			return merrors.Unauthorized(s.id, "password must not be empty")
		}

		// hardcoded check against service user
		if s.Config.ServiceUser.Username != "" &&
			strings.EqualFold(match[1], s.Config.ServiceUser.Username) &&
			match[2] == s.Config.ServiceUser.Password {
			out.Accounts = []*proto.Account{
				{
					Id:             "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
					AccountEnabled: true,
					PreferredName:  s.Config.ServiceUser.Username,
					DisplayName:    s.Config.ServiceUser.Username,
					UidNumber:      s.Config.ServiceUser.UID,
					GidNumber:      s.Config.ServiceUser.GID,
				},
			}
			return nil
		}
	}

	// only search for accounts
	tq := bleve.NewTermQuery("account")
	tq.SetField("bleve_type")

	query := bleve.NewConjunctionQuery(tq)

	if in.Query != "" {
		// parse the query like an odata filter
		var q *godata.GoDataFilterQuery
		if q, err = godata.ParseFilterString(in.Query); err != nil {
			s.log.Error().Err(err).Msg("could not parse query")
			return merrors.InternalServerError(s.id, "could not parse query: %v", err.Error())
		}

		// convert to bleve query
		bq, err := provider.BuildBleveQuery(q)
		if err != nil {
			s.log.Error().Err(err).Msg("could not build bleve query")
			return merrors.InternalServerError(s.id, "could not build bleve query: %v", err.Error())
		}
		query.AddQuery(bq)
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
		if err = s.repo.LoadAccount(ctx, hit.ID, a); err != nil {
			s.log.Error().Err(err).Str("account", hit.ID).Msg("could not load account, skipping")
			continue
		}
		var currentHash string
		if a.PasswordProfile != nil {
			currentHash = a.PasswordProfile.Password
		}

		s.debugLogAccount(a).Msg("found account")

		if password != "" {
			if a.PasswordProfile == nil {
				s.debugLogAccount(a).Msg("no password profile")
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
		if a.PasswordProfile != nil {
			a.PasswordProfile.Password = ""
		}

		out.Accounts = append(out.Accounts, a)
	}

	return
}

// GetAccount implements the AccountsServiceHandler interface
func (s Service) GetAccount(ctx context.Context, in *proto.GetAccountRequest, out *proto.Account) (err error) {
	if !s.hasAccountManagementPermissions(ctx) {
		return merrors.Forbidden(s.id, "no permission for GetAccount")
	}

	accLock.Lock()
	defer accLock.Unlock()
	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	if err = s.repo.LoadAccount(ctx, id, out); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "account not found: %v", err.Error())
		}

		s.log.Error().Err(err).Str("id", id).Msg("could not load account")
		return merrors.InternalServerError(s.id, "could not load account: %v", err.Error())
	}

	s.debugLogAccount(out).Msg("found account")

	// TODO add groups if requested
	// if in.FieldMask ...
	s.expandMemberOf(out)

	// remove password
	if out.PasswordProfile != nil {
		out.PasswordProfile.Password = ""
	}

	return
}

// CreateAccount implements the AccountsServiceHandler interface
func (s Service) CreateAccount(ctx context.Context, in *proto.CreateAccountRequest, out *proto.Account) (err error) {
	if !s.hasAccountManagementPermissions(ctx) {
		return merrors.Forbidden(s.id, "no permission for CreateAccount")
	}

	accLock.Lock()
	defer accLock.Unlock()
	var id string
	var acc = in.Account
	if acc == nil {
		return merrors.BadRequest(s.id, "account missing")
	}
	if acc.Id == "" {
		acc.Id = uuid.Must(uuid.NewV4()).String()
	}
	if !s.isValidUsername(acc.PreferredName) {
		return merrors.BadRequest(s.id, "preferred_name '%s' must be at least the local part of an email", acc.PreferredName)
	}
	if !s.isValidEmail(acc.Mail) {
		return merrors.BadRequest(s.id, "mail '%s' must be a valid email", acc.Mail)
	}

	if id, err = cleanupID(acc.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	if acc.PasswordProfile != nil {
		if acc.PasswordProfile.Password != "" {
			// encrypt password
			c := crypt.New(crypt.SHA512)
			if acc.PasswordProfile.Password, err = c.Generate([]byte(acc.PasswordProfile.Password), nil); err != nil {
				s.log.Error().Err(err).Str("id", id).Msg("could not hash password")
				return merrors.InternalServerError(s.id, "could not hash password: %v", err.Error())
			}
		}

		if err := passwordPoliciesValid(acc.PasswordProfile.PasswordPolicies); err != nil {
			return merrors.BadRequest(s.id, "%s", err)
		}
	}

	// extract group id
	// TODO groups should be ignored during create, use groups.AddMember? return error?

	// write and index account - note: don't do anything else in between!
	if err = s.repo.WriteAccount(ctx, acc); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not persist new account")
		s.debugLogAccount(acc).Msg("could not persist new account")
		return merrors.InternalServerError(s.id, "could not persist new account: %v", err.Error())
	}
	if err = s.indexAccount(acc.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not index new account: %v", err.Error())
	}
	s.log.Debug().Interface("account", acc).Msg("account after indexing")

	if acc.PasswordProfile != nil {
		acc.PasswordProfile.Password = ""
	}

	{
		out.Id = acc.Id
		out.Mail = acc.Mail
		out.PreferredName = acc.PreferredName
		out.AccountEnabled = acc.AccountEnabled
		out.DisplayName = acc.DisplayName
		out.OnPremisesSamAccountName = acc.OnPremisesSamAccountName
	}

	// TODO: assign user role to all new users for now, as create Account request does not have any role field
	if s.RoleService == nil {
		return merrors.InternalServerError(s.id, "could not assign role to account: roleService not configured")
	}
	if _, err = s.RoleService.AssignRoleToUser(ctx, &settings.AssignRoleToUserRequest{
		AccountUuid: acc.Id,
		RoleId:      settings_svc.BundleUUIDRoleUser,
	}); err != nil {
		return merrors.InternalServerError(s.id, "could not assign role to account: %v", err.Error())
	}

	return
}

// UpdateAccount implements the AccountsServiceHandler interface
// read only fields are ignored
// TODO how can we unset specific values? using the update mask
func (s Service) UpdateAccount(ctx context.Context, in *proto.UpdateAccountRequest, out *proto.Account) (err error) {
	if !s.hasAccountManagementPermissions(ctx) {
		return merrors.Forbidden(s.id, "no permission for UpdateAccount")
	}

	accLock.Lock()
	defer accLock.Unlock()
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

	if err = s.repo.LoadAccount(ctx, id, out); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "account not found: %v", err.Error())
		}

		s.log.Error().Err(err).Str("id", id).Msg("could not load account")
		return merrors.InternalServerError(s.id, "could not load account: %v", err.Error())

	}

	t := time.Now()
	tsnow := &timestamppb.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}

	validMask, err := validateUpdate(in.UpdateMask, updatableAccountPaths)
	if err != nil {
		return merrors.BadRequest(s.id, "%s", err)
	}

	if err := fieldmask_utils.StructToStruct(validMask, in.Account, out); err != nil {
		return merrors.InternalServerError(s.id, "%s", err)
	}

	if in.Account.PasswordProfile != nil {
		if out.PasswordProfile == nil {
			out.PasswordProfile = &proto.PasswordProfile{}
		}
		if in.Account.PasswordProfile.Password != "" {
			// encrypt password
			c := crypt.New(crypt.SHA512)
			if out.PasswordProfile.Password, err = c.Generate([]byte(in.Account.PasswordProfile.Password), nil); err != nil {
				in.Account.PasswordProfile.Password = ""
				s.log.Error().Err(err).Str("id", id).Msg("could not hash password")
				return merrors.InternalServerError(s.id, "could not hash password: %v", err.Error())
			}

			in.Account.PasswordProfile.Password = ""
		}

		if err := passwordPoliciesValid(in.Account.PasswordProfile.PasswordPolicies); err != nil {
			return merrors.BadRequest(s.id, "%s", err)
		}

		// lastPasswordChangeDateTime calculated, see password
		out.PasswordProfile.LastPasswordChangeDateTime = tsnow
	}

	// out.RefreshTokensValidFromDateTime TODO use to invalidate all existing sessions
	// out.SignInSessionsValidFromDateTime TODO use to invalidate all existing sessions

	// ... TODO on prem for sync

	if out.ExternalUserState != in.Account.ExternalUserState {
		out.ExternalUserState = in.Account.ExternalUserState
		out.ExternalUserStateChangeDateTime = tsnow
	}

	if err = s.repo.WriteAccount(ctx, out); err != nil {
		s.log.Error().Err(err).Str("id", out.Id).Msg("could not persist updated account")
		return merrors.InternalServerError(s.id, "could not persist updated account: %v", err.Error())
	}

	if err = s.indexAccount(id); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not index new account")
		return merrors.InternalServerError(s.id, "could not index updated account: %v", err.Error())
	}

	// remove password
	if out.PasswordProfile != nil {
		out.PasswordProfile.Password = ""
	}

	return
}

// whitelist of all paths/fields which can be updated by clients
var updatableAccountPaths = map[string]struct{}{
	"AccountEnabled":                   {},
	"IsResourceAccount":                {},
	"Identities":                       {},
	"DisplayName":                      {},
	"PreferredName":                    {},
	"UidNumber":                        {},
	"GidNumber":                        {},
	"Description":                      {},
	"Mail":                             {}, // read only?,
	"PasswordProfile.Password":         {},
	"PasswordProfile.PasswordPolicies": {},
	"PasswordProfile.ForceChangePasswordNextSignIn":        {},
	"PasswordProfile.ForceChangePasswordNextSignInWithMfa": {},
	"OnPremisesSyncEnabled":                                {},
	"OnPremisesSamAccountName":                             {},
}

// DeleteAccount implements the AccountsServiceHandler interface
func (s Service) DeleteAccount(ctx context.Context, in *proto.DeleteAccountRequest, out *empty.Empty) (err error) {
	if !s.hasAccountManagementPermissions(ctx) {
		return merrors.Forbidden(s.id, "no permission for DeleteAccount")
	}

	accLock.Lock()
	defer accLock.Unlock()
	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	a := &proto.Account{}
	if err = s.repo.LoadAccount(ctx, id, a); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "account not found: %v", err.Error())
		}

		s.log.Error().Err(err).Str("id", id).Msg("could not load account")
		return merrors.InternalServerError(s.id, "could not load account: %v", err.Error())
	}

	// delete member relationship in groups
	for i := range a.MemberOf {
		err = s.RemoveMember(ctx, &proto.RemoveMemberRequest{
			GroupId:   a.MemberOf[i].Id,
			AccountId: id,
		}, a.MemberOf[i])
		if err != nil {
			s.log.Error().Err(err).Str("accountid", id).Str("groupid", a.MemberOf[i].Id).Msg("could not remove group member, skipping")
		}
	}

	if err = s.repo.DeleteAccount(ctx, id); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "account not found: %v", err.Error())
		}

		s.log.Error().Err(err).Str("id", id).Str("accountId", id).Msg("could not remove account")
		return merrors.InternalServerError(s.id, "could not remove account: %v", err.Error())
	}

	if err = s.index.Delete(id); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("accountId", id).Msg("could not remove account from index")
		return merrors.InternalServerError(s.id, "could not remove account from index: %v", err.Error())
	}

	s.log.Info().Str("id", id).Msg("deleted account")
	return
}

// We want to allow email addresses as usernames so they show up when using them in ACLs on storages that allow intergration with our glauth LDAP service
// so we are adding a few restrictions from https://stackoverflow.com/questions/6949667/what-are-the-real-rules-for-linux-usernames-on-centos-6-and-rhel-6
// names should not start with numbers
var usernameRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]*(@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)*$")

func (s Service) isValidUsername(e string) bool {
	if len(e) < 1 && len(e) > 254 {
		return false
	}
	return usernameRegex.MatchString(e)
}

// regex from https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#valid-e-mail-address
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (s Service) isValidEmail(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

const (
	policyDisableStrongPassword     = "DisableStrongPassword"
	policyDisablePasswordExpiration = "DisablePasswordExpiration"
)

func passwordPoliciesValid(policies []string) error {
	for _, v := range policies {
		if v != policyDisableStrongPassword && v != policyDisablePasswordExpiration {
			return fmt.Errorf("invalid password-policy %s", v)
		}
	}

	return nil
}

// validateUpdate takes a update field-mask and validates it against a whitelist of updatable paths.
// Returns a FieldFilter on success which can be passed to the fieldmask_utils..StructToStruct. An error is returned
// if the mask tries to update no whitelisted fields.
//
// Given an empty or nil mask we assume that the client wants to update all whitelisted fields.
//
func validateUpdate(mask *field_mask.FieldMask, updatablePaths map[string]struct{}) (fieldmask_utils.FieldFilterContainer, error) {
	nop := func(s string) string { return s }
	// Assume that the client wants to update all updatable path if
	// no field-mask is given, so we create a mask with all paths
	if mask == nil || len(mask.Paths) == 0 {
		paths := make([]string, 0, len(updatablePaths))
		for fieldName := range updatablePaths {
			paths = append(paths, fieldName)
		}

		return fieldmask_utils.MaskFromPaths(paths, nop)
	}

	// Check that only allowed fields are updated
	for _, v := range mask.Paths {
		if _, ok := updatablePaths[v]; !ok {
			return nil, fmt.Errorf("can not update field %s, either unknown or readonly", v)
		}
	}

	return fieldmask_utils.MaskFromPaths(mask.Paths, nop)
}

// debugLogAccount returns a debug-log event with detailed account-info, and filtered password data
func (s Service) debugLogAccount(a *proto.Account) *zerolog.Event {
	return s.log.Debug().Fields(map[string]interface{}{
		"Id":                           a.Id,
		"Mail":                         a.Mail,
		"DisplayName":                  a.DisplayName,
		"AccountEnabled":               a.AccountEnabled,
		"IsResourceAccount":            a.IsResourceAccount,
		"Identities":                   a.Identities,
		"PreferredName":                a.PreferredName,
		"UidNumber":                    a.UidNumber,
		"GidNumber":                    a.GidNumber,
		"Description":                  a.Description,
		"OnPremisesSyncEnabled":        a.OnPremisesSyncEnabled,
		"OnPremisesSamAccountName":     a.OnPremisesSamAccountName,
		"OnPremisesUserPrincipalName":  a.OnPremisesUserPrincipalName,
		"OnPremisesSecurityIdentifier": a.OnPremisesSecurityIdentifier,
		"OnPremisesDistinguishedName":  a.OnPremisesDistinguishedName,
		"OnPremisesLastSyncDateTime":   a.OnPremisesLastSyncDateTime,
		"MemberOf":                     a.MemberOf,
		"CreatedDateTime":              a.CreatedDateTime,
		"DeletedDateTime":              a.DeletedDateTime,
	})
}

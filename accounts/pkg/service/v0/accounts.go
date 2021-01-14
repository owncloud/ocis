package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/owncloud/ocis/ocis-pkg/cache"
	"golang.org/x/crypto/bcrypt"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/log"

	"github.com/gofrs/uuid"
	p "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/accounts/pkg/storage"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/roles"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	settings_svc "github.com/owncloud/ocis/settings/pkg/service/v0"
	"github.com/rs/zerolog"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// passwordValidCache caches basic auth password validations
var passwordValidCache = cache.NewCache(1024)

// passwordValidCacheExpiration defines the entry lifetime
const passwordValidCacheExpiration = 10 * time.Minute

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

func (s Service) hasAccountManagementPermissions(ctx context.Context) bool {
	// get roles from context
	roleIDs, ok := roles.ReadRoleIDsFromContext(ctx)
	if !ok {
		/**
		* FIXME: with this we are skipping permission checks on all requests that are coming in without roleIDs in the
		* metadata context. This is a huge security impairment, as that's the case not only for grpc requests but also
		* for unauthenticated http requests and http requests coming in without hitting the ocis-proxy first.
		 */
		// TODO add system role for internal requests.
		// - at least the proxy needs to look up account info
		// - glauth needs to make bind requests
		// tracked as OCIS-454
		return true
	}

	// check if permission is present in roles of the authenticated account
	return s.RoleManager.FindPermissionByID(ctx, roleIDs, AccountManagementPermissionID) != nil
}

func (s Service) hasSelfManagementPermissions(ctx context.Context) bool {
	// get roles from context
	roleIDs, ok := roles.ReadRoleIDsFromContext(ctx)
	if !ok {
		return false
	}

	// check if permission is present in roles of the authenticated account
	return s.RoleManager.FindPermissionByID(ctx, roleIDs, SelfManagementPermissionID) != nil
}

// serviceUserToIndex temporarily adds a service user to the index, which is supposed to be removed before the lock on the handler function is released
func (s Service) serviceUserToIndex() (teardownServiceUser func()) {
	if s.Config.ServiceUser.Username != "" && s.Config.ServiceUser.UUID != "" {
		_, err := s.index.Add(s.getInMemoryServiceUser())
		if err != nil {
			s.log.Logger.Err(err).Msg("service user was configured but failed to be added to the index")
		} else {
			return func() {
				_ = s.index.Delete(s.getInMemoryServiceUser())
			}
		}
	}
	return func() {}
}

func (s Service) getInMemoryServiceUser() proto.Account {
	return proto.Account{
		AccountEnabled:           true,
		Id:                       s.Config.ServiceUser.UUID,
		PreferredName:            s.Config.ServiceUser.Username,
		OnPremisesSamAccountName: s.Config.ServiceUser.Username,
		DisplayName:              s.Config.ServiceUser.Username,
		UidNumber:                s.Config.ServiceUser.UID,
		GidNumber:                s.Config.ServiceUser.GID,
	}
}

// ListAccounts implements the AccountsServiceHandler interface
// the query contains account properties
func (s Service) ListAccounts(ctx context.Context, in *proto.ListAccountsRequest, out *proto.ListAccountsResponse) (err error) {
	hasSelf := s.hasSelfManagementPermissions(ctx)
	hasManagement := s.hasAccountManagementPermissions(ctx)
	if !hasSelf && !hasManagement {
		return merrors.Forbidden(s.id, "no permission for ListAccounts")
	}
	onlySelf := hasSelf && !hasManagement

	teardownServiceUser := s.serviceUserToIndex()
	defer teardownServiceUser()
	match, authRequest := getAuthQueryMatch(in.Query)
	if authRequest {
		password := match[2]
		if len(password) == 0 {
			return merrors.Unauthorized(s.id, "account not found or invalid credentials")
		}

		ids, err := s.index.FindBy(&proto.Account{}, "OnPremisesSamAccountName", match[1])
		if err != nil || len(ids) > 1 {
			return merrors.Unauthorized(s.id, "account not found or invalid credentials")
		}
		if len(ids) == 0 {
			ids, err = s.index.FindBy(&proto.Account{}, "Mail", match[1])
			if err != nil || len(ids) != 1 {
				return merrors.Unauthorized(s.id, "account not found or invalid credentials")
			}
		}

		a := &proto.Account{}
		err = s.repo.LoadAccount(ctx, ids[0], a)
		if err != nil || a.PasswordProfile == nil || len(a.PasswordProfile.Password) == 0 {
			return merrors.Unauthorized(s.id, "account not found or invalid credentials")
		}

		// isPasswordValid uses bcrypt.CompareHashAndPassword which is slow by design.
		// if every request that matches authQuery regex needs to do this step over and over again,
		// this is secure but also slow. In this implementation we keep it same secure but increase the speed.
		//
		// flow:
		// - request comes in
		// - it creates a sha256 based on found account PasswordProfile.LastPasswordChangeDateTime and requested password (v)
		// - it checks if the cache already contains an entry that matches found account Id // account PasswordProfile.LastPasswordChangeDateTime (k)
		// - if no entry exists it runs the bcrypt.CompareHashAndPassword as before and if everything is ok it stores the
		//   result by the (k) as key and (v) as value. If not it errors
		// - if a entry is found it checks if the given value matches (v). If it doesnt match, the cache entry gets removed
		//   and it errors.
		{
			var suspicious bool

			kh := sha256.New()
			kh.Write([]byte(a.Id))
			k := hex.EncodeToString(kh.Sum([]byte(a.PasswordProfile.LastPasswordChangeDateTime.String())))

			vh := sha256.New()
			vh.Write([]byte(a.PasswordProfile.Password))
			v := vh.Sum([]byte(password))

			e := passwordValidCache.Get(k)

			if e == nil {
				suspicious = !isPasswordValid(s.log, a.PasswordProfile.Password, password)
			} else if !bytes.Equal(e.V.([]byte), v) {
				suspicious = true
			}

			if suspicious {
				passwordValidCache.Unset(k)
				return merrors.Unauthorized(s.id, "account not found or invalid credentials")
			}

			if e == nil {
				passwordValidCache.Set(k, v, time.Now().Add(passwordValidCacheExpiration))
			}
		}

		a.PasswordProfile.Password = ""
		out.Accounts = []*proto.Account{a}

		return nil
	}

	if onlySelf {
		// limit list to own account id
		if aid, ok := metadata.Get(ctx, middleware.AccountID); ok {
			in.Query = "id eq '" + aid + "'"
		} else {
			return merrors.InternalServerError(s.id, "account id not in context")
		}
	}

	if in.Query == "" {
		err = s.repo.LoadAccounts(ctx, &out.Accounts)
		if err != nil {
			s.log.Err(err).Msg("failed to load all accounts from storage")
			return merrors.InternalServerError(s.id, "failed to load all accounts")
		}
		for i := range out.Accounts {
			a := out.Accounts[i]

			// TODO add groups only if requested
			// if in.FieldMask ...
			s.expandMemberOf(a)

			if a.PasswordProfile != nil {
				a.PasswordProfile.Password = ""
			}
		}
		return nil
	}

	searchResults, err := s.findAccountsByQuery(ctx, in.Query)
	out.Accounts = make([]*proto.Account, 0, len(searchResults))

	for _, hit := range searchResults {
		a := &proto.Account{}
		if hit == s.Config.ServiceUser.UUID {
			acc := s.getInMemoryServiceUser()
			a = &acc
		} else if err = s.repo.LoadAccount(ctx, hit, a); err != nil {
			s.log.Error().Err(err).Str("account", hit).Msg("could not load account, skipping")
			continue
		}

		s.debugLogAccount(a).Msg("found account")

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

func (s Service) findAccountsByQuery(ctx context.Context, query string) ([]string, error) {
	return s.index.Query(&proto.Account{}, query)
}

// GetAccount implements the AccountsServiceHandler interface
func (s Service) GetAccount(ctx context.Context, in *proto.GetAccountRequest, out *proto.Account) (err error) {
	hasSelf := s.hasSelfManagementPermissions(ctx)
	hasManagement := s.hasAccountManagementPermissions(ctx)
	if !hasSelf && !hasManagement {
		return merrors.Forbidden(s.id, "no permission for GetAccount")
	}
	onlySelf := hasSelf && !hasManagement

	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	if onlySelf {
		// limit get to own account id
		if aid, ok := metadata.Get(ctx, middleware.AccountID); ok {
			if id != aid {
				return merrors.Forbidden(s.id, "no permission for GetAccount of another user")
			}
		} else {
			return merrors.InternalServerError(s.id, "account id not in context")
		}
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

	var id string

	if in.Account == nil {
		return merrors.InternalServerError(s.id, "invalid account: empty")
	}

	p.Merge(out, in.Account)

	if out.Id == "" {
		out.Id = uuid.Must(uuid.NewV4()).String()
	}
	if err = validateAccount(s.id, out); err != nil {
		return err
	}

	if id, err = cleanupID(out.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	exists, err := s.accountExists(ctx, out.PreferredName, out.Mail, out.Id)
	if err != nil {
		return merrors.InternalServerError(s.id, "could not check if account exists: %v", err.Error())
	}
	if exists {
		return merrors.Conflict(s.id, "account already exists")
	}

	if out.PasswordProfile != nil {
		if out.PasswordProfile.Password != "" {
			// encrypt password
			hashed, err := bcrypt.GenerateFromPassword([]byte(in.Account.PasswordProfile.Password), s.Config.Server.HashDifficulty)
			if err != nil {
				s.log.Error().Err(err).Str("id", id).Msg("could not hash password")
				return merrors.InternalServerError(s.id, "could not hash password: %v", err.Error())
			}
			out.PasswordProfile.Password = string(hashed)
			in.Account.PasswordProfile.Password = ""
		}

		if err := passwordPoliciesValid(out.PasswordProfile.PasswordPolicies); err != nil {
			return merrors.BadRequest(s.id, "%s", err)
		}
	}

	// extract group id
	// TODO groups should be ignored during create, use groups.AddMember? return error?

	// write and index account - note: don't do anything else in between!
	if err = s.repo.WriteAccount(ctx, out); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not persist new account")
		s.debugLogAccount(out).Msg("could not persist new account")
		return merrors.InternalServerError(s.id, "could not persist new account: %v", err.Error())
	}
	indexResults, err := s.index.Add(out)
	if err != nil {
		s.rollbackCreateAccount(ctx, out)
		return merrors.Conflict(s.id, "Account already exists %v", err.Error())

	}
	s.log.Debug().Interface("account", out).Msg("account after indexing")

	for _, r := range indexResults {
		if r.Field == "UidNumber" {
			id, err := strconv.Atoi(path.Base(r.Value))
			if err != nil {
				s.rollbackCreateAccount(ctx, out)
				return err
			}
			out.UidNumber = int64(id)
			break
		}
	}

	if out.GidNumber == 0 {
		out.GidNumber = userDefaultGID
	}

	r := proto.ListGroupsResponse{}
	err = s.ListGroups(ctx, &proto.ListGroupsRequest{}, &r)
	if err != nil {
		// rollback account creation
		return err
	}

	for _, group := range r.Groups {
		if group.GidNumber == out.GidNumber {
			out.MemberOf = append(out.MemberOf, group)
		}
	}
	//acc.MemberOf = append(acc.MemberOf, &group)
	if err := s.repo.WriteAccount(context.Background(), out); err != nil {
		return err
	}

	if out.PasswordProfile != nil {
		out.PasswordProfile.Password = ""
	}

	// TODO: assign user role to all new users for now, as create Account request does not have any role field
	if s.RoleService == nil {
		return merrors.InternalServerError(s.id, "could not assign role to account: roleService not configured")
	}
	if _, err = s.RoleService.AssignRoleToUser(ctx, &settings.AssignRoleToUserRequest{
		AccountUuid: out.Id,
		RoleId:      settings_svc.BundleUUIDRoleUser,
	}); err != nil {
		return merrors.InternalServerError(s.id, "could not assign role to account: %v", err.Error())
	}

	return
}

// rollbackCreateAccount tries to rollback changes made by `CreateAccount` if parts of it failed.
func (s Service) rollbackCreateAccount(ctx context.Context, acc *proto.Account) {
	err := s.index.Delete(acc)
	if err != nil {
		s.log.Err(err).Msg("failed to rollback account from indices")
	}
	err = s.repo.DeleteAccount(ctx, acc.Id)
	if err != nil {
		s.log.Err(err).Msg("failed to rollback account from repo")
	}
}

// UpdateAccount implements the AccountsServiceHandler interface
// read only fields are ignored
// TODO how can we unset specific values? using the update mask
func (s Service) UpdateAccount(ctx context.Context, in *proto.UpdateAccountRequest, out *proto.Account) (err error) {
	hasSelf := s.hasSelfManagementPermissions(ctx)
	hasManagement := s.hasAccountManagementPermissions(ctx)
	if !hasSelf && !hasManagement {
		return merrors.Forbidden(s.id, "no permission for UpdateAccount")
	}
	onlySelf := hasSelf && !hasManagement

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

	if onlySelf {
		// limit update to own account id
		if aid, ok := metadata.Get(ctx, middleware.AccountID); ok {
			if id != aid {
				return merrors.Forbidden(s.id, "no permission to UpdateAccount of another user")
			}
		} else {
			return merrors.InternalServerError(s.id, "account id not in context")
		}
	}

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

	var validMask fieldmask_utils.FieldFilterContainer
	if onlySelf {
		if validMask, err = validateUpdate(in.UpdateMask, selfUpdatableAccountPaths); err != nil {
			return merrors.BadRequest(s.id, "%s", err)
		}
	} else {
		if validMask, err = validateUpdate(in.UpdateMask, updatableAccountPaths); err != nil {
			return merrors.BadRequest(s.id, "%s", err)
		}
	}

	if _, exists := validMask.Filter("PreferredName"); exists {
		if err = validateAccountPreferredName(s.id, in.Account); err != nil {
			return err
		}
	}
	if _, exists := validMask.Filter("OnPremisesSamAccountName"); exists {
		if err = validateAccountOnPremisesSamAccountName(s.id, in.Account); err != nil {
			return err
		}
	}
	if _, exists := validMask.Filter("Mail"); exists {
		if in.Account.Mail != "" {
			if err = validateAccountEmail(s.id, in.Account); err != nil {
				return err
			}
		}
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
			hashed, err := bcrypt.GenerateFromPassword([]byte(in.Account.PasswordProfile.Password), s.Config.Server.HashDifficulty)
			if err != nil {
				in.Account.PasswordProfile.Password = ""
				s.log.Error().Err(err).Str("id", id).Msg("could not hash password")
				return merrors.InternalServerError(s.id, "could not hash password: %v", err.Error())
			}
			out.PasswordProfile.Password = string(hashed)
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

	// We need to reload the old account state to be able to compute the update
	old := &proto.Account{}
	if err = s.repo.LoadAccount(ctx, id, old); err != nil {
		s.log.Error().Err(err).Str("id", out.Id).Msg("could not load old account representation during update, maybe the account got deleted meanwhile?")
		return merrors.InternalServerError(s.id, "could not load current account for update: %v", err.Error())
	}

	if err = s.repo.WriteAccount(ctx, out); err != nil {
		s.log.Error().Err(err).Str("id", out.Id).Msg("could not persist updated account")
		return merrors.InternalServerError(s.id, "could not persist updated account: %v", err.Error())
	}

	if err = s.index.Update(old, out); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not index new account")
		return merrors.InternalServerError(s.id, "could not index updated account: %v", err.Error())
	}

	// remove password
	if out.PasswordProfile != nil {
		out.PasswordProfile.Password = ""
	}

	return
}

// whitelist of all paths/fields which can be updated by users themself
var selfUpdatableAccountPaths = map[string]struct{}{
	"DisplayName":              {},
	"Description":              {},
	"Mail":                     {}, // read only?,
	"PasswordProfile.Password": {},
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

	if err = s.index.Delete(a); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("accountId", id).Msg("could not remove account from index")
		return merrors.InternalServerError(s.id, "could not remove account from index: %v", err.Error())
	}

	s.log.Info().Str("id", id).Msg("deleted account")
	return
}

func validateAccount(serviceID string, a *proto.Account) error {
	if err := validateAccountPreferredName(serviceID, a); err != nil {
		return err
	}
	if err := validateAccountOnPremisesSamAccountName(serviceID, a); err != nil {
		return err
	}
	if err := validateAccountEmail(serviceID, a); err != nil {
		return err
	}
	return nil
}

func validateAccountPreferredName(serviceID string, a *proto.Account) error {
	if !isValidUsername(a.PreferredName) {
		return merrors.BadRequest(serviceID, "preferred_name '%s' must be at least the local part of an email", a.PreferredName)
	}
	return nil
}

func validateAccountOnPremisesSamAccountName(serviceID string, a *proto.Account) error {
	if !isValidUsername(a.OnPremisesSamAccountName) {
		return merrors.BadRequest(serviceID, "on_premises_sam_account_name '%s' must be at least the local part of an email", a.OnPremisesSamAccountName)
	}
	return nil
}

func validateAccountEmail(serviceID string, a *proto.Account) error {
	if !isValidEmail(a.Mail) {
		return merrors.BadRequest(serviceID, "mail '%s' must be a valid email", a.Mail)
	}
	return nil
}

// We want to allow email addresses as usernames so they show up when using them in ACLs on storages that allow intergration with our glauth LDAP service
// so we are adding a few restrictions from https://stackoverflow.com/questions/6949667/what-are-the-real-rules-for-linux-usernames-on-centos-6-and-rhel-6
// names should not start with numbers
var usernameRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]*(@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)*$")

func isValidUsername(e string) bool {
	if len(e) < 1 && len(e) > 254 {
		return false
	}
	return usernameRegex.MatchString(e)
}

// regex from https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#valid-e-mail-address
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValidEmail(e string) bool {
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

func (s Service) accountExists(ctx context.Context, username, mail, id string) (exists bool, err error) {
	var ids []string
	ids, err = s.index.FindBy(&proto.Account{}, "preferred_name", username)
	if err != nil {
		return false, err
	}
	if len(ids) > 0 {
		return true, nil
	}

	ids, err = s.index.FindBy(&proto.Account{}, "on_premises_sam_account_name", username)
	if err != nil {
		return false, err
	}
	if len(ids) > 0 {
		return true, nil
	}

	ids, err = s.index.FindBy(&proto.Account{}, "mail", mail)
	if err != nil {
		return false, err
	}
	if len(ids) > 0 {
		return true, nil
	}

	a := &proto.Account{}
	err = s.repo.LoadAccount(ctx, id, a)
	if err == nil {
		return true, nil
	}
	if !storage.IsNotFoundErr(err) {
		return true, err
	}
	return false, nil
}

func getAuthQueryMatch(query string) (match []string, authRequest bool) {
	match = authQuery.FindStringSubmatch(query)
	return match, len(match) == 3
}

func isPasswordValid(logger log.Logger, hash string, pwd string) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error().Err(fmt.Errorf("%s", r)).Str("hash", hash).Msg("password lib panicked")
		}
	}()

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}

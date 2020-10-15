package service

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
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
	if !s.hasAccountManagementPermissions(ctx) {
		return merrors.Forbidden(s.id, "no permission for ListAccounts")
	}

	accLock.Lock()
	defer accLock.Unlock()
	var password string
	var searchResults []string

	teardownServiceUser := s.serviceUserToIndex()
	defer teardownServiceUser()

	// check if this looks like an auth request
	match := authQuery.FindStringSubmatch(in.Query)
	if len(match) == 3 {
		in.Query = fmt.Sprintf("on_premises_sam_account_name eq '%s'", match[1]) // todo fetch email? make query configurable
		password = match[2]
		if password == "" {
			return merrors.Unauthorized(s.id, "password must not be empty")
		}

		searchResults, err = s.index.FindBy(&proto.Account{}, "OnPremisesSamAccountName", match[1])
		if err != nil {
			return err
		}
	}

	var onPremQuery = regexp.MustCompile(`^on_premises_sam_account_name eq '(.*)'$`) // TODO how is ' escaped in the password?
	match = onPremQuery.FindStringSubmatch(in.Query)
	if len(match) == 2 {
		searchResults, err = s.index.FindBy(&proto.Account{}, "OnPremisesSamAccountName", match[1])
	}

	var mailQuery = regexp.MustCompile(`^mail eq '(.*)'$`)
	match = mailQuery.FindStringSubmatch(in.Query)
	if len(match) == 2 {
		searchResults, err = s.index.FindBy(&proto.Account{}, "Mail", match[1])
	}

	// startswith(on_premises_sam_account_name,'mar') or startswith(display_name,'mar') or startswith(mail,'mar')
	var searchQuery = regexp.MustCompile(`^startswith\(on_premises_sam_account_name,'(.*)'\) or startswith\(display_name,'(.*)'\) or startswith\(mail,'(.*)'\)$`)

	match = searchQuery.FindStringSubmatch(in.Query)
	if len(match) == 4 {
		resSam, _ := s.index.FindByPartial(&proto.Account{}, "OnPremisesSamAccountName", match[1]+"*")
		resDisp, _ := s.index.FindByPartial(&proto.Account{}, "DisplayName", match[2]+"*")
		resMail, _ := s.index.FindByPartial(&proto.Account{}, "Mail", match[3]+"*")

		searchResults = append(resSam, append(resDisp, resMail...)...)
		searchResults = unique(searchResults)

	}

	// id eq 'marie' or on_premises_sam_account_name eq 'marie'
	// id eq 'f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c' or on_premises_sam_account_name eq 'f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c'
	var idOrQuery = regexp.MustCompile(`^id eq '(.*)' or on_premises_sam_account_name eq '(.*)'$`)
	match = idOrQuery.FindStringSubmatch(in.Query)
	if len(match) == 3 {
		qID, qSam := match[1], match[2]
		tmp := &proto.Account{}
		_ = s.repo.LoadAccount(ctx, qID, tmp)
		searchResults, err = s.index.FindBy(&proto.Account{}, "OnPremisesSamAccountName", qSam)

		if tmp.Id != "" {
			searchResults = append(searchResults, tmp.Id)
		}

		searchResults = unique(searchResults)

	}

	if in.Query == "" {
		searchResults, _ = s.index.FindByPartial(&proto.Account{}, "Mail", "*")
	}

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
	if err = validateAccount(s.id, *acc); err != nil {
		return err
	}

	if id, err = cleanupID(acc.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	exists, err := s.accountExists(ctx, acc.PreferredName, acc.Mail, acc.Id)
	if err != nil {
		return merrors.InternalServerError(s.id, "could not check if account exists: %v", err.Error())
	}
	if exists {
		return merrors.BadRequest(s.id, "account already exists")
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
	indexResults, err := s.index.Add(acc)
	if err != nil {
		s.rollbackCreateAccount(ctx, acc)
		return merrors.BadRequest(s.id, "Account already exists %v", err.Error())

	}
	s.log.Debug().Interface("account", acc).Msg("account after indexing")

	for _, r := range indexResults {
		if r.Field == "UidNumber" {
			id, err := strconv.Atoi(path.Base(r.Value))
			if err != nil {
				s.rollbackCreateAccount(ctx, acc)
				return err
			}
			acc.UidNumber = int64(id)
			break
		}
	}

	if in.Account.GidNumber == 0 {
		in.Account.GidNumber = userDefaultGID
	}

	r := proto.ListGroupsResponse{}
	err = s.ListGroups(ctx, &proto.ListGroupsRequest{}, &r)
	if err != nil {
		// rollback account creation
		return err
	}

	for _, group := range r.Groups {
		if group.GidNumber == in.Account.GidNumber {
			in.Account.MemberOf = append(in.Account.MemberOf, group)
		}
	}
	//acc.MemberOf = append(acc.MemberOf, &group)
	if err := s.repo.WriteAccount(context.Background(), acc); err != nil {
		return err
	}

	if acc.PasswordProfile != nil {
		acc.PasswordProfile.Password = ""
	}

	*out = *acc

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

	if err = validateAccount(s.id, *in.Account); err != nil {
		return err
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

	if err = s.index.Delete(a); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("accountId", id).Msg("could not remove account from index")
		return merrors.InternalServerError(s.id, "could not remove account from index: %v", err.Error())
	}

	s.log.Info().Str("id", id).Msg("deleted account")
	return
}

func validateAccount(serviceID string, a proto.Account) error {
	if !isValidUsername(a.PreferredName) {
		return merrors.BadRequest(serviceID, "preferred_name '%s' must be at least the local part of an email", a.PreferredName)
	}
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

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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

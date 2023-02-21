package identity

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/CiscoM31/godata"
	"github.com/go-ldap/ldap/v3"
	"github.com/gofrs/uuid"
	"github.com/libregraph/idm/pkg/ldapdn"
	libregraph "github.com/owncloud/libre-graph-api-go"
	oldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	"golang.org/x/exp/slices"
)

const (
	_givenNameAttribute        = "givenname"
	_surNameAttribute          = "sn"
	_ldapGroupOfNamesAttribute = "(objectClass=groupOfNames)"
	_ldapGroupMemberAttribute  = "member"
)

// DisableUserMechanismType is used instead of directly using the string values from the configuration.
type DisableUserMechanismType int64

// The different DisableMechanism* constants are used for managing the enabling/disabling of users.
const (
	DisableMechanismNone DisableUserMechanismType = iota
	DisableMechanismAttribute
	DisableMechanismGroup
)

var mechanismMap = map[string]DisableUserMechanismType{
	"":          DisableMechanismNone,
	"none":      DisableMechanismNone,
	"attribute": DisableMechanismAttribute,
	"group":     DisableMechanismGroup,
}

type LDAP struct {
	useServerUUID   bool
	writeEnabled    bool
	refintEnabled   bool
	usePwModifyExOp bool

	userBaseDN       string
	userFilter       string
	userObjectClass  string
	userScope        int
	userAttributeMap userAttributeMap

	disableUserMechanism    DisableUserMechanismType
	localUserDisableGroupDN string

	groupBaseDN       string
	groupFilter       string
	groupObjectClass  string
	groupScope        int
	groupAttributeMap groupAttributeMap

	educationConfig educationConfig

	logger *log.Logger
	conn   ldap.Client
}

type userAttributeMap struct {
	displayName    string
	id             string
	mail           string
	userName       string
	givenName      string
	surname        string
	accountEnabled string
}

type ldapAttributeValues map[string][]string

// ParseDisableMechanismType checks that the configuration option for how to disable users is correct.
func ParseDisableMechanismType(disableMechanism string) (DisableUserMechanismType, error) {
	disableMechanism = strings.ToLower(disableMechanism)
	t, ok := mechanismMap[disableMechanism]
	if !ok {
		return -1, errors.New("invalid configuration option for disable user mechanism")
	}

	return t, nil
}

func NewLDAPBackend(lc ldap.Client, config config.LDAP, logger *log.Logger) (*LDAP, error) {
	if config.UserDisplayNameAttribute == "" || config.UserIDAttribute == "" ||
		config.UserEmailAttribute == "" || config.UserNameAttribute == "" {
		return nil, errors.New("invalid user attribute mappings")
	}
	uam := userAttributeMap{
		displayName:    config.UserDisplayNameAttribute,
		id:             config.UserIDAttribute,
		mail:           config.UserEmailAttribute,
		userName:       config.UserNameAttribute,
		accountEnabled: config.UserEnabledAttribute,
		givenName:      _givenNameAttribute,
		surname:        _surNameAttribute,
	}

	if config.GroupNameAttribute == "" || config.GroupIDAttribute == "" {
		return nil, errors.New("invalid group attribute mappings")
	}
	gam := groupAttributeMap{
		name:         config.GroupNameAttribute,
		id:           config.GroupIDAttribute,
		member:       _ldapGroupMemberAttribute,
		memberSyntax: "dn",
	}

	var userScope, groupScope int
	var err error
	if userScope, err = stringToScope(config.UserSearchScope); err != nil {
		return nil, fmt.Errorf("error configuring user scope: %w", err)
	}

	if groupScope, err = stringToScope(config.GroupSearchScope); err != nil {
		return nil, fmt.Errorf("error configuring group scope: %w", err)
	}

	var educationConfig educationConfig
	if educationConfig, err = newEducationConfig(config); err != nil {
		return nil, fmt.Errorf("error setting up education resource config: %w", err)
	}

	disableMechanismType, err := ParseDisableMechanismType(config.DisableUserMechanism)
	if err != nil {
		return nil, fmt.Errorf("error configuring disable user mechanism: %w", err)
	}

	return &LDAP{
		useServerUUID:           config.UseServerUUID,
		usePwModifyExOp:         config.UsePasswordModExOp,
		userBaseDN:              config.UserBaseDN,
		userFilter:              config.UserFilter,
		userObjectClass:         config.UserObjectClass,
		userScope:               userScope,
		userAttributeMap:        uam,
		groupBaseDN:             config.GroupBaseDN,
		groupFilter:             config.GroupFilter,
		groupObjectClass:        config.GroupObjectClass,
		groupScope:              groupScope,
		groupAttributeMap:       gam,
		educationConfig:         educationConfig,
		disableUserMechanism:    disableMechanismType,
		localUserDisableGroupDN: config.LdapDisabledUsersGroupDN,
		logger:                  logger,
		conn:                    lc,
		writeEnabled:            config.WriteEnabled,
		refintEnabled:           config.RefintEnabled,
	}, nil
}

// CreateUser implements the Backend Interface. It converts the libregraph.User into an
// LDAP User Entry (using the inetOrgPerson LDAP Objectclass) add adds that to the
// configured LDAP server
func (i *LDAP) CreateUser(ctx context.Context, user libregraph.User) (*libregraph.User, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("CreateUser")
	if !i.writeEnabled {
		return nil, ErrReadOnly
	}

	if user.AccountEnabled != nil && i.disableUserMechanism == DisableMechanismNone {
		return nil, errors.New("accountEnabled option not compatible with backend disable user mechanism")
	}

	ar, err := i.userToAddRequest(user)
	if err != nil {
		return nil, err
	}

	if err := i.conn.Add(ar); err != nil {
		var lerr *ldap.Error
		logger.Debug().Err(err).Msg("error adding user")
		if errors.As(err, &lerr) {
			if lerr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				err = errorcode.New(errorcode.NameAlreadyExists, lerr.Error())
			}
		}
		return nil, err
	}

	if i.usePwModifyExOp && user.PasswordProfile != nil && user.PasswordProfile.Password != nil {
		if err := i.updateUserPassowrd(ctx, ar.DN, user.PasswordProfile.GetPassword()); err != nil {
			return nil, err
		}
	}

	// Read	back user from LDAP to get the generated UUID
	e, err := i.getUserByDN(ar.DN)
	if err != nil {
		return nil, err
	}
	return i.createUserModelFromLDAP(e), nil
}

// DeleteUser implements the Backend Interface. It permanently deletes a User identified
// by name or id from the LDAP server
func (i *LDAP) DeleteUser(ctx context.Context, nameOrID string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("DeleteUser")
	if !i.writeEnabled {
		return ErrReadOnly
	}
	e, err := i.getLDAPUserByNameOrID(nameOrID)
	if err != nil {
		return err
	}
	dr := ldap.DelRequest{DN: e.DN}
	if err = i.conn.Del(&dr); err != nil {
		return err
	}

	if !i.refintEnabled {
		// Find all the groups that this user was a member of and remove it from there
		groupEntries, err := i.getLDAPGroupsByFilter(fmt.Sprintf("(%s=%s)", i.groupAttributeMap.member, e.DN), true, false)
		if err != nil {
			return err
		}
		for _, group := range groupEntries {
			logger.Debug().Str("group", group.DN).Str("user", e.DN).Msg("Cleaning up group membership")

			if mr, err := i.removeEntryByDNAndAttributeFromEntry(group, e.DN, i.groupAttributeMap.member); err == nil {
				if err = i.conn.Modify(mr); err != nil {
					// Errors when deleting the memberships are only logged as warnings but not returned
					// to the user as we already successfully deleted the users itself
					logger.Warn().Str("group", group.DN).Str("user", e.DN).Err(err).Msg("failed to remove member")
				}
			}
		}
	}
	return nil
}

// UpdateUser implements the Backend Interface for the LDAP Backend
func (i *LDAP) UpdateUser(ctx context.Context, nameOrID string, user libregraph.User) (*libregraph.User, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("UpdateUser")
	if !i.writeEnabled {
		return nil, ErrReadOnly
	}
	e, err := i.getLDAPUserByNameOrID(nameOrID)
	if err != nil {
		return nil, err
	}

	var updateNeeded bool

	// Don't allow updates of the ID
	if user.Id != nil && *user.Id != "" {
		if e.GetEqualFoldAttributeValue(i.userAttributeMap.id) != *user.Id {
			return nil, errorcode.New(errorcode.NotAllowed, "changing the UserId is not allowed")
		}
	}
	if user.OnPremisesSamAccountName != nil && *user.OnPremisesSamAccountName != "" {
		if eu := e.GetEqualFoldAttributeValue(i.userAttributeMap.userName); eu != *user.OnPremisesSamAccountName {
			e, err = i.changeUserName(ctx, e.DN, eu, user.GetOnPremisesSamAccountName())
			if err != nil {
				return nil, err
			}
		}
	}

	mr := ldap.ModifyRequest{DN: e.DN}
	if user.DisplayName != nil && *user.DisplayName != "" {
		if e.GetEqualFoldAttributeValue(i.userAttributeMap.displayName) != *user.DisplayName {
			mr.Replace(i.userAttributeMap.displayName, []string{*user.DisplayName})
			updateNeeded = true
		}
	}
	if user.Mail != nil && *user.Mail != "" {
		if e.GetEqualFoldAttributeValue(i.userAttributeMap.mail) != *user.Mail {
			mr.Replace(i.userAttributeMap.mail, []string{*user.Mail})
			updateNeeded = true
		}
	}
	if user.PasswordProfile != nil && user.PasswordProfile.Password != nil && *user.PasswordProfile.Password != "" {
		if i.usePwModifyExOp {
			if err := i.updateUserPassowrd(ctx, e.DN, user.PasswordProfile.GetPassword()); err != nil {
				return nil, err
			}
		} else {
			// password are hashed server side there is no need to check if the new password
			// is actually different from the old one.
			mr.Replace("userPassword", []string{*user.PasswordProfile.Password})
			updateNeeded = true
		}
	}

	// Behavior of enabling/disabling of users depends on the "disableUserMechanism" config option:
	//
	// "attribute": For the upstream user management service which modifies accountEnabled on the user entry
	// "group": Makes it possible for local admins to disable users by adding them to a special group
	if user.AccountEnabled != nil {
		switch i.disableUserMechanism {
		case DisableMechanismNone:
			return nil, errors.New("accountEnabled option not compatible with backend disable user mechanism")
		case DisableMechanismAttribute:
			boolString := strings.ToUpper(strconv.FormatBool(user.GetAccountEnabled()))
			ldapValue := e.GetEqualFoldAttributeValue(i.userAttributeMap.accountEnabled)
			if ldapValue != "" {
				mr.Replace(i.userAttributeMap.accountEnabled, []string{boolString})
			} else {
				mr.Add(i.userAttributeMap.accountEnabled, []string{boolString})
			}
			updateNeeded = true
		case DisableMechanismGroup:
			if user.GetAccountEnabled() {
				err = i.enableUser(logger, e.DN)
			} else {
				err = i.disableUser(logger, e.DN)
			}

			if err != nil {
				return nil, err
			}
			updateNeeded = true
		}
	}

	if updateNeeded {
		if err := i.conn.Modify(&mr); err != nil {
			return nil, err
		}
	}

	// Read	back user from LDAP to get the generated UUID
	e, err = i.getUserByDN(e.DN)
	if err != nil {
		return nil, err
	}

	returnUser := i.createUserModelFromLDAP(e)

	// To avoid an ldap lookup for group membership, set the enabled flag to same as input value
	// since this would have been updated with group membership from the input anyway.
	if user.AccountEnabled != nil && i.disableUserMechanism != DisableMechanismNone {
		returnUser.AccountEnabled = user.AccountEnabled
	}

	return returnUser, nil
}

func (i *LDAP) getUserByDN(dn string) (*ldap.Entry, error) {
	attrs := []string{
		i.userAttributeMap.displayName,
		i.userAttributeMap.id,
		i.userAttributeMap.mail,
		i.userAttributeMap.userName,
		i.userAttributeMap.surname,
		i.userAttributeMap.givenName,
		i.userAttributeMap.accountEnabled,
	}

	filter := fmt.Sprintf("(objectClass=%s)", i.userObjectClass)

	if i.userFilter != "" {
		filter = fmt.Sprintf("(&%s(%s))", filter, i.userFilter)
	}

	return i.getEntryByDN(dn, attrs, filter)
}

func (i *LDAP) getEntryByDN(dn string, attrs []string, filter string) (*ldap.Entry, error) {
	if filter == "" {
		filter = "(objectclass=*)"
	}

	searchRequest := ldap.NewSearchRequest(
		dn, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 1, 0, false,
		filter,
		attrs,
		nil,
	)

	i.logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("getEntryByDN")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		i.logger.Error().Err(err).Str("backend", "ldap").Str("dn", dn).Msg("Search ldap by DN failed")
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}
	if len(res.Entries) == 0 {
		return nil, ErrNotFound
	}

	return res.Entries[0], nil
}

func (i *LDAP) searchLDAPEntryByFilter(basedn string, attrs []string, filter string) (*ldap.Entry, error) {
	if filter == "" {
		filter = "(objectclass=*)"
	}

	searchRequest := ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases, 1, 0, false,
		filter,
		attrs,
		nil,
	)

	i.logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("getEntryByFilter")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		i.logger.Error().Err(err).Str("backend", "ldap").Str("dn", basedn).Str("filter", filter).Msg("Search user by filter failed")
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}
	if len(res.Entries) == 0 {
		return nil, ErrNotFound
	}

	return res.Entries[0], nil
}

func (i *LDAP) getLDAPUserByID(id string) (*ldap.Entry, error) {
	id = ldap.EscapeFilter(id)
	filter := fmt.Sprintf("(%s=%s)", i.userAttributeMap.id, id)
	return i.getLDAPUserByFilter(filter)
}

func (i *LDAP) getLDAPUserByNameOrID(nameOrID string) (*ldap.Entry, error) {
	nameOrID = ldap.EscapeFilter(nameOrID)
	filter := fmt.Sprintf("(|(%s=%s)(%s=%s))", i.userAttributeMap.userName, nameOrID, i.userAttributeMap.id, nameOrID)
	return i.getLDAPUserByFilter(filter)
}

func (i *LDAP) getLDAPUserByFilter(filter string) (*ldap.Entry, error) {
	filter = fmt.Sprintf("(&%s(objectClass=%s)%s)", i.userFilter, i.userObjectClass, filter)
	attrs := []string{
		i.userAttributeMap.displayName,
		i.userAttributeMap.id,
		i.userAttributeMap.mail,
		i.userAttributeMap.userName,
		i.userAttributeMap.surname,
		i.userAttributeMap.givenName,
		i.userAttributeMap.accountEnabled,
	}
	return i.searchLDAPEntryByFilter(i.userBaseDN, attrs, filter)
}

// GetUser implements the Backend Interface.
func (i *LDAP) GetUser(ctx context.Context, nameOrID string, oreq *godata.GoDataRequest) (*libregraph.User, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetUser")

	e, err := i.getLDAPUserByNameOrID(nameOrID)
	if err != nil {
		return nil, err
	}

	u := i.createUserModelFromLDAP(e)
	if u == nil {
		return nil, ErrNotFound
	}

	if i.disableUserMechanism != DisableMechanismNone {
		userEnabled, err := i.UserEnabled(e)
		if err == nil {
			u.AccountEnabled = &userEnabled
		}
	}

	exp, err := GetExpandValues(oreq.Query)
	if err != nil {
		return nil, err
	}

	if slices.Contains(exp, "memberOf") {
		userGroups, err := i.getGroupsForUser(e.DN)
		if err != nil {
			return nil, err
		}
		u.MemberOf = i.groupsFromLDAPEntries(userGroups)
	}
	return u, nil
}

// GetUsers implements the Backend Interface.
func (i *LDAP) GetUsers(ctx context.Context, oreq *godata.GoDataRequest) ([]*libregraph.User, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetUsers")

	search, err := GetSearchValues(oreq.Query)
	if err != nil {
		return nil, err
	}

	exp, err := GetExpandValues(oreq.Query)
	if err != nil {
		return nil, err
	}

	var userFilter string
	if search != "" {
		search = ldap.EscapeFilter(search)
		userFilter = fmt.Sprintf(
			"(|(%s=%s*)(%s=%s*)(%s=%s*))",
			i.userAttributeMap.userName, search,
			i.userAttributeMap.mail, search,
			i.userAttributeMap.displayName, search,
		)
	}
	userFilter = fmt.Sprintf("(&%s(objectClass=%s)%s)", i.userFilter, i.userObjectClass, userFilter)
	searchRequest := ldap.NewSearchRequest(
		i.userBaseDN, i.userScope, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{
			i.userAttributeMap.displayName,
			i.userAttributeMap.id,
			i.userAttributeMap.mail,
			i.userAttributeMap.userName,
			i.userAttributeMap.surname,
			i.userAttributeMap.givenName,
			i.userAttributeMap.accountEnabled,
		},
		nil,
	)
	logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("GetUsers")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}

	users := make([]*libregraph.User, 0, len(res.Entries))
	usersEnabledState, err := i.usersEnabledState(res.Entries)
	if err != nil {
		return nil, err
	}

	for _, e := range res.Entries {
		u := i.createUserModelFromLDAP(e)
		// Skip invalid LDAP users
		if u == nil {
			continue
		}

		if i.disableUserMechanism != DisableMechanismNone {
			userEnabled := usersEnabledState[e.DN]
			u.AccountEnabled = &userEnabled
		}

		if slices.Contains(exp, "memberOf") {
			userGroups, err := i.getGroupsForUser(e.DN)
			if err != nil {
				return nil, err
			}
			u.MemberOf = i.groupsFromLDAPEntries(userGroups)
		}
		users = append(users, u)
	}
	return users, nil
}

func (i *LDAP) changeUserName(ctx context.Context, dn, originalUserName, newUserName string) (*ldap.Entry, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)

	newDN := fmt.Sprintf("%s=%s", i.userAttributeMap.userName, newUserName)

	logger.Debug().Str("originalDN", dn).Str("newDN", newDN).Msg("Modifying DN")
	mrdn := ldap.NewModifyDNRequest(dn, newDN, true, "")

	if err := i.conn.ModifyDN(mrdn); err != nil {
		var lerr *ldap.Error
		logger.Debug().Str("originalDN", dn).Str("newDN", newDN).Err(err).Msg("Failed to modify DN")
		if errors.As(err, &lerr) {
			if lerr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				err = errorcode.New(errorcode.NameAlreadyExists, lerr.Error())
			}
		}
		return nil, err
	}

	parsed, err := ldap.ParseDN(dn)
	if err != nil {
		return nil, err
	}

	newFullDN, err := replaceDN(parsed, newDN)
	if err != nil {
		return nil, err
	}

	u, err := i.getUserByDN(newFullDN)
	if err != nil {
		return nil, err
	}

	if !i.refintEnabled {
		groups, err := i.getGroupsForUser(dn)
		if err != nil {
			return nil, err
		}
		for _, g := range groups {
			logger.Debug().Str("originalDN", dn).Str("newDN", u.DN).Str("group", g.DN).Msg("Changing member in group")
			err = i.renameMemberInGroup(ctx, g, dn, u.DN)
			// This could leave the groups in an inconsistent state, might be a good idea to
			// add a defer that changes everything back on error. Ideally, this entire function
			// should be atomic, but LDAP doesn't support that.
			if err != nil {
				return nil, err
			}
		}
	}

	return u, nil
}

func (i *LDAP) renameMemberInGroup(ctx context.Context, group *ldap.Entry, oldMember, newMember string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("oldMember", oldMember).Str("newMember", newMember).Msg("replacing group member")
	mr := ldap.NewModifyRequest(group.DN, nil)
	mr.Delete(i.groupAttributeMap.member, []string{oldMember})
	mr.Add(i.groupAttributeMap.member, []string{newMember})
	if err := i.conn.Modify(mr); err != nil {
		var lerr *ldap.Error
		if errors.As(err, &lerr) {
			if lerr.ResultCode == ldap.LDAPResultNoSuchObject {
				logger.Warn().Str("group", group.DN).Msg("Group no longer exists")
				return nil
			} else if lerr.ResultCode == ldap.LDAPResultNoSuchAttribute {
				logger.Warn().
					Str("oldMember", oldMember).
					Str("newMember", newMember).
					Str("groupDN", group.DN).
					Msg("member attribute not found, this probably means that the server has refint enabled, please configure the OCIS to respect that.")
				return nil
			}
		}
		return err
	}
	return nil
}

func (i *LDAP) updateUserPassowrd(ctx context.Context, dn, password string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("updateUserPassowrd")
	pwMod := ldap.PasswordModifyRequest{
		UserIdentity: dn,
		NewPassword:  password,
	}
	// Note: We can ignore the result message here, as it were only relevant if we requested
	// the server to generate a new Password
	_, err := i.conn.PasswordModify(&pwMod)
	if err != nil {
		var lerr *ldap.Error
		logger.Debug().Err(err).Msg("error setting password for user")
		if errors.As(err, &lerr) {
			if lerr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				err = errorcode.New(errorcode.NameAlreadyExists, lerr.Error())
			}
		}
	}
	return err
}

func (i *LDAP) createUserModelFromLDAP(e *ldap.Entry) *libregraph.User {
	if e == nil {
		return nil
	}

	opsan := e.GetEqualFoldAttributeValue(i.userAttributeMap.userName)
	id := e.GetEqualFoldAttributeValue(i.userAttributeMap.id)
	givenName := e.GetEqualFoldAttributeValue(i.userAttributeMap.givenName)
	surname := e.GetEqualFoldAttributeValue(i.userAttributeMap.surname)

	if id != "" && opsan != "" {
		return &libregraph.User{
			DisplayName:              pointerOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.displayName)),
			Mail:                     pointerOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.mail)),
			OnPremisesSamAccountName: &opsan,
			Id:                       &id,
			GivenName:                &givenName,
			Surname:                  &surname,
			AccountEnabled:           booleanOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.accountEnabled)),
		}
	}
	i.logger.Warn().Str("dn", e.DN).Msg("Invalid User. Missing username or id attribute")
	return nil
}

func (i *LDAP) userToLDAPAttrValues(user libregraph.User) (map[string][]string, error) {
	attrs := map[string][]string{
		i.userAttributeMap.displayName: {user.GetDisplayName()},
		i.userAttributeMap.userName:    {user.GetOnPremisesSamAccountName()},
		i.userAttributeMap.mail:        {user.GetMail()},
		"objectClass":                  {"inetOrgPerson", "organizationalPerson", "person", "top", "ownCloudUser"},
		"cn":                           {user.GetOnPremisesSamAccountName()},
	}

	if !i.useServerUUID {
		attrs["owncloudUUID"] = []string{uuid.Must(uuid.NewV4()).String()}
	}

	if user.AccountEnabled != nil {
		attrs[i.userAttributeMap.accountEnabled] = []string{
			strings.ToUpper(strconv.FormatBool(*user.AccountEnabled)),
		}
	}

	// inetOrgPerson requires "sn" to be set. Set it to the Username if
	// Surname is not set in the Request
	var sn string
	if user.Surname != nil && *user.Surname != "" {
		sn = *user.Surname
	} else {
		sn = *user.OnPremisesSamAccountName
	}
	attrs[i.userAttributeMap.surname] = []string{sn}

	// When we get a givenName, we set the attribute.
	if givenName := user.GetGivenName(); givenName != "" {
		attrs[i.userAttributeMap.givenName] = []string{givenName}
	}

	if !i.usePwModifyExOp && user.PasswordProfile != nil && user.PasswordProfile.Password != nil {
		// Depending on the LDAP server implementation this might cause the
		// password to be stored in cleartext in the LDAP database. Using the
		// "Password Modify LDAP Extended Operation" is recommended.
		attrs["userPassword"] = []string{*user.PasswordProfile.Password}
	}
	return attrs, nil
}

func (i *LDAP) getUserAttrTypes() []string {
	return []string{
		i.userAttributeMap.displayName,
		i.userAttributeMap.userName,
		i.userAttributeMap.mail,
		i.userAttributeMap.surname,
		i.userAttributeMap.givenName,
		"objectClass",
		"cn",
		"owncloudUUID",
		"userPassword",
		i.userAttributeMap.accountEnabled,
	}
}

func (i *LDAP) getUserLDAPDN(user libregraph.User) string {
	return fmt.Sprintf("uid=%s,%s", oldap.EscapeDNAttributeValue(*user.OnPremisesSamAccountName), i.userBaseDN)
}

func (i *LDAP) userToAddRequest(user libregraph.User) (*ldap.AddRequest, error) {
	ar := ldap.NewAddRequest(i.getUserLDAPDN(user), nil)

	attrMap, err := i.userToLDAPAttrValues(user)
	if err != nil {
		return nil, err
	}

	for _, attrType := range i.getUserAttrTypes() {
		if values, ok := attrMap[attrType]; ok {
			ar.Attribute(attrType, values)
		}
	}
	return ar, nil
}

func pointerOrNil(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func booleanOrNil(val string) *bool {
	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		return nil
	}

	return &boolValue
}

func stringToScope(scope string) (int, error) {
	var s int
	switch scope {
	case "sub":
		s = ldap.ScopeWholeSubtree
	case "one":
		s = ldap.ScopeSingleLevel
	case "base":
		s = ldap.ScopeBaseObject
	default:
		return 0, fmt.Errorf("invalid Scope '%s'", scope)
	}
	return s, nil
}

// removeEntryByDNAndAttributeFromEntry creates a request to remove a single member entry by attribute and DN from an ldap entry
func (i *LDAP) removeEntryByDNAndAttributeFromEntry(entry *ldap.Entry, dn string, attribute string) (*ldap.ModifyRequest, error) {
	nOldDN, err := ldapdn.ParseNormalize(dn)
	if err != nil {
		return nil, err
	}
	entries := entry.GetEqualFoldAttributeValues(attribute)
	found := false
	for _, entry := range entries {
		if entry == "" {
			continue
		}
		if nEntry, err := ldapdn.ParseNormalize(entry); err != nil {
			// We couldn't parse the entry value as a DN. Let's keep it
			// as it is but log a warning
			i.logger.Warn().Str("entryDN", entry).Err(err).Msg("Couldn't parse DN")
			continue
		} else {
			if nEntry == nOldDN {
				found = true
			}
		}
	}
	if !found {
		i.logger.Debug().Str("backend", "ldap").Str("entry", entry.DN).Str("target", dn).
			Msg("The target is not an entry in the attribute list")
		return nil, ErrNotFound
	}

	mr := ldap.ModifyRequest{DN: entry.DN}
	if len(entries) == 1 {
		mr.Add(attribute, []string{""})
	}
	mr.Delete(attribute, []string{dn})
	return &mr, nil
}

// expandLDAPAttributeEntries reads an attribute from an ldap entry and expands to users
func (i *LDAP) expandLDAPAttributeEntries(ctx context.Context, e *ldap.Entry, attribute string) ([]*ldap.Entry, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("ExpandLDAPAttributeEntries")
	result := []*ldap.Entry{}

	for _, entryDN := range e.GetEqualFoldAttributeValues(attribute) {
		if entryDN == "" {
			continue
		}
		logger.Debug().Str("entryDN", entryDN).Msg("lookup")
		ue, err := i.getUserByDN(entryDN)
		if err != nil {
			// Ignore errors when reading a specific entry fails, just log them and continue
			logger.Debug().Err(err).Str("entry", entryDN).Msg("error reading attribute member entry")
			continue
		}
		result = append(result, ue)
	}

	return result, nil
}

func replaceDN(fullDN *ldap.DN, newDN string) (string, error) {
	if len(fullDN.RDNs) == 0 {
		return "", fmt.Errorf("Can't operate on an empty dn")
	}

	if len(fullDN.RDNs) == 1 {
		return newDN, nil
	}

	for _, part := range fullDN.RDNs[1:] {
		newDN += "," + part.String()
	}

	return newDN, nil
}

// CreateLDAPGroupByDN is a helper method specifically intended for creating a "system" group
// for managing locally disabled users on service startup
func (i *LDAP) CreateLDAPGroupByDN(dn string) error {
	ar := ldap.NewAddRequest(dn, nil)

	attrs := map[string][]string{
		"objectClass": {"groupOfNames", "top"},
		"member":      {""},
	}

	for attrType, values := range attrs {
		ar.Attribute(attrType, values)
	}

	return i.conn.Add(ar)
}

func (i *LDAP) disableUser(logger log.Logger, userDN string) (err error) {
	group, err := i.getEntryByDN(i.localUserDisableGroupDN, []string{_ldapGroupMemberAttribute}, _ldapGroupOfNamesAttribute)

	if err != nil {
		return err
	}

	mr := ldap.ModifyRequest{DN: group.DN}
	mr.Add(_ldapGroupMemberAttribute, []string{userDN})

	err = i.conn.Modify(&mr)
	var lerr *ldap.Error
	if errors.As(err, &lerr) {
		// If the user is already in the group, just log a message and return
		if lerr.ResultCode == ldap.LDAPResultAttributeOrValueExists {
			logger.Info().Msg("User already in group for disabled users")
			return nil
		}
	}

	return err
}

func (i *LDAP) enableUser(logger log.Logger, userDN string) (err error) {
	group, err := i.getEntryByDN(i.localUserDisableGroupDN, []string{_ldapGroupMemberAttribute}, _ldapGroupOfNamesAttribute)

	if err != nil {
		return err
	}

	mr := ldap.ModifyRequest{DN: group.DN}
	mr.Delete(_ldapGroupMemberAttribute, []string{userDN})

	err = i.conn.Modify(&mr)
	var lerr *ldap.Error
	if errors.As(err, &lerr) {
		// If the user is not in the group, just log a message and return
		if lerr.ResultCode == ldap.LDAPResultNoSuchAttribute {
			logger.Info().Msg("User was not in group for disabled users")
			return nil
		}
	}

	return err
}

func (i *LDAP) userEnabledByAttribute(user *ldap.Entry) bool {
	enabledAttribute := booleanOrNil(user.GetEqualFoldAttributeValue(i.userAttributeMap.accountEnabled))

	if enabledAttribute == nil {
		return true
	}

	return *enabledAttribute
}

func (i *LDAP) usersEnabledStateFromGroup(users []string) (usersEnabledState map[string]bool, err error) {
	group, err := i.getEntryByDN(i.localUserDisableGroupDN, []string{_ldapGroupMemberAttribute}, _ldapGroupOfNamesAttribute)

	if err != nil {
		return nil, err
	}

	usersEnabledState = make(map[string]bool, len(users))
	for _, user := range users {
		usersEnabledState[user] = true
	}

	for _, memberDN := range group.GetEqualFoldAttributeValues(_ldapGroupMemberAttribute) {
		usersEnabledState[memberDN] = false
	}

	return usersEnabledState, err
}

// UserEnabled returns if a user is enabled. This can depend on a flag in the user entry or group membership
func (i *LDAP) UserEnabled(user *ldap.Entry) (bool, error) {
	usersEnabledState, err := i.usersEnabledState([]*ldap.Entry{user})

	if err != nil {
		return false, err
	}

	return usersEnabledState[user.DN], nil
}

func (i *LDAP) usersEnabledState(users []*ldap.Entry) (usersEnabledState map[string]bool, err error) {
	usersEnabledState = make(map[string]bool, len(users))
	keys := make([]string, len(users))
	for index, user := range users {
		usersEnabledState[user.DN] = true
		keys[index] = user.DN
	}

	switch i.disableUserMechanism {
	case DisableMechanismAttribute:
		for _, user := range users {
			usersEnabledState[user.DN] = i.userEnabledByAttribute(user)
		}

	case DisableMechanismGroup:
		userDisabledGroupState, err := i.usersEnabledStateFromGroup(keys)

		if err != nil {
			return nil, err
		}

		for _, user := range keys {
			usersEnabledState[user] = userDisabledGroupState[user]
		}
	}

	return usersEnabledState, nil
}

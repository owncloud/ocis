package identity

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/gofrs/uuid"
	libregraph "github.com/owncloud/libre-graph-api-go"
	oldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	"golang.org/x/exp/slices"
)

const (
	givenNameAttribute = "givenname"
	surNameAttribute   = "sn"
)

type LDAP struct {
	useServerUUID   bool
	writeEnabled    bool
	usePwModifyExOp bool

	userBaseDN       string
	userFilter       string
	userObjectClass  string
	userScope        int
	userAttributeMap userAttributeMap

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
	displayName string
	id          string
	mail        string
	userName    string
	givenName   string
	surname     string
}

type ldapAttributeValues map[string][]string

func NewLDAPBackend(lc ldap.Client, config config.LDAP, logger *log.Logger) (*LDAP, error) {
	if config.UserDisplayNameAttribute == "" || config.UserIDAttribute == "" ||
		config.UserEmailAttribute == "" || config.UserNameAttribute == "" {
		return nil, errors.New("invalid user attribute mappings")
	}
	uam := userAttributeMap{
		displayName: config.UserDisplayNameAttribute,
		id:          config.UserIDAttribute,
		mail:        config.UserEmailAttribute,
		userName:    config.UserNameAttribute,
		givenName:   givenNameAttribute,
		surname:     surNameAttribute,
	}

	if config.GroupNameAttribute == "" || config.GroupIDAttribute == "" {
		return nil, errors.New("invalid group attribute mappings")
	}
	gam := groupAttributeMap{
		name:         config.GroupNameAttribute,
		id:           config.GroupIDAttribute,
		member:       "member",
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

	return &LDAP{
		useServerUUID:     config.UseServerUUID,
		usePwModifyExOp:   config.UsePasswordModExOp,
		userBaseDN:        config.UserBaseDN,
		userFilter:        config.UserFilter,
		userObjectClass:   config.UserObjectClass,
		userScope:         userScope,
		userAttributeMap:  uam,
		groupBaseDN:       config.GroupBaseDN,
		groupFilter:       config.GroupFilter,
		groupObjectClass:  config.GroupObjectClass,
		groupScope:        groupScope,
		groupAttributeMap: gam,
		educationConfig:   educationConfig,
		logger:            logger,
		conn:              lc,
		writeEnabled:      config.WriteEnabled,
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

	// Find all the groups that this user was a member of and remove it from there
	groupEntries, err := i.getLDAPGroupsByFilter(fmt.Sprintf("(%s=%s)", i.groupAttributeMap.member, e.DN), true, false)
	if err != nil {
		return err
	}
	for _, group := range groupEntries {
		logger.Debug().Str("group", group.DN).Str("user", e.DN).Msg("Cleaning up group membership")

		if mr, err := i.removeMemberFromGroupEntry(group, e.DN); err == nil {
			if err = i.conn.Modify(mr); err != nil {
				// Errors when deleting the memberships are only logged as warnings but not returned
				// to the user as we already successfully deleted the users itself
				logger.Warn().Str("group", group.DN).Str("user", e.DN).Err(err).Msg("failed to remove member")
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
	// TODO: In order to allow updating the user name we'd need to issue a ModRDN operation
	// As we currently using uid as the naming Attribute for the user entries. (Do we even
	// want to allow changing the user name?). For now just disallow it.
	if user.OnPremisesSamAccountName != nil && *user.OnPremisesSamAccountName != "" {
		if e.GetEqualFoldAttributeValue(i.userAttributeMap.userName) != *user.OnPremisesSamAccountName {
			return nil, errorcode.New(errorcode.NotSupported, "changing the user name is currently not supported")
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
	// TODO implement account disabled/enabled

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
	return i.createUserModelFromLDAP(e), nil
}

func (i *LDAP) getUserByDN(dn string) (*ldap.Entry, error) {
	attrs := []string{
		i.userAttributeMap.displayName,
		i.userAttributeMap.id,
		i.userAttributeMap.mail,
		i.userAttributeMap.userName,
		i.userAttributeMap.surname,
		i.userAttributeMap.givenName,
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
		i.logger.Error().Err(err).Str("backend", "ldap").Str("dn", dn).Msg("Search user by DN failed")
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
	}
	return i.searchLDAPEntryByFilter(i.userBaseDN, attrs, filter)
}

func (i *LDAP) GetUser(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.User, error) {
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
	sel := strings.Split(queryParam.Get("$select"), ",")
	exp := strings.Split(queryParam.Get("$expand"), ",")
	if slices.Contains(sel, "memberOf") || slices.Contains(exp, "memberOf") {
		userGroups, err := i.getGroupsForUser(e.DN)
		if err != nil {
			return nil, err
		}
		u.MemberOf = i.groupsFromLDAPEntries(userGroups)
	}
	return u, nil
}

func (i *LDAP) GetUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.User, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetUsers")

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
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

	for _, e := range res.Entries {
		sel := strings.Split(queryParam.Get("$select"), ",")
		exp := strings.Split(queryParam.Get("$expand"), ",")
		u := i.createUserModelFromLDAP(e)
		// Skip invalid LDAP users
		if u == nil {
			continue
		}
		if slices.Contains(sel, "memberOf") || slices.Contains(exp, "memberOf") {
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
		"objectClass":                  {"inetOrgPerson", "organizationalPerson", "person", "top"},
		"cn":                           {user.GetOnPremisesSamAccountName()},
	}

	if !i.useServerUUID {
		attrs["owncloudUUID"] = []string{uuid.Must(uuid.NewV4()).String()}
		attrs["objectClass"] = append(attrs["objectClass"], "owncloud")
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

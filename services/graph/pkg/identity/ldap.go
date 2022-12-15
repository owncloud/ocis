package identity

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/gofrs/uuid"
	ldapdn "github.com/libregraph/idm/pkg/ldapdn"
	libregraph "github.com/owncloud/libre-graph-api-go"
	oldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	"golang.org/x/exp/slices"
)

var (
	errReadOnly = errorcode.New(errorcode.NotAllowed, "server is configured read-only")
	errNotFound = errorcode.New(errorcode.ItemNotFound, "not found")
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
}

type groupAttributeMap struct {
	name         string
	id           string
	member       string
	memberSyntax string
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
		return nil, errReadOnly
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
		return errReadOnly
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

		if mr, err := i.removeMemberFromGroupEntry(group, e.DN); err == nil && mr != nil {
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
		return nil, errReadOnly
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
	}

	filter := fmt.Sprintf("(objectClass=%s)", i.userObjectClass)

	if i.userFilter != "" {
		filter = fmt.Sprintf("(&%s(%s))", filter, i.userFilter)
	}

	return i.getEntryByDN(dn, attrs, filter)
}

func (i *LDAP) getGroupByDN(dn string) (*ldap.Entry, error) {
	attrs := []string{
		i.groupAttributeMap.id,
		i.groupAttributeMap.name,
	}
	filter := fmt.Sprintf("(objectClass=%s)", i.groupObjectClass)

	if i.groupFilter != "" {
		filter = fmt.Sprintf("(&%s(%s))", filter, i.groupFilter)
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
		return nil, errNotFound
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
		return nil, errNotFound
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
		return nil, errNotFound
	}
	sel := strings.Split(queryParam.Get("$select"), ",")
	exp := strings.Split(queryParam.Get("$expand"), ",")
	if slices.Contains(sel, "memberOf") || slices.Contains(exp, "memberOf") {
		userGroups, err := i.getGroupsForUser(e.DN)
		if err != nil {
			return nil, err
		}
		if len(userGroups) > 0 {
			groups := i.groupsFromLDAPEntries(userGroups)
			u.MemberOf = groups
		}
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
			if len(userGroups) > 0 {
				expand := ldap.EscapeFilter(queryParam.Get("expand"))
				if expand == "" {
					expand = "false"
				}
				groups := i.groupsFromLDAPEntries(userGroups)
				u.MemberOf = groups
			}
		}
		users = append(users, u)
	}
	return users, nil
}

func (i *LDAP) getGroupsForUser(dn string) ([]*ldap.Entry, error) {
	groupFilter := fmt.Sprintf(
		"(%s=%s)",
		i.groupAttributeMap.member, dn,
	)
	userGroups, err := i.getLDAPGroupsByFilter(groupFilter, false, false)
	if err != nil {
		return nil, err
	}
	return userGroups, nil
}

func (i *LDAP) GetGroup(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.Group, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetGroup")
	e, err := i.getLDAPGroupByNameOrID(nameOrID, true)
	if err != nil {
		return nil, err
	}
	sel := strings.Split(queryParam.Get("$select"), ",")
	exp := strings.Split(queryParam.Get("$expand"), ",")
	var g *libregraph.Group
	if g = i.createGroupModelFromLDAP(e); g == nil {
		return nil, errorcode.New(errorcode.ItemNotFound, "not found")
	}
	if slices.Contains(sel, "members") || slices.Contains(exp, "members") {
		members, err := i.expandLDAPGroupMembers(ctx, e)
		if err != nil {
			return nil, err
		}
		if len(members) > 0 {
			m := make([]libregraph.User, 0, len(members))
			for _, ue := range members {
				if u := i.createUserModelFromLDAP(ue); u != nil {
					m = append(m, *u)
				}
			}
			g.Members = m
		}
	}
	return g, nil
}

func (i *LDAP) getLDAPGroupByID(id string, requestMembers bool) (*ldap.Entry, error) {
	id = ldap.EscapeFilter(id)
	filter := fmt.Sprintf("(%s=%s)", i.groupAttributeMap.id, id)
	return i.getLDAPGroupByFilter(filter, requestMembers)
}

func (i *LDAP) getLDAPGroupByNameOrID(nameOrID string, requestMembers bool) (*ldap.Entry, error) {
	nameOrID = ldap.EscapeFilter(nameOrID)
	filter := fmt.Sprintf("(|(%s=%s)(%s=%s))", i.groupAttributeMap.name, nameOrID, i.groupAttributeMap.id, nameOrID)
	return i.getLDAPGroupByFilter(filter, requestMembers)
}

func (i *LDAP) getLDAPGroupByFilter(filter string, requestMembers bool) (*ldap.Entry, error) {
	e, err := i.getLDAPGroupsByFilter(filter, requestMembers, true)
	if err != nil {
		return nil, err
	}
	if len(e) == 0 {
		return nil, errorcode.New(errorcode.ItemNotFound, "not found")
	}

	return e[0], nil
}

// Search for LDAP Groups matching the specified filter, if requestMembers is true the groupMemberShip
// attribute will be part of the result attributes. The LDAP filter is combined with the configured groupFilter
// resulting in a filter like "(&(LDAP.groupFilter)(objectClass=LDAP.groupObjectClass)(<filter_from_args>))"
func (i *LDAP) getLDAPGroupsByFilter(filter string, requestMembers, single bool) ([]*ldap.Entry, error) {
	attrs := []string{
		i.groupAttributeMap.name,
		i.groupAttributeMap.id,
	}

	if requestMembers {
		attrs = append(attrs, i.groupAttributeMap.member)
	}

	sizelimit := 0
	if single {
		sizelimit = 1
	}
	searchRequest := ldap.NewSearchRequest(
		i.groupBaseDN, i.groupScope, ldap.NeverDerefAliases, sizelimit, 0, false,
		fmt.Sprintf("(&%s(objectClass=%s)%s)", i.groupFilter, i.groupObjectClass, filter),
		attrs,
		nil,
	)
	i.logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("getLDAPGroupsByFilter")
	res, err := i.conn.Search(searchRequest)

	if err != nil {
		var errmsg string
		if lerr, ok := err.(*ldap.Error); ok {
			if lerr.ResultCode == ldap.LDAPResultSizeLimitExceeded {
				errmsg = fmt.Sprintf("too many results searching for group '%s'", filter)
				i.logger.Debug().Str("backend", "ldap").Err(lerr).Msg(errmsg)
			}
		}
		return nil, errorcode.New(errorcode.ItemNotFound, errmsg)
	}
	return res.Entries, nil
}

// removeMemberFromGroupEntry creates an LDAP Modify request (not sending it)
// that would update the supplied entry to remove the specified member from the
// group
func (i *LDAP) removeMemberFromGroupEntry(group *ldap.Entry, memberDN string) (*ldap.ModifyRequest, error) {
	nOldMemberDN, err := ldapdn.ParseNormalize(memberDN)
	if err != nil {
		return nil, err
	}
	members := group.GetEqualFoldAttributeValues(i.groupAttributeMap.member)
	found := false
	for _, member := range members {
		if member == "" {
			continue
		}
		if nMember, err := ldapdn.ParseNormalize(member); err != nil {
			// We couldn't parse the member value as a DN. Let's keep it
			// as it is but log a warning
			i.logger.Warn().Str("memberDN", member).Err(err).Msg("Couldn't parse DN")
			continue
		} else {
			if nMember == nOldMemberDN {
				found = true
			}
		}
	}
	if !found {
		i.logger.Debug().Str("backend", "ldap").Str("groupdn", group.DN).Str("member", memberDN).
			Msg("The target is not a member of the group")
		return nil, nil
	}

	mr := ldap.ModifyRequest{DN: group.DN}
	if len(members) == 1 {
		mr.Add(i.groupAttributeMap.member, []string{""})
	}
	mr.Delete(i.groupAttributeMap.member, []string{memberDN})
	return &mr, nil
}

func (i *LDAP) GetGroups(ctx context.Context, queryParam url.Values) ([]*libregraph.Group, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetGroups")

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
	}

	var expandMembers bool
	sel := strings.Split(queryParam.Get("$select"), ",")
	exp := strings.Split(queryParam.Get("$expand"), ",")
	if slices.Contains(sel, "members") || slices.Contains(exp, "members") {
		expandMembers = true
	}

	var groupFilter string
	if search != "" {
		search = ldap.EscapeFilter(search)
		groupFilter = fmt.Sprintf(
			"(|(%s=%s*)(%s=%s*))",
			i.groupAttributeMap.name, search,
			i.groupAttributeMap.id, search,
		)
	}
	groupFilter = fmt.Sprintf("(&%s(objectClass=%s)%s)", i.groupFilter, i.groupObjectClass, groupFilter)

	groupAttrs := []string{
		i.groupAttributeMap.name,
		i.groupAttributeMap.id,
	}
	if expandMembers {
		groupAttrs = append(groupAttrs, i.groupAttributeMap.member)
	}

	searchRequest := ldap.NewSearchRequest(
		i.groupBaseDN, i.groupScope, ldap.NeverDerefAliases, 0, 0, false,
		groupFilter,
		groupAttrs,
		nil,
	)
	logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("GetGroups")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}

	groups := make([]*libregraph.Group, 0, len(res.Entries))

	var g *libregraph.Group
	for _, e := range res.Entries {
		if g = i.createGroupModelFromLDAP(e); g == nil {
			continue
		}
		if expandMembers {
			members, err := i.expandLDAPGroupMembers(ctx, e)
			if err != nil {
				return nil, err
			}
			if len(members) > 0 {
				m := make([]libregraph.User, 0, len(members))
				for _, ue := range members {
					if u := i.createUserModelFromLDAP(ue); u != nil {
						m = append(m, *u)
					}
				}
				g.Members = m
			}
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// GetGroupMembers implements the Backend Interface for the LDAP Backend
func (i *LDAP) GetGroupMembers(ctx context.Context, groupID string) ([]*libregraph.User, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetGroupMembers")
	e, err := i.getLDAPGroupByNameOrID(groupID, true)
	if err != nil {
		return nil, err
	}

	memberEntries, err := i.expandLDAPGroupMembers(ctx, e)
	result := make([]*libregraph.User, 0, len(memberEntries))
	if err != nil {
		return nil, err
	}
	for _, member := range memberEntries {
		if u := i.createUserModelFromLDAP(member); u != nil {
			result = append(result, u)
		}
	}

	return result, nil
}

func (i *LDAP) expandLDAPGroupMembers(ctx context.Context, e *ldap.Entry) ([]*ldap.Entry, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("expandLDAPGroupMembers")
	result := []*ldap.Entry{}

	for _, memberDN := range e.GetEqualFoldAttributeValues(i.groupAttributeMap.member) {
		if memberDN == "" {
			continue
		}
		logger.Debug().Str("memberDN", memberDN).Msg("lookup")
		ue, err := i.getUserByDN(memberDN)
		if err != nil {
			// Ignore errors when reading a specific member fails, just log them and continue
			logger.Debug().Err(err).Str("member", memberDN).Msg("error reading group member")
			continue
		}
		result = append(result, ue)
	}

	return result, nil
}

// CreateGroup implements the Backend Interface for the LDAP Backend
// It is currently restricted to managing groups based on the "groupOfNames" ObjectClass.
// As "groupOfNames" requires a "member" Attribute to be present. Empty Groups (groups
// without a member) a represented by adding an empty DN as the single member.
func (i *LDAP) CreateGroup(ctx context.Context, group libregraph.Group) (*libregraph.Group, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("create group")
	if !i.writeEnabled {
		return nil, errorcode.New(errorcode.NotAllowed, "server is configured read-only")
	}
	ar := ldap.AddRequest{
		DN: fmt.Sprintf("cn=%s,%s", oldap.EscapeDNAttributeValue(*group.DisplayName), i.groupBaseDN),
		Attributes: []ldap.Attribute{
			{
				Type: i.groupAttributeMap.name,
				Vals: []string{*group.DisplayName},
			},
			// This is a crutch to allow groups without members for LDAP Server's which
			// that apply strict Schema checking. The RFCs define "member/uniqueMember"
			// as required attribute for groupOfNames/groupOfUniqueNames. So we
			// add an empty string (which is a valid DN) as the initial member.
			// It will be replace once real members are added.
			// We might wanna use the newer, but not so broadly used "groupOfMembers"
			// objectclass (RFC2307bis-02) where "member" is optional.
			{
				Type: i.groupAttributeMap.member,
				Vals: []string{""},
			},
		},
	}

	// TODO make group objectclass configurable to support e.g. posixGroup, groupOfUniqueNames, groupOfMembers?}
	objectClasses := []string{"groupOfNames", "top"}

	if !i.useServerUUID {
		ar.Attribute("owncloudUUID", []string{uuid.Must(uuid.NewV4()).String()})
		objectClasses = append(objectClasses, "owncloud")
	}
	ar.Attribute("objectClass", objectClasses)

	if err := i.conn.Add(&ar); err != nil {
		return nil, err
	}

	// Read	back group from LDAP to get the generated UUID
	e, err := i.getGroupByDN(ar.DN)
	if err != nil {
		return nil, err
	}
	return i.createGroupModelFromLDAP(e), nil
}

// DeleteGroup implements the Backend Interface.
func (i *LDAP) DeleteGroup(ctx context.Context, id string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("DeleteGroup")
	if !i.writeEnabled {
		return errorcode.New(errorcode.NotAllowed, "server is configured read-only")
	}
	e, err := i.getLDAPGroupByID(id, false)
	if err != nil {
		return err
	}
	dr := ldap.DelRequest{DN: e.DN}
	if err = i.conn.Del(&dr); err != nil {
		return err
	}
	return nil
}

// AddMembersToGroup implements the Backend Interface for the LDAP backend.
// Currently it is limited to adding Users as Group members. Adding other groups
// as members is not yet implemented
func (i *LDAP) AddMembersToGroup(ctx context.Context, groupID string, memberIDs []string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("AddMembersToGroup")
	ge, err := i.getLDAPGroupByID(groupID, true)
	if err != nil {
		return err
	}

	mr := ldap.ModifyRequest{DN: ge.DN}
	// Handle empty groups (using the empty member attribute)
	current := ge.GetEqualFoldAttributeValues(i.groupAttributeMap.member)
	if len(current) == 1 && current[0] == "" {
		mr.Delete(i.groupAttributeMap.member, []string{""})
	}

	// Create a Set of current members for faster lookups
	currentSet := make(map[string]struct{}, len(current))
	for _, currentMember := range current {
		// We can ignore any empty member value here
		if currentMember == "" {
			continue
		}
		nCurrentMember, err := ldapdn.ParseNormalize(currentMember)
		if err != nil {
			// We couldn't parse the member value as a DN. Let's skip it, but log a warning
			logger.Warn().Str("memberDN", currentMember).Err(err).Msg("Couldn't parse DN")
			continue
		}
		currentSet[nCurrentMember] = struct{}{}
	}

	var newMemberDNs []string
	for _, memberID := range memberIDs {
		me, err := i.getLDAPUserByID(memberID)
		if err != nil {
			return err
		}
		nDN, err := ldapdn.ParseNormalize(me.DN)
		if err != nil {
			logger.Error().Str("new member", me.DN).Err(err).Msg("Couldn't parse DN")
			return err
		}
		if _, present := currentSet[nDN]; !present {
			newMemberDNs = append(newMemberDNs, me.DN)
		} else {
			logger.Debug().Str("memberDN", me.DN).Msg("Member already present in group. Skipping")
		}
	}

	if len(newMemberDNs) > 0 {
		mr.Add(i.groupAttributeMap.member, newMemberDNs)

		if err := i.conn.Modify(&mr); err != nil {
			return err
		}
	}
	return nil
}

// RemoveMemberFromGroup implements the Backend Interface.
func (i *LDAP) RemoveMemberFromGroup(ctx context.Context, groupID string, memberID string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("RemoveMemberFromGroup")
	ge, err := i.getLDAPGroupByID(groupID, true)
	if err != nil {
		logger.Debug().Str("backend", "ldap").Str("groupID", groupID).Msg("Error looking up group")
		return err
	}
	me, err := i.getLDAPUserByID(memberID)
	if err != nil {
		logger.Debug().Str("backend", "ldap").Str("memberID", memberID).Msg("Error looking up group member")
		return err
	}
	logger.Debug().Str("backend", "ldap").Str("groupdn", ge.DN).Str("member", me.DN).Msg("remove member")

	if mr, err := i.removeMemberFromGroupEntry(ge, me.DN); err == nil && mr != nil {
		return i.conn.Modify(mr)
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

	if id != "" && opsan != "" {
		return &libregraph.User{
			DisplayName:              pointerOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.displayName)),
			Mail:                     pointerOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.mail)),
			OnPremisesSamAccountName: &opsan,
			Id:                       &id,
		}
	}
	i.logger.Warn().Str("dn", e.DN).Msg("Invalid User. Missing username or id attribute")
	return nil
}

func (i *LDAP) createGroupModelFromLDAP(e *ldap.Entry) *libregraph.Group {
	name := e.GetEqualFoldAttributeValue(i.groupAttributeMap.name)
	id := e.GetEqualFoldAttributeValue(i.groupAttributeMap.id)

	if id != "" && name != "" {
		return &libregraph.Group{
			DisplayName: &name,
			Id:          &id,
		}
	}
	i.logger.Warn().Str("dn", e.DN).Msg("Group is missing name or id")
	return nil
}

func (i *LDAP) groupsFromLDAPEntries(e []*ldap.Entry) []libregraph.Group {
	groups := make([]libregraph.Group, 0, len(e))
	for _, g := range e {
		if grp := i.createGroupModelFromLDAP(g); grp != nil {
			groups = append(groups, *grp)
		}
	}
	return groups
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
	attrs["sn"] = []string{sn}

	if !i.usePwModifyExOp && user.PasswordProfile != nil && user.PasswordProfile.Password != nil {
		// Depending on the LDAP server implementation this might cause the
		// password to be stored in cleartext in the LDAP database. Using the
		// "Password Modify LDAP Extended Operation" is recommended.
		attrs["userPassword"] = []string{*user.PasswordProfile.Password}
	}
	return attrs, nil
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
	for attrType, values := range attrMap {
		ar.Attribute(attrType, values)
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

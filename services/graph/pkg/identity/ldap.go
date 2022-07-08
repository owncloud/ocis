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
		logger:            logger,
		conn:              lc,
		writeEnabled:      config.WriteEnabled,
	}, nil
}

// CreateUser implements the Backend Interface. It converts the libregraph.User into an
// LDAP User Entry (using the inetOrgPerson LDAP Objectclass) add adds that to the
// configured LDAP server
func (i *LDAP) CreateUser(ctx context.Context, user libregraph.User) (*libregraph.User, error) {
	if !i.writeEnabled {
		return nil, errReadOnly
	}
	ar := ldap.AddRequest{
		DN: fmt.Sprintf("uid=%s,%s", oldap.EscapeDNAttributeValue(*user.OnPremisesSamAccountName), i.userBaseDN),
		Attributes: []ldap.Attribute{
			// inetOrgPerson requires "cn"
			{
				Type: "cn",
				Vals: []string{*user.OnPremisesSamAccountName},
			},
			{
				Type: i.userAttributeMap.mail,
				Vals: []string{*user.Mail},
			},
			{
				Type: i.userAttributeMap.userName,
				Vals: []string{*user.OnPremisesSamAccountName},
			},
			{
				Type: i.userAttributeMap.displayName,
				Vals: []string{*user.DisplayName},
			},
		},
	}

	objectClasses := []string{"inetOrgPerson", "organizationalPerson", "person", "top"}

	if !i.usePwModifyExOp && user.PasswordProfile != nil && user.PasswordProfile.Password != nil {
		// Depending on the LDAP server implementation this might cause the
		// password to be stored in cleartext in the LDAP database. Using the
		// "Password Modify LDAP Extended Operation" is recommended.
		ar.Attribute("userPassword", []string{*user.PasswordProfile.Password})
	}
	if !i.useServerUUID {
		ar.Attribute("owncloudUUID", []string{uuid.Must(uuid.NewV4()).String()})
		objectClasses = append(objectClasses, "owncloud")
	}
	ar.Attribute("objectClass", objectClasses)

	// inetOrgPerson requires "sn" to be set. Set it to the Username if
	// Surname is not set in the Request
	var sn string
	if user.Surname != nil && *user.Surname != "" {
		sn = *user.Surname
	} else {
		sn = *user.OnPremisesSamAccountName
	}
	ar.Attribute("sn", []string{sn})

	if err := i.conn.Add(&ar); err != nil {
		var lerr *ldap.Error
		i.logger.Debug().Err(err).Msg("error adding user")
		if errors.As(err, &lerr) {
			if lerr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				err = errorcode.New(errorcode.NameAlreadyExists, lerr.Error())
			}
		}
		return nil, err
	}

	if i.usePwModifyExOp && user.PasswordProfile != nil && user.PasswordProfile.Password != nil {
		if err := i.updateUserPassowrd(ar.DN, user.PasswordProfile.GetPassword()); err != nil {
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
		i.logger.Debug().Str("group", group.DN).Str("user", e.DN).Msg("Cleaning up group membership")

		if mr, err := i.removeMemberFromGroupEntry(group, e.DN); err == nil && mr != nil {
			if err = i.conn.Modify(mr); err != nil {
				// Errors when deleting the memberships are only logged as warnings but not returned
				// to the user as we already successfully deleted the users itself
				i.logger.Warn().Str("group", group.DN).Str("user", e.DN).Err(err).Msg("failed to remove member")
			}
		}
	}
	return nil
}

// UpdateUser implements the Backend Interface for the LDAP Backend
func (i *LDAP) UpdateUser(ctx context.Context, nameOrID string, user libregraph.User) (*libregraph.User, error) {
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
			if err := i.updateUserPassowrd(e.DN, user.PasswordProfile.GetPassword()); err != nil {
				return nil, err
			}
		} else {
			// password are hashed server side there is no need to check if the new password
			// is actually different from the old one.
			mr.Replace("userPassword", []string{*user.PasswordProfile.Password})
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
	return i.createUserModelFromLDAP(e), nil
}

func (i *LDAP) getUserByDN(dn string) (*ldap.Entry, error) {
	attrs := []string{
		i.userAttributeMap.displayName,
		i.userAttributeMap.id,
		i.userAttributeMap.mail,
		i.userAttributeMap.userName,
	}
	return i.getEntryByDN(dn, attrs)
}

func (i *LDAP) getGroupByDN(dn string) (*ldap.Entry, error) {
	attrs := []string{
		i.groupAttributeMap.id,
		i.groupAttributeMap.name,
	}
	return i.getEntryByDN(dn, attrs)
}

func (i *LDAP) getEntryByDN(dn string, attrs []string) (*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		dn, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 1, 0, false,
		"(objectclass=*)",
		attrs,
		nil,
	)

	i.logger.Debug().Str("backend", "ldap").Str("dn", dn).Msg("Search user by DN")
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
	searchRequest := ldap.NewSearchRequest(
		i.userBaseDN, i.userScope, ldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(&%s(objectClass=%s)%s)", i.userFilter, i.userObjectClass, filter),
		[]string{
			i.userAttributeMap.displayName,
			i.userAttributeMap.id,
			i.userAttributeMap.mail,
			i.userAttributeMap.userName,
		},
		nil,
	)
	i.logger.Debug().Str("backend", "ldap").Msgf("Search %s", i.userBaseDN)
	res, err := i.conn.Search(searchRequest)

	if err != nil {
		var errmsg string
		if lerr, ok := err.(*ldap.Error); ok {
			if lerr.ResultCode == ldap.LDAPResultSizeLimitExceeded {
				errmsg = fmt.Sprintf("too many results searching for user '%s'", filter)
				i.logger.Debug().Str("backend", "ldap").Err(lerr).
					Str("userfilter", filter).Msg("too many results searching for user")
			}
		}
		return nil, errorcode.New(errorcode.ItemNotFound, errmsg)
	}
	if len(res.Entries) == 0 {
		return nil, errNotFound
	}

	return res.Entries[0], nil
}

func (i *LDAP) GetUser(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.User, error) {
	i.logger.Debug().Str("backend", "ldap").Msg("GetUser")
	e, err := i.getLDAPUserByNameOrID(nameOrID)
	if err != nil {
		return nil, err
	}
	sel := strings.Split(queryParam.Get("$select"), ",")
	exp := strings.Split(queryParam.Get("$expand"), ",")
	u := i.createUserModelFromLDAP(e)
	if slices.Contains(sel, "memberOf") || slices.Contains(exp, "memberOf") {
		userGroups, err := i.getGroupsForUser(e.DN)
		if err != nil {
			return nil, err
		}
		if len(userGroups) > 0 {
			groups := make([]libregraph.Group, 0, len(userGroups))
			for _, g := range userGroups {
				groups = append(groups, *i.createGroupModelFromLDAP(g))
			}
			u.MemberOf = groups
		}
	}
	return u, nil
}

func (i *LDAP) GetUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.User, error) {
	i.logger.Debug().Str("backend", "ldap").Msg("GetUsers")

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
	i.logger.Debug().Str("backend", "ldap").Msgf("Search %s", i.userBaseDN)
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}

	users := make([]*libregraph.User, 0, len(res.Entries))

	for _, e := range res.Entries {
		sel := strings.Split(queryParam.Get("$select"), ",")
		exp := strings.Split(queryParam.Get("$expand"), ",")
		u := i.createUserModelFromLDAP(e)
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
				groups := make([]libregraph.Group, 0, len(userGroups))
				for _, g := range userGroups {
					groups = append(groups, *i.createGroupModelFromLDAP(g))
				}
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
	i.logger.Debug().Str("backend", "ldap").Msg("GetGroup")
	e, err := i.getLDAPGroupByNameOrID(nameOrID, true)
	if err != nil {
		return nil, err
	}
	sel := strings.Split(queryParam.Get("$select"), ",")
	exp := strings.Split(queryParam.Get("$expand"), ",")
	g := i.createGroupModelFromLDAP(e)
	if slices.Contains(sel, "members") || slices.Contains(exp, "members") {
		members, err := i.GetGroupMembers(ctx, *g.Id)
		if err != nil {
			return nil, err
		}
		if len(members) > 0 {
			m := make([]libregraph.User, 0, len(members))
			for _, u := range members {
				m = append(m, *u)
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
	i.logger.Debug().Str("backend", "ldap").Msgf("Search %s", i.groupBaseDN)
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
	i.logger.Debug().Str("backend", "ldap").Msg("GetGroups")

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
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
	searchRequest := ldap.NewSearchRequest(
		i.groupBaseDN, i.groupScope, ldap.NeverDerefAliases, 0, 0, false,
		groupFilter,
		[]string{
			i.groupAttributeMap.name,
			i.groupAttributeMap.id,
		},
		nil,
	)
	i.logger.Debug().Str("backend", "ldap").Str("Base", i.groupBaseDN).Str("filter", groupFilter).Msg("ldap search")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}

	groups := make([]*libregraph.Group, 0, len(res.Entries))

	for _, e := range res.Entries {
		sel := strings.Split(queryParam.Get("$select"), ",")
		exp := strings.Split(queryParam.Get("$expand"), ",")
		g := i.createGroupModelFromLDAP(e)
		if slices.Contains(sel, "members") || slices.Contains(exp, "members") {
			members, err := i.GetGroupMembers(ctx, *g.Id)
			if err != nil {
				return nil, err
			}
			if len(members) > 0 {
				m := make([]libregraph.User, 0, len(members))
				for _, u := range members {
					m = append(m, *u)
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
	e, err := i.getLDAPGroupByNameOrID(groupID, true)
	if err != nil {
		return nil, err
	}

	result := []*libregraph.User{}

	for _, memberDN := range e.GetEqualFoldAttributeValues(i.groupAttributeMap.member) {
		if memberDN == "" {
			continue
		}
		i.logger.Debug().Str("memberDN", memberDN).Msg("lookup")
		ue, err := i.getUserByDN(memberDN)
		if err != nil {
			// Ignore errors when reading a specific member fails, just log them and continue
			i.logger.Warn().Err(err).Str("member", memberDN).Msg("error reading group member")
			continue
		}
		result = append(result, i.createUserModelFromLDAP(ue))
	}

	return result, nil
}

// CreateGroup implements the Backend Interface for the LDAP Backend
// It is currently restricted to managing groups based on the "groupOfNames" ObjectClass.
// As "groupOfNames" requires a "member" Attribute to be present. Empty Groups (groups
// without a member) a represented by adding an empty DN as the single member.
func (i *LDAP) CreateGroup(ctx context.Context, group libregraph.Group) (*libregraph.Group, error) {
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
			i.logger.Warn().Str("memberDN", currentMember).Err(err).Msg("Couldn't parse DN")
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
			i.logger.Error().Str("new member", me.DN).Err(err).Msg("Couldn't parse DN")
			return err
		}
		if _, present := currentSet[nDN]; !present {
			newMemberDNs = append(newMemberDNs, me.DN)
		} else {
			i.logger.Debug().Str("memberDN", me.DN).Msg("Member already present in group. Skipping")
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
	ge, err := i.getLDAPGroupByID(groupID, true)
	if err != nil {
		i.logger.Warn().Str("backend", "ldap").Str("groupID", groupID).Msg("Error looking up group")
		return err
	}
	me, err := i.getLDAPUserByID(memberID)
	if err != nil {
		i.logger.Warn().Str("backend", "ldap").Str("memberID", memberID).Msg("Error looking up group member")
		return err
	}
	i.logger.Debug().Str("backend", "ldap").Str("groupdn", ge.DN).Str("member", me.DN).Msg("remove member")

	if mr, err := i.removeMemberFromGroupEntry(ge, me.DN); err == nil && mr != nil {
		return i.conn.Modify(mr)
	}
	return nil
}

func (i *LDAP) updateUserPassowrd(dn, password string) error {
	pwMod := ldap.PasswordModifyRequest{
		UserIdentity: dn,
		NewPassword:  password,
	}
	// Note: We can ignore the result message here, as it were only relevant if we requested
	// the server to generate a new Password
	_, err := i.conn.PasswordModify(&pwMod)
	if err != nil {
		var lerr *ldap.Error
		i.logger.Debug().Err(err).Msg("error setting password for user")
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
	return &libregraph.User{
		DisplayName:              pointerOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.displayName)),
		Mail:                     pointerOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.mail)),
		OnPremisesSamAccountName: pointerOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.userName)),
		Id:                       pointerOrNil(e.GetEqualFoldAttributeValue(i.userAttributeMap.id)),
	}
}

func (i *LDAP) createGroupModelFromLDAP(e *ldap.Entry) *libregraph.Group {
	return &libregraph.Group{
		DisplayName: pointerOrNil(e.GetEqualFoldAttributeValue(i.groupAttributeMap.name)),
		Id:          pointerOrNil(e.GetEqualFoldAttributeValue(i.groupAttributeMap.id)),
	}
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

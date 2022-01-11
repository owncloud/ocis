package identity

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

type LDAP struct {
	userBaseDN       string
	userFilter       string
	userScope        int
	userAttributeMap userAttributeMap

	groupBaseDN       string
	groupFilter       string
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
	name string
	id   string
}

func NewLDAPBackend(lc ldap.Client, config config.LDAP, logger *log.Logger) (*LDAP, error) {
	if config.UserDisplayNameAttribute == "" || config.UserIDAttribute == "" ||
		config.UserEmailAttribute == "" || config.UserNameAttribute == "" {
		return nil, fmt.Errorf("Invalid user attribute mappings")
	}
	uam := userAttributeMap{
		displayName: config.UserDisplayNameAttribute,
		id:          config.UserIDAttribute,
		mail:        config.UserEmailAttribute,
		userName:    config.UserNameAttribute,
	}

	if config.GroupNameAttribute == "" || config.GroupIDAttribute == "" {
		return nil, fmt.Errorf("Invalid group attribute mappings")
	}
	gam := groupAttributeMap{
		name: config.GroupNameAttribute,
		id:   config.GroupIDAttribute,
	}

	var userScope, groupScope int
	var err error
	if userScope, err = stringToScope(config.UserSearchScope); err != nil {
		return nil, fmt.Errorf("Error configuring user scope: %w", err)
	}

	if groupScope, err = stringToScope(config.GroupSearchScope); err != nil {
		return nil, fmt.Errorf("Error configuring group scope: %w", err)
	}

	return &LDAP{
		userBaseDN:        config.UserBaseDN,
		userFilter:        config.UserFilter,
		userScope:         userScope,
		userAttributeMap:  uam,
		groupBaseDN:       config.GroupBaseDN,
		groupFilter:       config.GroupFilter,
		groupScope:        groupScope,
		groupAttributeMap: gam,
		logger:            logger,
		conn:              lc,
	}, nil
}

// CreateUser implements the Backend Interface. It converts the libregraph.User into an
// LDAP User Entry (using the inetOrgPerson LDAP Objectclass) add adds that to the
// configured LDAP server
func (i *LDAP) CreateUser(ctx context.Context, user libregraph.User) (*libregraph.User, error) {

	ar := ldap.AddRequest{
		DN: fmt.Sprintf("uid=%s,%s", *user.OnPremisesSamAccountName, i.userBaseDN),
		Attributes: []ldap.Attribute{
			{
				Type: "objectClass",
				Vals: []string{"inetOrgPerson", "organizationalPerson", "person", "top"},
			},
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

	if user.PasswordProfile != nil && user.PasswordProfile.Password != nil {
		// TODO? This relies to the LDAP server to properly hash the password.
		// We might want to add support for the Password Modify LDAP Extended
		// Operation for servers that implement it. (Or implement client-side
		// hashing here.
		ar.Attribute("userPassword", []string{*user.PasswordProfile.Password})
	}

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
		return nil, err
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
	e, err := i.getLDAPUserByNameOrID(nameOrID)
	if err != nil {
		return err
	}
	dr := ldap.DelRequest{DN: e.DN}
	if err = i.conn.Del(&dr); err != nil {
		return err
	}
	return nil
}

func (i *LDAP) getUserByDN(dn string) (*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		dn, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 1, 0, false,
		"(objectclass=*)",
		[]string{
			i.userAttributeMap.displayName,
			i.userAttributeMap.id,
			i.userAttributeMap.mail,
			i.userAttributeMap.userName,
		},
		nil,
	)

	i.logger.Debug().Str("backend", "ldap").Str("dn", dn).Msg("Search user by DN")
	res, err := i.conn.Search(searchRequest)

	if err != nil {
		i.logger.Error().Err(err).Str("backend", "ldap").Str("dn", dn).Msg("Search user by DN failed")
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}
	if len(res.Entries) == 0 {
		return nil, errorcode.New(errorcode.ItemNotFound, "not found")
	}

	return res.Entries[0], nil
}

func (i *LDAP) getLDAPUserByNameOrID(nameOrID string) (*ldap.Entry, error) {
	nameOrID = ldap.EscapeFilter(nameOrID)
	searchRequest := ldap.NewSearchRequest(
		i.userBaseDN, i.userScope, ldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(&%s(|(%s=%s)(%s=%s)))", i.userFilter, i.userAttributeMap.userName, nameOrID, i.userAttributeMap.id, nameOrID),
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
				errmsg = fmt.Sprintf("too many results searching for user '%s'", nameOrID)
				i.logger.Debug().Str("backend", "ldap").Err(lerr).
					Str("user", nameOrID).Msg("too many results searching for user")
			}
		}
		return nil, errorcode.New(errorcode.ItemNotFound, errmsg)
	}
	if len(res.Entries) == 0 {
		return nil, errorcode.New(errorcode.ItemNotFound, "not found")
	}

	return res.Entries[0], nil
}

func (i *LDAP) GetUser(ctx context.Context, nameOrID string) (*libregraph.User, error) {
	i.logger.Debug().Str("backend", "ldap").Msg("GetUser")
	e, err := i.getLDAPUserByNameOrID(nameOrID)
	if err != nil {
		return nil, err
	}
	return i.createUserModelFromLDAP(e), nil
}

func (i *LDAP) GetUsers(ctx context.Context, queryParam url.Values) ([]*libregraph.User, error) {
	i.logger.Debug().Str("backend", "ldap").Msg("GetUsers")

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
	}
	userFilter := i.userFilter
	if search != "" {
		search = ldap.EscapeFilter(search)
		userFilter = fmt.Sprintf(
			"(&(%s)(|(%s=%s*)(%s=%s*)(%s=%s*)))",
			userFilter,
			i.userAttributeMap.userName, search,
			i.userAttributeMap.mail, search,
			i.userAttributeMap.displayName, search,
		)
	}
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
		users = append(users, i.createUserModelFromLDAP(e))
	}
	return users, nil
}

func (i *LDAP) GetGroup(ctx context.Context, groupID string) (*libregraph.Group, error) {
	i.logger.Debug().Str("backend", "ldap").Msg("GetGroup")
	groupID = ldap.EscapeFilter(groupID)
	searchRequest := ldap.NewSearchRequest(
		i.groupBaseDN, i.groupScope, ldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(&%s(|(%s=%s)(%s=%s)))", i.groupFilter, i.groupAttributeMap.name, groupID, i.groupAttributeMap.id, groupID),
		[]string{
			i.groupAttributeMap.name,
			i.groupAttributeMap.id,
		},
		nil,
	)
	i.logger.Debug().Str("backend", "ldap").Msgf("Search %s", i.groupBaseDN)
	res, err := i.conn.Search(searchRequest)

	if err != nil {
		var errmsg string
		if lerr, ok := err.(*ldap.Error); ok {
			if lerr.ResultCode == ldap.LDAPResultSizeLimitExceeded {
				errmsg = fmt.Sprintf("too many results searching for group '%s'", groupID)
				i.logger.Debug().Str("backend", "ldap").Err(lerr).Msg(errmsg)
			}
		}
		return nil, errorcode.New(errorcode.ItemNotFound, errmsg)
	}
	if len(res.Entries) == 0 {
		return nil, errorcode.New(errorcode.ItemNotFound, "not found")
	}

	return i.createGroupModelFromLDAP(res.Entries[0]), nil
}

func (i *LDAP) GetGroups(ctx context.Context, queryParam url.Values) ([]*libregraph.Group, error) {
	i.logger.Debug().Str("backend", "ldap").Msg("GetGroups")

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
	}
	groupFilter := i.groupFilter
	if search != "" {
		search = ldap.EscapeFilter(search)
		groupFilter = fmt.Sprintf(
			"(&(%s)(|(%s=%s*)(%s=%s*)))",
			groupFilter,
			i.groupAttributeMap.name, search,
			i.groupAttributeMap.id, search,
		)
	}
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
		groups = append(groups, i.createGroupModelFromLDAP(e))
	}
	return groups, nil
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
		DisplayName:              pointerOrNil(e.GetEqualFoldAttributeValue(i.groupAttributeMap.name)),
		OnPremisesSamAccountName: pointerOrNil(e.GetEqualFoldAttributeValue(i.groupAttributeMap.name)),
		Id:                       pointerOrNil(e.GetEqualFoldAttributeValue(i.groupAttributeMap.id)),
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
		return 0, fmt.Errorf("Invalid Scope '%s'", scope)
	}
	return s, nil
}

package glauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/asim/go-micro/v3/metadata"
	"github.com/glauth/glauth/pkg/config"
	"github.com/glauth/glauth/pkg/handler"
	"github.com/glauth/glauth/pkg/stats"
	ber "github.com/nmcclain/asn1-ber"
	"github.com/nmcclain/ldap"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
)

type queryType string

const (
	usersQuery  queryType = "users"
	groupsQuery queryType = "groups"
)

type ocisHandler struct {
	as          accounts.AccountsService
	gs          accounts.GroupsService
	log         log.Logger
	basedn      string
	nameFormat  string
	groupFormat string
	rbid        string
}

func (h ocisHandler) Bind(bindDN, bindSimplePw string, conn net.Conn) (ldap.LDAPResultCode, error) {
	bindDN = strings.ToLower(bindDN)
	baseDN := strings.ToLower("," + h.basedn)

	h.log.Debug().
		Str("handler", "ocis").
		Str("binddn", bindDN).
		Str("basedn", h.basedn).
		Interface("src", conn.RemoteAddr()).
		Msg("Bind request")

	stats.Frontend.Add("bind_reqs", 1)

	// parse the bindDN - ensure that the bindDN ends with the BaseDN
	if !strings.HasSuffix(bindDN, baseDN) {
		h.log.Error().
			Str("handler", "ocis").
			Str("binddn", bindDN).
			Str("basedn", h.basedn).
			Interface("src", conn.RemoteAddr()).
			Msg("BindDN not part of our BaseDN")
		return ldap.LDAPResultInvalidCredentials, nil
	}
	parts := strings.Split(strings.TrimSuffix(bindDN, baseDN), ",")
	if len(parts) > 2 {
		h.log.Error().
			Str("handler", "ocis").
			Str("binddn", bindDN).
			Int("numparts", len(parts)).
			Interface("src", conn.RemoteAddr()).
			Msg("BindDN should have only one or two parts")
		return ldap.LDAPResultInvalidCredentials, nil
	}
	userName := strings.TrimPrefix(parts[0], "cn=")

	// TODO make glauth context aware
	ctx := context.Background()

	// use a session with the bound user?
	roleIDs, err := json.Marshal([]string{h.rbid})
	if err != nil {
		h.log.Error().
			Err(err).
			Str("handler", "ocis").
			Msg("could not marshal roleid json")
		return ldap.LDAPResultOperationsError, nil
	}
	ctx = metadata.Set(ctx, middleware.RoleIDs, string(roleIDs))

	// check password
	res, err := h.as.ListAccounts(ctx, &accounts.ListAccountsRequest{
		//Query: fmt.Sprintf("username eq '%s'", username),
		// TODO this allows looking up users when you know the username using basic auth
		// adding the password to the query is an option but sending this over the wire a la scim seems ugly
		// but to set passwords our accounts need it anyway
		Query: fmt.Sprintf("login eq '%s' and password eq '%s'", userName, bindSimplePw),
	})
	if err != nil || len(res.Accounts) == 0 {
		h.log.Error().
			Err(err).
			Str("handler", "ocis").
			Str("username", userName).
			Str("binddn", bindDN).
			Interface("src", conn.RemoteAddr()).
			Msg("Login failed")
		return ldap.LDAPResultInvalidCredentials, nil
	}

	stats.Frontend.Add("bind_successes", 1)
	h.log.Debug().
		Str("handler", "ocis").
		Str("binddn", bindDN).
		Interface("src", conn.RemoteAddr()).
		Msg("Bind success")
	return ldap.LDAPResultSuccess, nil
}

func (h ocisHandler) Search(bindDN string, searchReq ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {
	bindDN = strings.ToLower(bindDN)
	baseDN := strings.ToLower("," + h.basedn)
	searchBaseDN := strings.ToLower(searchReq.BaseDN)
	h.log.Debug().
		Str("handler", "ocis").
		Str("binddn", bindDN).
		Str("basedn", h.basedn).
		Str("filter", searchReq.Filter).
		Interface("src", conn.RemoteAddr()).
		Msg("Search request")
	stats.Frontend.Add("search_reqs", 1)

	// validate the user is authenticated and has appropriate access
	if len(bindDN) < 1 {
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultInsufficientAccessRights,
		}, fmt.Errorf("search error: Anonymous BindDN not allowed %s", bindDN)
	}
	if !strings.HasSuffix(bindDN, baseDN) {
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultInsufficientAccessRights,
		}, fmt.Errorf("search error: BindDN %s not in our BaseDN %s", bindDN, h.basedn)
	}
	if !strings.HasSuffix(searchBaseDN, h.basedn) {
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultInsufficientAccessRights,
		}, fmt.Errorf("search error: search BaseDN %s is not in our BaseDN %s", searchBaseDN, h.basedn)
	}

	var qtype queryType = ""
	query := ""
	var code ldap.LDAPResultCode
	var err error
	if searchReq.Filter == "(&)" { // see Absolute True and False Filters in https://tools.ietf.org/html/rfc4526#section-2
		query = ""
	} else {
		var cf *ber.Packet
		cf, err = ldap.CompileFilter(searchReq.Filter)
		if err != nil {
			h.log.Error().
				Err(err).
				Str("handler", "ocis").
				Str("binddn", bindDN).
				Str("basedn", h.basedn).
				Str("filter", searchReq.Filter).
				Interface("src", conn.RemoteAddr()).
				Msg("could not compile filter")
			return ldap.ServerSearchResult{
				ResultCode: ldap.LDAPResultOperationsError,
			}, fmt.Errorf("Search Error: error compiling filter: %s, error: %s", searchReq.Filter, err.Error())
		}
		qtype, query, code, err = parseFilter(cf)
		if err != nil {
			return ldap.ServerSearchResult{
				ResultCode: code,
			}, fmt.Errorf("Search Error: error parsing filter: %s, error: %s", searchReq.Filter, err.Error())
		}

		// check if the searchBaseDN already has a username and add it to the query
		parts := strings.Split(strings.TrimSuffix(searchBaseDN, baseDN), ",")
		if len(parts) > 0 && strings.HasPrefix(parts[0], "cn=") {
			if len(query) > 0 {
				query += " AND "
			}
			query += fmt.Sprintf("on_premises_sam_account_name eq '%s'", escapeValue(strings.TrimPrefix(parts[0], "cn=")))
		}
	}

	// TODO make glauth context aware
	ctx := context.Background()

	// use a session with the bound user?
	roleIDs, err := json.Marshal([]string{h.rbid})
	if err != nil {
		h.log.Error().
			Err(err).
			Str("handler", "ocis").
			Msg("could not marshal roleid json")
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultOperationsError,
		}, nil
	}
	ctx = metadata.Set(ctx, middleware.RoleIDs, string(roleIDs))

	entries := []*ldap.Entry{}
	h.log.Debug().
		Str("handler", "ocis").
		Str("binddn", bindDN).
		Str("basedn", h.basedn).
		Str("filter", searchReq.Filter).
		Str("qtype", string(qtype)).
		Str("query", query).
		Msg("parsed query")
	switch qtype {
	case usersQuery:
		accounts, err := h.as.ListAccounts(ctx, &accounts.ListAccountsRequest{
			Query: query,
		})
		if err != nil {
			h.log.Error().
				Err(err).
				Str("handler", "ocis").
				Str("binddn", bindDN).
				Str("basedn", h.basedn).
				Str("filter", searchReq.Filter).
				Str("query", query).
				Interface("src", conn.RemoteAddr()).
				Msg("Could not list accounts")

			return ldap.ServerSearchResult{
				ResultCode: ldap.LDAPResultOperationsError,
			}, fmt.Errorf("search error: error listing users")
		}
		entries = append(entries, h.mapAccounts(accounts.Accounts)...)
	case groupsQuery:
		groups, err := h.gs.ListGroups(ctx, &accounts.ListGroupsRequest{
			Query: query,
		})
		if err != nil {
			h.log.Error().
				Err(err).
				Str("handler", "ocis").
				Str("binddn", bindDN).
				Str("basedn", h.basedn).
				Str("filter", searchReq.Filter).
				Str("query", query).
				Interface("src", conn.RemoteAddr()).
				Msg("Could not list groups")

			return ldap.ServerSearchResult{
				ResultCode: ldap.LDAPResultOperationsError,
			}, fmt.Errorf("search error: error listing groups")
		}
		entries = append(entries, h.mapGroups(groups.Groups)...)
	}

	stats.Frontend.Add("search_successes", 1)
	h.log.Debug().
		Str("handler", "ocis").
		Int("num_entries", len(entries)).
		Str("binddn", bindDN).
		Str("basedn", h.basedn).
		Str("filter", searchReq.Filter).
		Interface("src", conn.RemoteAddr()).
		Msg("AP: Search OK")

	return ldap.ServerSearchResult{
		Entries:    entries,
		Referrals:  []string{},
		Controls:   []ldap.Control{},
		ResultCode: ldap.LDAPResultSuccess,
	}, nil
}

func attribute(name string, values ...string) *ldap.EntryAttribute {
	return &ldap.EntryAttribute{
		Name:   name,
		Values: values,
	}
}

func (h ocisHandler) mapAccounts(accounts []*accounts.Account) []*ldap.Entry {
	entries := make([]*ldap.Entry, 0, len(accounts))
	for i := range accounts {
		attrs := []*ldap.EntryAttribute{
			attribute("objectClass", "posixAccount", "inetOrgPerson", "organizationalPerson", "Person", "top"),
			attribute("cn", accounts[i].PreferredName),
			attribute("uid", accounts[i].PreferredName),
			attribute("sn", accounts[i].PreferredName),
			attribute("homeDirectory", ""),
			attribute("ownCloudUUID", accounts[i].Id), // see https://github.com/butonic/owncloud-ldap-schema/blob/master/owncloud.schema#L28-L34
		}
		if accounts[i].DisplayName != "" {
			attrs = append(attrs, attribute("displayName", accounts[i].DisplayName))
		}
		if accounts[i].Mail != "" {
			attrs = append(attrs, attribute("mail", accounts[i].Mail))
		}
		if accounts[i].UidNumber != 0 { // TODO no root?
			attrs = append(attrs, attribute("uidnumber", strconv.FormatInt(accounts[i].UidNumber, 10)))
		}
		if accounts[i].GidNumber != 0 {
			attrs = append(attrs, attribute("gidnumber", strconv.FormatInt(accounts[i].GidNumber, 10)))
		}
		if accounts[i].Description != "" {
			attrs = append(attrs, attribute("description", accounts[i].Description))
		}

		dn := fmt.Sprintf("%s=%s,%s=%s,%s",
			h.nameFormat,
			accounts[i].PreferredName,
			h.groupFormat,
			"users",
			h.basedn,
		)
		entries = append(entries, &ldap.Entry{DN: dn, Attributes: attrs})
	}
	return entries
}

func (h ocisHandler) mapGroups(groups []*accounts.Group) []*ldap.Entry {
	entries := make([]*ldap.Entry, 0, len(groups))
	for i := range groups {
		attrs := []*ldap.EntryAttribute{
			attribute("objectClass", "posixGroup", "groupOfNames", "top"),
			attribute("cn", groups[i].OnPremisesSamAccountName),
			attribute("ownCloudUUID", groups[i].Id), // see https://github.com/butonic/owncloud-ldap-schema/blob/master/owncloud.schema#L28-L34
		}
		if groups[i].DisplayName != "" {
			attrs = append(attrs, attribute("displayName", groups[i].DisplayName))
		}
		if groups[i].GidNumber != 0 {
			attrs = append(attrs, attribute("gidnumber", strconv.FormatInt(groups[i].GidNumber, 10)))
		}
		if groups[i].Description != "" {
			attrs = append(attrs, attribute("description", groups[i].Description))
		}

		dn := fmt.Sprintf("%s=%s,%s=%s,%s",
			h.nameFormat,
			groups[i].OnPremisesSamAccountName,
			h.groupFormat,
			"groups",
			h.basedn,
		)

		memberUids := make([]string, len(groups[i].Members))
		for j := range groups[i].Members {
			memberUids[j] = groups[i].Members[j].PreferredName
		}
		attrs = append(attrs, attribute("memberuid", memberUids...))
		entries = append(entries, &ldap.Entry{DN: dn, Attributes: attrs})
	}
	return entries
}

// LDAP filters might ask for groups and users at the same time, eg.
// (|
//   (&(objectClass=posixaccount)(cn=einstein))
//   (&(objectClass=posixgroup)(cn=users))
// )

// (&(objectClass=posixaccount)(objectClass=posixgroup))
// qtype is one of
// "" not determined
// "users"
// "groups"
func parseFilter(f *ber.Packet) (queryType, string, ldap.LDAPResultCode, error) {
	var qtype queryType
	var q string
	var code ldap.LDAPResultCode
	var err error
	switch ldap.FilterMap[f.Tag] {
	case "Present":
		if len(f.Children) != 0 {
			return "", "", ldap.LDAPResultOperationsError, fmt.Errorf("equality match must have no children, got %+v", f)
		}
		attribute := strings.ToLower(f.Data.String())

		if attribute == "objectclass" {
			// TODO implement proper present odata query, for now fall back to listing users
			return "users", q, code, err
		}
		return qtype, q, ldap.LDAPResultUnwillingToPerform, fmt.Errorf("%s filter match for %s not implemented", ldap.FilterMap[f.Tag], attribute)
	case "Equality Match":
		if len(f.Children) != 2 {
			return "", "", ldap.LDAPResultOperationsError, fmt.Errorf("equality match must have exactly two children")
		}
		attribute := strings.ToLower(f.Children[0].Value.(string))
		value := f.Children[1].Value.(string)

		// replace attributes
		switch attribute {
		case "objectclass":
			switch strings.ToLower(value) {
			case "posixaccount", "shadowaccount", "users", "person", "inetorgperson", "organizationalperson":
				qtype = usersQuery
			case "posixgroup", "groups":
				qtype = groupsQuery
			case "*":
				// TODO not implemented yet
				qtype = usersQuery
			default:
				qtype = ""
			}
		case "ownclouduuid":
			q = fmt.Sprintf("id eq '%s'", escapeValue(value))
		case "cn", "uid":
			// on_premises_sam_account_name is indexed using the lowercase analyzer in ocis-accounts
			// TODO use "tolower(on_premises_sam_account_name) eq '%s'" to be clear about the case insensitive comparison
			q = fmt.Sprintf("on_premises_sam_account_name eq '%s'", escapeValue(value))
		case "mail":
			q = fmt.Sprintf("mail eq '%s'", escapeValue(value))
		case "displayname":
			q = fmt.Sprintf("display_name eq '%s'", escapeValue(value))
		case "uidnumber":
			if i, err := strconv.ParseUint(value, 10, 64); err != nil {
				code = ldap.LDAPResultInvalidAttributeSyntax
			} else {
				q = fmt.Sprintf("uid_number eq %d", i)
			}
		case "gidnumber":
			if i, err := strconv.ParseUint(value, 10, 64); err != nil {
				code = ldap.LDAPResultInvalidAttributeSyntax
			} else {
				q = fmt.Sprintf("gid_number eq %d", i)
			}
		case "description":
			q = fmt.Sprintf("description eq '%s'", escapeValue(value))
		default:
			code = ldap.LDAPResultUndefinedAttributeType
			err = fmt.Errorf("unrecognized assertion type '%s' in filter item", attribute)
		}
		return qtype, q, code, err
	case "Substrings":
		if len(f.Children) != 2 {
			return "", "", ldap.LDAPResultOperationsError, fmt.Errorf("substrings filter must have exactly two children")
		}
		attribute := strings.ToLower(f.Children[0].Value.(string))
		if len(f.Children[1].Children) != 1 {
			return "", "", ldap.LDAPResultUnwillingToPerform, fmt.Errorf("substrings filter only supports prefix match")
		}
		value := f.Children[1].Children[0].Value.(string)

		// replace attributes
		switch attribute {
		case "objectclass":
			switch strings.ToLower(value) {
			case "posixaccount", "shadowaccount", "users", "person", "inetorgperson", "organizationalperson":
				qtype = usersQuery
			case "posixgroup", "groups":
				qtype = groupsQuery
			default:
				qtype = ""
			}
		case "ownclouduuid":
			q = fmt.Sprintf("startswith(id,'%s')", escapeValue(value))
		case "cn", "uid":
			// on_premises_sam_account_name is indexed using the lowercase analyzer in ocis-accounts
			// TODO use "tolower(on_premises_sam_account_name) eq '%s'" to be clear about the case insensitive comparison
			q = fmt.Sprintf("startswith(on_premises_sam_account_name,'%s')", escapeValue(value))
		case "mail":
			q = fmt.Sprintf("startswith(mail,'%s')", escapeValue(value))
		case "displayname":
			q = fmt.Sprintf("startswith(display_name,'%s')", escapeValue(value))
		case "description":
			q = fmt.Sprintf("startswith(description,'%s')", escapeValue(value))
		default:
			code = ldap.LDAPResultUndefinedAttributeType
			err = fmt.Errorf("unrecognized assertion type '%s' in filter item", attribute)
		}
		return qtype, q, code, err
	case "And", "Or":
		subQueries := []string{}
		for i := range f.Children {
			var subQuery string
			var qt queryType
			qt, subQuery, code, err = parseFilter(f.Children[i])
			if err != nil {
				return "", "", code, err
			}
			if qtype == "" {
				qtype = qt
			} else if qt != "" && qt != qtype {
				return "", "", ldap.LDAPResultUnwillingToPerform, fmt.Errorf("mixing user and group filters not supported")
			}
			if subQuery != "" {
				subQueries = append(subQueries, subQuery)
			}
		}
		return qtype, strings.Join(subQueries, " "+strings.ToLower(ldap.FilterMap[f.Tag])+" "), ldap.LDAPResultSuccess, nil
	case "Not":
		if len(f.Children) != 1 {
			return "", "", ldap.LDAPResultOperationsError, fmt.Errorf("not filter match must have exactly one child")
		}
		qtype, subQuery, code, err := parseFilter(f.Children[0])
		if err != nil {
			return "", "", code, err
		}
		if subQuery != "" {
			q = fmt.Sprintf("not %s", subQuery)
		}
		return qtype, q, code, nil
	}
	return qtype, q, ldap.LDAPResultUnwillingToPerform, fmt.Errorf("%s filter not implemented", ldap.FilterMap[f.Tag])
}

// escapeValue escapes all special characters in the value
func escapeValue(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func (h ocisHandler) Close(boundDN string, conn net.Conn) error {
	stats.Frontend.Add("closes", 1)
	return nil
}

// Add is not yet supported for the ocis backend
func (h ocisHandler) Add(boundDN string, req ldap.AddRequest, conn net.Conn) (result ldap.LDAPResultCode, err error) {
	return ldap.LDAPResultInsufficientAccessRights, nil
}

// Modify is not yet supported for the ocis backend
func (h ocisHandler) Modify(boundDN string, req ldap.ModifyRequest, conn net.Conn) (result ldap.LDAPResultCode, err error) {
	return ldap.LDAPResultInsufficientAccessRights, nil
}

// Delete is not yet supported for the ocis backend
func (h ocisHandler) Delete(boundDN string, deleteDN string, conn net.Conn) (result ldap.LDAPResultCode, err error) {
	return ldap.LDAPResultInsufficientAccessRights, nil
}

// FindUser with the given username
func (h ocisHandler) FindUser(userName string) (found bool, user config.User, err error) {
	return false, config.User{}, nil
}

// NewOCISHandler implements a glauth backend with ocis-accounts as the datasource
func NewOCISHandler(opts ...Option) handler.Handler {
	options := newOptions(opts...)

	handler := ocisHandler{
		log:         options.Logger,
		as:          options.AccountsService,
		gs:          options.GroupsService,
		basedn:      options.BaseDN,
		nameFormat:  options.NameFormat,
		groupFormat: options.GroupFormat,
		rbid:        options.RoleBundleUUID,
	}
	return handler
}

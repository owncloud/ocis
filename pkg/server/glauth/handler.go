package glauth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/glauth/glauth/pkg/config"
	"github.com/glauth/glauth/pkg/handler"
	"github.com/glauth/glauth/pkg/stats"
	ber "github.com/nmcclain/asn1-ber"
	"github.com/nmcclain/ldap"
	accounts "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-pkg/v2/log"
)

type ocisHandler struct {
	as  accounts.AccountsService
	log log.Logger
	cfg *config.Config
}

func (h ocisHandler) Bind(bindDN, bindSimplePw string, conn net.Conn) (ldap.LDAPResultCode, error) {
	bindDN = strings.ToLower(bindDN)
	baseDN := strings.ToLower("," + h.cfg.Backend.BaseDN)

	h.log.Debug().Str("binddn", bindDN).Str("basedn", h.cfg.Backend.BaseDN).Interface("src", conn.RemoteAddr()).Msg("Bind request")

	stats.Frontend.Add("bind_reqs", 1)

	// parse the bindDN - ensure that the bindDN ends with the BaseDN
	if !strings.HasSuffix(bindDN, baseDN) {
		h.log.Error().Str("binddn", bindDN).Str("basedn", h.cfg.Backend.BaseDN).Interface("src", conn.RemoteAddr()).Msg("BindDN not part of our BaseDN")
		return ldap.LDAPResultInvalidCredentials, nil
	}
	parts := strings.Split(strings.TrimSuffix(bindDN, baseDN), ",")
	if len(parts) > 2 {
		h.log.Error().Str("binddn", bindDN).Int("numparts", len(parts)).Interface("src", conn.RemoteAddr()).Msg("BindDN should have only one or two parts")
		return ldap.LDAPResultInvalidCredentials, nil
	}
	userName := strings.TrimPrefix(parts[0], "cn=")

	// check password
	_, err := h.as.ListAccounts(context.TODO(), &accounts.ListAccountsRequest{
		//Query: fmt.Sprintf("username eq '%s'", username),
		// TODO this allows lookung up users when you know the username using basic auth
		// adding the password to the query is an option but sending the sover the wira a la scim seems ugly
		// but to set passwords our accounts need it anyway
		Query: fmt.Sprintf("login eq '%s' and password eq '%s'", userName, bindSimplePw),
	})
	if err != nil {
		h.log.Error().Str("username", userName).Str("binddn", bindDN).Interface("src", conn.RemoteAddr()).Msg("Login failed")
		return ldap.LDAPResultInvalidCredentials, nil
	}

	stats.Frontend.Add("bind_successes", 1)
	h.log.Debug().Str("binddn", bindDN).Interface("src", conn.RemoteAddr()).Msg("Bind success")
	return ldap.LDAPResultSuccess, nil
}

func (h ocisHandler) Search(bindDN string, searchReq ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {
	bindDN = strings.ToLower(bindDN)
	baseDN := strings.ToLower("," + h.cfg.Backend.BaseDN)
	searchBaseDN := strings.ToLower(searchReq.BaseDN)
	h.log.Debug().Str("binddn", bindDN).Str("basedn", h.cfg.Backend.BaseDN).Str("filter", searchReq.Filter).Interface("src", conn.RemoteAddr()).Msg("Search request")
	stats.Frontend.Add("search_reqs", 1)

	// validate the user is authenticated and has appropriate access
	if len(bindDN) < 1 {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("search error: Anonymous BindDN not allowed %s", bindDN)
	}
	if !strings.HasSuffix(bindDN, baseDN) {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("search error: BindDN %s not in our BaseDN %s", bindDN, h.cfg.Backend.BaseDN)
	}
	if !strings.HasSuffix(searchBaseDN, h.cfg.Backend.BaseDN) {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("search error: search BaseDN %s is not in our BaseDN %s", searchBaseDN, h.cfg.Backend.BaseDN)
	}

	qtype := ""
	query := ""
	var err error
	if searchReq.Filter == "(&)" { // see Absolute True and False Filters in https://tools.ietf.org/html/rfc4526#section-2
		query = ""
	} else {
		var cf *ber.Packet
		cf, err = ldap.CompileFilter(searchReq.Filter)
		if err != nil {
			return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: error parsing filter: %s", searchReq.Filter)
		}
		qtype, query, err = parseFilter(cf)
		if err != nil {
			return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: error parsing filter: %s", searchReq.Filter)
		}
	}

	entries := []*ldap.Entry{}
	h.log.Debug().Str("binddn", bindDN).Str("basedn", h.cfg.Backend.BaseDN).Str("filter", searchReq.Filter).Str("qtype", qtype).Str("query", query).Msg("parsed query")
	if qtype == "users" {
		accounts, err := h.as.ListAccounts(context.TODO(), &accounts.ListAccountsRequest{
			Query: query,
		})
		if err != nil {
			h.log.Error().Err(err).Str("binddn", bindDN).Str("basedn", h.cfg.Backend.BaseDN).Str("filter", searchReq.Filter).Str("query", query).Interface("src", conn.RemoteAddr()).Msg("Could not list accounts")
			return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, errors.New("search error: error getting users")
		}
		for i := range accounts.Accounts {
			attrs := []*ldap.EntryAttribute{
				{Name: "objectClass", Values: []string{"posixAccount", "inetOrgPerson", "organizationalPerson", "Person", "top"}},
				{Name: "cn", Values: []string{accounts.Accounts[i].PreferredName}},
				{Name: "uid", Values: []string{accounts.Accounts[i].PreferredName}},
				{Name: "sn", Values: []string{accounts.Accounts[i].PreferredName}}, // must be set for a valid person
			}
			if accounts.Accounts[i].DisplayName != "" {
				attrs = append(attrs, &ldap.EntryAttribute{Name: "displayName", Values: []string{accounts.Accounts[i].DisplayName}})
			}
			if accounts.Accounts[i].Mail != "" {
				attrs = append(attrs, &ldap.EntryAttribute{Name: "mail", Values: []string{accounts.Accounts[i].Mail}})
			}
			if accounts.Accounts[i].UidNumber != 0 { // TODO no root?
				attrs = append(attrs, &ldap.EntryAttribute{Name: "uidnumber", Values: []string{strconv.FormatInt(accounts.Accounts[i].UidNumber, 10)}})
			}
			if accounts.Accounts[i].GidNumber != 0 {
				attrs = append(attrs, &ldap.EntryAttribute{Name: "gidnumber", Values: []string{strconv.FormatInt(accounts.Accounts[i].GidNumber, 10)}})
			}
			if accounts.Accounts[i].Description != "" {
				attrs = append(attrs, &ldap.EntryAttribute{Name: "description", Values: []string{accounts.Accounts[i].Description}})
			}

			dn := fmt.Sprintf("%s=%s,%s=%s,%s", h.cfg.Backend.NameFormat, accounts.Accounts[i].PreferredName, h.cfg.Backend.GroupFormat, "users", h.cfg.Backend.BaseDN)
			entries = append(entries, &ldap.Entry{DN: dn, Attributes: attrs})
		}
	}

	stats.Frontend.Add("search_successes", 1)
	h.log.Debug().Str("binddn", bindDN).Str("basedn", h.cfg.Backend.BaseDN).Str("filter", searchReq.Filter).Interface("src", conn.RemoteAddr()).Msg("AP: Search OK")
	return ldap.ServerSearchResult{Entries: entries, Referrals: []string{}, Controls: []ldap.Control{}, ResultCode: ldap.LDAPResultSuccess}, nil
}

// LDAP filters might ask for grouips and users at the same time, eg.
// (|
//   (&(objectClass=posixaccount)(cn=einstein))
//   (&(objectClass=posixgroup)(cn=users))
// )

// (&(objectClass=posixaccount)(objectClass=posixgroup))
// qtype is one of
// "" not determined
// "users"
// "groups"
func parseFilter(f *ber.Packet) (qtype string, q string, err error) {
	switch ldap.FilterMap[f.Tag] {
	case "Equality Match":
		if len(f.Children) != 2 {
			return "", "", errors.New("equality match must have only two children")
		}
		attribute := strings.ToLower(f.Children[0].Value.(string))
		value := f.Children[1].Value.(string)

		// replace attributes
		switch attribute {
		case "objectclass":
			switch value {
			case "posixaccount", "shadowaccount", "users", "person", "inetorgperson", "organizationalperson":
				qtype = "users"
			case "posixgroup", "groups":
				qtype = "groups"
			default:
				qtype = ""
			}
			return qtype, "", nil
		case "cn", "uid":
			return "", fmt.Sprintf("preferred_name eq '%s'", strings.ReplaceAll(value, "'", "''")), nil
		case "mail":
			return "", fmt.Sprintf("mail eq '%s'", strings.ReplaceAll(value, "'", "''")), nil
		case "displayname":
			return "", fmt.Sprintf("display_name eq '%s'", strings.ReplaceAll(value, "'", "''")), nil
		case "uidnumber":
			return "", fmt.Sprintf("uid_number eq '%s'", strings.ReplaceAll(value, "'", "''")), nil
		case "gidnumber":
			return "", fmt.Sprintf("gid_number eq '%s'", strings.ReplaceAll(value, "'", "''")), nil
		case "description":
			return "", fmt.Sprintf("description eq '%s'", strings.ReplaceAll(value, "'", "''")), nil
		}

		return "", "", fmt.Errorf("filter by %s not implmented", attribute)
	case "And":
		subQueries := []string{}
		for _, child := range f.Children {
			var subQuery string
			var qt string
			qt, subQuery, err = parseFilter(child)
			if err != nil {
				return "", "", err
			}
			if qtype == "" {
				qtype = qt
			} else if qt != "" && qt != qtype {
				return "", "", fmt.Errorf("mixing user and group filters not supported")
			}
			if subQuery != "" {
				subQueries = append(subQueries, subQuery)
			}
		}
		return qtype, strings.Join(subQueries, " and "), nil
	case "Or":
		subQueries := []string{}
		for _, child := range f.Children {
			var subQuery string
			var qt string
			qt, subQuery, err = parseFilter(child)
			if err != nil {
				return "", "", err
			}
			if qtype == "" {
				qtype = qt
			} else if qt != "" && qt != qtype {
				return "", "", fmt.Errorf("mixing user and group filters not supported")
			}
			if subQuery != "" {
				subQueries = append(subQueries, subQuery)
			}
		}
		return qtype, strings.Join(subQueries, " or "), nil
	case "Not":
		if len(f.Children) != 1 {
			return "", "", errors.New("not filter must have only one child")
		}
		qtype, subQuery, err := parseFilter(f.Children[0])
		if err != nil {
			return "", "", err
		}
		if subQuery != "" {
			q = fmt.Sprintf("not %s", subQuery)
		}
		return qtype, q, nil
	}
	return
}

func (h ocisHandler) Close(boundDN string, conn net.Conn) error {
	stats.Frontend.Add("closes", 1)
	return nil
}

// NewOCISHandler implements a glauth backend with ocis-accounts as tdhe datasource
func NewOCISHandler(opts ...Option) handler.Handler {
	options := newOptions(opts...)

	handler := ocisHandler{
		log: options.Logger,
		cfg: options.Config,
		as:  options.AccountsService,
	}
	return handler
}

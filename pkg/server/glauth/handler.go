package glauth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/glauth/glauth/pkg/config"
	"github.com/glauth/glauth/pkg/handler"
	"github.com/glauth/glauth/pkg/stats"
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
	// return all users in the config file - the LDAP library will filter results for us
	entries := []*ldap.Entry{}
	filterEntity, err := ldap.GetFilterObjectClass(searchReq.Filter)
	if err != nil {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("search error: error parsing filter: %s", searchReq.Filter)
	}

	switch filterEntity {
	default:
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("search error: unhandled filter type: %s [%s]", filterEntity, searchReq.Filter)
	case "posixgroup":
		/*
			groups, err := session.getGroups()
			if err != nil {
				return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, errors.New("search error: error getting groups")
			}
			for _, g := range groups {
				attrs := []*ldap.EntryAttribute{}
				attrs = append(attrs, &ldap.EntryAttribute{Name: "cn", Values: []string{*g.ID}})
				attrs = append(attrs, &ldap.EntryAttribute{Name: "description", Values: []string{fmt.Sprintf("%s from ownCloud", *g.ID)}})
				//			attrs = append(attrs, &ldap.EntryAttribute{"gidNumber", []string{fmt.Sprintf("%d", g.UnixID)}})
				attrs = append(attrs, &ldap.EntryAttribute{Name: "objectClass", Values: []string{"posixGroup"}})
				if g.Members != nil {
					members := make([]string, len(g.Members))
					for i, v := range g.Members {
						members[i] = *v.ID
					}

					attrs = append(attrs, &ldap.EntryAttribute{Name: "memberUid", Values: members})
				}
				dn := fmt.Sprintf("cn=%s,%s=groups,%s", *g.ID, h.cfg.Backend.GroupFormat, h.cfg.Backend.BaseDN)
				entries = append(entries, &ldap.Entry{DN: dn, Attributes: attrs})
			}
		*/
	case "posixaccount", "":
		userName := ""
		if searchBaseDN != strings.ToLower(h.cfg.Backend.BaseDN) {
			parts := strings.Split(strings.TrimSuffix(searchBaseDN, baseDN), ",")
			if len(parts) >= 1 {
				userName = strings.TrimPrefix(parts[0], "cn=")
			}
		}
		accounts, err := h.as.ListAccounts(context.TODO(), &accounts.ListAccountsRequest{
			Query: fmt.Sprintf("preferred_name eq '%s'", strings.ReplaceAll(userName, "'", "''")),
		})
		if err != nil {
			h.log.Error().Err(err).Str("username", userName).Interface("src", conn.RemoteAddr()).Msg("Could not list accounts")
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

func (h ocisHandler) Close(boundDN string, conn net.Conn) error {
	stats.Frontend.Add("closes", 1)
	return nil
}

// helper functions
func connID(conn net.Conn) string {
	h := sha256.New()
	h.Write([]byte(conn.LocalAddr().String() + conn.RemoteAddr().String()))
	sha := fmt.Sprintf("% x", h.Sum(nil))
	return string(sha)
}

func NewOCISHandler(opts ...Option) handler.Handler {
	options := newOptions(opts...)

	handler := ocisHandler{
		log: options.Logger,
		cfg: options.Config,
		as:  options.AccountsService,
	}
	return handler
}

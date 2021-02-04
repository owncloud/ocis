package svc

import (
	"errors"

	"github.com/go-ldap/ldap/v3"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func (g Graph) ldapGetSingleEntry(baseDn string, filter string) (*ldap.Entry, error) {
	conn, err := g.initLdap()
	if err != nil {
		return nil, err
	}
	result, err := g.ldapSearch(conn, filter, baseDn)
	if err != nil {
		return nil, err
	}
	if len(result.Entries) == 0 {
		return nil, errors.New("resource not found")
	}
	return result.Entries[0], nil
}

func (g Graph) initLdap() (*ldap.Conn, error) {
	g.logger.Info().Msgf("Dialing ldap %s://%s", g.config.Ldap.Network, g.config.Ldap.Address)
	con, err := ldap.Dial(g.config.Ldap.Network, g.config.Ldap.Address)

	if err != nil {
		return nil, err
	}

	if err := con.Bind(g.config.Ldap.UserName, g.config.Ldap.Password); err != nil {
		return nil, err
	}
	return con, nil
}

func (g Graph) ldapSearch(con *ldap.Conn, filter string, baseDN string) (*ldap.SearchResult, error) {
	search := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"dn",
			"uid",
			"givenname",
			"mail",
			"displayname",
			"entryuuid",
			"sn",
			"cn",
		},
		nil,
	)

	return con.Search(search)
}

func createUserModelFromLDAP(entry *ldap.Entry) *msgraph.User {
	displayName := entry.GetAttributeValue("displayname")
	givenName := entry.GetAttributeValue("givenname")
	mail := entry.GetAttributeValue("mail")
	surName := entry.GetAttributeValue("sn")
	id := entry.GetAttributeValue("entryuuid")
	return &msgraph.User{
		DisplayName: &displayName,
		GivenName:   &givenName,
		Surname:     &surName,
		Mail:        &mail,
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &id,
			},
		},
	}
}

func createGroupModelFromLDAP(entry *ldap.Entry) *msgraph.Group {
	id := entry.GetAttributeValue("entryuuid")
	displayName := entry.GetAttributeValue("cn")

	return &msgraph.Group{
		DisplayName: &displayName,
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &id,
			},
		},
	}
}

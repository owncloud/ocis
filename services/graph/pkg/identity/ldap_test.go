package identity

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/test-go/testify/mock"
)

func getMockedBackend(l ldap.Client, lc config.LDAP, logger *log.Logger) (*LDAP, error) {
	return NewLDAPBackend(l, lconfig, logger)
}

var lconfig = config.LDAP{
	UserBaseDN:               "ou=people,dc=test",
	UserObjectClass:          "inetOrgPerson",
	UserSearchScope:          "sub",
	UserFilter:               "",
	UserDisplayNameAttribute: "displayname",
	UserIDAttribute:          "entryUUID",
	UserEmailAttribute:       "mail",
	UserNameAttribute:        "uid",

	GroupBaseDN:        "ou=groups,dc=test",
	GroupObjectClass:   "groupOfNames",
	GroupSearchScope:   "sub",
	GroupFilter:        "",
	GroupNameAttribute: "cn",
	GroupIDAttribute:   "entryUUID",
}

var userEntry = ldap.NewEntry("uid=user",
	map[string][]string{
		"uid":         {"user"},
		"displayname": {"DisplayName"},
		"mail":        {"user@example"},
		"entryuuid":   {"abcd-defg"},
	})
var invalidUserEntry = ldap.NewEntry("uid=user",
	map[string][]string{
		"uid":         {"invalid"},
		"displayname": {"DisplayName"},
		"mail":        {"user@example"},
	})
var groupEntry = ldap.NewEntry("cn=group",
	map[string][]string{
		"cn":        {"group"},
		"entryuuid": {"abcd-defg"},
		"member": {
			"uid=user,ou=people,dc=test",
			"uid=invalid,ou=people,dc=test",
		},
	})
var invalidGroupEntry = ldap.NewEntry("cn=invalid",
	map[string][]string{
		"cn": {"invalid"},
	})

var logger = log.NewLogger(log.Level("debug"))

func TestNewLDAPBackend(t *testing.T) {

	l := &mocks.Client{}

	tc := lconfig
	tc.UserDisplayNameAttribute = ""
	if _, err := NewLDAPBackend(l, tc, &logger); err == nil {
		t.Error("Should fail with incomplete user attr config")
	}

	tc = lconfig
	tc.GroupIDAttribute = ""
	if _, err := NewLDAPBackend(l, tc, &logger); err == nil {
		t.Errorf("Should fail with incomplete group	config")
	}

	tc = lconfig
	tc.UserSearchScope = ""
	if _, err := NewLDAPBackend(l, tc, &logger); err == nil {
		t.Errorf("Should fail with invalid user search scope")
	}

	tc = lconfig
	tc.GroupSearchScope = ""
	if _, err := NewLDAPBackend(l, tc, &logger); err == nil {
		t.Errorf("Should fail with invalid group search scope")
	}

	if _, err := NewLDAPBackend(l, lconfig, &logger); err != nil {
		t.Errorf("Should fail with invalid group search scope")
	}

}

func TestCreateUserModelFromLDAP(t *testing.T) {
	l := &mocks.Client{}
	logger := log.NewLogger(log.Level("debug"))

	b, _ := NewLDAPBackend(l, lconfig, &logger)
	if user := b.createUserModelFromLDAP(nil); user != nil {
		t.Errorf("createUserModelFromLDAP should return on nil Entry")
	}
	user := b.createUserModelFromLDAP(userEntry)
	if user == nil {
		t.Error("Converting a valid LDAP Entry should succeed")
	} else {
		if *user.OnPremisesSamAccountName != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.userName) {
			t.Errorf("Error creating msGraph User from LDAP Entry: %v != %v", user.OnPremisesSamAccountName, pointerOrNil(userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.userName)))
		}
		if *user.Mail != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.mail) {
			t.Errorf("Error creating msGraph User from LDAP Entry: %s != %s", *user.Mail, userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.mail))
		}
		if *user.DisplayName != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.displayName) {
			t.Errorf("Error creating msGraph User from LDAP Entry: %s != %s", *user.DisplayName, userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.displayName))
		}
		if *user.Id != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.id) {
			t.Errorf("Error creating msGraph User from LDAP Entry: %s != %s", *user.Id, userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.id))
		}
	}
}

func TestGetUser(t *testing.T) {
	// Mock a Sizelimit Error
	lm := &mocks.Client{}
	lm.On("Search", mock.Anything).
		Return(
			nil, ldap.NewError(ldap.LDAPResultSizeLimitExceeded, errors.New("mock")))
	b, _ := getMockedBackend(lm, lconfig, &logger)

	queryParamExpand := url.Values{
		"$expand": []string{"memberOf"},
	}
	queryParamSelect := url.Values{
		"$select": []string{"memberOf"},
	}
	_, err := b.GetUser(context.Background(), "fred", nil)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	_, err = b.GetUser(context.Background(), "fred", queryParamExpand)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	_, err = b.GetUser(context.Background(), "fred", queryParamSelect)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	// Mock an empty Search Result
	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	_, err = b.GetUser(context.Background(), "fred", nil)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	_, err = b.GetUser(context.Background(), "fred", queryParamExpand)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	_, err = b.GetUser(context.Background(), "fred", queryParamSelect)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	// Mock a valid Search Result
	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{
				Entries: []*ldap.Entry{userEntry},
			},
			nil)

	b, _ = getMockedBackend(lm, lconfig, &logger)
	u, err := b.GetUser(context.Background(), "user", nil)
	if err != nil {
		t.Errorf("Expected GetUser to succeed. Got %s", err.Error())
	} else if *u.Id != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.id) {
		t.Errorf("Expected GetUser to return a valid user")
	}

	u, err = b.GetUser(context.Background(), "user", queryParamExpand)
	if err != nil {
		t.Errorf("Expected GetUser to succeed. Got %s", err.Error())
	} else if *u.Id != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.id) {
		t.Errorf("Expected GetUser to return a valid user")
	}

	u, err = b.GetUser(context.Background(), "user", queryParamSelect)
	if err != nil {
		t.Errorf("Expected GetUser to succeed. Got %s", err.Error())
	} else if *u.Id != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.id) {
		t.Errorf("Expected GetUser to return a valid user")
	}

	// Mock invalid Search Result
	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{
				Entries: []*ldap.Entry{invalidUserEntry},
			},
			nil)

	b, _ = getMockedBackend(lm, lconfig, &logger)
	u, err = b.GetUser(context.Background(), "invalid", nil)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}
}

func TestGetUsers(t *testing.T) {
	// Mock a Sizelimit Error
	lm := &mocks.Client{}
	lm.On("Search", mock.Anything).Return(nil, ldap.NewError(ldap.LDAPResultOperationsError, errors.New("mock")))

	b, _ := getMockedBackend(lm, lconfig, &logger)
	_, err := b.GetUsers(context.Background(), url.Values{})
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).Return(&ldap.SearchResult{}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err := b.GetUsers(context.Background(), url.Values{})
	if err != nil {
		t.Errorf("Expected success, got '%s'", err.Error())
	} else if g == nil || len(g) != 0 {
		t.Errorf("Expected zero length user slice")
	}
}

func TestGetGroup(t *testing.T) {
	// Mock a Sizelimit Error
	lm := &mocks.Client{}
	lm.On("Search", mock.Anything).Return(nil, ldap.NewError(ldap.LDAPResultSizeLimitExceeded, errors.New("mock")))

	queryParamExpand := url.Values{
		"$expand": []string{"members"},
	}
	queryParamSelect := url.Values{
		"$select": []string{"members"},
	}
	b, _ := getMockedBackend(lm, lconfig, &logger)
	_, err := b.GetGroup(context.Background(), "group", nil)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}
	_, err = b.GetGroup(context.Background(), "group", queryParamExpand)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}
	_, err = b.GetGroup(context.Background(), "group", queryParamSelect)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	// Mock an empty Search Result
	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).Return(&ldap.SearchResult{}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	_, err = b.GetGroup(context.Background(), "group", nil)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}
	_, err = b.GetGroup(context.Background(), "group", queryParamExpand)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}
	_, err = b.GetGroup(context.Background(), "group", queryParamSelect)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	// Mock an invalid Search Result
	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).Return(&ldap.SearchResult{
		Entries: []*ldap.Entry{invalidGroupEntry},
	}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err := b.GetGroup(context.Background(), "group", nil)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}
	g, err = b.GetGroup(context.Background(), "group", queryParamExpand)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}
	g, err = b.GetGroup(context.Background(), "group", queryParamSelect)
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	// Mock a valid	Search Result
	lm = &mocks.Client{}
	sr1 := &ldap.SearchRequest{
		BaseDN:     "ou=groups,dc=test",
		Scope:      2,
		SizeLimit:  1,
		Filter:     "(&(objectClass=groupOfNames)(|(cn=group)(entryUUID=group)))",
		Attributes: []string{"cn", "entryUUID", "member"},
		Controls:   []ldap.Control(nil),
	}
	sr2 := &ldap.SearchRequest{
		BaseDN:     "uid=user,ou=people,dc=test",
		SizeLimit:  1,
		Filter:     "(objectclass=*)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid"},
		Controls:   []ldap.Control(nil),
	}
	sr3 := &ldap.SearchRequest{
		BaseDN:     "uid=invalid,ou=people,dc=test",
		SizeLimit:  1,
		Filter:     "(objectclass=*)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid"},
		Controls:   []ldap.Control(nil),
	}

	lm.On("Search", sr1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{groupEntry}}, nil)
	lm.On("Search", sr2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{userEntry}}, nil)
	lm.On("Search", sr3).Return(&ldap.SearchResult{Entries: []*ldap.Entry{invalidUserEntry}}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err = b.GetGroup(context.Background(), "group", nil)
	if err != nil {
		t.Errorf("Expected GetGroup to succeed. Got %s", err.Error())
	} else if *g.Id != groupEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id) {
		t.Errorf("Expected GetGroup to return a valid group")
	}
	g, err = b.GetGroup(context.Background(), "group", queryParamExpand)
	switch {
	case err != nil:
		t.Errorf("Expected GetGroup to succeed. Got %s", err.Error())
	case g.GetId() != groupEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id):
		t.Errorf("Expected GetGroup to return a valid group")
	case len(g.Members) != 1:
		t.Errorf("Expected GetGroup with expand to return one member")
	case g.Members[0].GetId() != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.id):
		t.Errorf("Expected GetGroup with expand to return correct member")
	}
	g, err = b.GetGroup(context.Background(), "group", queryParamSelect)
	switch {
	case err != nil:
		t.Errorf("Expected GetGroup to succeed. Got %s", err.Error())
	case g.GetId() != groupEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id):
		t.Errorf("Expected GetGroup to return a valid group")
	case len(g.Members) != 1:
		t.Errorf("Expected GetGroup with expand to return one member")
	case g.Members[0].GetId() != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.id):
		t.Errorf("Expected GetGroup with expand to return correct member")
	}
}

func TestGetGroups(t *testing.T) {
	queryParamExpand := url.Values{
		"$expand": []string{"members"},
	}
	queryParamSelect := url.Values{
		"$select": []string{"members"},
	}
	lm := &mocks.Client{}
	lm.On("Search", mock.Anything).Return(nil, ldap.NewError(ldap.LDAPResultOperationsError, errors.New("mock")))
	b, _ := getMockedBackend(lm, lconfig, &logger)
	_, err := b.GetGroups(context.Background(), url.Values{})
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).Return(&ldap.SearchResult{}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err := b.GetGroups(context.Background(), url.Values{})
	if err != nil {
		t.Errorf("Expected success, got '%s'", err.Error())
	} else if g == nil || len(g) != 0 {
		t.Errorf("Expected zero length user slice")
	}

	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).Return(&ldap.SearchResult{
		Entries: []*ldap.Entry{groupEntry},
	}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err = b.GetGroups(context.Background(), url.Values{})
	if err != nil {
		t.Errorf("Expected GetGroup to succeed. Got %s", err.Error())
	} else if *g[0].Id != groupEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id) {
		t.Errorf("Expected GetGroup to return a valid group")
	}

	// Mock a valid	Search Result with expanded group members
	lm = &mocks.Client{}
	sr1 := &ldap.SearchRequest{
		BaseDN:     "ou=groups,dc=test",
		Scope:      2,
		Filter:     "(&(objectClass=groupOfNames))",
		Attributes: []string{"cn", "entryUUID", "member"},
		Controls:   []ldap.Control(nil),
	}
	sr2 := &ldap.SearchRequest{
		BaseDN:     "uid=user,ou=people,dc=test",
		SizeLimit:  1,
		Filter:     "(objectclass=*)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid"},
		Controls:   []ldap.Control(nil),
	}
	sr3 := &ldap.SearchRequest{
		BaseDN:     "uid=invalid,ou=people,dc=test",
		SizeLimit:  1,
		Filter:     "(objectclass=*)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid"},
		Controls:   []ldap.Control(nil),
	}

	for _, param := range []url.Values{queryParamSelect, queryParamExpand} {
		lm.On("Search", sr1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{groupEntry}}, nil)
		lm.On("Search", sr2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{userEntry}}, nil)
		lm.On("Search", sr3).Return(&ldap.SearchResult{Entries: []*ldap.Entry{invalidUserEntry}}, nil)
		b, _ = getMockedBackend(lm, lconfig, &logger)
		g, err = b.GetGroups(context.Background(), param)
		switch {
		case err != nil:
			t.Errorf("Expected GetGroup to succeed. Got %s", err.Error())
		case g[0].GetId() != groupEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id):
			t.Errorf("Expected GetGroup to return a valid group")
		case len(g[0].Members) != 1:
			t.Errorf("Expected GetGroup to return group with one member")
		case g[0].Members[0].GetId() != userEntry.GetEqualFoldAttributeValue(b.userAttributeMap.id):
			t.Errorf("Expected GetGroup to return group with correct member")
		}
	}
}

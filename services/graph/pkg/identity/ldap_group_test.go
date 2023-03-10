package identity

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/test-go/testify/mock"
)

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
		Filter:     "(objectClass=inetOrgPerson)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid", "sn", "givenname", "userEnabledAttribute", "userTypeAttribute"},
		Controls:   []ldap.Control(nil),
	}
	sr3 := &ldap.SearchRequest{
		BaseDN:     "uid=invalid,ou=people,dc=test",
		SizeLimit:  1,
		Filter:     "(objectClass=inetOrgPerson)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid", "sn", "givenname", "userEnabledAttribute", "userTypeAttribute"},
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
		Filter:     "(objectClass=inetOrgPerson)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid", "sn", "givenname", "userEnabledAttribute", "userTypeAttribute"},
		Controls:   []ldap.Control(nil),
	}
	sr3 := &ldap.SearchRequest{
		BaseDN:     "uid=invalid,ou=people,dc=test",
		SizeLimit:  1,
		Filter:     "(objectClass=inetOrgPerson)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid", "sn", "givenname", "userEnabledAttribute", "userTypeAttribute"},
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

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
	UserBaseDN:               "dc=test",
	UserSearchScope:          "sub",
	UserFilter:               "filter",
	UserDisplayNameAttribute: "displayname",
	UserIDAttribute:          "entryUUID",
	UserEmailAttribute:       "mail",
	UserNameAttribute:        "uid",

	GroupBaseDN:        "dc=test",
	GroupSearchScope:   "sub",
	GroupFilter:        "filter",
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
var groupEntry = ldap.NewEntry("cn=group",
	map[string][]string{
		"cn":        {"group"},
		"entryuuid": {"abcd-defg"},
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
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
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
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
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

	// Mock a valid	Search Result
	lm = &mocks.Client{}
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
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
}

func TestGetUsers(t *testing.T) {
	// Mock a Sizelimit Error
	lm := &mocks.Client{}
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, ldap.NewError(ldap.LDAPResultOperationsError, errors.New("mock")))

	b, _ := getMockedBackend(lm, lconfig, &logger)
	_, err := b.GetUsers(context.Background(), url.Values{})
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	lm = &mocks.Client{}
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&ldap.SearchResult{}, nil)
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
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, ldap.NewError(ldap.LDAPResultSizeLimitExceeded, errors.New("mock")))

	queryParamExpand := url.Values{
		"$expand": []string{"memberOf"},
	}
	queryParamSelect := url.Values{
		"$select": []string{"memberOf"},
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
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&ldap.SearchResult{}, nil)
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

	// Mock a valid	Search Result
	lm = &mocks.Client{}
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&ldap.SearchResult{
			Entries: []*ldap.Entry{groupEntry},
		}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err := b.GetGroup(context.Background(), "group", nil)
	if err != nil {
		t.Errorf("Expected GetGroup to succeed. Got %s", err.Error())
	} else if *g.Id != groupEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id) {
		t.Errorf("Expected GetGroup to return a valid group")
	}
	g, err = b.GetGroup(context.Background(), "group", queryParamExpand)
	if err != nil {
		t.Errorf("Expected GetGroup to succeed. Got %s", err.Error())
	} else if *g.Id != groupEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id) {
		t.Errorf("Expected GetGroup to return a valid group")
	}
	g, err = b.GetGroup(context.Background(), "group", queryParamSelect)
	if err != nil {
		t.Errorf("Expected GetGroup to succeed. Got %s", err.Error())
	} else if *g.Id != groupEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id) {
		t.Errorf("Expected GetGroup to return a valid group")
	}
}

func TestGetGroups(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, ldap.NewError(ldap.LDAPResultOperationsError, errors.New("mock")))
	b, _ := getMockedBackend(lm, lconfig, &logger)
	_, err := b.GetGroups(context.Background(), url.Values{})
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	lm = &mocks.Client{}
	lm.On("Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&ldap.SearchResult{}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err := b.GetGroups(context.Background(), url.Values{})
	if err != nil {
		t.Errorf("Expected success, got '%s'", err.Error())
	} else if g == nil || len(g) != 0 {
		t.Errorf("Expected zero length user slice")
	}
}

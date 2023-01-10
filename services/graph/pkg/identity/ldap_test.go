package identity

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/mock"
)

func getMockedBackend(l ldap.Client, lc config.LDAP, logger *log.Logger) (*LDAP, error) {
	return NewLDAPBackend(l, lc, logger)
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

	WriteEnabled: true,
}

var userEntry = ldap.NewEntry("uid=user",
	map[string][]string{
		"uid":         {"user"},
		"displayname": {"DisplayName"},
		"mail":        {"user@example"},
		"entryuuid":   {"abcd-defg"},
		"sn":          {"surname"},
		"givenname":   {"givenName"},
	})

var invalidUserEntry = ldap.NewEntry("uid=user",
	map[string][]string{
		"uid":         {"invalid"},
		"displayname": {"DisplayName"},
		"mail":        {"user@example"},
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

func TestCreateUser(t *testing.T) {
	l := &mocks.Client{}
	l.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{
				Entries: []*ldap.Entry{userEntry},
			},
			nil)
	l.On("Add", mock.Anything).Return(nil)
	logger := log.NewLogger(log.Level("debug"))

	displayName := "DisplayName"
	mail := "user@example"
	userName := "user"
	surname := "surname"
	givenName := "givenName"

	user := libregraph.NewUser()
	user.SetDisplayName(displayName)
	user.SetMail(mail)
	user.SetOnPremisesSamAccountName(userName)
	user.SetSurname(surname)
	user.SetGivenName(givenName)

	c := lconfig
	c.UseServerUUID = true
	b, _ := NewLDAPBackend(l, c, &logger)

	newUser, err := b.CreateUser(context.Background(), *user)
	assert.Nil(t, err)
	assert.Equal(t, displayName, newUser.GetDisplayName())
	assert.Equal(t, mail, newUser.GetMail())
	assert.Equal(t, userName, newUser.GetOnPremisesSamAccountName())
	assert.Equal(t, givenName, newUser.GetGivenName())
	assert.Equal(t, surname, newUser.GetSurname())
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

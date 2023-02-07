package identity

import (
	"context"
	"testing"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/mock"
)

var eduUserEntry = ldap.NewEntry("uid=user,ou=people,dc=test",
	map[string][]string{
		"uid":         {"testuser"},
		"displayname": {"Test User"},
		"mail":        {"user@example"},
		"entryuuid":   {"abcd-defg"},
		"userClass":   {"student"},
		"oCExternalIdentity": {
			"$ http://idp $ testuser",
			"xxx $ http://idpnew $ xxxxx-xxxxx-xxxxx",
		},
	})
var eduUserEntryWithSchool = ldap.NewEntry("uid=user,ou=people,dc=test",
	map[string][]string{
		"uid":              {"testuser"},
		"displayname":      {"Test User"},
		"mail":             {"user@example"},
		"entryuuid":        {"abcd-defg"},
		"userClass":        {"student"},
		"ocMemberOfSchool": {"abcd-defg"},
		"oCExternalIdentity": {
			"$ http://idp $ testuser",
			"xxx $ http://idpnew $ xxxxx-xxxxx-xxxxx",
		},
	})

var sr1 *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=people,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=ocEducationUser)(|(uid=abcd-defg)(entryUUID=abcd-defg)))",
	Attributes: []string{"displayname", "entryUUID", "mail", "uid", "oCExternalIdentity", "userClass", "ocMemberOfSchool"},
	Controls:   []ldap.Control(nil),
}
var sr2 *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=people,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=ocEducationUser)(|(uid=xxxx-xxxx)(entryUUID=xxxx-xxxx)))",
	Attributes: []string{"displayname", "entryUUID", "mail", "uid", "oCExternalIdentity", "userClass", "ocMemberOfSchool"},
	Controls:   []ldap.Control(nil),
}

func TestCreateEducationUser(t *testing.T) {
	lm := &mocks.Client{}
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	//assert.NotEqual(t, "", b.educationConfig.schoolObjectClass)
	lm.On("Add", mock.Anything).Return(nil)

	lm.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{
				Entries: []*ldap.Entry{
					eduUserEntry,
				},
			},
			nil)
	user := libregraph.NewEducationUser()
	user.SetDisplayName("Test User")
	user.SetOnPremisesSamAccountName("testuser")
	user.SetMail("testuser@example.org")
	user.SetPrimaryRole("student")
	eduUser, err := b.CreateEducationUser(context.Background(), *user)
	lm.AssertNumberOfCalls(t, "Add", 1)
	lm.AssertNumberOfCalls(t, "Search", 1)
	assert.NotNil(t, eduUser)
	assert.Nil(t, err)
	assert.Equal(t, eduUser.GetDisplayName(), user.GetDisplayName())
	assert.Equal(t, eduUser.GetOnPremisesSamAccountName(), user.GetOnPremisesSamAccountName())
	assert.Equal(t, "abcd-defg", eduUser.GetId())
	assert.Equal(t, eduUser.GetPrimaryRole(), user.GetPrimaryRole())
}

func TestDeleteEducationUser(t *testing.T) {
	lm := &mocks.Client{}

	lm.On("Search", sr1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntry}}, nil)
	lm.On("Search", sr2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	dr1 := &ldap.DelRequest{
		DN: "uid=user,ou=people,dc=test",
	}
	lm.On("Del", dr1).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.DeleteEducationUser(context.Background(), "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 1)
	lm.AssertNumberOfCalls(t, "Del", 1)
	assert.Nil(t, err)

	err = b.DeleteEducationUser(context.Background(), "xxxx-xxxx")
	lm.AssertNumberOfCalls(t, "Search", 2)
	lm.AssertNumberOfCalls(t, "Del", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "itemNotFound", err.Error())
}

func TestGetEducationUser(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", sr1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntry}}, nil)
	lm.On("Search", sr2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	user, err := b.GetEducationUser(context.Background(), "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 1)
	assert.Nil(t, err)
	assert.Equal(t, "Test User", user.GetDisplayName())
	assert.Equal(t, "abcd-defg", user.GetId())

	_, err = b.GetEducationUser(context.Background(), "xxxx-xxxx")
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	assert.Equal(t, "itemNotFound", err.Error())
}

func TestGetEducationUsers(t *testing.T) {
	lm := &mocks.Client{}
	sr := &ldap.SearchRequest{
		BaseDN:     "ou=people,dc=test",
		Scope:      2,
		SizeLimit:  0,
		Filter:     "(objectClass=ocEducationUser)",
		Attributes: []string{"displayname", "entryUUID", "mail", "uid", "oCExternalIdentity", "userClass", "ocMemberOfSchool"},
		Controls:   []ldap.Control(nil),
	}
	lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntry}}, nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	_, err = b.GetEducationUsers(context.Background())
	lm.AssertNumberOfCalls(t, "Search", 1)
	assert.Nil(t, err)
}

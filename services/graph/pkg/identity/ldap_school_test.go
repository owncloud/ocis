package identity

import (
	"context"
	"testing"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/mock"
)

var eduConfig = config.LDAP{
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

	WriteEnabled:              true,
	EducationResourcesEnabled: true,
}

var schoolEntry = ldap.NewEntry("ou=Test School",
	map[string][]string{
		"ou":                      {"Test School"},
		"ocEducationSchoolNumber": {"0123"},
		"owncloudUUID":            {"abcd-defg"},
	})

func TestCreateSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Add", mock.Anything).
		Return(nil)

	lm.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{
				Entries: []*ldap.Entry{schoolEntry},
			},
			nil)

	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	assert.NotEqual(t, "", b.educationConfig.schoolObjectClass)
	school := libregraph.NewEducationSchool()
	school.SetDisplayName("Test School")
	school.SetSchoolNumber("0123")
	school.SetId("abcd-defg")
	res_school, err := b.CreateSchool(context.Background(), *school)
	lm.AssertNumberOfCalls(t, "Add", 1)
	lm.AssertNumberOfCalls(t, "Search", 1)
	assert.Nil(t, err)
	assert.NotNil(t, res_school)
	assert.Equal(t, res_school.GetDisplayName(), school.GetDisplayName())
	assert.Equal(t, res_school.GetId(), school.GetId())
	assert.Equal(t, res_school.GetSchoolNumber(), school.GetSchoolNumber())
}

func TestDeleteSchool(t *testing.T) {
	lm := &mocks.Client{}
	sr1 := &ldap.SearchRequest{
		BaseDN:     "",
		Scope:      2,
		SizeLimit:  1,
		Filter:     "(&(objectClass=ocEducationSchool)(owncloudUUID=abcd-defg))",
		Attributes: []string{"ou", "owncloudUUID", "ocEducationSchoolNumber"},
		Controls:   []ldap.Control(nil),
	}
	sr2 := &ldap.SearchRequest{
		BaseDN:     "",
		Scope:      2,
		SizeLimit:  1,
		Filter:     "(&(objectClass=ocEducationSchool)(owncloudUUID=xxxx-xxxx))",
		Attributes: []string{"ou", "owncloudUUID", "ocEducationSchoolNumber"},
		Controls:   []ldap.Control(nil),
	}
	lm.On("Search", sr1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", sr2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	dr1 := &ldap.DelRequest{
		DN: "ou=Test School",
	}
	lm.On("Del", dr1).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)

	err = b.DeleteSchool(context.Background(), "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 1)
	lm.AssertNumberOfCalls(t, "Del", 1)
	assert.Nil(t, err)

	err = b.DeleteSchool(context.Background(), "xxxx-xxxx")
	lm.AssertNumberOfCalls(t, "Search", 2)
	lm.AssertNumberOfCalls(t, "Del", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "itemNotFound", err.Error())
}

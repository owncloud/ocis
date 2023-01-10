package identity

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/mock"
)

var classEntry = ldap.NewEntry("ocEducationExternalId=Math0123",
	map[string][]string{
		"cn":                    {"Math"},
		"ocEducationExternalId": {"Math0123"},
		"ocEducationClassType":  {"course"},
		"entryUUID":             {"abcd-defg"},
	})
var classEntryWithMember = ldap.NewEntry("ocEducationExternalId=Math0123",
	map[string][]string{
		"cn":                    {"Math"},
		"ocEducationExternalId": {"Math0123"},
		"ocEducationClassType":  {"course"},
		"entryUUID":             {"abcd-defg"},
		"member":                {"uid=user"},
	})

func TestCreateEducationClass(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Add", mock.Anything).
		Return(nil)

	lm.On("Search", mock.Anything).
		Return(
			&ldap.SearchResult{
				Entries: []*ldap.Entry{classEntry},
			},
			nil)

	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	assert.NotEqual(t, "", b.educationConfig.classObjectClass)
	class := libregraph.NewEducationClass("Math", "course")
	class.SetExternalId("Math0123")
	class.SetId("abcd-defg")
	res_class, err := b.CreateEducationClass(context.Background(), *class)
	lm.AssertNumberOfCalls(t, "Add", 1)
	lm.AssertNumberOfCalls(t, "Search", 1)
	assert.Nil(t, err)
	assert.NotNil(t, res_class)
	assert.Equal(t, res_class.GetDisplayName(), class.GetDisplayName())
	assert.Equal(t, res_class.GetId(), class.GetId())
	assert.Equal(t, res_class.GetExternalId(), class.GetExternalId())
	assert.Equal(t, res_class.GetClassification(), class.GetClassification())
}

func TestGetEducationClasses(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", mock.Anything).Return(nil, ldap.NewError(ldap.LDAPResultOperationsError, errors.New("mock")))
	b, _ := getMockedBackend(lm, lconfig, &logger)
	_, err := b.GetEducationClasses(context.Background(), url.Values{})
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).Return(&ldap.SearchResult{}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err := b.GetEducationClasses(context.Background(), url.Values{})
	if err != nil {
		t.Errorf("Expected success, got '%s'", err.Error())
	} else if g == nil || len(g) != 0 {
		t.Errorf("Expected zero length user slice")
	}

	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).Return(&ldap.SearchResult{
		Entries: []*ldap.Entry{classEntry},
	}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err = b.GetEducationClasses(context.Background(), url.Values{})
	if err != nil {
		t.Errorf("Expected GetEducationClasses to succeed. Got %s", err.Error())
	} else if *g[0].Id != classEntry.GetEqualFoldAttributeValue(b.groupAttributeMap.id) {
		t.Errorf("Expected GetEducationClasses to return a valid group")
	}
}

func TestGetEducationClass(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		filter               string
		expectedItemNotFound bool
	}{
		{
			name:                 "Test search class using id",
			id:                   "abcd-defg",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=abcd-defg)(ocEducationExternalId=abcd-defg)))",
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search class using unknown Id",
			id:                   "xxxx-xxxx",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=xxxx-xxxx)(ocEducationExternalId=xxxx-xxxx)))",
			expectedItemNotFound: true,
		},
		{
			name:                 "Test search class using external ID",
			id:                   "Math0123",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=Math0123)(ocEducationExternalId=Math0123)))",
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search school using unknown externalID",
			id:                   "Unknown3210",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=Unknown3210)(ocEducationExternalId=Unknown3210)))",
			expectedItemNotFound: true,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		sr := &ldap.SearchRequest{
			BaseDN:     "ou=groups,dc=test",
			Scope:      2,
			SizeLimit:  1,
			Filter:     tt.filter,
			Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId"},
			Controls:   []ldap.Control(nil),
		}
		if tt.expectedItemNotFound {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
		} else {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntry}}, nil)
		}

		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		class, err := b.GetEducationClass(context.Background(), tt.id, nil)
		lm.AssertNumberOfCalls(t, "Search", 1)

		if tt.expectedItemNotFound {
			assert.NotNil(t, err)
			assert.Equal(t, "itemNotFound", err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, "Math", class.GetDisplayName())
			assert.Equal(t, "abcd-defg", class.GetId())
			assert.Equal(t, "Math0123", class.GetExternalId())
		}
	}
}

func TestDeleteEducationClass(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		filter               string
		expectedItemNotFound bool
	}{
		{
			name:                 "Test search class using id",
			id:                   "abcd-defg",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=abcd-defg)(ocEducationExternalId=abcd-defg)))",
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search class using unknown Id",
			id:                   "xxxx-xxxx",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=xxxx-xxxx)(ocEducationExternalId=xxxx-xxxx)))",
			expectedItemNotFound: true,
		},
		{
			name:                 "Test search class using external ID",
			id:                   "Math0123",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=Math0123)(ocEducationExternalId=Math0123)))",
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search school using unknown externalID",
			id:                   "Unknown3210",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=Unknown3210)(ocEducationExternalId=Unknown3210)))",
			expectedItemNotFound: true,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		sr := &ldap.SearchRequest{
			BaseDN:     "ou=groups,dc=test",
			Scope:      2,
			SizeLimit:  1,
			Filter:     tt.filter,
			Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId"},
			Controls:   []ldap.Control(nil),
		}
		if tt.expectedItemNotFound {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
		} else {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntry}}, nil)
		}
		dr := &ldap.DelRequest{
			DN: "ocEducationExternalId=Math0123",
		}
		lm.On("Del", dr).Return(nil)

		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		err = b.DeleteEducationClass(context.Background(), tt.id)
		lm.AssertNumberOfCalls(t, "Search", 1)

		if tt.expectedItemNotFound {
			lm.AssertNumberOfCalls(t, "Del", 0)
			assert.NotNil(t, err)
			assert.Equal(t, "itemNotFound", err.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetEducationClassMembers(t *testing.T) {
	tests := []struct {
		name                 string
		id                   string
		filter               string
		expectedItemNotFound bool
	}{
		{
			name:                 "Test search class using id",
			id:                   "abcd-defg",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=abcd-defg)(ocEducationExternalId=abcd-defg)))",
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search class using unknown Id",
			id:                   "xxxx-xxxx",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=xxxx-xxxx)(ocEducationExternalId=xxxx-xxxx)))",
			expectedItemNotFound: true,
		},
		{
			name:                 "Test search class using external ID",
			id:                   "Math0123",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=Math0123)(ocEducationExternalId=Math0123)))",
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search school using unknown externalID",
			id:                   "Unknown3210",
			filter:               "(&(objectClass=ocEducationClass)(|(entryUUID=Unknown3210)(ocEducationExternalId=Unknown3210)))",
			expectedItemNotFound: true,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		user_sr := &ldap.SearchRequest{
			BaseDN:     "uid=user",
			Scope:      0,
			SizeLimit:  1,
			Filter:     "(objectClass=inetOrgPerson)",
			Attributes: []string{"displayname", "entryUUID", "mail", "uid", "sn", "givenname"},
			Controls:   []ldap.Control(nil),
		}
		lm.On("Search", user_sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{userEntry}}, nil)
		sr := &ldap.SearchRequest{
			BaseDN:     "ou=groups,dc=test",
			Scope:      2,
			SizeLimit:  1,
			Filter:     tt.filter,
			Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId", "member"},
			Controls:   []ldap.Control(nil),
		}
		if tt.expectedItemNotFound {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
		} else {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntryWithMember}}, nil)
		}

		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		users, err := b.GetEducationClassMembers(context.Background(), tt.id)

		if tt.expectedItemNotFound {
			lm.AssertNumberOfCalls(t, "Search", 1)
			assert.NotNil(t, err)
			assert.Equal(t, "itemNotFound", err.Error())
		} else {
			lm.AssertNumberOfCalls(t, "Search", 2)
			assert.Nil(t, err)
			assert.Equal(t, len(users), 1)
		}
	}
}

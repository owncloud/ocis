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
	UserEnabledAttribute:     "userEnabledAttribute",
	UserTypeAttribute:        "userTypeAttribute",

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

var schoolEntry1 = ldap.NewEntry("ou=Test School1",
	map[string][]string{
		"ou":                      {"Test School1"},
		"ocEducationSchoolNumber": {"0042"},
		"owncloudUUID":            {"hijk-defg"},
	})

var (
	filterSchoolSearchByIdExisting        = "(&(objectClass=ocEducationSchool)(|(owncloudUUID=abcd-defg)(ocEducationSchoolNumber=abcd-defg)))"
	filterSchoolSearchByIdNonexistant     = "(&(objectClass=ocEducationSchool)(|(owncloudUUID=xxxx-xxxx)(ocEducationSchoolNumber=xxxx-xxxx)))"
	filterSchoolSearchByNumberExisting    = "(&(objectClass=ocEducationSchool)(|(owncloudUUID=0123)(ocEducationSchoolNumber=0123)))"
	filterSchoolSearchByNumberNonexistant = "(&(objectClass=ocEducationSchool)(|(owncloudUUID=3210)(ocEducationSchoolNumber=3210)))"
)

func TestCreateEducationSchool(t *testing.T) {
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
	res_school, err := b.CreateEducationSchool(context.Background(), *school)
	lm.AssertNumberOfCalls(t, "Add", 1)
	lm.AssertNumberOfCalls(t, "Search", 1)
	assert.Nil(t, err)
	assert.NotNil(t, res_school)
	assert.Equal(t, res_school.GetDisplayName(), school.GetDisplayName())
	assert.Equal(t, res_school.GetId(), school.GetId())
	assert.Equal(t, res_school.GetSchoolNumber(), school.GetSchoolNumber())
}

func TestUpdateEducationSchoolOperation(t *testing.T) {
	testSchoolName := "A name"
	testSchoolNumber := "1234"
	tests := []struct {
		name              string
		displayName       string
		schoolNumber      string
		expectedOperation schoolUpdateOperation
	}{
		{
			name:              "Test using school with both number and name",
			displayName:       testSchoolName,
			schoolNumber:      testSchoolNumber,
			expectedOperation: tooManyValues,
		},
		{
			name:              "Test with unchanged number",
			schoolNumber:      testSchoolNumber,
			expectedOperation: schoolUnchanged,
		},
		{
			name:              "Test with unchanged name",
			displayName:       testSchoolName,
			expectedOperation: schoolUnchanged,
		},
		{
			name:              "Test new name",
			displayName:       "Something new",
			expectedOperation: displayNameUpdated,
		},
		{
			name:              "Test new number",
			schoolNumber:      "9876",
			expectedOperation: schoolNumberUpdated,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		displayName := "A name"
		schoolNumber := "1234"

		currentSchool := libregraph.EducationSchool{
			DisplayName:  &displayName,
			SchoolNumber: &schoolNumber,
		}

		schoolUpdate := libregraph.EducationSchool{
			DisplayName:  &tt.displayName,
			SchoolNumber: &tt.schoolNumber,
		}

		operation := b.UpdateEducationSchoolOperation(schoolUpdate, currentSchool)
		assert.Equal(t, tt.expectedOperation, operation)
	}
}

func TestDeleteEducationSchool(t *testing.T) {
	tests := []struct {
		name                 string
		numberOrId           string
		filter               string
		expectedItemNotFound bool
	}{
		{
			name:                 "Test delete school using schoolId",
			numberOrId:           "abcd-defg",
			filter:               filterSchoolSearchByIdExisting,
			expectedItemNotFound: false,
		},
		{
			name:                 "Test delete school using unknown schoolId",
			numberOrId:           "xxxx-xxxx",
			filter:               filterSchoolSearchByIdNonexistant,
			expectedItemNotFound: true,
		},
		{
			name:                 "Test delete school using schoolNumber",
			numberOrId:           "0123",
			filter:               filterSchoolSearchByNumberExisting,
			expectedItemNotFound: false,
		},
		{
			name:                 "Test delete school using unknown schoolNumber",
			numberOrId:           "3210",
			filter:               filterSchoolSearchByNumberNonexistant,
			expectedItemNotFound: true,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		sr := &ldap.SearchRequest{
			BaseDN:     "",
			Scope:      2,
			SizeLimit:  1,
			Filter:     tt.filter,
			Attributes: []string{"ou", "owncloudUUID", "ocEducationSchoolNumber"},
			Controls:   []ldap.Control(nil),
		}
		if tt.expectedItemNotFound {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
		} else {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
		}
		dr := &ldap.DelRequest{
			DN: "ou=Test School",
		}
		lm.On("Del", dr).Return(nil)

		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		err = b.DeleteEducationSchool(context.Background(), tt.numberOrId)
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

func TestGetEducationSchool(t *testing.T) {
	tests := []struct {
		name                 string
		numberOrId           string
		filter               string
		expectedItemNotFound bool
	}{
		{
			name:                 "Test search school using schoolId",
			numberOrId:           "abcd-defg",
			filter:               filterSchoolSearchByIdExisting,
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search school using unknown schoolId",
			numberOrId:           "xxxx-xxxx",
			filter:               filterSchoolSearchByIdNonexistant,
			expectedItemNotFound: true,
		},
		{
			name:                 "Test search school using schoolNumber",
			numberOrId:           "0123",
			filter:               filterSchoolSearchByNumberExisting,
			expectedItemNotFound: false,
		},
		{
			name:                 "Test search school using unknown schoolNumber",
			numberOrId:           "3210",
			filter:               filterSchoolSearchByNumberNonexistant,
			expectedItemNotFound: true,
		},
	}

	for _, tt := range tests {
		lm := &mocks.Client{}
		sr := &ldap.SearchRequest{
			BaseDN:     "",
			Scope:      2,
			SizeLimit:  1,
			Filter:     tt.filter,
			Attributes: []string{"ou", "owncloudUUID", "ocEducationSchoolNumber"},
			Controls:   []ldap.Control(nil),
		}
		if tt.expectedItemNotFound {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
		} else {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
		}

		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		school, err := b.GetEducationSchool(context.Background(), tt.numberOrId)
		lm.AssertNumberOfCalls(t, "Search", 1)

		if tt.expectedItemNotFound {
			assert.NotNil(t, err)
			assert.Equal(t, "itemNotFound", err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, "Test School", school.GetDisplayName())
			assert.Equal(t, "abcd-defg", school.GetId())
			assert.Equal(t, "0123", school.GetSchoolNumber())
		}
	}
}

func TestGetEducationSchools(t *testing.T) {
	lm := &mocks.Client{}
	sr1 := &ldap.SearchRequest{
		BaseDN:     "",
		Scope:      2,
		SizeLimit:  0,
		Filter:     "(objectClass=ocEducationSchool)",
		Attributes: []string{"ou", "owncloudUUID", "ocEducationSchoolNumber"},
		Controls:   []ldap.Control(nil),
	}
	lm.On("Search", sr1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry, schoolEntry1}}, nil)
	//	lm.On("Search", sr2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	_, err = b.GetEducationSchools(context.Background())
	lm.AssertNumberOfCalls(t, "Search", 1)
	assert.Nil(t, err)
}

var schoolByIDSearch1 *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "",
	Scope:      2,
	SizeLimit:  1,
	Filter:     filterSchoolSearchByIdExisting,
	Attributes: []string{"ou", "owncloudUUID", "ocEducationSchoolNumber"},
	Controls:   []ldap.Control(nil),
}

var schoolByNumberSearch *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "",
	Scope:      2,
	SizeLimit:  1,
	Filter:     filterSchoolSearchByNumberExisting,
	Attributes: []string{"ou", "owncloudUUID", "ocEducationSchoolNumber"},
	Controls:   []ldap.Control(nil),
}

var userByIDSearch1 *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=people,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=ocEducationUser)(|(uid=abcd-defg)(entryUUID=abcd-defg)))",
	Attributes: eduUserAttrs,
	Controls:   []ldap.Control(nil),
}

var userByIDSearch2 *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=people,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=ocEducationUser)(|(uid=does-not-exist)(entryUUID=does-not-exist)))",
	Attributes: eduUserAttrs,
	Controls:   []ldap.Control(nil),
}

var userToSchoolModRequest *ldap.ModifyRequest = &ldap.ModifyRequest{
	DN: "uid=user,ou=people,dc=test",
	Changes: []ldap.Change{
		{
			Operation: ldap.AddAttribute,
			Modification: ldap.PartialAttribute{
				Type: "ocMemberOfSchool",
				Vals: []string{"abcd-defg"},
			},
		},
	},
}

var userFromSchoolModRequest *ldap.ModifyRequest = &ldap.ModifyRequest{
	DN: "uid=user,ou=people,dc=test",
	Changes: []ldap.Change{
		{
			Operation: ldap.DeleteAttribute,
			Modification: ldap.PartialAttribute{
				Type: "ocMemberOfSchool",
				Vals: []string{"abcd-defg"},
			},
		},
	},
}

var classToSchoolModRequest *ldap.ModifyRequest = &ldap.ModifyRequest{
	DN: "ocEducationExternalId=Math0123",
	Changes: []ldap.Change{
		{
			Operation: ldap.AddAttribute,
			Modification: ldap.PartialAttribute{
				Type: "ocMemberOfSchool",
				Vals: []string{"abcd-defg"},
			},
		},
	},
}

var classFromSchoolModRequest *ldap.ModifyRequest = &ldap.ModifyRequest{
	DN: "ocEducationExternalId=Math0123",
	Changes: []ldap.Change{
		{
			Operation: ldap.DeleteAttribute,
			Modification: ldap.PartialAttribute{
				Type: "ocMemberOfSchool",
				Vals: []string{"abcd-defg"},
			},
		},
	},
}

func TestAddUsersToEducationSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", userByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntry}}, nil)
	lm.On("Search", userByIDSearch2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	lm.On("Modify", userToSchoolModRequest).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.AddUsersToEducationSchool(context.Background(), "abcd-defg", []string{"does-not-exist"})
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	err = b.AddUsersToEducationSchool(context.Background(), "abcd-defg", []string{"abcd-defg", "does-not-exist"})
	lm.AssertNumberOfCalls(t, "Search", 5)
	assert.NotNil(t, err)
	err = b.AddUsersToEducationSchool(context.Background(), "abcd-defg", []string{"abcd-defg"})
	lm.AssertNumberOfCalls(t, "Search", 7)
	assert.Nil(t, err)
	// try to add by school number (instead or id)
	err = b.AddUsersToEducationSchool(context.Background(), "0123", []string{"abcd-defg"})
	lm.AssertNumberOfCalls(t, "Search", 9)
	assert.Nil(t, err)
}

func TestRemoveMemberFromEducationSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", userByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntryWithSchool}}, nil)
	lm.On("Search", userByIDSearch2).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	lm.On("Modify", userFromSchoolModRequest).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.RemoveUserFromEducationSchool(context.Background(), "abcd-defg", "does-not-exist")
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	assert.Equal(t, "itemNotFound", err.Error())
	err = b.RemoveUserFromEducationSchool(context.Background(), "abcd-defg", "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 4)
	lm.AssertNumberOfCalls(t, "Modify", 1)
	// try to remove by school number (instead or id)
	err = b.RemoveUserFromEducationSchool(context.Background(), "0123", "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 6)
	lm.AssertNumberOfCalls(t, "Modify", 2)
	assert.Nil(t, err)
}

var usersBySchoolIDSearch *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=people,dc=test",
	Scope:      2,
	SizeLimit:  0,
	Filter:     "(&(objectClass=ocEducationUser)(ocMemberOfSchool=abcd-defg))",
	Attributes: eduUserAttrs,
	Controls:   []ldap.Control(nil),
}

func TestGetEducationSchoolUsers(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", usersBySchoolIDSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{eduUserEntryWithSchool}}, nil)
	b, _ := getMockedBackend(lm, eduConfig, &logger)
	users, err := b.GetEducationSchoolUsers(context.Background(), "abcd-defg")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	users, err = b.GetEducationSchoolUsers(context.Background(), "0123")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

var classesBySchoolIDSearch *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=groups,dc=test",
	Scope:      2,
	SizeLimit:  0,
	Filter:     "(&(objectClass=ocEducationClass)(ocMemberOfSchool=abcd-defg))",
	Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId", "ocMemberOfSchool", "ocEducationTeacherMember"},
	Controls:   []ldap.Control(nil),
}

func TestGetEducationSchoolClasses(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", classesBySchoolIDSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntry}}, nil)
	b, _ := getMockedBackend(lm, eduConfig, &logger)
	users, err := b.GetEducationSchoolClasses(context.Background(), "abcd-defg")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	users, err = b.GetEducationSchoolClasses(context.Background(), "0123")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

var classesByUUIDSearchNotFound *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=groups,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=ocEducationClass)(|(entryUUID=does-not-exist)(ocEducationExternalId=does-not-exist)))",
	Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId", "ocMemberOfSchool", "ocEducationTeacherMember"},
	Controls:   []ldap.Control(nil),
}

var classesByUUIDSearchFound *ldap.SearchRequest = &ldap.SearchRequest{
	BaseDN:     "ou=groups,dc=test",
	Scope:      2,
	SizeLimit:  1,
	Filter:     "(&(objectClass=ocEducationClass)(|(entryUUID=abcd-defg)(ocEducationExternalId=abcd-defg)))",
	Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId", "ocMemberOfSchool", "ocEducationTeacherMember"},
	Controls:   []ldap.Control(nil),
}

func TestAddClassesToEducationSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", classesByUUIDSearchNotFound).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	lm.On("Search", classesByUUIDSearchFound).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Modify", classToSchoolModRequest).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.AddClassesToEducationSchool(context.Background(), "abcd-defg", []string{"does-not-exist"})
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	err = b.AddClassesToEducationSchool(context.Background(), "abcd-defg", []string{"abcd-defg", "does-not-exist"})
	lm.AssertNumberOfCalls(t, "Search", 5)
	assert.NotNil(t, err)
	err = b.AddClassesToEducationSchool(context.Background(), "abcd-defg", []string{"abcd-defg"})
	lm.AssertNumberOfCalls(t, "Search", 7)
	assert.Nil(t, err)
	// try to add by school number (instead or id)
	err = b.AddClassesToEducationSchool(context.Background(), "0123", []string{"abcd-defg"})
	lm.AssertNumberOfCalls(t, "Search", 9)
	assert.Nil(t, err)
}

func TestRemoveClassFromEducationSchool(t *testing.T) {
	lm := &mocks.Client{}
	lm.On("Search", schoolByIDSearch1).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", schoolByNumberSearch).Return(&ldap.SearchResult{Entries: []*ldap.Entry{schoolEntry}}, nil)
	lm.On("Search", classesByUUIDSearchFound).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntryWithSchool}}, nil)
	lm.On("Search", classesByUUIDSearchNotFound).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	lm.On("Modify", classFromSchoolModRequest).Return(nil)
	b, err := getMockedBackend(lm, eduConfig, &logger)
	assert.Nil(t, err)
	err = b.RemoveClassFromEducationSchool(context.Background(), "abcd-defg", "does-not-exist")
	lm.AssertNumberOfCalls(t, "Search", 2)
	assert.NotNil(t, err)
	assert.Equal(t, "itemNotFound", err.Error())
	err = b.RemoveClassFromEducationSchool(context.Background(), "abcd-defg", "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 4)
	lm.AssertNumberOfCalls(t, "Modify", 1)
	// try to remove by school number (instead or id)
	err = b.RemoveClassFromEducationSchool(context.Background(), "0123", "abcd-defg")
	lm.AssertNumberOfCalls(t, "Search", 6)
	lm.AssertNumberOfCalls(t, "Modify", 2)
	assert.Nil(t, err)
}

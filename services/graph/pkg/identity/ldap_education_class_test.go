package identity

import (
	"context"
	"errors"
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

var classEntryWithSchool = ldap.NewEntry("ocEducationExternalId=Math0123",
	map[string][]string{
		"cn":                    {"Math"},
		"ocEducationExternalId": {"Math0123"},
		"ocEducationClassType":  {"course"},
		"entryUUID":             {"abcd-defg"},
		"ocMemberOfSchool":      {"abcd-defg"},
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
	_, err := b.GetEducationClasses(context.Background())
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	lm = &mocks.Client{}
	lm.On("Search", mock.Anything).Return(&ldap.SearchResult{}, nil)
	b, _ = getMockedBackend(lm, lconfig, &logger)
	g, err := b.GetEducationClasses(context.Background())
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
	g, err = b.GetEducationClasses(context.Background())
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
			Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId", "ocMemberOfSchool", "ocEducationTeacherMember"},
			Controls:   []ldap.Control(nil),
		}
		if tt.expectedItemNotFound {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
		} else {
			lm.On("Search", sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{classEntry}}, nil)
		}

		b, err := getMockedBackend(lm, eduConfig, &logger)
		assert.Nil(t, err)

		class, err := b.GetEducationClass(context.Background(), tt.id)
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
			Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId", "ocMemberOfSchool", "ocEducationTeacherMember"},
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
			Attributes: []string{"displayname", "entryUUID", "mail", "uid", "sn", "givenname", "userEnabledAttribute", "userTypeAttribute"},
			Controls:   []ldap.Control(nil),
		}
		lm.On("Search", user_sr).Return(&ldap.SearchResult{Entries: []*ldap.Entry{userEntry}}, nil)
		sr := &ldap.SearchRequest{
			BaseDN:     "ou=groups,dc=test",
			Scope:      2,
			SizeLimit:  1,
			Filter:     tt.filter,
			Attributes: []string{"cn", "entryUUID", "ocEducationClassType", "ocEducationExternalId", "ocMemberOfSchool", "ocEducationTeacherMember", "member"},
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

func TestLDAP_UpdateEducationClass(t *testing.T) {
	externalIDs := []string{"Math3210"}
	changeString := "xxxx-xxxx"
	type args struct {
		id    string
		class libregraph.EducationClass
	}
	type modifyData struct {
		arg *ldap.ModifyRequest
		ret error
	}
	type modifyDNData struct {
		arg *ldap.ModifyDNRequest
		ret error
	}
	type searchData struct {
		res *ldap.SearchResult
		err error
	}
	tests := []struct {
		name         string
		args         args
		modifyDNData modifyDNData
		modifyData   modifyData
		searchData   searchData
		assertion    func(assert.TestingT, error, ...interface{}) bool
	}{
		{
			name: "Change name",
			args: args{
				id: "abcd-defg",
				class: libregraph.EducationClass{
					DisplayName: "Math-2",
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool { return assert.Nil(tt, err) },
			modifyData: modifyData{
				arg: &ldap.ModifyRequest{
					DN: "ocEducationExternalId=Math0123",
					Changes: []ldap.Change{
						{
							Operation: ldap.ReplaceAttribute,
							Modification: ldap.PartialAttribute{
								Type: "cn",
								Vals: []string{"Math-2"},
							},
						},
					},
				},
			},
			modifyDNData: modifyDNData{
				arg: &ldap.ModifyDNRequest{},
				ret: nil,
			},
			searchData: searchData{
				res: &ldap.SearchResult{
					Entries: []*ldap.Entry{classEntry},
				},
				err: nil,
			},
		},
		{
			name: "Change external ID",
			args: args{
				id: "abcd-defg",
				class: libregraph.EducationClass{
					ExternalId: &externalIDs[0],
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool { return assert.Nil(tt, err) },
			modifyData: modifyData{
				arg: &ldap.ModifyRequest{},
			},
			modifyDNData: modifyDNData{
				arg: &ldap.ModifyDNRequest{
					DN:           "ocEducationExternalId=Math0123",
					NewRDN:       "ocEducationExternalId=Math3210",
					DeleteOldRDN: true,
					NewSuperior:  "",
				},
				ret: nil,
			},
			searchData: searchData{
				res: &ldap.SearchResult{
					Entries: []*ldap.Entry{classEntry},
				},
				err: nil,
			},
		},
		{
			name: "Change both name and external ID",
			args: args{
				id: "abcd-defg",
				class: libregraph.EducationClass{
					DisplayName: "Math-2",
					ExternalId:  &externalIDs[0],
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool { return assert.Nil(tt, err) },
			modifyData: modifyData{
				arg: &ldap.ModifyRequest{
					DN: "ocEducationExternalId=Math3210,ou=groups,dc=test",
					Changes: []ldap.Change{
						{
							Operation: ldap.ReplaceAttribute,
							Modification: ldap.PartialAttribute{
								Type: "cn",
								Vals: []string{"Math-2"},
							},
						},
					},
				},
			},
			modifyDNData: modifyDNData{
				arg: &ldap.ModifyDNRequest{
					DN:           "ocEducationExternalId=Math0123",
					NewRDN:       "ocEducationExternalId=Math3210",
					DeleteOldRDN: true,
					NewSuperior:  "",
				},
				ret: nil,
			},
			searchData: searchData{
				res: &ldap.SearchResult{
					Entries: []*ldap.Entry{classEntry},
				},
				err: nil,
			},
		},
		{
			name: "Check error: attempt at changing ID",
			args: args{
				id: "abcd-defg",
				class: libregraph.EducationClass{
					Id: &changeString,
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool { return assert.Error(tt, err) },
			modifyData: modifyData{
				arg: &ldap.ModifyRequest{},
			},
			modifyDNData: modifyDNData{
				arg: &ldap.ModifyDNRequest{},
				ret: nil,
			},
			searchData: searchData{
				res: &ldap.SearchResult{
					Entries: []*ldap.Entry{classEntry},
				},
				err: nil,
			},
		},
		{
			name: "Check error: attempt at changing description",
			args: args{
				id: "abcd-defg",
				class: libregraph.EducationClass{
					Description: &changeString,
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool { return assert.Error(tt, err) },
			modifyData: modifyData{
				arg: &ldap.ModifyRequest{},
			},
			modifyDNData: modifyDNData{
				arg: &ldap.ModifyDNRequest{},
				ret: nil,
			},
			searchData: searchData{
				res: &ldap.SearchResult{
					Entries: []*ldap.Entry{classEntry},
				},
				err: nil,
			},
		},
		{
			name: "Check error: attempt at changing classification",
			args: args{
				id: "abcd-defg",
				class: libregraph.EducationClass{
					Classification: changeString,
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool { return assert.Error(tt, err) },
			modifyData: modifyData{
				arg: &ldap.ModifyRequest{},
			},
			modifyDNData: modifyDNData{
				arg: &ldap.ModifyDNRequest{},
				ret: nil,
			},
			searchData: searchData{
				res: &ldap.SearchResult{
					Entries: []*ldap.Entry{classEntry},
				},
				err: nil,
			},
		},
		{
			name: "Check error: attempt at changing members",
			args: args{
				id: "abcd-defg",
				class: libregraph.EducationClass{
					Members: []libregraph.User{*libregraph.NewUser()},
				},
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool { return assert.Error(tt, err) },
			modifyData: modifyData{
				arg: &ldap.ModifyRequest{},
			},
			modifyDNData: modifyDNData{
				arg: &ldap.ModifyDNRequest{},
				ret: nil,
			},
			searchData: searchData{
				res: &ldap.SearchResult{
					Entries: []*ldap.Entry{classEntry},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := &mocks.Client{}
			b, err := getMockedBackend(lm, eduConfig, &logger)
			if err != nil {
				panic(err)
			}

			lm.On("Modify", tt.modifyData.arg).Return(tt.modifyData.ret)
			lm.On("ModifyDN", tt.modifyDNData.arg).Return(tt.modifyDNData.ret)
			lm.On("Search", mock.Anything).Return(tt.searchData.res, tt.searchData.err)

			ctx := context.Background()

			_, err = b.UpdateEducationClass(ctx, tt.args.id, tt.args.class)
			tt.assertion(t, err)
		})
	}
}

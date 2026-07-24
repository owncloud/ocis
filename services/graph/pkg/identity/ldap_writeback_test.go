package identity

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAttrsFromAddRequest(t *testing.T) {
	ar := ldap.NewAddRequest("uid=alice,ou=people,dc=test", nil)
	ar.Attribute("uid", []string{"alice"})
	ar.Attribute("displayname", []string{"Alice Example"})
	ar.Attribute("mail", []string{"alice@example.org", "a@example.org"})

	attrs := attrsFromAddRequest(ar)

	assert.Equal(t, []string{"alice"}, attrs["uid"])
	assert.Equal(t, []string{"Alice Example"}, attrs["displayname"])
	assert.Equal(t, []string{"alice@example.org", "a@example.org"}, attrs["mail"])

	// round-trips through ldap.NewEntry with case-insensitive lookup
	e := ldap.NewEntry(ar.DN, attrs)
	assert.Equal(t, ar.DN, e.DN)
	assert.Equal(t, "alice", e.GetEqualFoldAttributeValue("UID"))
	assert.Equal(t, "Alice Example", e.GetEqualFoldAttributeValue("displayName"))
	assert.Equal(t, []string{"alice@example.org", "a@example.org"}, e.GetEqualFoldAttributeValues("mail"))
}

// TestAttrsFromAddRequestDoesNotAliasRequest guards against attrsFromAddRequest
// returning slices that alias ar.Attributes[].Vals: mutating the returned map must
// never change the AddRequest's own values, and vice versa.
func TestAttrsFromAddRequestDoesNotAliasRequest(t *testing.T) {
	ar := ldap.NewAddRequest("uid=alice,ou=people,dc=test", nil)
	ar.Attribute("mail", []string{"alice@example.org"})

	attrs := attrsFromAddRequest(ar)
	attrs["mail"][0] = "mutated@example.org"

	assert.Equal(t, []string{"alice@example.org"}, ar.Attributes[0].Vals)
}

func TestApplyModifyToEntry(t *testing.T) {
	base := ldap.NewEntry("uid=alice,ou=people,dc=test", map[string][]string{
		"uid":         {"alice"},
		"displayname": {"Alice Example"},
		"mail":        {"alice@example.org"},
		"member":      {"cn=a", "cn=b"},
	})

	t.Run("Replace overwrites existing attribute", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: base.DN}
		mr.Replace("displayname", []string{"Alice New"})

		got := applyModifyToEntry(base, mr)
		assert.Equal(t, "Alice New", got.GetEqualFoldAttributeValue("displayname"))
	})

	t.Run("Replace adds a not-yet-present attribute", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: base.DN}
		mr.Replace("givenname", []string{"Alice"})

		got := applyModifyToEntry(base, mr)
		assert.Equal(t, "Alice", got.GetEqualFoldAttributeValue("givenname"))
	})

	t.Run("Add appends to existing attribute", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: base.DN}
		mr.Add("member", []string{"cn=c"})

		got := applyModifyToEntry(base, mr)
		assert.Equal(t, []string{"cn=a", "cn=b", "cn=c"}, got.GetEqualFoldAttributeValues("member"))
	})

	t.Run("Delete whole attribute", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: base.DN}
		mr.Delete("mail", []string{})

		got := applyModifyToEntry(base, mr)
		assert.Empty(t, got.GetEqualFoldAttributeValues("mail"))
	})

	t.Run("Delete specific value", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: base.DN}
		mr.Delete("member", []string{"cn=a"})

		got := applyModifyToEntry(base, mr)
		assert.Equal(t, []string{"cn=b"}, got.GetEqualFoldAttributeValues("member"))
	})

	t.Run("Delete last remaining value drops the attribute", func(t *testing.T) {
		single := ldap.NewEntry(base.DN, map[string][]string{
			"member": {"cn=a"},
		})
		mr := &ldap.ModifyRequest{DN: single.DN}
		mr.Delete("member", []string{"cn=a"})

		got := applyModifyToEntry(single, mr)
		assert.Empty(t, got.GetEqualFoldAttributeValues("member"))
		for _, a := range got.Attributes {
			assert.NotEqual(t, "member", strings.ToLower(a.Name), "attribute must be dropped, not left empty")
		}
	})

	t.Run("nil ModifyRequest returns base unchanged", func(t *testing.T) {
		got := applyModifyToEntry(base, nil)
		assert.Equal(t, base.DN, got.DN)
		assert.Equal(t, "Alice Example", got.GetEqualFoldAttributeValue("displayname"))
		assert.Equal(t, []string{"cn=a", "cn=b"}, got.GetEqualFoldAttributeValues("member"))
	})

	t.Run("nil base returns nil", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: "uid=alice,ou=people,dc=test"}
		mr.Replace("displayname", []string{"Alice New"})

		got := applyModifyToEntry(nil, mr)
		assert.Nil(t, got)
	})

	t.Run("case-insensitive match creates no duplicate", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: base.DN}
		mr.Replace("DisplayName", []string{"Alice Caps"})

		got := applyModifyToEntry(base, mr)
		assert.Equal(t, "Alice Caps", got.GetEqualFoldAttributeValue("displayname"))
		// exactly one attribute entry whose name folds to displayname
		count := 0
		for _, a := range got.Attributes {
			if strings.EqualFold(a.Name, "displayname") {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})

	t.Run("base is never mutated", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: base.DN}
		mr.Replace("displayname", []string{"Mutated"})
		mr.Add("member", []string{"cn=z"})
		mr.Delete("mail", []string{})

		_ = applyModifyToEntry(base, mr)
		assert.Equal(t, "Alice Example", base.GetEqualFoldAttributeValue("displayname"))
		assert.Equal(t, []string{"cn=a", "cn=b"}, base.GetEqualFoldAttributeValues("member"))
		assert.Equal(t, "alice@example.org", base.GetEqualFoldAttributeValue("mail"))
	})

	t.Run("Replace populates ByteValues for raw accessors", func(t *testing.T) {
		mr := &ldap.ModifyRequest{DN: base.DN}
		mr.Replace("displayname", []string{"Raw Value"})

		got := applyModifyToEntry(base, mr)
		assert.Equal(t, []byte("Raw Value"), got.GetEqualFoldRawAttributeValue("displayname"))
	})
}

// TestCreateUserSynthesizesWhenNotServerUUID asserts that with useServerUUID=false
// CreateUser builds the response model from the AddRequest instead of reading the
// entry back — i.e. no Search is issued after the Add.
func TestCreateUserSynthesizesWhenNotServerUUID(t *testing.T) {
	displayName := "DisplayName"
	mail := "user@example"
	userName := "user"
	surname := "surname"
	givenName := "givenName"
	userType := "Member"

	l := &mocks.Client{}
	// Add succeeds; Search must NOT be called (no read-back).
	l.On("Add", mock.Anything).Return(nil)
	logger := log.NewLogger(log.Level("debug"))

	user := libregraph.NewUser(displayName, userName)
	user.SetMail(mail)
	user.SetSurname(surname)
	user.SetGivenName(givenName)
	user.SetAccountEnabled(true)
	user.SetUserType(userType)

	c := lconfig
	c.UseServerUUID = false
	b, err := NewLDAPBackend(l, c, &logger, "")
	assert.Nil(t, err)

	newUser, err := b.CreateUser(context.Background(), *user)
	assert.Nil(t, err)
	assert.Equal(t, displayName, newUser.GetDisplayName())
	assert.Equal(t, mail, newUser.GetMail())
	assert.Equal(t, userName, newUser.GetOnPremisesSamAccountName())
	assert.Equal(t, givenName, newUser.GetGivenName())
	assert.Equal(t, surname, newUser.GetSurname())
	assert.True(t, newUser.GetAccountEnabled())
	assert.Equal(t, userType, newUser.GetUserType())
	// oCIS generated the id and wrote it into the Add; it must be present on the model.
	assert.NotEmpty(t, newUser.GetId())

	l.AssertNotCalled(t, "Search", mock.Anything)
}

// TestCreateUserSynthesizedModelEqualsReadBack is the regression guard: the model
// CreateUser returns from the synthesize path (useServerUUID=false) must equal the
// model the unchanged builder produces from the entry a read-back would have returned
// — i.e. the server's echo of exactly what oCIS wrote. Eliminating the read-back must
// not change the CreateUser response for any populated field.
func TestCreateUserSynthesizedModelEqualsReadBack(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	user := libregraph.NewUser("DisplayName", "user")
	user.SetMail("user@example")
	user.SetSurname("surname")
	user.SetGivenName("givenName")
	user.SetAccountEnabled(true)
	user.SetUserType("Member")

	// Capture the AddRequest so we can reproduce the entry a read-back would return.
	var written *ldap.AddRequest
	l := &mocks.Client{}
	l.On("Add", mock.Anything).Run(func(args mock.Arguments) {
		written = args.Get(0).(*ldap.AddRequest)
	}).Return(nil)

	c := lconfig
	c.UseServerUUID = false
	b, err := NewLDAPBackend(l, c, &logger, "")
	assert.Nil(t, err)

	synthUser, err := b.CreateUser(context.Background(), *user)
	assert.Nil(t, err)

	// Ground truth: what the old read-back path would have built. A base-scoped search
	// on the DN returns exactly the attributes that were written, so the read-back entry
	// is ldap.NewEntry over the captured AddRequest.
	readBackEntry := ldap.NewEntry(written.DN, attrsFromAddRequest(written))
	readBackModel := b.createUserModelFromLDAP(readBackEntry)

	assert.Equal(t, readBackModel, synthUser)
}

// TestSynthesizedEntryPopulatesInstancesWithoutSearch is the regression guard for the
// attribute-sourced multi-instance fields. With instanceURLTemplate +
// crossInstanceReferenceTemplate configured and instanceMapperEnabled=false, building a
// user model from a synthesized entry (as the create/update fold paths do) must populate
// user.Instances / CrossInstanceReference AND must not issue a live Search — the codepath
// only searches when instanceMapperEnabled=true, which is mutually exclusive with the
// "no second LDAP call" guarantee (see the design spec's Testing section).
func TestSynthesizedEntryPopulatesInstancesWithoutSearch(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	c := lconfig
	c.UserMemberAttribute = "ocMemberOfInstance"
	c.InstanceURLTemplate = "https://{{.Instancename}}.example.org"
	c.CrossInstanceReferenceTemplate = "{{.Username}}@{{.Instancename}}"
	c.InstanceMapperEnabled = false

	l := &mocks.Client{}
	// A non-empty instanceID is required for the instance templates to be parsed.
	b, err := NewLDAPBackend(l, c, &logger, "instanceA")
	assert.Nil(t, err)

	// Synthesized entry as attrsFromAddRequest + ldap.NewEntry would produce, carrying a
	// member-of-instance value so the Instances loop runs.
	entry := ldap.NewEntry("uid=user,ou=people,dc=test", map[string][]string{
		"uid":                {"user"},
		"displayname":        {"DisplayName"},
		"mail":               {"user@example"},
		"entryUUID":          {"abcd-defg"},
		"ocMemberOfInstance": {"instanceA"},
	})

	model := b.createUserModelFromLDAP(entry)
	assert.NotNil(t, model)
	assert.NotEmpty(t, model.Instances)
	assert.Equal(t, "https://instanceA.example.org", model.Instances[0].GetUrl())
	assert.NotNil(t, model.CrossInstanceReference)
	assert.Equal(t, "user@instanceA", model.GetCrossInstanceReference())

	// instanceMapperEnabled=false → getInstance returns the value verbatim, no Search.
	l.AssertNotCalled(t, "Search", mock.Anything)
}

// TestCreateGroupSynthesizesWhenNotServerUUID asserts CreateGroup builds its response
// from the AddRequest (no read-back) when oCIS generated the id itself.
func TestCreateGroupSynthesizesWhenNotServerUUID(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	// oCIS writes the generated group id under the configured GroupIDAttribute, so the
	// synthesized entry must carry the id under that attribute for the model builder.
	c := lconfig
	c.UseServerUUID = false
	c.GroupIDAttribute = "owncloudUUID"

	var written *ldap.AddRequest
	l := &mocks.Client{}
	l.On("Add", mock.Anything).Run(func(args mock.Arguments) {
		written = args.Get(0).(*ldap.AddRequest)
	}).Return(nil)

	b, err := NewLDAPBackend(l, c, &logger, "")
	assert.Nil(t, err)

	group := libregraph.NewGroup()
	group.SetDisplayName("mygroup")

	newGroup, err := b.CreateGroup(context.Background(), *group)
	assert.Nil(t, err)
	assert.Equal(t, "mygroup", newGroup.GetDisplayName())
	assert.NotEmpty(t, newGroup.GetId())
	l.AssertNotCalled(t, "Search", mock.Anything)

	// parity with the read-back-derived model
	readBackModel := b.createGroupModelFromLDAP(ldap.NewEntry(written.DN, attrsFromAddRequest(written)))
	assert.Equal(t, readBackModel, newGroup)
}

// TestCreateEducationUserSynthesizesWhenNotServerUUID asserts CreateEducationUser
// synthesizes its response (no read-back) when oCIS generated the id itself.
func TestCreateEducationUserSynthesizesWhenNotServerUUID(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	c := eduConfig
	c.UseServerUUID = false

	var written *ldap.AddRequest
	l := &mocks.Client{}
	l.On("Add", mock.Anything).Run(func(args mock.Arguments) {
		written = args.Get(0).(*ldap.AddRequest)
	}).Return(nil)

	b, err := getMockedBackend(l, c, &logger)
	assert.Nil(t, err)

	user := libregraph.NewEducationUser()
	user.SetDisplayName("Test User")
	user.SetOnPremisesSamAccountName("testuser")
	user.SetMail("testuser@example.org")
	user.SetPrimaryRole("student")
	user.SetUserType("Member")
	user.SetAccountEnabled(false)
	user.SetExternalID("ext-ernal-id")

	eduUser, err := b.CreateEducationUser(context.Background(), *user)
	assert.Nil(t, err)
	assert.Equal(t, "Test User", eduUser.GetDisplayName())
	assert.Equal(t, "testuser", eduUser.GetOnPremisesSamAccountName())
	assert.NotEmpty(t, eduUser.GetId())
	l.AssertNotCalled(t, "Search", mock.Anything)

	readBackModel := b.createEducationUserModelFromLDAP(ldap.NewEntry(written.DN, attrsFromAddRequest(written)))
	assert.Equal(t, readBackModel, eduUser)
}

// TestCreateEducationClassSynthesizesWhenNotServerUUID asserts CreateEducationClass
// synthesizes its response (no read-back) when oCIS generated the id itself.
func TestCreateEducationClassSynthesizesWhenNotServerUUID(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	c := eduConfig
	c.UseServerUUID = false
	c.GroupIDAttribute = "owncloudUUID"

	var written *ldap.AddRequest
	l := &mocks.Client{}
	l.On("Add", mock.Anything).Run(func(args mock.Arguments) {
		written = args.Get(0).(*ldap.AddRequest)
	}).Return(nil)

	b, err := getMockedBackend(l, c, &logger)
	assert.Nil(t, err)

	class := libregraph.NewEducationClass()
	class.SetExternalId("Math0123")
	class.SetDisplayName("Math")
	class.SetClassification("course")

	resClass, err := b.CreateEducationClass(context.Background(), *class)
	assert.Nil(t, err)
	assert.Equal(t, "Math", resClass.GetDisplayName())
	assert.Equal(t, "Math0123", resClass.GetExternalId())
	assert.Equal(t, "course", resClass.GetClassification())
	assert.NotEmpty(t, resClass.GetId())
	l.AssertNotCalled(t, "Search", mock.Anything)

	readBackModel := b.createEducationClassModelFromLDAP(ldap.NewEntry(written.DN, attrsFromAddRequest(written)))
	assert.Equal(t, readBackModel, resClass)
}

// TestCreateEducationSchoolSynthesizesWhenNotServerUUID asserts CreateEducationSchool
// synthesizes its response (no read-back) when oCIS generated the id itself. The
// pre-create duplicate-number Search still fires; only the read-back is eliminated.
func TestCreateEducationSchoolSynthesizesWhenNotServerUUID(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	c := eduConfig
	c.UseServerUUID = false

	var written *ldap.AddRequest
	l := &mocks.Client{}
	l.On("Add", mock.Anything).Run(func(args mock.Arguments) {
		written = args.Get(0).(*ldap.AddRequest)
	}).Return(nil)
	// The duplicate-schoolNumber lookup before Add is a subtree search that returns
	// "not found" so create proceeds. No base-scoped read-back must follow the Add.
	l.On("Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.Scope == ldap.ScopeWholeSubtree
	})).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)

	b, err := getMockedBackend(l, c, &logger)
	assert.Nil(t, err)

	school := libregraph.NewEducationSchool()
	school.SetDisplayName("Test School")
	school.SetSchoolNumber("0123")

	resSchool, err := b.CreateEducationSchool(context.Background(), *school)
	assert.Nil(t, err)
	assert.Equal(t, "Test School", resSchool.GetDisplayName())
	assert.Equal(t, "0123", resSchool.GetSchoolNumber())
	assert.NotEmpty(t, resSchool.GetId())
	// no base-scoped read-back on the created school DN
	l.AssertNotCalled(t, "Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.Scope == ldap.ScopeBaseObject
	}))

	readBackModel := b.createSchoolModelFromLDAP(ldap.NewEntry(written.DN, attrsFromAddRequest(written)))
	assert.Equal(t, readBackModel, resSchool)
}

// TestUpdateUserFoldsNoReadBack asserts a non-rename UpdateUser does not read the
// entry back: it pre-reads once, Modifies, and folds the ModifyRequest onto the
// pre-read entry to build the response.
func TestUpdateUserFoldsNoReadBack(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	preRead := ldap.NewEntry("uid=user,ou=people,dc=test", map[string][]string{
		"uid":               {"user"},
		"displayname":       {"Old Name"},
		"mail":              {"old@example.org"},
		"entryUUID":         {"abcd-defg"},
		"userTypeAttribute": {"Member"},
	})

	l := &mocks.Client{}
	// Pre-read (subtree search by name/id) returns the entry once.
	l.On("Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.Scope == ldap.ScopeWholeSubtree
	})).Return(&ldap.SearchResult{Entries: []*ldap.Entry{preRead}}, nil).Once()
	l.On("Modify", mock.Anything).Return(nil)

	b, err := getMockedBackend(l, lconfig, &logger)
	assert.Nil(t, err)

	id := "abcd-defg"
	newMail := "new@example.org"
	update := libregraph.UserUpdate{Mail: &newMail}
	update.Id = &id

	got, err := b.UpdateUser(context.Background(), "user", update)
	assert.Nil(t, err)
	assert.Equal(t, newMail, got.GetMail())
	assert.Equal(t, "Old Name", got.GetDisplayName())
	assert.Equal(t, "abcd-defg", got.GetId())

	// No base-scoped read-back after the Modify.
	l.AssertNotCalled(t, "Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.Scope == ldap.ScopeBaseObject
	}))
}

// TestUpdateEducationUserFoldsNoReadBack asserts a non-rename UpdateEducationUser
// folds the ModifyRequest onto the pre-read entry rather than reading it back.
func TestUpdateEducationUserFoldsNoReadBack(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	preRead := ldap.NewEntry("uid=testuser,ou=people,dc=test", map[string][]string{
		"uid":               {"testuser"},
		"displayname":       {"Test User"},
		"mail":              {"old@example.org"},
		"entryUUID":         {"abcd-defg"},
		"userClass":         {"student"},
		"userTypeAttribute": {"Member"},
	})

	l := &mocks.Client{}
	// Pre-read by name/id (subtree) returns the entry; no base-scoped read-back after.
	l.On("Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.Scope == ldap.ScopeWholeSubtree
	})).Return(&ldap.SearchResult{Entries: []*ldap.Entry{preRead}}, nil).Once()
	l.On("Modify", mock.Anything).Return(nil)

	b, err := getMockedBackend(l, eduConfig, &logger)
	assert.Nil(t, err)

	newMail := "new@example.org"
	user := libregraph.NewEducationUser()
	user.SetMail(newMail)

	got, err := b.UpdateEducationUser(context.Background(), "testuser", *user)
	assert.Nil(t, err)
	assert.Equal(t, newMail, got.GetMail())
	assert.Equal(t, "Test User", got.GetDisplayName())
	assert.Equal(t, "abcd-defg", got.GetId())

	l.AssertNotCalled(t, "Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.Scope == ldap.ScopeBaseObject
	}))
}

// TestUpdateEducationSchoolFoldsNoReadBack asserts UpdateEducationSchool folds the
// applied property changes onto the pre-read entry rather than reading it back.
func TestUpdateEducationSchoolFoldsNoReadBack(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	preRead := ldap.NewEntry("ou=Test School", map[string][]string{
		"ou":                      {"Test School"},
		"ocEducationSchoolNumber": {"0123"},
		"owncloudUUID":            {"abcd-defg"},
	})

	l := &mocks.Client{}
	// Single pre-read; the termination Replace is folded on. No second Search.
	l.On("Search", mock.Anything).Return(&ldap.SearchResult{Entries: []*ldap.Entry{preRead}}, nil).Once()
	l.On("Modify", mock.Anything).Return(nil)

	b, err := getMockedBackend(l, eduConfig, &logger)
	assert.Nil(t, err)

	school := libregraph.NewEducationSchool()
	terminationTime := time.Date(2042, time.January, 31, 12, 0, 0, 0, time.UTC)
	school.SetTerminationDate(terminationTime)

	got, err := b.UpdateEducationSchool(context.Background(), "abcd-defg", *school)
	assert.Nil(t, err)
	assert.Equal(t, "Test School", got.GetDisplayName())
	assert.Equal(t, "abcd-defg", got.GetId())
	assert.Equal(t, "0123", got.GetSchoolNumber())
	// the folded termination date is reflected in the response
	assert.True(t, got.HasTerminationDate())
	assert.True(t, terminationTime.Equal(got.GetTerminationDate()))

	l.AssertNumberOfCalls(t, "Search", 1)
}

// TestUpdateEducationClassFoldPreservesClassificationAndExternalID guards the Site 9
// regression: folding must not drop classification/externalID from the response, even
// on an update that only changes displayName. It also asserts no read-back after Modify.
func TestUpdateEducationClassFoldPreservesClassificationAndExternalID(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	// Full education-class entry as getEducationClassByID would return it.
	preRead := ldap.NewEntry("ocEducationExternalId=Math0123,ou=groups,dc=test", map[string][]string{
		"cn":                    {"Math"},
		"entryUUID":             {"abcd-defg"},
		"ocEducationExternalId": {"Math0123"},
		"ocEducationClassType":  {"course"},
	})

	l := &mocks.Client{}
	l.On("Search", mock.Anything).Return(&ldap.SearchResult{Entries: []*ldap.Entry{preRead}}, nil).Once()
	l.On("Modify", mock.Anything).Return(nil)

	b, err := getMockedBackend(l, eduConfig, &logger)
	assert.Nil(t, err)

	newName := "Math-2"
	class := libregraph.EducationClass{DisplayName: &newName}

	got, err := b.UpdateEducationClass(context.Background(), "abcd-defg", class)
	assert.Nil(t, err)
	assert.Equal(t, "Math-2", got.GetDisplayName())
	assert.Equal(t, "abcd-defg", got.GetId())
	// classification and externalID must survive a displayName-only update
	assert.Equal(t, "course", got.GetClassification())
	assert.Equal(t, "Math0123", got.GetExternalId())

	l.AssertNotCalled(t, "Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.Scope == ldap.ScopeBaseObject
	}))
}

// TestCreateGroupSynthesizesWithNonDefaultIDAttribute guards the fix for the
// hardcoded "owncloudUUID" id write: the generated id must be stored under the
// configured GroupIDAttribute, otherwise the synthesized model has no id and
// createGroupModelFromLDAP returns nil (CreateGroup would return (nil, nil)).
func TestCreateGroupSynthesizesWithNonDefaultIDAttribute(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	c := lconfig
	c.UseServerUUID = false
	c.GroupIDAttribute = "entryUUID" // non-default: previously the id was written to owncloudUUID only

	var written *ldap.AddRequest
	l := &mocks.Client{}
	l.On("Add", mock.Anything).Run(func(args mock.Arguments) {
		written = args.Get(0).(*ldap.AddRequest)
	}).Return(nil)

	b, err := NewLDAPBackend(l, c, &logger, "")
	assert.Nil(t, err)

	group := libregraph.NewGroup()
	group.SetDisplayName("mygroup")

	newGroup, err := b.CreateGroup(context.Background(), *group)
	assert.Nil(t, err)
	assert.NotNil(t, newGroup)
	assert.Equal(t, "mygroup", newGroup.GetDisplayName())
	// the id was written under the configured attribute and survives synthesis
	assert.NotEmpty(t, newGroup.GetId())
	assert.NotEmpty(t, written.Attributes)
	// the AddRequest carries the id under the configured attribute, not "owncloudUUID"
	var idFromConfigured, idFromHardcoded string
	for _, a := range written.Attributes {
		switch a.Type {
		case c.GroupIDAttribute:
			if len(a.Vals) > 0 {
				idFromConfigured = a.Vals[0]
			}
		case "owncloudUUID":
			if len(a.Vals) > 0 {
				idFromHardcoded = a.Vals[0]
			}
		}
	}
	assert.NotEmpty(t, idFromConfigured, "id must be written under the configured GroupIDAttribute")
	assert.Empty(t, idFromHardcoded, "id must not be written under a hardcoded owncloudUUID")
	assert.Equal(t, idFromConfigured, newGroup.GetId())
	l.AssertNotCalled(t, "Search", mock.Anything)
}

// TestCreateEducationClassSynthesizesWithNonDefaultIDAttribute guards the same
// hardcoded-id fix on the education-class create path, which reuses
// groupToLDAPAttrValues.
func TestCreateEducationClassSynthesizesWithNonDefaultIDAttribute(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	c := eduConfig
	c.UseServerUUID = false
	c.GroupIDAttribute = "entryUUID" // non-default

	l := &mocks.Client{}
	l.On("Add", mock.Anything).Return(nil)

	b, err := getMockedBackend(l, c, &logger)
	assert.Nil(t, err)

	class := libregraph.NewEducationClass()
	class.SetExternalId("Math0123")
	class.SetDisplayName("Math")
	class.SetClassification("course")

	resClass, err := b.CreateEducationClass(context.Background(), *class)
	assert.Nil(t, err)
	assert.NotNil(t, resClass)
	assert.Equal(t, "Math", resClass.GetDisplayName())
	// the id survives synthesis under the configured attribute
	assert.NotEmpty(t, resClass.GetId())
	l.AssertNotCalled(t, "Search", mock.Anything)
}

// TestGetGroupByNameUsesGroupNameAttribute guards the fix for the name-only lookup
// branch of getLDAPGroupByNameOrID, which previously filtered by the user name
// attribute (uid) instead of the group name attribute (cn), so a group could never
// be found by name. That else-branch is reached when filterEscapeUUID errors, which
// happens for a non-UUID name when the id attribute is an octet string. The lookup
// is exercised directly so the assertion isolates the filter, independent of model
// building.
func TestGetGroupByNameUsesGroupNameAttribute(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))

	c := lconfig
	// Octet-string group IDs require server-assigned UUIDs; this also makes
	// filterEscapeUUID error on a non-UUID name, reaching the name-only else branch.
	c.UseServerUUID = true
	c.GroupIDIsOctetString = true

	var groupSearch *ldap.SearchRequest
	l := &mocks.Client{}
	l.On("Search", mock.Anything).Run(func(args mock.Arguments) {
		groupSearch = args.Get(0).(*ldap.SearchRequest)
	}).Return(&ldap.SearchResult{Entries: []*ldap.Entry{groupEntry}}, nil)

	b, err := getMockedBackend(l, c, &logger)
	assert.Nil(t, err)

	_, err = b.getLDAPGroupByNameOrID("group", false)
	assert.Nil(t, err)
	assert.NotNil(t, groupSearch)
	// the filter must use the group name attribute (cn), not the user name attribute (uid)
	assert.Contains(t, groupSearch.Filter, "("+b.groupAttributeMap.name+"=group)")
	assert.NotContains(t, groupSearch.Filter, b.userAttributeMap.userName+"=group")
}

// TestNewLDAPBackendRejectsOctetStringWithoutServerUUID guards the fail-fast check:
// octet-string ID attributes (server-assigned, e.g. AD objectGUID) are incompatible
// with UseServerUUID=false, where oCIS generates string UUIDs. The combination must
// be rejected at startup rather than silently producing corrupt IDs.
func TestNewLDAPBackendRejectsOctetStringWithoutServerUUID(t *testing.T) {
	logger := log.NewLogger(log.Level("debug"))
	l := &mocks.Client{}

	t.Run("user octet-string without server UUID is rejected", func(t *testing.T) {
		c := lconfig
		c.UseServerUUID = false
		c.UserIDIsOctetString = true
		_, err := NewLDAPBackend(l, c, &logger, "")
		assert.Error(t, err)
	})

	t.Run("group octet-string without server UUID is rejected", func(t *testing.T) {
		c := lconfig
		c.UseServerUUID = false
		c.GroupIDIsOctetString = true
		_, err := NewLDAPBackend(l, c, &logger, "")
		assert.Error(t, err)
	})

	t.Run("octet-string with server UUID is allowed", func(t *testing.T) {
		c := lconfig
		c.UseServerUUID = true
		c.UserIDIsOctetString = true
		c.GroupIDIsOctetString = true
		_, err := NewLDAPBackend(l, c, &logger, "")
		assert.NoError(t, err)
	})

	t.Run("string IDs without server UUID is allowed", func(t *testing.T) {
		c := lconfig
		c.UseServerUUID = false
		c.UserIDIsOctetString = false
		c.GroupIDIsOctetString = false
		_, err := NewLDAPBackend(l, c, &logger, "")
		assert.NoError(t, err)
	})
}

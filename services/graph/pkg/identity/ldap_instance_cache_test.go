package identity

import (
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var instanceMapperConfig = config.LDAP{
	UserBaseDN:               "ou=people,dc=test",
	UserObjectClass:          "inetOrgPerson",
	UserSearchScope:          "sub",
	UserDisplayNameAttribute: "displayname",
	UserIDAttribute:          "entryUUID",
	UserEmailAttribute:       "mail",
	UserNameAttribute:        "uid",

	GroupBaseDN:          "ou=groups,dc=test",
	GroupObjectClass:     "groupOfNames",
	GroupSearchScope:     "sub",
	GroupNameAttribute:   "cn",
	GroupMemberAttribute: "member",
	GroupIDAttribute:     "entryUUID",

	InstanceMapperEnabled:       true,
	InstanceMapperBaseDN:        "ou=instances,dc=test",
	InstanceMapperNameAttribute: "instanceName",
	InstanceMapperIDAttribute:   "instanceID",
}

func TestGetInstanceNameIsCached(t *testing.T) {
	lm := &mocks.Client{}
	b, err := getMockedBackend(lm, instanceMapperConfig, &logger)
	assert.Nil(t, err)

	instanceEntry := ldap.NewEntry("instanceID=instance-1,ou=instances,dc=test",
		map[string][]string{
			"instanceName": {"Instance One"},
		})
	lm.On("Search", mock.Anything).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{instanceEntry}}, nil)

	name1, err1 := b.getInstanceName("instance-1")
	assert.Nil(t, err1)
	assert.Equal(t, "Instance One", name1)

	name2, err2 := b.getInstanceName("instance-1")
	assert.Nil(t, err2)
	assert.Equal(t, "Instance One", name2)

	lm.AssertNumberOfCalls(t, "Search", 1)
}

func TestGetInstanceNameCachesNotFound(t *testing.T) {
	lm := &mocks.Client{}
	b, err := getMockedBackend(lm, instanceMapperConfig, &logger)
	assert.Nil(t, err)

	lm.On("Search", mock.Anything).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)

	_, err1 := b.getInstanceName("stale-instance")
	assert.ErrorIs(t, err1, ErrNotFound)

	_, err2 := b.getInstanceName("stale-instance")
	assert.ErrorIs(t, err2, ErrNotFound)

	lm.AssertNumberOfCalls(t, "Search", 1)
}

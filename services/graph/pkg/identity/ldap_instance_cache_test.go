package identity

import (
	"testing"
	"time"

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
	InstanceMapperCacheTTL:      time.Minute,
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

func TestGetInstanceNameZeroTTLDisablesCache(t *testing.T) {
	lm := &mocks.Client{}
	cfg := instanceMapperConfig
	cfg.InstanceMapperCacheTTL = 0
	b, err := getMockedBackend(lm, cfg, &logger)
	assert.Nil(t, err)

	instanceEntry := ldap.NewEntry("instanceID=instance-1,ou=instances,dc=test",
		map[string][]string{
			"instanceName": {"Instance One"},
		})
	lm.On("Search", mock.Anything).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{instanceEntry}}, nil)

	for range 2 {
		name, err := b.getInstanceName("instance-1")
		assert.Nil(t, err)
		assert.Equal(t, "Instance One", name)
	}

	// A TTL of 0 means "don't cache", not "cache forever".
	lm.AssertNumberOfCalls(t, "Search", 2)
}

func TestGetInstanceNameCacheCapacityIsHonoured(t *testing.T) {
	lm := &mocks.Client{}
	cfg := instanceMapperConfig
	cfg.InstanceMapperCacheCapacity = 1
	b, err := getMockedBackend(lm, cfg, &logger)
	assert.Nil(t, err)

	lm.On("Search", mock.Anything).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)

	// Caching "a" then "b" evicts "a", because only one entry fits, so asking
	// for "a" again has to hit LDAP a second time.
	for _, instanceID := range []string{"a", "b", "a"} {
		_, err := b.getInstanceName(instanceID)
		assert.ErrorIs(t, err, ErrNotFound)
	}

	lm.AssertNumberOfCalls(t, "Search", 3)
}

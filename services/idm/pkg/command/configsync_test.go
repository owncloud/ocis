package command

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// sampleOcisConfig mirrors the relevant subset of an `ocis init` generated
// ocis.yaml, carrying the bind_password keys the reva/libregraph/idp service
// users are wired into.
const sampleOcisConfig = `token_manager:
  jwt_secret: keepme
idm:
  service_user_passwords:
    admin_password: adminpw
    idm_password: OLDIDM
    reva_password: OLDREVA
    idp_password: OLDIDP
idp:
  ldap:
    bind_password: OLDIDP
auth_basic:
  auth_providers:
    ldap:
      bind_password: OLDREVA
users:
  drivers:
    ldap:
      bind_password: OLDREVA
groups:
  drivers:
    ldap:
      bind_password: OLDREVA
graph:
  identity:
    ldap:
      bind_password: OLDIDM
`

func getYAMLPath(t *testing.T, data []byte, path ...string) string {
	t.Helper()
	var root yaml.Node
	require.NoError(t, yaml.Unmarshal(data, &root))
	node := root.Content[0]
	for _, key := range path {
		var next *yaml.Node
		for i := 0; i+1 < len(node.Content); i += 2 {
			if node.Content[i].Value == key {
				next = node.Content[i+1]
				break
			}
		}
		require.NotNilf(t, next, "path %v: key %q not found", path, key)
		node = next
	}
	return node.Value
}

func TestSyncServiceUserPasswordReva(t *testing.T) {
	out, updated, err := syncServiceUserPasswordInConfig([]byte(sampleOcisConfig), "reva", "NEWPW")
	require.NoError(t, err)

	// reva binds via auth_basic, users and groups, and is bootstrapped via
	// idm.service_user_passwords.reva_password — all must move to NEWPW.
	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "auth_basic", "auth_providers", "ldap", "bind_password"))
	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "users", "drivers", "ldap", "bind_password"))
	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "groups", "drivers", "ldap", "bind_password"))
	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "idm", "service_user_passwords", "reva_password"))

	// Unrelated service users must stay untouched.
	assert.Equal(t, "OLDIDM", getYAMLPath(t, out, "graph", "identity", "ldap", "bind_password"))
	assert.Equal(t, "OLDIDP", getYAMLPath(t, out, "idp", "ldap", "bind_password"))
	assert.Equal(t, "OLDIDM", getYAMLPath(t, out, "idm", "service_user_passwords", "idm_password"))

	// Unrelated keys and comments-free structure must survive the round-trip.
	assert.Equal(t, "keepme", getYAMLPath(t, out, "token_manager", "jwt_secret"))

	assert.ElementsMatch(t, []string{
		"auth_basic.auth_providers.ldap.bind_password",
		"users.drivers.ldap.bind_password",
		"groups.drivers.ldap.bind_password",
		"idm.service_user_passwords.reva_password",
	}, updated)
}

func TestSyncServiceUserPasswordLibregraph(t *testing.T) {
	out, updated, err := syncServiceUserPasswordInConfig([]byte(sampleOcisConfig), "libregraph", "NEWPW")
	require.NoError(t, err)

	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "graph", "identity", "ldap", "bind_password"))
	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "idm", "service_user_passwords", "idm_password"))
	// reva keys stay put
	assert.Equal(t, "OLDREVA", getYAMLPath(t, out, "users", "drivers", "ldap", "bind_password"))

	assert.ElementsMatch(t, []string{
		"graph.identity.ldap.bind_password",
		"idm.service_user_passwords.idm_password",
	}, updated)
}

func TestSyncServiceUserPasswordIdp(t *testing.T) {
	out, updated, err := syncServiceUserPasswordInConfig([]byte(sampleOcisConfig), "idp", "NEWPW")
	require.NoError(t, err)

	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "idp", "ldap", "bind_password"))
	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "idm", "service_user_passwords", "idp_password"))

	assert.ElementsMatch(t, []string{
		"idp.ldap.bind_password",
		"idm.service_user_passwords.idp_password",
	}, updated)
}

func TestSyncServiceUserPasswordUnknownUser(t *testing.T) {
	// An arbitrary (non service) user has no bind config to sync; the config
	// must be returned unchanged with no reported updates.
	out, updated, err := syncServiceUserPasswordInConfig([]byte(sampleOcisConfig), "einstein", "NEWPW")
	require.NoError(t, err)
	assert.Empty(t, updated)
	assert.Equal(t, strings.TrimSpace(sampleOcisConfig), strings.TrimSpace(string(out)))
}

func TestSyncServiceUserPasswordMissingKeys(t *testing.T) {
	// A slim config that omits some bind_password keys: only the present ones
	// are updated and reported, absent ones are silently skipped.
	const slim = `idm:
  service_user_passwords:
    reva_password: OLDREVA
`
	out, updated, err := syncServiceUserPasswordInConfig([]byte(slim), "reva", "NEWPW")
	require.NoError(t, err)
	assert.Equal(t, "NEWPW", getYAMLPath(t, out, "idm", "service_user_passwords", "reva_password"))
	assert.ElementsMatch(t, []string{"idm.service_user_passwords.reva_password"}, updated)
}

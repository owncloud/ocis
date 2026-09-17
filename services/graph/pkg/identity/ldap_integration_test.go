package identity

// Integration tests: full stack NewLDAPWithReconnect → NewLDAPBackend → libregraph/idm boltdb LDAP server.
// Uses the same server oCIS ships in production (services/idm). Regression baseline for
// GetUser, GetUsers, CreateUser, DeleteUser.

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/CiscoM31/godata"
	gldap "github.com/go-ldap/ldap/v3"
	"github.com/libregraph/idm/pkg/ldappassword"
	"github.com/libregraph/idm/pkg/ldapserver"
	"github.com/libregraph/idm/pkg/ldbbolt"
	"github.com/libregraph/idm/server/handler/boltdb"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	connldap "github.com/owncloud/reva/v2/pkg/utils/ldap"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	idmBaseDN  = "o=libregraph-idm"
	idmAdminDN = "uid=admin,ou=sysusers,o=libregraph-idm"
	idmAdminPW = "admin"
)

// idmConfig mirrors lconfig but with base DNs matching the libregraph/idm schema.
var idmConfig = lconfig

func init() {
	idmConfig.UserBaseDN = "ou=users,o=libregraph-idm"
	idmConfig.GroupBaseDN = "ou=groups,o=libregraph-idm"
	idmConfig.LdapDisabledUsersGroupDN = "cn=DisabledUsersGroup,ou=groups,o=libregraph-idm"
	idmConfig.WriteEnabled = true
	idmConfig.UseServerUUID = false // boltdb does not auto-generate entryUUID; client must supply it
}

// idmEntry builds a gldap.Entry matching idmConfig attribute names.
func idmEntry(uid, uuid string) *gldap.Entry {
	dn := fmt.Sprintf("uid=%s,%s", uid, idmConfig.UserBaseDN)
	return gldap.NewEntry(dn, map[string][]string{
		idmConfig.UserNameAttribute:        {uid},
		idmConfig.UserDisplayNameAttribute: {"Display " + uid},
		idmConfig.UserEmailAttribute:       {uid + "@example.com"},
		idmConfig.UserIDAttribute:          {uuid},
		"sn":                               {uid},
		"givenname":                        {uid},
		idmConfig.UserEnabledAttribute:     {"TRUE"},
		idmConfig.UserTypeAttribute:        {"Member"},
		"objectClass":                      {"inetOrgPerson", "organizationalPerson", "person", "top", "ownCloudUser"},
	})
}

// seedIDMDB seeds the boltdb with the OU structure, admin sysuser, and any additional entries.
func seedIDMDB(t *testing.T, dbfile string, entries ...*gldap.Entry) {
	t.Helper()

	bdb := &ldbbolt.LdbBolt{}
	require.NoError(t, bdb.Configure(logrus.New(), idmBaseDN, dbfile, nil))
	require.NoError(t, bdb.Initialize())
	defer bdb.Close()

	hashedPW, err := ldappassword.Hash(idmAdminPW, "{ARGON2}")
	require.NoError(t, err)

	// OU structure mirroring the real IDM base.ldif.tmpl
	base := []*gldap.Entry{
		gldap.NewEntry(idmBaseDN, map[string][]string{
			"o": {"libregraph-idm"}, "objectClass": {"organization"},
		}),
		gldap.NewEntry("ou=users,o=libregraph-idm", map[string][]string{
			"ou": {"users"}, "objectClass": {"organizationalUnit"},
		}),
		gldap.NewEntry("ou=sysusers,o=libregraph-idm", map[string][]string{
			"ou": {"sysusers"}, "objectClass": {"organizationalUnit"},
		}),
		gldap.NewEntry("ou=groups,o=libregraph-idm", map[string][]string{
			"ou": {"groups"}, "objectClass": {"organizationalUnit"},
		}),
		gldap.NewEntry(idmAdminDN, map[string][]string{
			"uid": {"admin"}, "objectClass": {"account", "simpleSecurityObject"},
			"userPassword": {hashedPW},
		}),
	}
	for _, e := range base {
		require.NoError(t, bdb.EntryPut(e), "seeding base entry %s", e.DN)
	}
	for _, e := range entries {
		require.NoError(t, bdb.EntryPut(e), "seeding entry %s", e.DN)
	}
}

// newIDMBackend starts the libregraph/idm boltdb LDAP server and returns a wired LDAP identity backend.
func newIDMBackend(t *testing.T, entries ...*gldap.Entry) *LDAP {
	t.Helper()

	dbfile := filepath.Join(t.TempDir(), "idm.db")
	seedIDMDB(t, dbfile, entries...)

	boltHandler, err := boltdb.NewBoltDBHandler(logrus.New(), dbfile, &boltdb.Options{
		BaseDN:  idmBaseDN,
		AdminDN: idmAdminDN,
	})
	require.NoError(t, err)

	srv := ldapserver.NewServer()
	srv.EnforceLDAP = true // boltdb does not apply LDAP filters; the server must
	ctx, cancel := context.WithCancel(context.Background())
	h := boltHandler.WithContext(ctx)
	srv.AddFunc("", h)
	srv.BindFunc("", h)
	srv.DeleteFunc("", h)
	srv.ModifyFunc("", h)
	srv.SearchFunc("", h)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	quit := make(chan bool)
	srv.QuitChannel(quit)
	go func() { _ = srv.Serve(ln) }()
	t.Cleanup(func() { cancel(); close(quit); ln.Close() })

	level := os.Getenv("OCIS_LOG_LEVEL")
	if level == "" {
		level = "error"
	}
	lgr := log.NewLogger(log.Level(level))
	conn := connldap.NewLDAPWithReconnect(connldap.Config{
		URI:           fmt.Sprintf("ldap://%s", ln.Addr().String()),
		BindDN:        idmAdminDN,
		BindPassword:  idmAdminPW,
		RetryMaxCount: 1,
	})
	conn.SetLogger(&lgr.Logger)

	backend, err := NewLDAPBackend(conn, idmConfig, &lgr, "", "")
	require.NoError(t, err)
	return backend
}

func odataReq(t *testing.T) *godata.GoDataRequest {
	t.Helper()
	req, err := godata.ParseRequest(context.Background(), "", url.Values{})
	require.NoError(t, err)
	return req
}

func TestIntegrationGetUser_HappyPath(t *testing.T) {
	b := newIDMBackend(t,
		idmEntry("alice", "aaaaaaaa-0000-0000-0000-000000000001"),
	)

	user, err := b.GetUser(context.Background(), "alice", odataReq(t))
	require.NoError(t, err)
	assert.Equal(t, "alice", user.GetOnPremisesSamAccountName())
	assert.Equal(t, "Display alice", user.GetDisplayName())
	assert.Equal(t, "alice@example.com", user.GetMail())
}

func TestIntegrationGetUsers_HappyPath(t *testing.T) {
	b := newIDMBackend(t,
		idmEntry("alice", "aaaaaaaa-0000-0000-0000-000000000001"),
		idmEntry("bob", "bbbbbbbb-0000-0000-0000-000000000002"),
	)

	users, err := b.GetUsers(context.Background(), odataReq(t))
	require.NoError(t, err)
	require.Len(t, users, 2)
	names := make(map[string]bool, 2)
	for _, u := range users {
		names[u.GetOnPremisesSamAccountName()] = true
	}
	assert.True(t, names["alice"])
	assert.True(t, names["bob"])
}

func TestIntegrationCreateUser_HappyPath(t *testing.T) {
	b := newIDMBackend(t) // empty — no pre-seeded users

	newUser := libregraph.NewUser("Charlie", "charlie")
	newUser.SetMail("charlie@example.com")
	newUser.SetSurname("Charlie")
	newUser.SetGivenName("Charlie")
	newUser.SetAccountEnabled(true)
	newUser.SetUserType("Member")

	created, err := b.CreateUser(context.Background(), *newUser)
	require.NoError(t, err)
	assert.Equal(t, "charlie", created.GetOnPremisesSamAccountName())
	assert.Equal(t, "Charlie", created.GetDisplayName())
	assert.Equal(t, "charlie@example.com", created.GetMail())
}

func TestIntegrationDeleteUser_HappyPath(t *testing.T) {
	entry := idmEntry("dave", "dddddddd-0000-0000-0000-000000000004")
	b := newIDMBackend(t, entry)

	// Confirm user exists first.
	_, err := b.GetUser(context.Background(), "dave", odataReq(t))
	require.NoError(t, err)

	err = b.DeleteUser(context.Background(), "dave")
	require.NoError(t, err)

	// Confirm gone.
	_, err = b.GetUser(context.Background(), "dave", odataReq(t))
	require.Error(t, err)
}

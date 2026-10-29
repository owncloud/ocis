package shared

import (
	"path/filepath"
	"testing"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
	"github.com/stretchr/testify/require"
)

func newTestCreds(t *testing.T) string {
	t.Helper()
	accountKP, err := nkeys.CreateAccount()
	require.NoError(t, err)
	userKP, err := nkeys.CreateUser()
	require.NoError(t, err)
	userPub, err := userKP.PublicKey()
	require.NoError(t, err)
	userJWT, err := jwt.NewUserClaims(userPub).Encode(accountKP)
	require.NoError(t, err)
	seed, err := userKP.Seed()
	require.NoError(t, err)
	creds, err := jwt.FormatUserConfig(userJWT, seed)
	require.NoError(t, err)
	return string(creds)
}

func TestCheckNatsCredentials(t *testing.T) {
	t.Run("no credentials", func(t *testing.T) {
		require.NoError(t, CheckNatsCredentials("proxy", "nats-js-kv", "", ""))
	})

	t.Run("valid content", func(t *testing.T) {
		require.NoError(t, CheckNatsCredentials("proxy", "nats-js-kv", "", newTestCreds(t)))
	})

	t.Run("partly written content", func(t *testing.T) {
		creds := newTestCreds(t)
		err := CheckNatsCredentials("proxy", "nats-js-kv", "", creds[:len(creds)/3])
		require.ErrorContains(t, err, "proxy")
	})

	t.Run("missing file", func(t *testing.T) {
		err := CheckNatsCredentials("proxy", "nats-js-kv", filepath.Join(t.TempDir(), "missing.creds"), "")
		require.ErrorContains(t, err, "proxy")
	})

	t.Run("store without nats", func(t *testing.T) {
		for _, store := range []string{"memory", "redis-sentinel", "noop", ""} {
			require.NoError(t, CheckNatsCredentials("proxy", store, filepath.Join(t.TempDir(), "missing.creds"), ""), store)
		}
	})

	t.Run("nats-js store", func(t *testing.T) {
		err := CheckNatsCredentials("proxy", "nats-js", filepath.Join(t.TempDir(), "missing.creds"), "")
		require.ErrorContains(t, err, "proxy")
	})
}

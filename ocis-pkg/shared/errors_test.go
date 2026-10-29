package shared

import (
	"path/filepath"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/natstest"
	"github.com/stretchr/testify/require"
)

func TestCheckNatsCredentials(t *testing.T) {
	t.Run("no credentials", func(t *testing.T) {
		require.NoError(t, CheckNatsCredentials("proxy", "nats-js-kv", "", ""))
	})

	t.Run("valid content", func(t *testing.T) {
		require.NoError(t, CheckNatsCredentials("proxy", "nats-js-kv", "", string(natstest.Creds(t))))
	})

	t.Run("partly written content", func(t *testing.T) {
		creds := string(natstest.Creds(t))
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

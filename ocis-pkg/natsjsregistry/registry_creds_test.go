package natsjsregistry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCreds(t *testing.T) {
	t.Run("no variables", func(t *testing.T) {
		credsFile, creds := getCreds()
		require.Empty(t, credsFile)
		require.Empty(t, creds)
	})

	t.Run("global variables", func(t *testing.T) {
		t.Setenv(_natsCredsFileEnv, "/global.creds")
		t.Setenv(_natsCredsEnv, "global")
		credsFile, creds := getCreds()
		require.Equal(t, "/global.creds", credsFile)
		require.Equal(t, "global", creds)
	})

	t.Run("registry file takes precedence over global variables", func(t *testing.T) {
		t.Setenv(_natsCredsFileEnv, "/global.creds")
		t.Setenv(_natsCredsEnv, "global")
		t.Setenv(_registryCredsFileEnv, "/registry.creds")
		credsFile, creds := getCreds()
		require.Equal(t, "/registry.creds", credsFile)
		require.Empty(t, creds)
	})

	t.Run("registry content takes precedence over global variables", func(t *testing.T) {
		t.Setenv(_natsCredsFileEnv, "/global.creds")
		t.Setenv(_registryCredsEnv, "registry")
		credsFile, creds := getCreds()
		require.Empty(t, credsFile)
		require.Equal(t, "registry", creds)
	})
}

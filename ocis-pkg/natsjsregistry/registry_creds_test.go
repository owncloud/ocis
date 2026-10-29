package natsjsregistry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/owncloud/ocis/v2/ocis-pkg/natstest"
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

	t.Run("registry username and password take precedence over global variables", func(t *testing.T) {
		t.Setenv(_natsCredsFileEnv, "/global.creds")
		t.Setenv(_natsCredsEnv, "global")
		t.Setenv(_registryUsernameEnv, "user")
		t.Setenv(_registryPasswordEnv, "pass")
		credsFile, creds := getCreds()
		require.Empty(t, credsFile)
		require.Empty(t, creds)
	})

	t.Run("registry username alone does not take precedence over global variables", func(t *testing.T) {
		t.Setenv(_natsCredsFileEnv, "/global.creds")
		t.Setenv(_registryUsernameEnv, "user")
		credsFile, creds := getCreds()
		require.Equal(t, "/global.creds", credsFile)
		require.Empty(t, creds)
	})

	t.Run("registry file takes precedence over registry username", func(t *testing.T) {
		t.Setenv(_registryCredsFileEnv, "/registry.creds")
		t.Setenv(_registryUsernameEnv, "user")
		credsFile, creds := getCreds()
		require.Equal(t, "/registry.creds", credsFile)
		require.Empty(t, creds)
	})
}

func TestCredentialsOptionChecksOnlyOnce(t *testing.T) {
	creds := natstest.Creds(t)
	credsFile := filepath.Join(t.TempDir(), "user.creds")
	require.NoError(t, os.WriteFile(credsFile, creds, 0o600))
	t.Setenv(_natsCredsFileEnv, credsFile)

	n := &storeregistry{}

	// first call: credentialsOption checks the file
	opt, err := n.credentialsOption()
	require.NoError(t, err)
	require.NotNil(t, opt)

	// a rotation writes the file partly
	require.NoError(t, os.WriteFile(credsFile, creds[:len(creds)/3], 0o600))

	// later call, for example from Init after a reconnect: no check, no error
	opt, err = n.credentialsOption()
	require.NoError(t, err)
	o := nats.GetDefaultOptions()
	require.NoError(t, opt(&o), "applying the option must not read the file")
	require.NotNil(t, o.UserJWT)
	require.NotNil(t, o.SignatureCB)

	// the connect attempt with the partly written file fails, NATS retries it
	_, err = o.SignatureCB([]byte("nonce"))
	require.Error(t, err)

	// a missing file during the rotation also fails only the connect attempt
	require.NoError(t, os.Remove(credsFile))
	_, err = o.UserJWT()
	require.Error(t, err)

	// after the rotation the callbacks use the new file
	newCreds := natstest.Creds(t)
	require.NoError(t, os.WriteFile(credsFile, newCreds, 0o600))
	userJWT, err := o.UserJWT()
	require.NoError(t, err)
	require.Contains(t, string(newCreds), userJWT)
	_, err = o.SignatureCB([]byte("nonce"))
	require.NoError(t, err)
}

func TestCredentialsOptionFirstCallFails(t *testing.T) {
	t.Setenv(_natsCredsFileEnv, filepath.Join(t.TempDir(), "missing.creds"))
	n := &storeregistry{}
	_, err := n.credentialsOption()
	require.Error(t, err)
}

func TestCredentialsOptionContent(t *testing.T) {
	t.Setenv(_natsCredsEnv, string(natstest.Creds(t)))
	n := &storeregistry{}
	for i := 0; i < 2; i++ {
		opt, err := n.credentialsOption()
		require.NoError(t, err)
		o := nats.GetDefaultOptions()
		require.NoError(t, opt(&o))
		_, err = o.SignatureCB([]byte("nonce"))
		require.NoError(t, err)
	}
}

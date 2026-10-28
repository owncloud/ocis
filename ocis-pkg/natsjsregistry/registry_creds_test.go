package natsjsregistry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
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

func newTestCreds(t *testing.T) []byte {
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
	return creds
}

func TestCredentialsOptionChecksOnlyOnce(t *testing.T) {
	creds := newTestCreds(t)
	credsFile := filepath.Join(t.TempDir(), "user.creds")
	require.NoError(t, os.WriteFile(credsFile, creds, 0o600))
	t.Setenv(_natsCredsFileEnv, credsFile)

	n := &storeregistry{}

	// first call: the file is checked
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
	newCreds := newTestCreds(t)
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
	t.Setenv(_natsCredsEnv, string(newTestCreds(t)))
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

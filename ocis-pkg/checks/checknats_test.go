package checks_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/owncloud/ocis/v2/ocis-pkg/checks"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/stretchr/testify/require"
)

// startOperatorServer starts a NATS server in operator (JWT) mode.
// It returns the client URL and the ".creds" content of a user.
func startOperatorServer(t *testing.T) (string, []byte) {
	t.Helper()

	operatorKP, err := nkeys.CreateOperator()
	require.NoError(t, err)
	operatorPub, err := operatorKP.PublicKey()
	require.NoError(t, err)
	operatorJWT, err := jwt.NewOperatorClaims(operatorPub).Encode(operatorKP)
	require.NoError(t, err)
	operator, err := jwt.DecodeOperatorClaims(operatorJWT)
	require.NoError(t, err)

	accountKP, err := nkeys.CreateAccount()
	require.NoError(t, err)
	accountPub, err := accountKP.PublicKey()
	require.NoError(t, err)
	accountJWT, err := jwt.NewAccountClaims(accountPub).Encode(operatorKP)
	require.NoError(t, err)

	resolver := &server.MemAccResolver{}
	require.NoError(t, resolver.Store(accountPub, accountJWT))

	s, err := server.NewServer(&server.Options{
		Host:             "127.0.0.1",
		Port:             -1,
		TrustedOperators: []*jwt.OperatorClaims{operator},
		AccountResolver:  resolver,
		NoLog:            true,
		NoSigs:           true,
	})
	require.NoError(t, err)
	go s.Start()
	t.Cleanup(s.Shutdown)
	require.True(t, s.ReadyForConnections(10*time.Second))

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

	return s.ClientURL(), creds
}

func TestNewNatsCheckWithConfig(t *testing.T) {
	url, creds := startOperatorServer(t)
	credsFile := filepath.Join(t.TempDir(), "user.creds")
	require.NoError(t, os.WriteFile(credsFile, creds, 0o600))

	t.Run("no credentials", func(t *testing.T) {
		require.Error(t, checks.NewNatsCheckWithConfig(url, stream.NatsConfig{})(context.Background()))
	})

	t.Run("username and password", func(t *testing.T) {
		require.Error(t, checks.NewNatsCheckWithConfig(url, stream.NatsConfig{
			AuthUsername: "user",
			AuthPassword: "pass",
		})(context.Background()))
	})

	t.Run("creds file", func(t *testing.T) {
		require.NoError(t, checks.NewNatsCheckWithConfig(url, stream.NatsConfig{
			CredsFile: credsFile,
		})(context.Background()))
	})

	t.Run("creds content", func(t *testing.T) {
		require.NoError(t, checks.NewNatsCheckWithConfig(url, stream.NatsConfig{
			Creds:        string(creds),
			AuthUsername: "ignored",
			AuthPassword: "ignored",
		})(context.Background()))
	})

	t.Run("invalid creds", func(t *testing.T) {
		err := checks.NewNatsCheckWithConfig(url, stream.NatsConfig{Creds: "not a creds file"})(context.Background())
		require.ErrorContains(t, err, "could not configure the nats connection")
	})

	t.Run("probes use the same connection", func(t *testing.T) {
		rotatedFile := filepath.Join(t.TempDir(), "user.creds")
		require.NoError(t, os.WriteFile(rotatedFile, creds, 0o600))
		check := checks.NewNatsCheckWithConfig(url, stream.NatsConfig{CredsFile: rotatedFile})
		require.NoError(t, check(context.Background()))

		// A rotation writes the file. The connection is open, so the probe does not read the file.
		require.NoError(t, os.WriteFile(rotatedFile, creds[:len(creds)/3], 0o600))
		require.NoError(t, check(context.Background()))
	})

	t.Run("probe retries invalid creds", func(t *testing.T) {
		laterFile := filepath.Join(t.TempDir(), "user.creds")
		check := checks.NewNatsCheckWithConfig(url, stream.NatsConfig{CredsFile: laterFile})
		require.ErrorContains(t, check(context.Background()), "could not configure the nats connection")

		require.NoError(t, os.WriteFile(laterFile, creds, 0o600))
		require.NoError(t, check(context.Background()))
	})
}

func TestNatsOptions(t *testing.T) {
	apply := func(t *testing.T, options []nats.Option) nats.Options {
		t.Helper()
		o := nats.GetDefaultOptions()
		for _, opt := range options {
			require.NoError(t, opt(&o))
		}
		return o
	}

	t.Run("no settings", func(t *testing.T) {
		options, err := checks.NatsOptions(stream.NatsConfig{})
		require.NoError(t, err)
		require.Empty(t, options)
	})

	t.Run("username and password", func(t *testing.T) {
		options, err := checks.NatsOptions(stream.NatsConfig{AuthUsername: "user", AuthPassword: "pass"})
		require.NoError(t, err)
		o := apply(t, options)
		require.Equal(t, "user", o.User)
		require.Equal(t, "pass", o.Password)
		require.Nil(t, o.UserJWT)
	})

	t.Run("tls", func(t *testing.T) {
		options, err := checks.NatsOptions(stream.NatsConfig{EnableTLS: true, TLSInsecure: true})
		require.NoError(t, err)
		o := apply(t, options)
		require.True(t, o.Secure)
		require.NotNil(t, o.TLSConfig)
		require.True(t, o.TLSConfig.InsecureSkipVerify)
	})
}

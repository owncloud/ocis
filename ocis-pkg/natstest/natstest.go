// Package natstest has helpers for tests with NATS JWT credentials.
package natstest

import (
	"testing"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nkeys"
	"github.com/stretchr/testify/require"
)

// Creds returns the content of a ".creds" file for a new user in a new account.
// No server knows the account, so use the content only for the checks of the format.
func Creds(t *testing.T) []byte {
	t.Helper()
	accountKP, err := nkeys.CreateAccount()
	require.NoError(t, err)
	return UserCreds(t, accountKP)
}

// UserCreds returns the content of a ".creds" file for a new user of the account accountKP.
func UserCreds(t *testing.T, accountKP nkeys.KeyPair) []byte {
	t.Helper()
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

// StartOperatorServer starts a NATS server in operator (JWT) mode.
// It returns the server and the ".creds" content of a user.
func StartOperatorServer(t *testing.T) (*server.Server, []byte) {
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

	return s, UserCreds(t, accountKP)
}

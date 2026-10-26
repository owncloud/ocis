// Package natsauth builds the JWT authentication options for NATS connections.
//
// A NATS server in operator (JWT) mode authenticates a client with a user JWT
// and an NKey seed. Both are stored in a ".creds" file. The package accepts the
// path of such a file or its content.
package natsauth

import (
	"errors"
	"fmt"
	"os"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

// CredentialsOption returns the nats.Option that authenticates a connection with
// the given JWT credentials. creds is the content of a ".creds" file and credsFile
// is the path of a ".creds" file. creds takes precedence over credsFile.
//
// A file is read again on each connect, so a rotated JWT is used without a restart.
// The content is fixed, so a rotated JWT needs a restart.
//
// CredentialsOption returns nil and no error when both values are empty.
// It returns an error when the credentials cannot be read or parsed.
// The error never contains the seed.
func CredentialsOption(credsFile, creds string) (nats.Option, error) {
	switch {
	case creds != "":
		if err := Validate([]byte(creds)); err != nil {
			return nil, fmt.Errorf("invalid nats creds: %w", err)
		}
		return nats.UserCredentialBytes([]byte(creds)), nil
	case credsFile != "":
		content, err := os.ReadFile(credsFile)
		if err != nil {
			return nil, fmt.Errorf("cannot read nats creds file: %w", err)
		}
		if err := Validate(content); err != nil {
			return nil, fmt.Errorf("invalid nats creds file %s: %w", credsFile, err)
		}
		return nats.UserCredentials(credsFile), nil
	default:
		return nil, nil
	}
}

// Check checks the JWT credentials like CredentialsOption. Use it in a constructor
// that can return an error, so that invalid credentials stop the service at startup.
// It returns no error when both values are empty.
func Check(credsFile, creds string) error {
	_, err := CredentialsOption(credsFile, creds)
	return err
}

// Validate checks that content is a ".creds" file with a user JWT and the
// matching user NKey seed.
func Validate(content []byte) error {
	rawJWT, err := nkeys.ParseDecoratedJWT(content)
	if err != nil {
		return fmt.Errorf("cannot parse user jwt: %w", err)
	}
	claims, err := jwt.DecodeUserClaims(rawJWT)
	if err != nil {
		return fmt.Errorf("cannot decode user jwt: %w", err)
	}

	kp, err := nkeys.ParseDecoratedUserNKey(content)
	if err != nil {
		return fmt.Errorf("cannot parse user nkey seed: %w", err)
	}
	defer kp.Wipe()
	pub, err := kp.PublicKey()
	if err != nil {
		return fmt.Errorf("cannot get public key of user nkey seed: %w", err)
	}
	if pub != claims.Subject {
		return errors.New("user nkey seed does not match the user jwt")
	}
	return nil
}

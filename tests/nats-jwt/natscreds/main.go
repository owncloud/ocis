// Command natscreds writes a NATS server config in operator (JWT) mode and
// the user ".creds" files for an oCIS test.
//
// Usage: go run ./tests/nats-jwt/natscreds -out <dir> [-port 9233]
//
// Output in <dir>:
//   - server.conf: nats-server config with JetStream and a memory resolver
//   - user.creds:  user of the account "ocis" with JetStream
//   - other.creds: user of the account "other" with JetStream
//
// The command creates new keys at each run. The repository never contains a seed.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

type account struct {
	pub  string
	jwt  string
	user []byte // content of the .creds file of a user of the account
}

func must[T any](v T, err error) T {
	if err != nil {
		fmt.Fprintln(os.Stderr, "natscreds:", err)
		os.Exit(1)
	}
	return v
}

// newAccount returns an account with JetStream and one user.
func newAccount(name string, operatorKP nkeys.KeyPair) account {
	kp := must(nkeys.CreateAccount())
	pub := must(kp.PublicKey())
	claims := jwt.NewAccountClaims(pub)
	claims.Name = name
	claims.Limits.JetStreamLimits = jwt.JetStreamLimits{MemoryStorage: -1, DiskStorage: -1, Streams: -1, Consumer: -1}

	userKP := must(nkeys.CreateUser())
	userPub := must(userKP.PublicKey())
	userClaims := jwt.NewUserClaims(userPub)
	userClaims.Name = name
	userJWT := must(userClaims.Encode(kp))

	return account{
		pub:  pub,
		jwt:  must(claims.Encode(operatorKP)),
		user: must(jwt.FormatUserConfig(userJWT, must(userKP.Seed()))),
	}
}

func main() {
	out := flag.String("out", ".", "output folder")
	port := flag.Int("port", 9233, "client port of the NATS server")
	flag.Parse()

	operatorKP := must(nkeys.CreateOperator())
	operatorPub := must(operatorKP.PublicKey())

	sysKP := must(nkeys.CreateAccount())
	sysPub := must(sysKP.PublicKey())
	sysClaims := jwt.NewAccountClaims(sysPub)
	sysClaims.Name = "SYS"
	sysJWT := must(sysClaims.Encode(operatorKP))

	ocis := newAccount("ocis", operatorKP)
	other := newAccount("other", operatorKP)

	operatorClaims := jwt.NewOperatorClaims(operatorPub)
	operatorClaims.Name = "ci"
	operatorClaims.SystemAccount = sysPub
	operatorJWT := must(operatorClaims.Encode(operatorKP))

	conf := fmt.Sprintf(`port: %d
jetstream {
  store_dir: "/data"
}
operator: %s
system_account: %s
resolver: MEMORY
resolver_preload: {
  %s: %s
  %s: %s
  %s: %s
}
`, *port, operatorJWT, sysPub, sysPub, sysJWT, ocis.pub, ocis.jwt, other.pub, other.jwt)

	must(0, os.MkdirAll(*out, 0o700))
	// server.conf has JWTs but no seed. The container user can differ from the runner user.
	must(0, os.WriteFile(filepath.Join(*out, "server.conf"), []byte(conf), 0o644)) //nolint:gosec
	must(0, os.WriteFile(filepath.Join(*out, "user.creds"), ocis.user, 0o600))
	must(0, os.WriteFile(filepath.Join(*out, "other.creds"), other.user, 0o600))
	fmt.Println("wrote server.conf, user.creds, and other.creds to", *out)
}

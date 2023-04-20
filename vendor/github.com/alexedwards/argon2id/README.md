# Argon2id

This package provides a convenience wrapper around Go's [argon2](https://pkg.go.dev/golang.org/x/crypto/argon2?tab=doc) implementation, making it simpler to securely hash and verify passwords using Argon2.

It enforces use of the Argon2id algorithm variant and cryptographically-secure random salts.

## Usage

```go
package main

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func main() {
	// CreateHash returns a Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	hash, err := argon2id.CreateHash("pa$$word", argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash("pa$$word", hash)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Match: %v", match)
}
```

### Changing the Parameters

When creating a hash you can and should configure the parameters to be suitable for the environment that the code is running in. The parameters are:

* Memory — The amount of memory used by the Argon2 algorithm (in kibibytes).
* Iterations — The number of iterations (or passes) over the memory.
* Parallelism — The number of threads (or lanes) used by the algorithm.
* Salt length — Length of the random salt. 16 bytes is recommended for password hashing.
* Key length — Length of the generated key (or password hash). 16 bytes or more is recommended.

The Memory and Iterations parameters control the computational cost of hashing the password. The higher these figures are, the greater the cost of generating the hash and the longer the runtime. It also follows that the greater the cost will be for any attacker trying to guess the password.

If the code is running on a machine with multiple cores, then you can decrease the runtime without reducing the cost by increasing the Parallelism parameter. This controls the number of threads that the work is spread across. Important note: Changing the value of the Parallelism parameter changes the hash output.

```go
params := &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  4,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}


hash, err := argon2id.CreateHash("pa$$word", params)
if err != nil {
	log.Fatal(err)
}
```

For guidance and an outline process for choosing appropriate parameters see https://tools.ietf.org/html/draft-irtf-cfrg-argon2-04#section-4.

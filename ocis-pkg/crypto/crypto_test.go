package crypto_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/ocis-pkg/crypto"
	"github.com/owncloud/ocis/ocis-pkg/log"

	. "github.com/onsi/ginkgo"
	cfg "github.com/owncloud/ocis/ocis-pkg/config"
)

var _ = Describe("Crypto", func() {
	var (
		userConfigDir string
		err           error
		config        = cfg.New()
	)

	BeforeEach(func() {
		userConfigDir, err = os.UserConfigDir()
		if err != nil {
			Fail(err.Error())
		}
		config.Proxy.HTTP.TLSKey = filepath.Join(userConfigDir, "ocis", "server.key")
		config.Proxy.HTTP.TLSCert = filepath.Join(userConfigDir, "ocis", "server.cert")
	})

	AfterEach(func() {
		if err := os.RemoveAll(filepath.Join(userConfigDir, "ocis")); err != nil {
			Fail(err.Error())
		}
	})

	// This little test should nail down the main functionality of this package, which is providing with a default location
	// for the key / certificate pair in case none is configured. Regardless of how the values ended in the configuration,
	// the side effects of GenCert is what we want to test.
	Describe("Creating key / certificate pair", func() {
		Context("For ocis-proxy in the location of the user config directory", func() {
			It(fmt.Sprintf("Creates the cert / key tuple in: %s", filepath.Join(userConfigDir, "ocis")), func() {
				if err := crypto.GenCert(config.Proxy.HTTP.TLSCert, config.Proxy.HTTP.TLSKey, log.NewLogger()); err != nil {
					Fail(err.Error())
				}

				if _, err := os.Stat(filepath.Join(userConfigDir, "ocis", "server.key")); err != nil {
					Fail("key not found at the expected location")
				}

				if _, err := os.Stat(filepath.Join(userConfigDir, "ocis", "server.cert")); err != nil {
					Fail("certificate not found at the expected location")
				}
			})
		})
	})
})

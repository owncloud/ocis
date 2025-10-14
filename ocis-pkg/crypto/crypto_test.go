package crypto_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"

	. "github.com/onsi/ginkgo/v2"
	cfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
)

var _ = Describe("Crypto", func() {
	var (
		userConfigDir string
		err           error
		config        = cfg.DefaultConfig()
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
				if err := crypto.GenCert(config.Proxy.HTTP.TLSCert, config.Proxy.HTTP.TLSKey, log.NopLogger()); err != nil {
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
	Describe("Creating a new cert pool", func() {
		var (
			crtOne string
			keyOne string
			crtTwo string
			keyTwo string
		)
		BeforeEach(func() {
			crtOne = filepath.Join(userConfigDir, "ocis/one.cert")
			keyOne = filepath.Join(userConfigDir, "ocis/one.key")
			crtTwo = filepath.Join(userConfigDir, "ocis/two.cert")
			keyTwo = filepath.Join(userConfigDir, "ocis/two.key")
			if err := crypto.GenCert(crtOne, keyOne, log.NopLogger()); err != nil {
				Fail(err.Error())
			}
			if err := crypto.GenCert(crtTwo, keyTwo, log.NopLogger()); err != nil {
				Fail(err.Error())
			}
		})
		It("handles one certificate", func() {
			f1, _ := os.Open(crtOne)
			defer f1.Close()

			c, err := crypto.NewCertPoolFromPEM(f1)
			if err != nil {
				Fail(err.Error())
			}
			if len(c.Subjects()) != 1 {
				Fail("expected 1 certificate in the cert pool")
			}
		})
		It("handles multiple certificates", func() {
			f1, _ := os.Open(crtOne)
			f2, _ := os.Open(crtTwo)
			defer f1.Close()
			defer f2.Close()

			c, err := crypto.NewCertPoolFromPEM(f1, f2)
			if err != nil {
				Fail(err.Error())
			}
			if len(c.Subjects()) != 2 {
				Fail("expected 2 certificates in the cert pool")
			}
		})
	})
})

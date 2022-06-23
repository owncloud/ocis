package ldap

import (
	"crypto/x509"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

const (
	caCheckRetries = 3
	caCheckSleep   = 2
)

func WaitForCA(log log.Logger, insecure bool, caCert string) error {
	if !insecure && caCert != "" {
		for i := 0; i < caCheckRetries; i++ {
			if _, err := os.Stat(caCert); err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			// Check if this actually is a CA cert. We need to retry here as well
			// as the file might exist already, but have no contents yet.
			certs := x509.NewCertPool()
			pemData, err := ioutil.ReadFile(caCert)
			if err != nil {
				log.Debug().Err(err).Str("LDAP CACert", caCert).Msg("Error reading CA")
			} else if !certs.AppendCertsFromPEM(pemData) {
				log.Debug().Str("LDAP CAcert", caCert).Msg("Failed to append CA to pool")
			} else {
				return nil
			}
			time.Sleep(caCheckSleep * time.Second)
			log.Warn().Str("LDAP CACert", caCert).Msgf("CA cert file is not ready yet. Waiting %d seconds for it to appear.", caCheckSleep)
		}
	}
	return nil
}

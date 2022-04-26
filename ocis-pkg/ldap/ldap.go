package ldap

import (
	"errors"
	"os"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

const _caTimeout = 5

func WaitForCA(log log.Logger, insecure bool, caCert string) error {
	if !insecure && caCert != "" {
		if _, err := os.Stat(caCert); errors.Is(err, os.ErrNotExist) {
			log.Warn().Str("LDAP CACert", caCert).Msgf("File does not exist. Waiting %d seconds for it to appear.", _caTimeout)
			time.Sleep(_caTimeout * time.Second)
			if _, err := os.Stat(caCert); errors.Is(err, os.ErrNotExist) {
				log.Warn().Str("LDAP CACert", caCert).Msgf("File still does not exist after Timeout")
				return err
			}
		}
	}
	return nil
}

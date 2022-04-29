package shared

import (
	"fmt"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func MissingMachineAuthApiKeyError(service string) error {
	return fmt.Errorf("machine_auth_api_key has not your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting OCIS_MACHINE_AUTH_API_KEY).\n",
		service, defaults.BaseConfigPath())
}

func MissingJWTTokenError(service string) error {
	return fmt.Errorf("jwt_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting OCIS_JWT_SECRET).\n",
		service, defaults.BaseConfigPath())
}

func MissingRevaTransferSecretError(service string) error {
	return fmt.Errorf("transfer_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting STORAGE_TRANSFER_SECRET).\n",
		service, defaults.BaseConfigPath())
}

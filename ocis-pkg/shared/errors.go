package shared

import (
	"fmt"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func MissingMachineAuthAPIKey(service string) error {
	return fmt.Errorf("The Machineauth API key has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingJWTToken(service string) error {
	return fmt.Errorf("jwt_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingRevaTransferSecret(service string) error {
	return fmt.Errorf("transfer_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingLDAPBindPassword(service string) error {
	return fmt.Errorf("bind_password has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingServiceUserPassword(service, serviceUser string) error {
	return fmt.Errorf("password of service user %s has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		serviceUser, service, defaults.BaseConfigPath())
}

func MissingSystemUserID(service string) error {
	return fmt.Errorf("The system user ID has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingSystemAuthAPIKey(service string) error {
	return fmt.Errorf("The system auth API key has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

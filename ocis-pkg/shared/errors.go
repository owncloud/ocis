package shared

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
)

func MissingMachineAuthApiKeyError(service string) error {
	return fmt.Errorf("The Machineauth API key has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingSystemUserApiKeyError(service string) error {
	return fmt.Errorf("The SystemUser API key has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingJWTTokenError(service string) error {
	return fmt.Errorf("The jwt_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingRevaTransferSecretError(service string) error {
	return fmt.Errorf("The transfer_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingLDAPBindPassword(service string) error {
	return fmt.Errorf("The ldap bind_password has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingServiceUserPassword(service, serviceUser string) error {
	return fmt.Errorf("The password of service user %s has not been set properly in your config for %s. "+
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

func MissingAdminUserID(service string) error {
	return fmt.Errorf("The admin user ID has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingServiceAccountID(service string) error {
	return fmt.Errorf("The service account id has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingServiceAccountSecret(service string) error {
	return fmt.Errorf("The service account secret has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingWOPISecretError(service string) error {
	return fmt.Errorf("The WOPI secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

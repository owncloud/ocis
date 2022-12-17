package shared

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
)

// MissingMachineAuthApiKeyError
// FIXME: nolint
// nolint: revive
func MissingMachineAuthApiKeyError(service string) error {
	return fmt.Errorf("The Machineauth API key has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

// MissingSystemUserApiKeyError
// FIXME: nolint
// nolint: revive
func MissingSystemUserApiKeyError(service string) error {
	return fmt.Errorf("The SystemUser API key has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

// MissingJWTTokenError
// FIXME: nolint
// nolint: revive
func MissingJWTTokenError(service string) error {
	return fmt.Errorf("The jwt_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

// MissingRevaTransferSecretError
// FIXME: nolint
// nolint: revive
func MissingRevaTransferSecretError(service string) error {
	return fmt.Errorf("The transfer_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

// MissingLDAPBindPassword
// FIXME: nolint
// nolint: revive
func MissingLDAPBindPassword(service string) error {
	return fmt.Errorf("The ldap bind_password has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

// MissingServiceUserPassword
// FIXME: nolint
// nolint: revive
func MissingServiceUserPassword(service, serviceUser string) error {
	return fmt.Errorf("The password of service user %s has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		serviceUser, service, defaults.BaseConfigPath())
}

// MissingSystemUserID
// FIXME: nolint
// nolint: revive
func MissingSystemUserID(service string) error {
	return fmt.Errorf("The system user ID has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

// MissingAdminUserID
// FIXME: nolint
// nolint: revive
func MissingAdminUserID(service string) error {
	return fmt.Errorf("The admin user ID has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by running ocis init or setting it manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

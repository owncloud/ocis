package shared

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/reva/v2/pkg/natsauth"
	revastore "github.com/owncloud/reva/v2/pkg/store"
)

func MissingMachineAuthApiKeyError(service string) error {
	return fmt.Errorf("The Machineauth API key has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingSystemUserApiKeyError(service string) error {
	return fmt.Errorf("The SystemUser API key has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingJWTTokenError(service string) error {
	return fmt.Errorf("The jwt_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingRevaTransferSecretError(service string) error {
	return fmt.Errorf("The transfer_secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingLDAPBindPassword(service string) error {
	return fmt.Errorf("The ldap bind_password has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingServiceUserPassword(service, serviceUser string) error {
	return fmt.Errorf("The password of service user %s has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		serviceUser, service, defaults.BaseConfigPath())
}

func MissingSystemUserID(service string) error {
	return fmt.Errorf("The system user ID has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingAdminUserID(service string) error {
	return fmt.Errorf("The admin user ID has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingServiceAccountID(service string) error {
	return fmt.Errorf("The service account id has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingServiceAccountSecret(service string) error {
	return fmt.Errorf("The service account secret has not been configured for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

func MissingWOPISecretError(service string) error {
	return fmt.Errorf("The WOPI secret has not been set properly in your config for %s. "+
		"Make sure your %s config contains the proper values "+
		"(e.g. by using 'ocis init --diff' and applying the patch or setting a value manually in "+
		"the config/corresponding environment variable).",
		service, defaults.BaseConfigPath())
}

// CheckNatsCredentials checks the NATS JWT credentials of a store at startup. The reva
// store cannot return an error for invalid credentials. It only logs the error, and each
// store operation fails later. It returns no error when both values are empty.
// Only the nats-js and nats-js-kv stores use the credentials. For other store types,
// it does not check the credentials.
func CheckNatsCredentials(service, store, credsFile, creds string) error {
	if store != revastore.TypeNatsJS && store != revastore.TypeNatsJSKV {
		return nil
	}
	if err := natsauth.Check(credsFile, creds); err != nil {
		return fmt.Errorf("The NATS credentials for %s are not valid: %w", service, err)
	}
	return nil
}

# Auth-Service

The ocis Auth Service is used to authenticate service accounts. Compared to normal accounts, service accounts are ocis internal only and not available as ordinary users like via LDAP.

## The `auth` Service Family

ocis uses serveral authentication services for different use cases. All services that start with `auth-` are part of the authentication service family. Each member authenticates requests with different scopes. As of now, these services exist:
  -   `auth-app` handles authentication of external 3rd party apps
  -   `auth-basic` handles basic authentication
  -   `auth-bearer` handles oidc authentication
  -   `auth-machine` handles interservice authentication when a user is impersonated
  -   `auth-service` handles interservice authentication when using service accounts

## Service Accounts

Service accounts are user accounts that are only used for inter service communication. The users have no personal space, do not show up in user lists and cannot login via the UI. Service accounts can be configured in the settings service. Only the `admin` service user is available for now. Additionally to the actions it can do via its role, all service users can stat all files on all spaces.

## Configuring Service Accounts

By using the envvars `OCIS_SERVICE_ACCOUNT_ID` and `OCIS_SERVICE_ACCOUNT_SECRET`, one can configure the ID and the secret of the service user. The secret can be rotated regulary to increase security. For activating a new secret, all services where the envvars are used need to be restarted. The secret is always and only stored in memory and never written into any persistant store. Though you can use any string for the service account, it is recommmended to use a UUIDv4 string.

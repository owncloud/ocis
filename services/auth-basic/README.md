# Auth-Basic

The oCIS Auth Basic service provides basic authentication for those clients who cannot handle OpenID Connect. This should only be enabled for tests and development.

The `auth-basic` service is responsible for validating authentication of incoming requests. To do so, it will use the configured `auth manager`, see the `Auth Managers` section. Only HTTP basic auth requests to ocis will involve the `auth-basic` service.

To enable `auth-basic`, you first must set `PROXY_ENABLE_BASIC_AUTH` to `true`.

## The `auth` Service Family

ocis uses serveral authentication services for different use cases. All services that start with `auth-` are part of the authentication service family. Each member authenticates requests with different scopes. As of now, these services exist:
  -   `auth-app` handles authentication of external 3rd party apps
  -   `auth-basic` handles basic authentication
  -   `auth-bearer` handles oidc authentication
  -   `auth-machine` handles interservice authentication when a user is impersonated
  -   `auth-service` handles interservice authentication when using service accounts

## Auth Managers

Since the `auth-basic` service does not do any validation itself, it needs to be configured with an authentication manager. One can use the `AUTH_BASIC_AUTH_MANAGER` environment variable to configure this. Currently only one auth manager is supported: `"ldap"`

### LDAP Auth Manager

Setting `AUTH_BASIC_AUTH_MANAGER` to `"ldap"` will configure the `auth-basic` service to use LDAP as auth manager. This is the recommended option for running in a production and testing environment. More details on how to configure LDAP with ocis can be found in the admin docs.

### Other Auth Managers

oCIS currently supports no other auth manager

## Scalability

When using `"ldap"` as auth manager, there is no persistance as requests will just be forwarded to the LDAP server. Therefore, multiple instances of the `auth-basic` service can be started without further configuration. Be aware, that other auth managers might not allow that.


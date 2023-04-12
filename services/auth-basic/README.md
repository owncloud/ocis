# Auth-Basic

The oCIS Auth Basic service provides basic authentication for those clients who cannot handle OpenID Connect. This should only be enabled for tests and development.

The `auth-basic` service is responsible for validating authentication of incoming requests. To do so, it will use the configured `auth manager`, see the `Auth Managers` section. Only HTTP basic auth requests to ocis will involve the `auth-basic` service.

To enable `auth-basic`, you first must set `PROXY_ENABLE_BASIC_AUTH` to `true`.

## Auth Managers

Since the `auth-basic` service does not do any validation itself, it needs to be configured with an authentication manager. One can use the `AUTH_BASIC_AUTH_MANAGER` environment variable to configure this. Currently only one auth manager is supported: `"ldap"`

### LDAP Auth Manager

Setting `AUTH_BASIC_AUTH_MANAGER` to `"ldap"` will configure the `auth-basic` service to use LDAP as auth manager. This is the recommended option for running in a production and testing environment. More details on how to configure LDAP with ocis can be found in the admin docs.

### Other Auth Managers

oCIS currently supports no other auth manager

## Scalability

When using `"ldap"` as auth manager, there is no persistance as requests will just be forwarded to the LDAP server. Therefore, multiple instances of the `auth-basic` service can be started without further configuration. Be aware, that other auth managers might not allow that.


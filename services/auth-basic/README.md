# Auth-Basic Service

The `auth-basic` service is responsible for validating authentication of incoming requests. To do so, it will use the configured `auth manager`, see the `Auth Managers` section. Only HTTP basic auth requests to ocis will involve the `auth-basic` service.

## Auth Managers

Since the `auth-basic` service does not do any validation itself, it needs to be configured with an authentication manager. One can use the `AUTH_BASIC_AUTH_PROVIDER` environment variable to configure this.

### LDAP Auth Manager

Setting `AUTH_BASIC_AUTH_PROVIDER` to `"ldap"` will configure the `auth-basic` service to use LDAP as auth manager. This is the recommended option for running in a production and testing environment. More details on how to configure LDAP with ocis can be found in the admin docs.

### Other Auth Managers

The possible auth mangers which can be selected are `"ldap"` and `"owncloudsql"`. Those are tested and usable though `"ldap"` is the recommend manager. Refer to the admin docs for additional information about those.

## Scalability

Scalability, just like memory and cpu consumption, are highly dependent on the configured auth manager. When using the recommended one (`"ldap"`) there is no persistance as requests will just be 
forwarded to the ldap-server. Therefore multiple instances of the `auth-basic` service can be started without further configuration. Be aware that other auth managers might not allow that. `"json"`
auth managers for example persist to the disc and can therefore not be scaled to multiple instances easily.


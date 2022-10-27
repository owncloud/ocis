# Auth-Basic Service

The `auth-basic` service is responsible for validating authentication of incoming requests. To do so it will forward the call to an `auth manager`
(see `Auth Managers` section). Almost every request to the server will involve the `auth-basic` service as access tokens need to be verified.

## Auth Managers

Since the `auth-basic` service does not do any validation itself, it needs to be configured with an auth manager. One can use the `AUTH_BASIC_AUTH_PROVIDER` envvar to configure this.

### LDAP Auth Manager

Setting `AUTH_BASIC_AUTH_PROVIDER` to `"ldap"` will configure the `auth-basic` service to use ldap as auth manager. This is the recommended option for productive and testing deployments. 
More details on how to configure ldap with ocis can be found in the admin docs.

### Other Auth Managers

There are a number of other possible auth mangers, including `"json"` and `"owncloudsql"`. Those are tested and usable, but we recommend using ldap. 
Refer to the admin docs for additional information about these.

## Scalability

Scalability, just like memory and cpu consumption, are highly dependent on the configured auth manager. When using the recommended one (`"ldap"`) there is no persistance as requests will just be 
forwarded to the ldap-server. Therefore multiple instances of the `auth-basic` service can be started without further configuration. Be aware that other auth managers might not allow that. `"json"`
auth managers for example persist to the disc and can therefore not be scaled to multiple instances easily.


# Auth-Bearer

The oCIS Auth Bearer service communicates with the configured OpenID Connect identity provider to authenticate requests. OpenID Connect is the default authentication mechanism for all clients: web, desktop and mobile. Basic auth is only used for testing and has to be explicity enabled.

## The `auth` Service Family

ocis uses serveral authentication services for different use cases. All services that start with `auth-` are part of the authentication service family. Each member authenticates requests with different scopes. As of now, these services exist:
  -   `auth-basic` handles basic authentication
  -   `auth-bearer` handles oidc authentication
  -   `auth-machine` handles interservice authentication when a user is impersonated
  -   `auth-service` handles interservice authentication when using service accounts
  -   `auth-app` handles authentication of external 3rd party apps

## Built in OpenID Connect Identity Provider

A default oCIS deployment will start a [built in OpenID Connect identity provider](https://github.com/owncloud/ocis/tree/master/services/idp) but can be configured to use an external one as well.

## Scalability

There is no persistance or caching. The proxy caches verified auth bearer tokens. Requests will be forwarded to the identity provider. Therefore, multiple instances of the `auth-bearer` service can be started without further configuration. Currently, the auth registry used by the gateway can only use a single instance of the service. To use more than one auth provider per deployment you need to scale the gateway.

This will change when we use the service registry in more places and use micro clients to select an instance of a service.

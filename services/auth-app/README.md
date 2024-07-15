# Auth-App

The auth-app service provides authentication for 3rd party apps.

## The `auth` Service Family

ocis uses serveral authentication services for different use cases. All services that start with `auth-` are part of the authentication service family. Each member authenticates requests with different scopes. As of now, these services exist:
  -   `auth-basic` handles basic authentication
  -   `auth-bearer` handles oidc authentication
  -   `auth-machine` handles interservice authentication when a user is impersonated
  -   `auth-service` handles interservice authentication when using service accounts
  -   `auth-app` handles authentication of external 3rd party apps

## Optional Service

This service is an optional service that will not run with default settings. To start use it, two envvars need to be set:
```bash
OCIS_ADD_RUN_SERVICES=auth-app # to start the service. Alternatively you can start the service explicitly via the command line.
PROXY_ENABLE_APP_AUTH=true # to allow app authentication. This envvar goes to the proxy service in case of a distributed environment.
```

## App Tokens

App Tokens are used to authenticate 3rd party apps. To be able to use an app token, one must first create a token via cli.

```bash
ocis auth-app create --user-name={user-name} --expiration={token-expiration}
```

Once generated, these tokens can be used to authenticate requests to the oCIS services. They can be passed in any request as `Basic Auth` header.

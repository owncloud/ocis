# Auth-Machine

The oCIS Auth Machine is used for interservice communication when using user impersonation.

ocis uses serveral authentication services for different use cases. All services that start with `auth-` are part of the authentication service family. Each member authenticates requests with different scopes. As of now, these services exist:
  -   `auth-app` handles authentication of external 3rd party apps
  -   `auth-basic` handles basic authentication
  -   `auth-bearer` handles oidc authentication
  -   `auth-machine` handles interservice authentication when a user is impersonated
  -   `auth-service` handles interservice authentication when using service accounts

## User Impersonation

When one ocis service is trying to talk to other ocis services, it needs to authenticate itself. To do so, it will impersonate a user using the `auth-machine` service. It will then act on behalf of this user. Any action will show up as action of this specific user, which gets visible when e.g. logged in the audit log.

## Deprecation

With the upcoming `auth-service` service, the `auth-machine` service will be used less frequently and is probably a candidate for deprecation.

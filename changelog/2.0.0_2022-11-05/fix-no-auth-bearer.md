Bugfix: Don't run auth-bearer service by default

We no longer start the auth-bearer service by default. This service is
currently unused and not required to run ocis. The equivalent functionality
to verify OpenID connect tokens and to mint reva tokes for OIDC authenticated
clients is currently implemented inside the oidc-auth middleware of the proxy.

https://github.com/owncloud/ocis/issues/4692

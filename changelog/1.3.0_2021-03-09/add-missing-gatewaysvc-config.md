Bugfix: Add missing gateway config

The auth provider `ldap` and `oidc` drivers now need to be able talk to the reva gateway. We added the `gatewayscv` to the config that is passed to reva.

https://github.com/owncloud/ocis/pull/1716

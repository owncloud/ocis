Change: default to running behind ocis-proxy

We changed the default configuration to integrate better with ocis.

- We use ocis-glauth as the default ldap server on port 9125 with base `dc=example,dc=org`.
- We use a dedicated technical `reva` user to make ldap binds
- Clients are supposed to use the ocis-proxy endpoint `https://localhost:9200`
- We removed unneeded ocis configuration from the frontend which no longer serves an oidc provider.
- We changed the default user OpaqueID attribute from `sub` to `preferred_username`. The latter is a claim populated by konnectd that can also be used by the reva ldap user manager to look up users by their OpaqueId

https://github.com/owncloud/ocis/ocis-revapull/113

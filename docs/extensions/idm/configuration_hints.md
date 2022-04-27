---
title: Configuration Hints
date: 2022-04-27:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/idm
geekdocFilePath: configuration_hints.md
geekdocCollapseSection: true
---

## TLS Server Certificates
By default IDM generates a self-signed certificate and key on startup to be
able to provide TLS protected services. The certificate is stored in
`idm/ldap.crt` inside the oCIS base data directory. The key is in
`idm/ldap.key` in the same directory. You configure custom a custom server
certificate by setting the `IDM_LDAPS_CERT` and `IDM_LDAPS_KEY`.

## Default / Demo Users
On startup IDM creates a set of default services users, that are needed
internally to provide other oCIS service access to IDM. These users are stored
in a separate subtree. The base DN of that subtree is:
`ou=sysusers,o=libregraph-idm`. The service users are:

* `uid=libregraph,ou=sysusers,o=libregraph-idm`: This is the only user with write
  access to the LDAP tree. It is used by the Graph service to lookup, create, delete
  modify users and groups.
* `uid=idp,ou=sysusers,o=libregraph-idm`: This user is used by the IDP service to
  perform user lookups for authentication.
* `uid=reva,ou=sysusers,o=libregraph-idm`: This user is used by the "reva" services
  "user, group and auth-basic.

IDM is also able to create [Demo Users](../../../ocis/getting-started/demo-users)
upon startup. 

## Access via LDAP command line tools
For testing purposes it is sometimes helpful to query IDM using the ldap
command line clients. To e.g. list all user can use this command:

```
ldapsearch -x -H ldaps://127.0.0.1:9235 -x -D uid=libregraph,ou=sysusers,o=libregraph-idm -w idm -b o=libregraph-idm objectclass=inetorgperson
```

When using the default configuration with the self-signed server certificate
you might need to switch of Certificate Validation the `LDAPTL_REQCERT` env
variable:

```
LDAPTLS_REQCERT=never ldapsearch -x -H ldaps://127.0.0.1:9235 -x -D uid=libregraph,ou=sysusers,o=libregraph-idm -w idm -b o=libregraph-idm objectclass=inetorgperson
```

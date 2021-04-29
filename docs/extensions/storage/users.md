---
title: "Users"
date: 2020-01-16T00:00:00+00:00
weight: 35
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: users.md
---

### Demo driver

This is a simple user driver for testing. It contains three users:
```
einstein:relativity
marie:radioactivity
richard:superfluidity
```
In order to use the `demo` driver you need to export the relevant environment variable:
```
export STORAGE_USERS_DRIVER=demo
```

### JSON driver

In order to switch from the `ldap` driver to JSON based users you need to export the relevant environment variables:
```
export STORAGE_USERS_DRIVER=json
export STORAGE_USERS_JSON=/path/to/users.json
```

For the format of the users.json have a look at the [reva examples](https://github.com/cs3org/reva/blob/master/examples/oc-phoenix/users.demo.json)

### LDAP driver

This is the default user driver.

If the below defaults don't match your environment change them accordingly:
```
export STORAGE_LDAP_HOSTNAME=localhost
export STORAGE_LDAP_PORT=9126
export STORAGE_LDAP_BASE_DN='dc=example,dc=org'
export STORAGE_LDAP_USERFILTER='(&(objectclass=posixAccount)(cn=%s))'
export STORAGE_LDAP_GROUPFILTER='(&(objectclass=posixGroup)(cn=%s))'
export STORAGE_LDAP_BIND_DN='cn=reva,ou=sysusers,dc=example,dc=org'
export STORAGE_LDAP_BIND_PASSWORD=reva
export STORAGE_LDAP_SCHEMA_UID=uid
export STORAGE_LDAP_SCHEMA_MAIL=mail
export STORAGE_LDAP_SCHEMA_DISPLAYNAME=sn
export STORAGE_LDAP_SCHEMA_CN=cn
```

Then restart the `bin/storage users` and `bin/storage auth-basic` services for the changes to take effect.

---
title: "Users"
date: 2020-01-16T00:00:00+00:00
anchor: "users"
weight: 35
---

### Demo driver

This is the default user driver.It contains three users:
```
einstein:relativity
marie:radioactivty
richard:superfluidity
```

### JSON driver

In order to switch from the `demo` driver to JSON based users you need to export the relevant environment variables:
```
export REVA_USERS_DRIVER=json
export REVA_USERS_JSON=/path/to/users.json
```

For the format of the users.json have a look at the [reva examples](https://github.com/cs3org/reva/blob/master/examples/separate/users.demo.json)

### LDAP driver

In order to switch from the `demo` driver to LDAP you need to export the relevant environment variable:
```
export REVA_USERS_DRIVER=ldap
```

If the below defaults don't match your environment change them accordingly:
```
export REVA_LDAP_HOSTNAME=localhost
export REVA_LDAP_PORT=636
export REVA_LDAP_BASE_DN='dc=owncloud,dc=com'
export REVA_LDAP_USERFILTER='(&(objectclass=posixAccount)(cn=%s))'
export REVA_LDAP_GROUPFILTER='(&(objectclass=posixGroup)(cn=%s))'
export REVA_LDAP_BIND_DN='cn=admin,dc=owncloud,dc=com'
export REVA_LDAP_BIND_PASSWORD=admin
export REVA_LDAP_SCHEMA_UID=uid
export REVA_LDAP_SCHEMA_MAIL=mail
export REVA_LDAP_SCHEMA_DISPLAYNAME=displayName
export REVA_LDAP_SCHEMA_CN=cn
```

Then restart the `bin/ocis-reva users` and `bin/ocis-reva auth-basic` services for the changes to take effect.

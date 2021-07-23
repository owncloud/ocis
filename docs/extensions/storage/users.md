---
title: "Users"
date: 2020-01-16T00:00:00+00:00
weight: 17
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: users.md
---

TODO add this to the storage overview? or is this a different part? That should be started as a separate service ? And documented elsewhere, eg. in the accounts?

### User and Group provisioning

In oc10 users are identified by a username, which cannot change, because it is used as a foreign key in several tables. For oCIS we are internally identifying users by a UUID, while using the username in the WebDAV and OCS APIs for backwards compatability. To distinguish this in the URLs we are using `<username>` instead of `<userid>`. You may have encountered `<userlayout>`, which refers to a template that can be configured to build several path segments by filling in user properties, e.g. the first  character of the username (`{{substr 0 1 .Username}}/{{.Username}}`), the identity provider (`{{.Id.Idp}}/{{.Username}}`) or the email (`{{.Mail}}`)

{{< hint warning >}}
Make no mistake, the [OCS Provisioning API](https://doc.owncloud.com/server/developer_manual/core/apis/provisioning-api.html) uses `userid` while it actually is the username, because it is what you use to login. 
{{< /hint >}}

We are currently working on adding [user management through the CS3 API](https://github.com/owncloud/ocis/pull/1930) to handle user and group provisioning (and deprovisioning).

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
export STORAGE_LDAP_USER_SCHEMA_UID=uid
export STORAGE_LDAP_USER_SCHEMA_MAIL=mail
export STORAGE_LDAP_USER_SCHEMA_DISPLAYNAME=sn
export STORAGE_LDAP_USER_SCHEMA_CN=cn
```

Then restart the `bin/storage users` and `bin/storage auth-basic` services for the changes to take effect.

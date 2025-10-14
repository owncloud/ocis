---
title: Resetting a lost administrator password
date: 2022-08-29:00:00+00:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/idm
geekdocFilePath: admin_password_reset.md
geekdocCollapseSection: true
---

## Resetting a lost administrator password
By default, when using oCIS with the builtin IDM an ad generates the
user `admin` (DN `uid=admin,ou=users,o=libregraph-idm`) if, for any
reason, the password of that user is lost, it can be reset using
the `resetpassword` sub-command:

```
ocis idm resetpassword
```

It will prompt for a new password and set the password of that user
accordingly. Note: As this command is accessing the idm database directly
will only work while ocis is not running and nothing else is accessing
database.

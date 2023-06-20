Bugfix: Don't connect to ldap on startup

This leads to misleading error messages. Instead we connect on first request

https://github.com/owncloud/ocis/pull/6565

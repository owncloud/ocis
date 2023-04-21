Bugfix: Fix retry handling for LDAP connections

We've fixed the handling of network issues (e.g. connection loss) during LDAP Write Operations
to correctly retry the request.

https://github.com/owncloud/ocis/issues/2974

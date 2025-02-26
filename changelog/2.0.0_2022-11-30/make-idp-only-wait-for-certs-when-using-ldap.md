Bugfix: Make IDP only wait for certs when using LDAP

When configuring cs3 as the backend the IDP no longer waits for an LDAP certificate to appear.

https://github.com/owncloud/ocis/pull/3965

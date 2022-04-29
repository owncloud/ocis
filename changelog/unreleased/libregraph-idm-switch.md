Change: Switched default configuration to use libregraph/idm

We switched the default configuration of oCIS to use the "idm" service (based
on libregraph/idm) as the standard source for user and group information. The
accounts and glauth services are no longer enabled by default and will be
removed with an upcoming release.

https://github.com/owncloud/ocis/pull/3331
https://github.com/owncloud/ocis/pull/3633

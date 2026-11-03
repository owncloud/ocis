Change: Remove the owncloudsql driver

We removed the deprecated `owncloudsql` driver from oCIS.

The services `storage-users`, `auth-basic`, `users`, `groups`, and `sharing` no
longer accept the value `owncloudsql` for their driver. The related
configuration sections and environment variables are removed.

https://github.com/owncloud/ocis/pull/13189

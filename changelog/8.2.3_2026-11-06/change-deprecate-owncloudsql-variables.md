Change: Deprecate the owncloudsql environment variables

We deprecated all environment variables of the `owncloudsql` driver in the
services `storage-users`, `auth-basic`, `users`, `groups`, and `sharing`. The
variables have the name part `OWNCLOUDSQL`. oCIS 9.0.0 removes the variables
together with the `owncloudsql` driver.

https://github.com/owncloud/ocis/pull/13193

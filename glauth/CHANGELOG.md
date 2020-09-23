# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-glauth unreleased.

[unreleased]: https://github.com/owncloud/ocis/glauth/compare/v0.5.0...master

## Summary

* Bugfix - Return invalid credentials when user was not found: [#30](https://github.com/owncloud/ocis/glauth/pull/30)
* Bugfix - Query numeric attribute values without quotes: [#28](https://github.com/owncloud/ocis/glauth/issues/28)
* Bugfix - Use searchBaseDN if already a user/group name: [#214](https://github.com/owncloud/product/issues/214)
* Bugfix - Fix LDAP substring startswith filters: [#31](https://github.com/owncloud/ocis/glauth/pull/31)

## Details

* Bugfix - Return invalid credentials when user was not found: [#30](https://github.com/owncloud/ocis/glauth/pull/30)

   We were relying on an error code of the ListAccounts call when the username and password was
   wrong. But the list will be empty if no user with the given login was found. So we also need to check
   if the list of accounts is empty.

   https://github.com/owncloud/ocis/glauth/pull/30


* Bugfix - Query numeric attribute values without quotes: [#28](https://github.com/owncloud/ocis/glauth/issues/28)

   Some LDAP properties like `uidnumber` and `gidnumber` are numeric. When an OS tries to look up a
   user it will not only try to lookup the user by username, but also by the `uidnumber`:
   `(&(objectclass=posixAccount)(uidnumber=20000))`. The accounts backend for glauth was
   sending that as a string query `uid_number eq '20000'` in the ListAccounts query. This PR
   changes that to `uid_number eq 20000`. The removed quotes allow the parser in ocis-accounts to
   identify the numeric literal.

   https://github.com/owncloud/ocis/glauth/issues/28
   https://github.com/owncloud/ocis/glauth/pull/29
   https://github.com/owncloud/ocis/accounts/pull/68


* Bugfix - Use searchBaseDN if already a user/group name: [#214](https://github.com/owncloud/product/issues/214)

   In case of the searchBaseDN already referencing a user or group, the search query was ignoring
   the user/group name entirely, because the searchBaseDN is not part of the LDAP filters. We
   fixed this by including an additional query part if the searchBaseDN contains a CN.

   https://github.com/owncloud/product/issues/214
   https://github.com/owncloud/ocis/glauth/pull/32


* Bugfix - Fix LDAP substring startswith filters: [#31](https://github.com/owncloud/ocis/glauth/pull/31)

   Filters like `(mail=mar*)` are currentld not parsed correctly, but they are used when
   searching for recipients. This PR correctly converts them to odata filters like
   `startswith(mail,'mar')`.

   https://github.com/owncloud/ocis/glauth/pull/31

# Changelog for [0.5.0] (2020-07-23)

The following sections list the changes in ocis-glauth 0.5.0.

[0.5.0]: https://github.com/owncloud/ocis/glauth/compare/v0.4.0...v0.5.0

## Summary

* Bugfix - Ignore case when comparing objectclass values: [#26](https://github.com/owncloud/ocis/glauth/pull/26)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#24](https://github.com/owncloud/ocis/glauth/pull/24)
* Enhancement - Handle ownCloudUUID attribute: [#27](https://github.com/owncloud/ocis/glauth/pull/27)
* Enhancement - Implement group queries: [#22](https://github.com/owncloud/ocis/glauth/issues/22)

## Details

* Bugfix - Ignore case when comparing objectclass values: [#26](https://github.com/owncloud/ocis/glauth/pull/26)

   The LDAP equality comparison is specified as case insensitive. We fixed the comparison for
   objectclass properties.

   https://github.com/owncloud/ocis/glauth/pull/26


* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#24](https://github.com/owncloud/ocis/glauth/pull/24)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis/glauth/pull/24


* Enhancement - Handle ownCloudUUID attribute: [#27](https://github.com/owncloud/ocis/glauth/pull/27)

   Clients can now query an accounts immutable id by using the [new `ownCloudUUID`
   attribute](https://github.com/butonic/owncloud-ldap-schema/blob/master/owncloud.schema#L28-L34).

   https://github.com/owncloud/ocis/glauth/pull/27


* Enhancement - Implement group queries: [#22](https://github.com/owncloud/ocis/glauth/issues/22)

   Refactored the handler and implemented group queries.

   https://github.com/owncloud/ocis/glauth/issues/22
   https://github.com/owncloud/ocis/glauth/pull/23

# Changelog for [0.4.0] (2020-03-18)

The following sections list the changes in ocis-glauth 0.4.0.

[0.4.0]: https://github.com/owncloud/ocis/glauth/compare/v0.2.0...v0.4.0

## Summary

* Enhancement - Configuration: [#11](https://github.com/owncloud/ocis/glauth/pull/11)
* Enhancement - Improve default settings: [#12](https://github.com/owncloud/ocis/glauth/pull/12)
* Enhancement - Generate temporary ldap certificates if LDAPS is enabled: [#12](https://github.com/owncloud/ocis/glauth/pull/12)
* Enhancement - Provide additional tls-endpoint: [#12](https://github.com/owncloud/ocis/glauth/pull/12)

## Details

* Enhancement - Configuration: [#11](https://github.com/owncloud/ocis/glauth/pull/11)

   Extensions should be responsible of configuring themselves. We use Viper for config loading
   from default paths. Environment variables **WILL** take precedence over config files.

   https://github.com/owncloud/ocis/glauth/pull/11


* Enhancement - Improve default settings: [#12](https://github.com/owncloud/ocis/glauth/pull/12)

   This helps achieve zero-config in single-binary.

   https://github.com/owncloud/ocis/glauth/pull/12


* Enhancement - Generate temporary ldap certificates if LDAPS is enabled: [#12](https://github.com/owncloud/ocis/glauth/pull/12)

   This change helps to achieve zero-configuration in single-binary mode.

   https://github.com/owncloud/ocis/glauth/pull/12


* Enhancement - Provide additional tls-endpoint: [#12](https://github.com/owncloud/ocis/glauth/pull/12)

   Ocis-glauth is now able to concurrently serve a encrypted and an unencrypted ldap-port.
   Please note that only SSL (no StarTLS) is supported at the moment.

   https://github.com/owncloud/ocis/glauth/pull/12

# Changelog for [0.2.0] (2020-03-17)

The following sections list the changes in ocis-glauth 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis/glauth/compare/v0.3.0...v0.2.0

## Summary

* Change - Default to config based user backend: [#6](https://github.com/owncloud/ocis/glauth/pull/6)

## Details

* Change - Default to config based user backend: [#6](https://github.com/owncloud/ocis/glauth/pull/6)

   We changed the default configuration to use the config file backend instead of the ownCloud
   backend.

   The config backend currently only has two hard coded users: demo and admin. To switch back to the
   ownCloud backend use `GLAUTH_BACKEND_DATASTORE=owncloud`

   https://github.com/owncloud/ocis/glauth/pull/6

# Changelog for [0.3.0] (2020-03-17)

The following sections list the changes in ocis-glauth 0.3.0.

[0.3.0]: https://github.com/owncloud/ocis/glauth/compare/v0.1.0...v0.3.0

## Summary

* Change - Use physicist demo users: [#5](https://github.com/owncloud/ocis/glauth/issues/5)

## Details

* Change - Use physicist demo users: [#5](https://github.com/owncloud/ocis/glauth/issues/5)

   Demo users like admin, demo and test don't allow you to tell a story. Which is why we changed the
   set of hard coded demo users to `einstein`, `marie` and `feynman`. You should know who they are.
   This also changes the ldap domain from `dc=owncloud,dc=com` to `dc=example,dc=org` because
   that is what these users use as their email domain. There are also `konnectd` and `reva` for
   technical purposes, eg. to allow konnectd and reva to bind to glauth.

   https://github.com/owncloud/ocis/glauth/issues/5

# Changelog for [0.1.0] (2020-02-28)

The following sections list the changes in ocis-glauth 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis/glauth/compare/178b6ccde34b64a88e8c14a9acb5857a4c6a3164...v0.1.0

## Summary

* Enhancement - Initial release of basic version: [#1](https://github.com/owncloud/ocis/glauth/pull/1)

## Details

* Enhancement - Initial release of basic version: [#1](https://github.com/owncloud/ocis/glauth/pull/1)

   Just prepare an initial basic version to provide a glauth service.

   https://github.com/owncloud/ocis/glauth/pull/1


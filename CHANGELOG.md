# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-glauth unreleased.

[unreleased]: https://github.com/owncloud/ocis-glauth/compare/v0.4.0...master

## Summary

* Enhancement - Implement group queries: [#22](https://github.com/owncloud/ocis-glauth/issues/22)

## Details

* Enhancement - Implement group queries: [#22](https://github.com/owncloud/ocis-glauth/issues/22)

   Refactored the handler and implemented group queries.

   https://github.com/owncloud/ocis-glauth/issues/22

# Changelog for [0.4.0] (2020-03-18)

The following sections list the changes in ocis-glauth 0.4.0.

[0.4.0]: https://github.com/owncloud/ocis-glauth/compare/v0.3.0...v0.4.0

## Summary

* Enhancement - Configuration: [#11](https://github.com/owncloud/ocis-glauth/pull/11)
* Enhancement - Improve default settings: [#12](https://github.com/owncloud/ocis-glauth/pull/12)
* Enhancement - Generate temporary ldap certificates if LDAPS is enabled: [#12](https://github.com/owncloud/ocis-glauth/pull/12)
* Enhancement - Provide additional tls-endpoint: [#12](https://github.com/owncloud/ocis-glauth/pull/12)

## Details

* Enhancement - Configuration: [#11](https://github.com/owncloud/ocis-glauth/pull/11)

   Extensions should be responsible of configuring themselves. We use Viper for config loading
   from default paths. Environment variables **WILL** take precedence over config files.

   https://github.com/owncloud/ocis-glauth/pull/11


* Enhancement - Improve default settings: [#12](https://github.com/owncloud/ocis-glauth/pull/12)

   This helps achieve zero-config in single-binary.

   https://github.com/owncloud/ocis-glauth/pull/12


* Enhancement - Generate temporary ldap certificates if LDAPS is enabled: [#12](https://github.com/owncloud/ocis-glauth/pull/12)

   This change helps to achieve zero-configuration in single-binary mode.

   https://github.com/owncloud/ocis-glauth/pull/12


* Enhancement - Provide additional tls-endpoint: [#12](https://github.com/owncloud/ocis-glauth/pull/12)

   Ocis-glauth is now able to concurrently serve a encrypted and an unencrypted ldap-port.
   Please note that only SSL (no StarTLS) is supported at the moment.

   https://github.com/owncloud/ocis-glauth/pull/12

# Changelog for [0.3.0] (2020-03-17)

The following sections list the changes in ocis-glauth 0.3.0.

[0.3.0]: https://github.com/owncloud/ocis-glauth/compare/v0.2.0...v0.3.0

## Summary

* Change - Use physicist demo users: [#5](https://github.com/owncloud/ocis-glauth/issues/5)

## Details

* Change - Use physicist demo users: [#5](https://github.com/owncloud/ocis-glauth/issues/5)

   Demo users like admin, demo and test don't allow you to tell a story. Which is why we changed the
   set of hard coded demo users to `einstein`, `marie` and `feynman`. You should know who they are.
   This also changes the ldap domain from `dc=owncloud,dc=com` to `dc=example,dc=org` because
   that is what these users use as their email domain. There are also `konnectd` and `reva` for
   technical purposes, eg. to allow konnectd and reva to bind to glauth.

   https://github.com/owncloud/ocis-glauth/issues/5

# Changelog for [0.2.0] (2020-03-17)

The following sections list the changes in ocis-glauth 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis-glauth/compare/v0.1.0...v0.2.0

## Summary

* Change - Default to config based user backend: [#6](https://github.com/owncloud/ocis-glauth/pull/6)

## Details

* Change - Default to config based user backend: [#6](https://github.com/owncloud/ocis-glauth/pull/6)

   We changed the default configuration to use the config file backend instead of the ownCloud
   backend.

   The config backend currently only has two hard coded users: demo and admin. To switch back to the
   ownCloud backend use `GLAUTH_BACKEND_DATASTORE=owncloud`

   https://github.com/owncloud/ocis-glauth/pull/6

# Changelog for [0.1.0] (2020-02-28)

The following sections list the changes in ocis-glauth 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-glauth/compare/178b6ccde34b64a88e8c14a9acb5857a4c6a3164...v0.1.0

## Summary

* Enhancement - Initial release of basic version: [#1](https://github.com/owncloud/ocis-glauth/pull/1)

## Details

* Enhancement - Initial release of basic version: [#1](https://github.com/owncloud/ocis-glauth/pull/1)

   Just prepare an initial basic version to provide a glauth service.

   https://github.com/owncloud/ocis-glauth/pull/1


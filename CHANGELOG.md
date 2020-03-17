# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-glauth unreleased.

[unreleased]: https://github.com/owncloud/ocis-glauth/compare/v0.1.0...master

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


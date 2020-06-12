# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-accounts unreleased.

[unreleased]: https://github.com/owncloud/ocis-accounts/compare/v0.1.1...master

## Summary

* Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis-accounts/pull/33)

## Details

* Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis-accounts/pull/33)

   We had a timezone setting in our profile settings bundle. As we're not dealing with a timezone
   yet it would be confusing for the user to have a timezone setting available. We removed it, until
   we have a timezone implementation available in ocis-web.

   https://github.com/owncloud/ocis-accounts/pull/33

# Changelog for [0.1.1] (2020-04-29)

The following sections list the changes in ocis-accounts 0.1.1.

[0.1.1]: https://github.com/owncloud/ocis-accounts/compare/v0.1.0...v0.1.1

## Summary

* Enhancement - Logging is configurable: [#24](https://github.com/owncloud/ocis-accounts/pull/24)

## Details

* Enhancement - Logging is configurable: [#24](https://github.com/owncloud/ocis-accounts/pull/24)

   ACCOUNTS_LOG_* env-vars or cli-flags can be used for logging configuration. See --help for
   more details.

   https://github.com/owncloud/ocis-accounts/pull/24

# Changelog for [0.1.0] (2020-03-18)

The following sections list the changes in ocis-accounts 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-accounts/compare/500e303cb544ed93d84153f01219d77eeee44929...v0.1.0

## Summary

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-accounts/issues/1)
* Enhancement - Configuration: [#15](https://github.com/owncloud/ocis-accounts/pull/15)

## Details

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-accounts/issues/1)

   Just prepared an initial basic version.

   https://github.com/owncloud/ocis-accounts/issues/1


* Enhancement - Configuration: [#15](https://github.com/owncloud/ocis-accounts/pull/15)

   Extensions should be responsible of configuring themselves. We use Viper for config loading
   from default paths. Environment variables **WILL** take precedence over config files.

   https://github.com/owncloud/ocis-accounts/pull/15


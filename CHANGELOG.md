# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-accounts unreleased.

[unreleased]: https://github.com/owncloud/ocis-accounts/compare/v0.1.1...master

## Summary

* Change - Align structure of this extension with other extensions: [#51](https://github.com/owncloud/ocis-accounts/pull/51)
* Change - Enable accounts on creation: [#43](https://github.com/owncloud/ocis-accounts/issues/43)
* Change - Pass around the correct logger throughout the code: [#41](https://github.com/owncloud/ocis-accounts/issues/41)
* Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis-accounts/pull/33)
* Enhancement - Update accounts API: [#30](https://github.com/owncloud/ocis-accounts/pull/30)
* Enhancement - Add simple user listing UI: [#51](https://github.com/owncloud/ocis-accounts/pull/51)

## Details

* Change - Align structure of this extension with other extensions: [#51](https://github.com/owncloud/ocis-accounts/pull/51)

   We aim to have a similar project structure for all our ocis extensions. This extension was
   different with regard to the structure of the server command and naming of some flag names.

   https://github.com/owncloud/ocis-accounts/pull/51


* Change - Enable accounts on creation: [#43](https://github.com/owncloud/ocis-accounts/issues/43)

   Accounts have been created with the account_enabled flag set to false. Now when they are
   created accounts will be enabled per default.

   https://github.com/owncloud/ocis-accounts/issues/43


* Change - Pass around the correct logger throughout the code: [#41](https://github.com/owncloud/ocis-accounts/issues/41)

   Pass around the logger to have consistent log formatting, log level, etc.

   https://github.com/owncloud/ocis-accounts/issues/41
   https://github.com/owncloud/ocis-accounts/pull/48


* Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis-accounts/pull/33)

   We had a timezone setting in our profile settings bundle. As we're not dealing with a timezone
   yet it would be confusing for the user to have a timezone setting available. We removed it, until
   we have a timezone implementation available in ocis-web.

   https://github.com/owncloud/ocis-accounts/pull/33


* Enhancement - Update accounts API: [#30](https://github.com/owncloud/ocis-accounts/pull/30)

   We updated the api to allow fetching users not onyl by UUID, but also by identity (OpenID issuer
   and subject) email, username and optionally a password.

   https://github.com/owncloud/ocis-accounts/pull/30


* Enhancement - Add simple user listing UI: [#51](https://github.com/owncloud/ocis-accounts/pull/51)

   We added an extension for ocis-web that shows a simple list of all existing users.

   https://github.com/owncloud/ocis-accounts/pull/51

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


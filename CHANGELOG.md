# Changes in unreleased

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


# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-phoenix unreleased.

[unreleased]: https://github.com/owncloud/ocis-phoenix/compare/v0.13.0...master

## Summary

* Bugfix - Fix external app URLs: [#218](https://github.com/owncloud/product/issues/218)

## Details

* Bugfix - Fix external app URLs: [#218](https://github.com/owncloud/product/issues/218)

   The URLs for the default set of external apps was hardcoded to localhost:9200. We fixed that by
   using relative paths instead.

   https://github.com/owncloud/product/issues/218
   https://github.com/owncloud/ocis-phoenix/pull/83

# Changelog for [0.13.0] (2020-08-25)

The following sections list the changes in ocis-phoenix 0.13.0.

[0.13.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.12.0...v0.13.0

## Summary

* Change - Update Phoenix: [#81](https://github.com/owncloud/ocis-phoenix/pull/81)

## Details

* Change - Update Phoenix: [#81](https://github.com/owncloud/ocis-phoenix/pull/81)

   Updated phoenix from v0.15.0 to v0.16.0

   https://github.com/owncloud/ocis-phoenix/pull/81
   https://github.com/owncloud/phoenix/releases/tag/v0.16.0

# Changelog for [0.12.0] (2020-08-19)

The following sections list the changes in ocis-phoenix 0.12.0.

[0.12.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.11.0...v0.12.0

## Summary

* Change - Enable Settings and Accounts apps by default: [#80](https://github.com/owncloud/ocis-phoenix/pull/80)
* Change - Update Phoenix: [#79](https://github.com/owncloud/ocis-phoenix/pull/79)

## Details

* Change - Enable Settings and Accounts apps by default: [#80](https://github.com/owncloud/ocis-phoenix/pull/80)

   The default ocis-web config now adds the frontend of ocis-accounts and ocis-settings to the
   builtin web config.

   https://github.com/owncloud/ocis-phoenix/pull/80

* Change - Update Phoenix: [#79](https://github.com/owncloud/ocis-phoenix/pull/79)

   Updated phoenix from v0.14.0 to v0.15.0

   https://github.com/owncloud/ocis-phoenix/pull/79
   https://github.com/owncloud/phoenix/releases/tag/v0.15.0

# Changelog for [0.11.0] (2020-08-17)

The following sections list the changes in ocis-phoenix 0.11.0.

[0.11.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.10.0...v0.11.0

## Summary

* Bugfix - Exit when assets or config are not found: [#76](https://github.com/owncloud/ocis-phoenix/pull/76)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#73](https://github.com/owncloud/ocis-phoenix/pull/73)
* Change - Hide searchbar by default: [#116](https://github.com/owncloud/product/issues/116)
* Change - Update Phoenix: [#78](https://github.com/owncloud/ocis-phoenix/pull/78)

## Details

* Bugfix - Exit when assets or config are not found: [#76](https://github.com/owncloud/ocis-phoenix/pull/76)

   When a non-existing assets folder is specified, there was only a warning log statement and the
   service served the builtin assets instead. It is safe to exit the service in such a scenario,
   instead of serving other assets than specified. We changed the log level to `Fatal` on
   non-existing assets. Similar for the web config, it was not failing on service level, but only
   showing an error in the web ui, wenn the specified config file could not be found. We changed the
   log level to `Fatal` as well.

   https://github.com/owncloud/ocis-phoenix/pull/76
   https://github.com/owncloud/ocis-phoenix/pull/77

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#73](https://github.com/owncloud/ocis-phoenix/pull/73)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-phoenix/pull/73

* Change - Hide searchbar by default: [#116](https://github.com/owncloud/product/issues/116)

   Since file search is not working at the moment we decided to hide the search bar by default.

   https://github.com/owncloud/product/issues/116
   https://github.com/owncloud/ocis-phoenix/pull/74

* Change - Update Phoenix: [#78](https://github.com/owncloud/ocis-phoenix/pull/78)

   Updated phoenix from v0.13.0 to v0.14.0

   https://github.com/owncloud/ocis-phoenix/pull/78

# Changelog for [0.10.0] (2020-07-17)

The following sections list the changes in ocis-phoenix 0.10.0.

[0.10.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.9.0...v0.10.0

## Summary

* Change - Update Phoenix: [#72](https://github.com/owncloud/ocis-phoenix/pull/72)

## Details

* Change - Update Phoenix: [#72](https://github.com/owncloud/ocis-phoenix/pull/72)

   Updated phoenix from v0.12.0 to v0.13.0

   https://github.com/owncloud/ocis-phoenix/pull/72

# Changelog for [0.9.0] (2020-07-10)

The following sections list the changes in ocis-phoenix 0.9.0.

[0.9.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.8.1...v0.9.0

## Summary

* Bugfix - Allow silent refresh of access token: [#69](https://github.com/owncloud/ocis-konnectd/issues/69)
* Change - Update Phoenix: [#70](https://github.com/owncloud/ocis-phoenix/pull/70)

## Details

* Bugfix - Allow silent refresh of access token: [#69](https://github.com/owncloud/ocis-konnectd/issues/69)

   Sets the `X-Frame-Options` header to `SAMEORIGIN` so the oidc client can refresh the token in
   an iframe.

   https://github.com/owncloud/ocis-konnectd/issues/69
   https://github.com/owncloud/ocis-phoenix/pull/69

* Change - Update Phoenix: [#70](https://github.com/owncloud/ocis-phoenix/pull/70)

   Updated phoenix from v0.11.1 to v0.12.0

   https://github.com/owncloud/ocis-phoenix/pull/70

# Changelog for [0.8.1] (2020-06-29)

The following sections list the changes in ocis-phoenix 0.8.1.

[0.8.1]: https://github.com/owncloud/ocis-phoenix/compare/v0.8.0...v0.8.1

## Summary

* Change - Update Phoenix: [#68](https://github.com/owncloud/ocis-phoenix/pull/68)

## Details

* Change - Update Phoenix: [#68](https://github.com/owncloud/ocis-phoenix/pull/68)

   Updated phoenix from v0.11.0 to v0.11.1

   https://github.com/owncloud/ocis-phoenix/pull/68

# Changelog for [0.8.0] (2020-06-26)

The following sections list the changes in ocis-phoenix 0.8.0.

[0.8.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.7.0...v0.8.0

## Summary

* Change - Update Phoenix: [#67](https://github.com/owncloud/ocis-phoenix/pull/67)

## Details

* Change - Update Phoenix: [#67](https://github.com/owncloud/ocis-phoenix/pull/67)

   Updated phoenix from v0.10.0 to v0.11.0

   https://github.com/owncloud/ocis-phoenix/pull/67

# Changelog for [0.7.0] (2020-05-26)

The following sections list the changes in ocis-phoenix 0.7.0.

[0.7.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.6.0...v0.7.0

## Summary

* Change - Update Phoenix: [#66](https://github.com/owncloud/ocis-phoenix/pull/66)

## Details

* Change - Update Phoenix: [#66](https://github.com/owncloud/ocis-phoenix/pull/66)

   Updated phoenix from v0.9.0 to v0.10.0

   https://github.com/owncloud/ocis-phoenix/pull/66

# Changelog for [0.6.0] (2020-04-28)

The following sections list the changes in ocis-phoenix 0.6.0.

[0.6.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.5.0...v0.6.0

## Summary

* Change - Update Phoenix: [#65](https://github.com/owncloud/ocis-phoenix/pull/65)

## Details

* Change - Update Phoenix: [#65](https://github.com/owncloud/ocis-phoenix/pull/65)

   Updated phoenix from v0.8.0 to v0.9.0

   https://github.com/owncloud/ocis-phoenix/pull/65

# Changelog for [0.5.0] (2020-04-14)

The following sections list the changes in ocis-phoenix 0.5.0.

[0.5.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.4.1...v0.5.0

## Summary

* Change - Update Phoenix: [#63](https://github.com/owncloud/ocis-phoenix/pull/63)

## Details

* Change - Update Phoenix: [#63](https://github.com/owncloud/ocis-phoenix/pull/63)

   Updated phoenix from v0.7.0 to v0.8.0

   https://github.com/owncloud/ocis-phoenix/pull/63

# Changelog for [0.4.1] (2020-04-01)

The following sections list the changes in ocis-phoenix 0.4.1.

[0.4.1]: https://github.com/owncloud/ocis-phoenix/compare/v0.4.0...v0.4.1

## Summary

* Bugfix - Create a new tag to fix v0.4.0: [#62](https://github.com/owncloud/ocis-phoenix/pull/62)

## Details

* Bugfix - Create a new tag to fix v0.4.0: [#62](https://github.com/owncloud/ocis-phoenix/pull/62)

   Release v0.4.0 is using the wrong assets. We fixed that by creating a new release.

   https://github.com/owncloud/ocis-phoenix/pull/62

# Changelog for [0.4.0] (2020-03-31)

The following sections list the changes in ocis-phoenix 0.4.0.

[0.4.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.2.0...v0.4.0

## Summary

* Change - Update Phoenix: [#60](https://github.com/owncloud/ocis-phoenix/pull/60)
* Enhancement - Configuration: [#57](https://github.com/owncloud/ocis-phoenix/pull/57)

## Details

* Change - Update Phoenix: [#60](https://github.com/owncloud/ocis-phoenix/pull/60)

   Updated phoenix from v0.6.0 to v0.7.0

   https://github.com/owncloud/ocis-phoenix/pull/60

* Enhancement - Configuration: [#57](https://github.com/owncloud/ocis-phoenix/pull/57)

   Extensions should be responsible of configuring themselves. We use Viper for config loading
   from default paths. Environment variables **WILL** take precedence over config files.

   https://github.com/owncloud/ocis-phoenix/pull/57

# Changelog for [0.2.0] (2020-03-17)

The following sections list the changes in ocis-phoenix 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.3.0...v0.2.0

## Summary

* Bugfix - Config file value not being read: [#45](https://github.com/owncloud/ocis-phoenix/pull/45)
* Enhancement - Update to Phoenix v0.5.0: [#43](https://github.com/owncloud/ocis-phoenix/issues/43)
* Enhancement - Update to Phoenix v0.6.0: [#53](https://github.com/owncloud/ocis-phoenix/pull/53)

## Details

* Bugfix - Config file value not being read: [#45](https://github.com/owncloud/ocis-phoenix/pull/45)

   There was a bug in which phoenix config is always set to the default values and the contents of the
   config file were actually ignored.

   https://github.com/owncloud/ocis-phoenix/issues/46
   https://github.com/owncloud/ocis-phoenix/issues/47
   https://github.com/owncloud/ocis-phoenix/pull/45

* Enhancement - Update to Phoenix v0.5.0: [#43](https://github.com/owncloud/ocis-phoenix/issues/43)

   Use the latest phoenix release

   https://github.com/owncloud/ocis-phoenix/issues/43

* Enhancement - Update to Phoenix v0.6.0: [#53](https://github.com/owncloud/ocis-phoenix/pull/53)

   Use the latest phoenix release

   https://github.com/owncloud/ocis-phoenix/pull/53

# Changelog for [0.3.0] (2020-03-17)

The following sections list the changes in ocis-phoenix 0.3.0.

[0.3.0]: https://github.com/owncloud/ocis-phoenix/compare/v0.1.0...v0.3.0

## Summary

* Change - Default to running behind ocis-proxy: [#55](https://github.com/owncloud/ocis-phoenix/pull/55)

## Details

* Change - Default to running behind ocis-proxy: [#55](https://github.com/owncloud/ocis-phoenix/pull/55)

   We changed the default configuration to integrate better with ocis.

   Clients are supposed to use the ocis-proxy endpoint `https://localhost:9200`

   https://github.com/owncloud/ocis-phoenix/pull/55

# Changelog for [0.1.0] (2020-02-03)

The following sections list the changes in ocis-phoenix 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-phoenix/compare/432c57c406a8421a20ba596818d95f816e2ef9c7...v0.1.0

## Summary

* Change - Initial release of basic version: [#3](https://github.com/owncloud/ocis-phoenix/issues/3)

## Details

* Change - Initial release of basic version: [#3](https://github.com/owncloud/ocis-phoenix/issues/3)

   Just prepared an initial basic version to serve Phoenix for the ownCloud Infinite Scale
   project. It just provides a minimal viable product to demonstrate the microservice pattern.

   https://github.com/owncloud/ocis-phoenix/issues/3


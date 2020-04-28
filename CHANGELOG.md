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


# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes for ocis-proxy unreleased.

[unreleased]: https://github.com/owncloud/ocis-proxy/compare/v0.2.1...master

## Summary

* Enhancement - Configurable OpenID Connect client: [#27](https://github.com/owncloud/ocis-proxy/pull/27)
* Enhancement - Add policy selectors: [#4](https://github.com/owncloud/ocis-proxy/issues/4)

## Details

* Enhancement - Configurable OpenID Connect client: [#27](https://github.com/owncloud/ocis-proxy/pull/27)

   The proxy will try to authenticate every request with the configured OIDC provider.

   See configs/proxy-example.oidc.json for an example-configuration.

   https://github.com/owncloud/ocis-proxy/pull/27


* Enhancement - Add policy selectors: [#4](https://github.com/owncloud/ocis-proxy/issues/4)

   "Static-Policy" can be configured to always select a specific policy. See:
   config/proxy-example.json.

   "Migration-Policy" selects policy depending on existence of the uid in the ocis-accounts
   service. See: config/proxy-example-migration.json

   https://github.com/owncloud/ocis-proxy/issues/4

# Changelog for [0.2.1] (2020-03-25)

The following sections list the changes for ocis-proxy 0.2.1.

[0.2.1]: https://github.com/owncloud/ocis-proxy/compare/v0.2.0...v0.2.1

## Summary

* Bugfix - Set TLS-Certificate correctly: [#25](https://github.com/owncloud/ocis-proxy/pull/25)

## Details

* Bugfix - Set TLS-Certificate correctly: [#25](https://github.com/owncloud/ocis-proxy/pull/25)

   https://github.com/owncloud/ocis-proxy/pull/25

# Changelog for [0.2.0] (2020-03-25)

The following sections list the changes for ocis-proxy 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis-proxy/compare/v0.1.0...v0.2.0

## Summary

* Change - Route requests based on regex or query parameters: [#21](https://github.com/owncloud/ocis-proxy/issues/21)
* Enhancement - Proxy client urls in default configuration: [#19](https://github.com/owncloud/ocis-proxy/issues/19)
* Enhancement - Make TLS-Cert configurable: [#14](https://github.com/owncloud/ocis-proxy/pull/14)

## Details

* Change - Route requests based on regex or query parameters: [#21](https://github.com/owncloud/ocis-proxy/issues/21)

   Some requests needed to be distinguished based on a pattern or a query parameter. We've
   implemented the functionality to route requests based on different conditions.

   https://github.com/owncloud/ocis-proxy/issues/21


* Enhancement - Proxy client urls in default configuration: [#19](https://github.com/owncloud/ocis-proxy/issues/19)

   Proxy /status.php and index.php/*

   https://github.com/owncloud/ocis-proxy/issues/19


* Enhancement - Make TLS-Cert configurable: [#14](https://github.com/owncloud/ocis-proxy/pull/14)

   Before a generates certificates on every start was used for dev purposes.

   https://github.com/owncloud/ocis-proxy/pull/14

# Changelog for [0.1.0] (2020-03-18)

The following sections list the changes for ocis-proxy 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-proxy/compare/500e303cb544ed93d84153f01219d77eeee44929...v0.1.0

## Summary

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-proxy/issues/1)
* Enhancement - Load Proxy Policies at Runtime: [#17](https://github.com/owncloud/ocis-proxy/issues/17)

## Details

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-proxy/issues/1)

   Just prepared an initial basic version.

   https://github.com/owncloud/ocis-proxy/issues/1


* Enhancement - Load Proxy Policies at Runtime: [#17](https://github.com/owncloud/ocis-proxy/issues/17)

   While a proxy without policies is of no use, the current state of ocis-proxy expects a config
   file either at an expected Viper location or specified via -- config-file flag. To ease
   deployments and ensure a working set of policies out of the box we need a series of defaults.

   https://github.com/owncloud/ocis-proxy/issues/17
   https://github.com/owncloud/ocis-proxy/pull/16


# Changelog for [0.6.0] (2020-08-17)

The following sections list the changes for ocis-proxy 0.6.0.

[0.6.0]: https://github.com/owncloud/ocis-proxy/compare/v0.5.0...v0.6.0

## Summary

* Bugfix - Enable new accounts by default: [#79](https://github.com/owncloud/ocis-proxy/pull/79)
* Bugfix - Lookup user by id for presigned URLs: [#85](https://github.com/owncloud/ocis-proxy/pull/85)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#78](https://github.com/owncloud/ocis-proxy/pull/78)
* Change - Add settings and ocs group routes: [#81](https://github.com/owncloud/ocis-proxy/pull/81)
* Change - Add route for user provisioning API in ocis-ocs: [#80](https://github.com/owncloud/ocis-proxy/pull/80)

## Details

* Bugfix - Enable new accounts by default: [#79](https://github.com/owncloud/ocis-proxy/pull/79)

   When new accounts are created, they also need to be enabled to be useable.

   https://github.com/owncloud/ocis-proxy/pull/79


* Bugfix - Lookup user by id for presigned URLs: [#85](https://github.com/owncloud/ocis-proxy/pull/85)

   Phoenix will send the `userid`, not the `username` as the `OC-Credential` for presigned URLs.
   This PR uses the new `ocisid` claim in the OIDC userinfo to pass the userid to the account
   middleware.

   https://github.com/owncloud/ocis/issues/436
   https://github.com/owncloud/ocis-proxy/pull/85
   https://github.com/owncloud/ocis-pkg/pull/50


* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#78](https://github.com/owncloud/ocis-proxy/pull/78)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-proxy/pull/78


* Change - Add settings and ocs group routes: [#81](https://github.com/owncloud/ocis-proxy/pull/81)

   Route settings requests and ocs group related requests to new services

   https://github.com/owncloud/ocis-proxy/pull/81


* Change - Add route for user provisioning API in ocis-ocs: [#80](https://github.com/owncloud/ocis-proxy/pull/80)

   We added a route to send requests on the user provisioning API endpoints to ocis-ocs.

   https://github.com/owncloud/ocis-proxy/pull/80

# Changelog for [0.5.0] (2020-07-23)

The following sections list the changes for ocis-proxy 0.5.0.

[0.5.0]: https://github.com/owncloud/ocis-proxy/compare/v0.4.0...v0.5.0

## Summary

* Bugfix - Provide token configuration from config: [#69](https://github.com/owncloud/ocis-proxy/pull/69)
* Bugfix - Provide token configuration from config: [#76](https://github.com/owncloud/ocis-proxy/pull/76)
* Change - Add OIDC config flags: [#66](https://github.com/owncloud/ocis-proxy/pull/66)
* Change - Mint new username property in the reva token: [#62](https://github.com/owncloud/ocis-proxy/pull/62)
* Enhancement - Add Accounts UI routes: [#65](https://github.com/owncloud/ocis-proxy/pull/65)
* Enhancement - Add option to disable TLS: [#71](https://github.com/owncloud/ocis-proxy/issues/71)
* Enhancement - Only send create home request if an account has been migrated: [#52](https://github.com/owncloud/ocis-proxy/issues/52)
* Enhancement - Create a root span on proxy that propagates down to consumers: [#64](https://github.com/owncloud/ocis-proxy/pull/64)
* Enhancement - Support signed URLs: [#73](https://github.com/owncloud/ocis-proxy/issues/73)

## Details

* Bugfix - Provide token configuration from config: [#69](https://github.com/owncloud/ocis-proxy/pull/69)

   Fixed a bug that causes the createHome middleware to crash if no configuration for the
   TokenManager is propagated.

   https://github.com/owncloud/ocis-proxy/pull/69


* Bugfix - Provide token configuration from config: [#76](https://github.com/owncloud/ocis-proxy/pull/76)

   Fixed a bug that causes the createHome middleware to crash if the createHome response has no
   Status set

   https://github.com/owncloud/ocis-proxy/pull/76


* Change - Add OIDC config flags: [#66](https://github.com/owncloud/ocis-proxy/pull/66)

   To authenticate requests with an oidc provider we added two environment variables: -
   `PROXY_OIDC_ISSUER="https://localhost:9200"` and - `PROXY_OIDC_INSECURE=true`

   This changes ocis-proxy to now load the oidc-middleware by default, requiring a bearer token
   and exchanging the email in the OIDC claims for an account id at the ocis-accounts service.

   Setting `PROXY_OIDC_ISSUER=""` will disable the OIDC middleware.

   https://github.com/owncloud/ocis-proxy/pull/66


* Change - Mint new username property in the reva token: [#62](https://github.com/owncloud/ocis-proxy/pull/62)

   An accounts username is now taken from the on_premises_sam_account_name property instead of
   the preferred_name. Furthermore the group name (also from on_premises_sam_account_name
   property) is now minted into the token as well.

   https://github.com/owncloud/ocis-proxy/pull/62


* Enhancement - Add Accounts UI routes: [#65](https://github.com/owncloud/ocis-proxy/pull/65)

   The accounts service has a ui that requires routing - `/api/v0/accounts` and - `/accounts.js`

   To http://localhost:9181

   https://github.com/owncloud/ocis-proxy/pull/65


* Enhancement - Add option to disable TLS: [#71](https://github.com/owncloud/ocis-proxy/issues/71)

   Can be used to disable TLS when the ocis-proxy is behind an TLS-Terminating reverse proxy.

   Env PROXY_TLS=false or --tls=false

   https://github.com/owncloud/ocis-proxy/issues/71
   https://github.com/owncloud/ocis-proxy/pull/72


* Enhancement - Only send create home request if an account has been migrated: [#52](https://github.com/owncloud/ocis-proxy/issues/52)

   This change adds a check if an account has been migrated by getting it from the ocis-accounts
   service. If no account is returned it means it hasn't been migrated.

   https://github.com/owncloud/ocis-proxy/issues/52
   https://github.com/owncloud/ocis-proxy/pull/63


* Enhancement - Create a root span on proxy that propagates down to consumers: [#64](https://github.com/owncloud/ocis-proxy/pull/64)

   In order to propagate and correctly associate a span with a request we need a root span that gets
   sent to other services.

   https://github.com/owncloud/ocis-proxy/pull/64


* Enhancement - Support signed URLs: [#73](https://github.com/owncloud/ocis-proxy/issues/73)

   We added a middleware that verifies signed urls as generated by the owncloud-sdk. This allows
   directly downloading large files with browsers instead of using `blob://` urls, which eats
   memory ...

   https://github.com/owncloud/ocis-proxy/issues/73
   https://github.com/owncloud/ocis-proxy/pull/75
   https://github.com/owncloud/ocis-ocs/pull/18
   https://github.com/owncloud/owncloud-sdk/pull/504

# Changelog for [0.4.0] (2020-06-25)

The following sections list the changes for ocis-proxy 0.4.0.

[0.4.0]: https://github.com/owncloud/ocis-proxy/compare/v0.3.1...v0.4.0

## Summary

* Bugfix - Accounts service response was ignored: [#43](https://github.com/owncloud/ocis-proxy/pull/43)
* Bugfix - Fix x-access-token in header: [#41](https://github.com/owncloud/ocis-proxy/pull/41)
* Change - Point /data endpoint to reva frontend: [#45](https://github.com/owncloud/ocis-proxy/pull/45)
* Change - Send autocreate home request to reva gateway: [#51](https://github.com/owncloud/ocis-proxy/pull/51)
* Change - Update to new accounts API: [#39](https://github.com/owncloud/ocis-proxy/issues/39)
* Enhancement - Retrieve Account UUID From User Claims: [#36](https://github.com/owncloud/ocis-proxy/pull/36)
* Enhancement - Create account if it doesn't exist in ocis-accounts: [#55](https://github.com/owncloud/ocis-proxy/issues/55)
* Enhancement - Disable keep-alive on server-side OIDC requests: [#268](https://github.com/owncloud/ocis/issues/268)
* Enhancement - Make jwt secret configurable: [#41](https://github.com/owncloud/ocis-proxy/pull/41)
* Enhancement - Respect account_enabled flag: [#53](https://github.com/owncloud/ocis-proxy/issues/53)

## Details

* Bugfix - Accounts service response was ignored: [#43](https://github.com/owncloud/ocis-proxy/pull/43)

   We fixed an error in the AccountUUID middleware that was responsible for ignoring an account
   uuid provided by the accounts service.

   https://github.com/owncloud/ocis-proxy/pull/43


* Bugfix - Fix x-access-token in header: [#41](https://github.com/owncloud/ocis-proxy/pull/41)

   We fixed setting the x-access-token in the request header, which was broken before.

   https://github.com/owncloud/ocis-proxy/pull/41
   https://github.com/owncloud/ocis-proxy/pull/46


* Change - Point /data endpoint to reva frontend: [#45](https://github.com/owncloud/ocis-proxy/pull/45)

   Adjusted example config files to point /data to the reva frontend.

   https://github.com/owncloud/ocis-proxy/pull/45


* Change - Send autocreate home request to reva gateway: [#51](https://github.com/owncloud/ocis-proxy/pull/51)

   Send autocreate home request to reva gateway

   https://github.com/owncloud/ocis-proxy/pull/51


* Change - Update to new accounts API: [#39](https://github.com/owncloud/ocis-proxy/issues/39)

   Update to new accounts API

   https://github.com/owncloud/ocis-proxy/issues/39


* Enhancement - Retrieve Account UUID From User Claims: [#36](https://github.com/owncloud/ocis-proxy/pull/36)

   OIDC Middleware can make use of uuidFromClaims to trade claims.Email for an account's UUID.
   For this, a general purpose cache was added that caches on a per-request basis, meaning
   whenever the request parameters match a set of keys, the cached value is returned, saving a
   round trip to the accounts service that otherwise would happen in every single request.

   https://github.com/owncloud/ocis-proxy/pull/36


* Enhancement - Create account if it doesn't exist in ocis-accounts: [#55](https://github.com/owncloud/ocis-proxy/issues/55)

   The accounts_uuid middleware tries to get the account from ocis-accounts. If it doens't exist
   there yet the proxy creates the account using the ocis-account api.

   https://github.com/owncloud/ocis-proxy/issues/55
   https://github.com/owncloud/ocis-proxy/issues/58


* Enhancement - Disable keep-alive on server-side OIDC requests: [#268](https://github.com/owncloud/ocis/issues/268)

   This should reduce file-descriptor counts

   https://github.com/owncloud/ocis/issues/268
   https://github.com/owncloud/ocis-proxy/pull/42
   https://github.com/cs3org/reva/pull/787


* Enhancement - Make jwt secret configurable: [#41](https://github.com/owncloud/ocis-proxy/pull/41)

   We added a config option for the reva token manager JWTSecret. It was hardcoded before and is now
   configurable.

   https://github.com/owncloud/ocis-proxy/pull/41


* Enhancement - Respect account_enabled flag: [#53](https://github.com/owncloud/ocis-proxy/issues/53)

   If the account returned by the accounts service has the account_enabled flag set to false, the
   proxy will return immediately with the status code unauthorized.

   https://github.com/owncloud/ocis-proxy/issues/53

# Changelog for [0.3.1] (2020-03-31)

The following sections list the changes for ocis-proxy 0.3.1.

[0.3.1]: https://github.com/owncloud/ocis-proxy/compare/v0.3.0...v0.3.1

## Summary

* Change - Update ocis-pkg: [#30](https://github.com/owncloud/ocis-proxy/pull/30)

## Details

* Change - Update ocis-pkg: [#30](https://github.com/owncloud/ocis-proxy/pull/30)

   We updated ocis-pkg from 2.0.2 to 2.2.0.

   https://github.com/owncloud/ocis-proxy/pull/30

# Changelog for [0.3.0] (2020-03-30)

The following sections list the changes for ocis-proxy 0.3.0.

[0.3.0]: https://github.com/owncloud/ocis-proxy/compare/v0.2.0...v0.3.0

## Summary

* Change - Insecure http-requests are now redirected to https: [#29](https://github.com/owncloud/ocis-proxy/pull/29)
* Enhancement - Configurable OpenID Connect client: [#27](https://github.com/owncloud/ocis-proxy/pull/27)
* Enhancement - Add policy selectors: [#4](https://github.com/owncloud/ocis-proxy/issues/4)

## Details

* Change - Insecure http-requests are now redirected to https: [#29](https://github.com/owncloud/ocis-proxy/pull/29)

   https://github.com/owncloud/ocis-proxy/pull/29


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

# Changelog for [0.2.0] (2020-03-25)

The following sections list the changes for ocis-proxy 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis-proxy/compare/v0.2.1...v0.2.0

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

# Changelog for [0.2.1] (2020-03-25)

The following sections list the changes for ocis-proxy 0.2.1.

[0.2.1]: https://github.com/owncloud/ocis-proxy/compare/v0.1.0...v0.2.1

## Summary

* Bugfix - Set TLS-Certificate correctly: [#25](https://github.com/owncloud/ocis-proxy/pull/25)

## Details

* Bugfix - Set TLS-Certificate correctly: [#25](https://github.com/owncloud/ocis-proxy/pull/25)

   https://github.com/owncloud/ocis-proxy/pull/25

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


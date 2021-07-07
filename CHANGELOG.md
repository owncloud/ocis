# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes for unreleased.

[unreleased]: https://github.com/owncloud/ocis/compare/v1.8.0...master

## Summary

* Bugfix - Panic when service fails to start: [#2252](https://github.com/owncloud/ocis/pull/2252)
* Enhancement - Runtime support for cherry picking extensions: [#2229](https://github.com/owncloud/ocis/pull/2229)
* Enhancement - Remove unnecessary Service.Init(): [#1705](https://github.com/owncloud/ocis/pull/1705)
* Enhancement - Update REVA to v1.9.1-0.20210628143859-9d29c36c0c3f: [#2227](https://github.com/owncloud/ocis/pull/2227)

## Details

* Bugfix - Panic when service fails to start: [#2252](https://github.com/owncloud/ocis/pull/2252)

   Tags: runtime

   When attempting to run a service through the runtime that is currently running and fails to
   start, a race condition still redirect os Interrupt signals to a closed channel.

   https://github.com/owncloud/ocis/pull/2252

* Enhancement - Runtime support for cherry picking extensions: [#2229](https://github.com/owncloud/ocis/pull/2229)

   Support for running certain extensions supervised via cli flags. Example usage:

   ``` > ocis server --extensions="proxy, idp, storage-metadata, accounts" ```

   https://github.com/owncloud/ocis/pull/2229

* Enhancement - Remove unnecessary Service.Init(): [#1705](https://github.com/owncloud/ocis/pull/1705)

   As it turns out oCIS already calls this method. Invoking it twice would end in accidentally
   resetting values.

   https://github.com/owncloud/ocis/pull/1705

* Enhancement - Update REVA to v1.9.1-0.20210628143859-9d29c36c0c3f: [#2227](https://github.com/owncloud/ocis/pull/2227)

   https://github.com/owncloud/ocis/pull/2227
# Changelog for [1.8.0] (2021-06-28)

The following sections list the changes for 1.8.0.

[1.8.0]: https://github.com/owncloud/ocis/compare/v1.7.0...v1.8.0

## Summary

* Bugfix - External storage registration used wrong config: [#2120](https://github.com/owncloud/ocis/pull/2120)
* Bugfix - Remove authentication from /status.php completely: [#2188](https://github.com/owncloud/ocis/pull/2188)
* Bugfix - Make webdav namespace configurable across services: [#2198](https://github.com/owncloud/ocis/pull/2198)
* Change - Update ownCloud Web to v3.3.0: [#2187](https://github.com/owncloud/ocis/pull/2187)
* Enhancement - Properly configure graph-explorer client registration: [#2118](https://github.com/owncloud/ocis/pull/2118)
* Enhancement - Use system default location to store TLS artefacts: [#2129](https://github.com/owncloud/ocis/pull/2129)
* Enhancement - Update REVA to v1.9: [#2205](https://github.com/owncloud/ocis/pull/2205)

## Details

* Bugfix - External storage registration used wrong config: [#2120](https://github.com/owncloud/ocis/pull/2120)

   The go-micro registry-singleton ignores the ocis configuration and defaults to mdns

   https://github.com/owncloud/ocis/pull/2120

* Bugfix - Remove authentication from /status.php completely: [#2188](https://github.com/owncloud/ocis/pull/2188)

   Despite requests without Authentication header being successful, requests with an invalid
   bearer token in the Authentication header were rejected in the proxy with an 401
   unauthenticated. Now the Authentication header is completely ignored for the /status.php
   route.

   https://github.com/owncloud/client/issues/8538
   https://github.com/owncloud/ocis/pull/2188

* Bugfix - Make webdav namespace configurable across services: [#2198](https://github.com/owncloud/ocis/pull/2198)

   The WebDAV namespace is used across various services, but it was previously hardcoded in some
   of the services. This PR uses the same environment variable to set the config correctly across
   the services.

   https://github.com/owncloud/ocis/pull/2198

* Change - Update ownCloud Web to v3.3.0: [#2187](https://github.com/owncloud/ocis/pull/2187)

   Tags: web

   We updated ownCloud Web to v3.3.0. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/2187
   https://github.com/owncloud/web/releases/tag/v3.3.0

* Enhancement - Properly configure graph-explorer client registration: [#2118](https://github.com/owncloud/ocis/pull/2118)

   The client registration in the `identifier-registration.yaml` for the graph-explorer
   didn't contain `redirect_uris` nor `origins`. Both were added to prevent exploitation.

   https://github.com/owncloud/ocis/pull/2118

* Enhancement - Use system default location to store TLS artefacts: [#2129](https://github.com/owncloud/ocis/pull/2129)

   This used to default to the current location of the binary, which is not ideal after a first run as
   it leaves traces behind. It now uses the system's location for artefacts with the help of
   https://golang.org/pkg/os/#UserConfigDir.

   https://github.com/owncloud/ocis/pull/2129

* Enhancement - Update REVA to v1.9: [#2205](https://github.com/owncloud/ocis/pull/2205)

   This update includes * [set Content-Type
   correctly](https://github.com/cs3org/reva/pull/1750) * [Return file checksum
   available from the metadata for the EOS
   driver](https://github.com/cs3org/reva/pull/1755) * [Sort share entries
   alphabetically](https://github.com/cs3org/reva/pull/1772) * [Initial work on the
   owncloudsql driver](https://github.com/cs3org/reva/pull/1710) * [Add user ID cache
   warmup to EOS storage driver](https://github.com/cs3org/reva/pull/1774) * [Use
   UidNumber and GidNumber fields in User
   objects](https://github.com/cs3org/reva/pull/1573) * [EOS GRPC
   interface](https://github.com/cs3org/reva/pull/1471) * [switch
   references](https://github.com/cs3org/reva/pull/1721) * [remove user's uuid from
   trashbin file key](https://github.com/cs3org/reva/pull/1793) * [fix restore behavior of
   the trashbin API](https://github.com/cs3org/reva/pull/1795) * [eosfs: add arbitrary
   metadata support](https://github.com/cs3org/reva/pull/1811)

   https://github.com/owncloud/ocis/pull/2205
   https://github.com/owncloud/ocis/pull/2210
# Changelog for [1.7.0] (2021-06-04)

The following sections list the changes for 1.7.0.

[1.7.0]: https://github.com/owncloud/ocis/compare/v1.6.0...v1.7.0

## Summary

* Bugfix - Change the groups index to be case sensitive: [#2109](https://github.com/owncloud/ocis/pull/2109)
* Change - Update ownCloud Web to v3.2.0: [#2096](https://github.com/owncloud/ocis/pull/2096)
* Enhancement - Enable the s3ng storage driver: [#1886](https://github.com/owncloud/ocis/pull/1886)
* Enhancement - Color contrasts on IDP/OIDC login pages: [#2088](https://github.com/owncloud/ocis/pull/2088)
* Enhancement - Announce user profile picture capability: [#2036](https://github.com/owncloud/ocis/pull/2036)
* Enhancement - Update reva to v1.7.1-0.20210531093513-b74a2b156af6: [#2104](https://github.com/owncloud/ocis/pull/2104)

## Details

* Bugfix - Change the groups index to be case sensitive: [#2109](https://github.com/owncloud/ocis/pull/2109)

   Groups are considered to be case sensitive. The index must handle them case sensitive too
   otherwise we will have undeterministic behavior while editing or deleting groups.

   https://github.com/owncloud/ocis/pull/2109

* Change - Update ownCloud Web to v3.2.0: [#2096](https://github.com/owncloud/ocis/pull/2096)

   Tags: web

   We updated ownCloud Web to v3.2.0. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/2096
   https://github.com/owncloud/web/releases/tag/v3.2.0

* Enhancement - Enable the s3ng storage driver: [#1886](https://github.com/owncloud/ocis/pull/1886)

   We made it possible to use the new s3ng storage driver by adding according commandline flags and
   environment variables.

   https://github.com/owncloud/ocis/pull/1886

* Enhancement - Color contrasts on IDP/OIDC login pages: [#2088](https://github.com/owncloud/ocis/pull/2088)

   We have updated the color contrasts on the IDP pages in order to improve accessibility.

   https://github.com/owncloud/ocis/pull/2088

* Enhancement - Announce user profile picture capability: [#2036](https://github.com/owncloud/ocis/pull/2036)

   Added a new capability (through https://github.com/cs3org/reva/pull/1694) to prevent the
   web frontend from fetching (nonexistent) user avatar profile pictures which added latency &
   console errors.

   https://github.com/owncloud/ocis/pull/2036

* Enhancement - Update reva to v1.7.1-0.20210531093513-b74a2b156af6: [#2104](https://github.com/owncloud/ocis/pull/2104)

   This reva update includes: * [fix move in the owncloud storage
   driver](https://github.com/cs3org/reva/pull/1696) * [add checksum header to the tus
   preflight response](https://github.com/cs3org/reva/pull/1702) * [Add reliability
   calculations support to Mentix](https://github.com/cs3org/reva/pull/1649) * [fix
   response format when accepting shares](https://github.com/cs3org/reva/pull/1724) *
   [Datatx createtransfershare](https://github.com/cs3org/reva/pull/1725)

   https://github.com/owncloud/ocis/issues/2102
   https://github.com/owncloud/ocis/pull/2104
# Changelog for [1.6.0] (2021-05-12)

The following sections list the changes for 1.6.0.

[1.6.0]: https://github.com/owncloud/ocis/compare/v1.5.0...v1.6.0

## Summary

* Bugfix - Fix STORAGE_METADATA_ROOT default value override: [#1956](https://github.com/owncloud/ocis/pull/1956)
* Bugfix - Stop the supervisor if a service fails to start: [#1963](https://github.com/owncloud/ocis/pull/1963)
* Change - Update ownCloud Web to v3.1.0: [#2045](https://github.com/owncloud/ocis/pull/2045)
* Enhancement - Added dictionary files: [#2003](https://github.com/owncloud/ocis/pull/2003)
* Enhancement - Introduce login form with h1 tag for screen readers only: [#1991](https://github.com/owncloud/ocis/pull/1991)
* Enhancement - User Deprovisioning for the OCS API: [#1962](https://github.com/owncloud/ocis/pull/1962)
* Enhancement - Support thumbnails for txt files: [#1988](https://github.com/owncloud/ocis/pull/1988)
* Enhancement - Update reva to v1.7.1-0.20210430154404-69bd21f2cc97: [#2010](https://github.com/owncloud/ocis/pull/2010)
* Enhancement - Update reva to v1.7.1-0.20210507160327-e2c3841d0dbc: [#2044](https://github.com/owncloud/ocis/pull/2044)
* Enhancement - Use oc-select: [#1979](https://github.com/owncloud/ocis/pull/1979)
* Enhancement - Set SameSite settings to Strict for Web: [#2019](https://github.com/owncloud/ocis/pull/2019)

## Details

* Bugfix - Fix STORAGE_METADATA_ROOT default value override: [#1956](https://github.com/owncloud/ocis/pull/1956)

   The way the value was being set ensured that it was NOT being overridden where it should have
   been. This patch ensures the correct loading order of values.

   https://github.com/owncloud/ocis/pull/1956

* Bugfix - Stop the supervisor if a service fails to start: [#1963](https://github.com/owncloud/ocis/pull/1963)

   Steps to make the supervisor fail:

   `PROXY_HTTP_ADDR=0.0.0.0:9144 bin/ocis server`

   https://github.com/owncloud/ocis/pull/1963

* Change - Update ownCloud Web to v3.1.0: [#2045](https://github.com/owncloud/ocis/pull/2045)

   Tags: web

   We updated ownCloud Web to v3.1.0. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/2045
   https://github.com/owncloud/web/releases/tag/v3.1.0

* Enhancement - Added dictionary files: [#2003](https://github.com/owncloud/ocis/pull/2003)

   Added the dictionary.js file for package settings and accounts which contains strings that
   should be synced to transifex but not exist in the UI directly.

   https://github.com/owncloud/ocis/pull/2003

* Enhancement - Introduce login form with h1 tag for screen readers only: [#1991](https://github.com/owncloud/ocis/pull/1991)

   https://github.com/owncloud/ocis/pull/1991

* Enhancement - User Deprovisioning for the OCS API: [#1962](https://github.com/owncloud/ocis/pull/1962)

   Use the CS3 API and Reva to deprovision users completely.

   Two new environment variables introduced: ``` OCS_IDM_ADDRESS OCS_STORAGE_USERS_DRIVER
   ```

   `OCS_IDM_ADDRESS` is also an alias for `OCIS_URL`; allows the OCS service to mint jwt tokens
   for the authenticated user that will be read by the reva authentication middleware.

   `OCS_STORAGE_USERS_DRIVER` determines how a user is deprovisioned. This kind of behavior is
   needed since every storage driver deals with deleting differently.

   https://github.com/owncloud/ocis/pull/1962

* Enhancement - Support thumbnails for txt files: [#1988](https://github.com/owncloud/ocis/pull/1988)

   Implemented support for thumbnails for txt files in the thumbnails service.

   https://github.com/owncloud/ocis/pull/1988

* Enhancement - Update reva to v1.7.1-0.20210430154404-69bd21f2cc97: [#2010](https://github.com/owncloud/ocis/pull/2010)

  * Fix recycle to different locations (https://github.com/cs3org/reva/pull/1541)
  * Fix user share as grantee in json backend (https://github.com/cs3org/reva/pull/1650)
  * Introduce named services (https://github.com/cs3org/reva/pull/1509)
  * Improve json marshalling of share protobuf messages (https://github.com/cs3org/reva/pull/1655)
  * Cache resources from share getter methods in OCS (https://github.com/cs3org/reva/pull/1643)
  * Fix public file shares (https://github.com/cs3org/reva/pull/1666)

   https://github.com/owncloud/ocis/pull/2010

* Enhancement - Update reva to v1.7.1-0.20210507160327-e2c3841d0dbc: [#2044](https://github.com/owncloud/ocis/pull/2044)

  * Add user profile picture to capabilities (https://github.com/cs3org/reva/pull/1694)
  * Mint scope-based access tokens for RBAC (https://github.com/cs3org/reva/pull/1669)
  * Add cache warmup strategy for OCS resource infos (https://github.com/cs3org/reva/pull/1664)
  * Filter shares based on type in OCS (https://github.com/cs3org/reva/pull/1683)

   https://github.com/owncloud/ocis/pull/2044

* Enhancement - Use oc-select: [#1979](https://github.com/owncloud/ocis/pull/1979)

   Replace oc-drop with oc select in settings

   https://github.com/owncloud/ocis/pull/1979

* Enhancement - Set SameSite settings to Strict for Web: [#2019](https://github.com/owncloud/ocis/pull/2019)

   Changed SameSite settings to Strict for Web to prevent warnings in Firefox

   https://github.com/owncloud/ocis/pull/2019
# Changelog for [1.5.0] (2021-04-21)

The following sections list the changes for 1.5.0.

[1.5.0]: https://github.com/owncloud/ocis/compare/v1.4.0...v1.5.0

## Summary

* Bugfix - Fixes "unaligned 64-bit atomic operation" panic on 32-bit ARM: [#1888](https://github.com/owncloud/ocis/pull/1888)
* Change - Make Protobuf package names unique: [#1875](https://github.com/owncloud/ocis/pull/1875)
* Change - Update ownCloud Web to v3.0.0: [#1938](https://github.com/owncloud/ocis/pull/1938)
* Enhancement - Change default path for thumbnails: [#1892](https://github.com/owncloud/ocis/pull/1892)
* Enhancement - Parse config on supervised mode with run subcommand: [#1931](https://github.com/owncloud/ocis/pull/1931)
* Enhancement - Update ODS in accounts & settings extension: [#1934](https://github.com/owncloud/ocis/pull/1934)
* Enhancement - Add config for public share SQL driver: [#1916](https://github.com/owncloud/ocis/pull/1916)
* Enhancement - Remove dead runtime code: [#1923](https://github.com/owncloud/ocis/pull/1923)
* Enhancement - Add option to reading registry rules from json file: [#1917](https://github.com/owncloud/ocis/pull/1917)
* Enhancement - Update reva to v1.6.1-0.20210414111318-a4b5148cbfb2: [#1872](https://github.com/owncloud/ocis/pull/1872)

## Details

* Bugfix - Fixes "unaligned 64-bit atomic operation" panic on 32-bit ARM: [#1888](https://github.com/owncloud/ocis/pull/1888)

   Sync/cache had uint64s that were not 64-bit aligned causing panics on 32-bit systems during
   atomic access

   https://github.com/owncloud/ocis/issues/1887
   https://github.com/owncloud/ocis/pull/1888

* Change - Make Protobuf package names unique: [#1875](https://github.com/owncloud/ocis/pull/1875)

   Introduce unique `package` and `go_package` names for our Protobuf definitions

   https://github.com/owncloud/ocis/pull/1875

* Change - Update ownCloud Web to v3.0.0: [#1938](https://github.com/owncloud/ocis/pull/1938)

   Tags: web

   We updated ownCloud Web to v3.0.0. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/1938
   https://github.com/owncloud/web/releases/tag/v3.0.0

* Enhancement - Change default path for thumbnails: [#1892](https://github.com/owncloud/ocis/pull/1892)

   Changes the default path for thumbnails from `<os tmp dir>/ocis-thumbnails` to
   `/var/tmp/ocis/thumbnails`

   https://github.com/owncloud/ocis/issues/1891
   https://github.com/owncloud/ocis/pull/1892

* Enhancement - Parse config on supervised mode with run subcommand: [#1931](https://github.com/owncloud/ocis/pull/1931)

   Currenntly it is not possible to parse a single config file from an extension when running on
   supervised mode.

   https://github.com/owncloud/ocis/pull/1931

* Enhancement - Update ODS in accounts & settings extension: [#1934](https://github.com/owncloud/ocis/pull/1934)

   The accounts and settings extensions were updated to reflect the latest changes in the
   ownCloud design system. In addition, a couple of quick wins in terms of accessibility are
   included.

   https://github.com/owncloud/ocis/pull/1934

* Enhancement - Add config for public share SQL driver: [#1916](https://github.com/owncloud/ocis/pull/1916)

   https://github.com/owncloud/ocis/pull/1916

* Enhancement - Remove dead runtime code: [#1923](https://github.com/owncloud/ocis/pull/1923)

   When moving from the old runtime to the new one there were lots of files left behind that are
   essentially dead code and should be removed. The original code lives here
   github.com/refs/pman/ if someone finds it interesting to read.

   https://github.com/owncloud/ocis/pull/1923

* Enhancement - Add option to reading registry rules from json file: [#1917](https://github.com/owncloud/ocis/pull/1917)

   https://github.com/owncloud/ocis/pull/1917

* Enhancement - Update reva to v1.6.1-0.20210414111318-a4b5148cbfb2: [#1872](https://github.com/owncloud/ocis/pull/1872)

  * enforce quota (https://github.com/cs3org/reva/pull/1557)
  * Make additional info attribute configureable (https://github.com/cs3org/reva/pull/1588)
  * check ENOTDIR for readlink (https://github.com/cs3org/reva/pull/1597)
  * Add wrappers for EOS and EOS Home storage drivers (https://github.com/cs3org/reva/pull/1624)
  * eos: fixes for enabling file sharing (https://github.com/cs3org/reva/pull/1619)
  * implement checksums in the owncloud storage driver (https://github.com/cs3org/reva/pull/1629)

   https://github.com/owncloud/ocis/pull/1872
# Changelog for [1.4.0] (2021-03-30)

The following sections list the changes for 1.4.0.

[1.4.0]: https://github.com/owncloud/ocis/compare/v1.3.0...v1.4.0

## Summary

* Bugfix - Fix thumbnail generation for jpegs: [#1785](https://github.com/owncloud/ocis/pull/1785)
* Change - Update ownCloud Web to v2.1.0: [#1870](https://github.com/owncloud/ocis/pull/1870)
* Enhancement - Add focus to input elements on login page: [#1792](https://github.com/owncloud/ocis/pull/1792)
* Enhancement - Improve accessibility to input elements on login page: [#1794](https://github.com/owncloud/ocis/pull/1794)
* Enhancement - Add new build targets: [#1824](https://github.com/owncloud/ocis/pull/1824)
* Enhancement - Clarify expected failures: [#1790](https://github.com/owncloud/ocis/pull/1790)
* Enhancement - Replace special character in login page title with a regular minus: [#1813](https://github.com/owncloud/ocis/pull/1813)
* Enhancement - File Logging: [#1816](https://github.com/owncloud/ocis/pull/1816)
* Enhancement - Runtime Hostname and Port are now configurable: [#1822](https://github.com/owncloud/ocis/pull/1822)
* Enhancement - Generate thumbnails for .gif files: [#1791](https://github.com/owncloud/ocis/pull/1791)
* Enhancement - Tracing Refactor: [#1819](https://github.com/owncloud/ocis/pull/1819)
* Enhancement - Update reva to v1.6.1-0.20210326165326-e8a00d9b2368: [#1683](https://github.com/owncloud/ocis/pull/1683)

## Details

* Bugfix - Fix thumbnail generation for jpegs: [#1785](https://github.com/owncloud/ocis/pull/1785)

   Images with the extension `.jpeg` were not properly supported.

   https://github.com/owncloud/ocis/issues/1490
   https://github.com/owncloud/ocis/pull/1785

* Change - Update ownCloud Web to v2.1.0: [#1870](https://github.com/owncloud/ocis/pull/1870)

   Tags: web

   We updated ownCloud Web to v2.1.0. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/1870
   https://github.com/owncloud/web/releases/tag/v2.1.0

* Enhancement - Add focus to input elements on login page: [#1792](https://github.com/owncloud/ocis/pull/1792)

   https://github.com/owncloud/web/issues/4322
   https://github.com/owncloud/ocis/pull/1792

* Enhancement - Improve accessibility to input elements on login page: [#1794](https://github.com/owncloud/ocis/pull/1794)

   https://github.com/owncloud/web/issues/4319
   https://github.com/owncloud/ocis/pull/1794
   https://github.com/owncloud/ocis/pull/1811

* Enhancement - Add new build targets: [#1824](https://github.com/owncloud/ocis/pull/1824)

   Make build target `build` used to build a binary twice, the second occurrence having symbols
   for debugging. We split this step in two and added `build-all` and `build-debug` targets.

   - `build-all` now behaves as the previous `build` target, it will generate 2 binaries, one for
   debug. - `build-debug` will build a single binary for debugging.

   https://github.com/owncloud/ocis/pull/1824

* Enhancement - Clarify expected failures: [#1790](https://github.com/owncloud/ocis/pull/1790)

   Some features, while covered by the ownCloud 10 acceptance tests, will not be implmented for
   now: - blacklisted / ignored files, because ocis does not need to blacklist `.htaccess` files -
   `OC-LazyOps` support was [removed from the
   clients](https://github.com/owncloud/client/pull/8398). We are thinking about [a state
   machine for uploads to properly solve that scenario and also list the state of files in progress
   in the web ui](https://github.com/owncloud/ocis/issues/214). The expected failures
   files now have a dedicated _Won't fix_ section for these items.

   https://github.com/owncloud/ocis/issues/214
   https://github.com/owncloud/ocis/pull/1790
   https://github.com/owncloud/client/pull/8398

* Enhancement - Replace special character in login page title with a regular minus: [#1813](https://github.com/owncloud/ocis/pull/1813)

   https://github.com/owncloud/ocis/pull/1813

* Enhancement - File Logging: [#1816](https://github.com/owncloud/ocis/pull/1816)

   When running supervised, support for configuring all logs to a single log file:
   `OCIS_LOG_FILE=/Users/foo/bar/ocis.log MICRO_REGISTRY=etcd bin/ocis server`

   Supports directing log from single extensions to a log file:
   `PROXY_LOG_FILE=/Users/foo/bar/proxy.log MICRO_REGISTRY=etcd bin/ocis proxy`

   https://github.com/owncloud/ocis/pull/1816

* Enhancement - Runtime Hostname and Port are now configurable: [#1822](https://github.com/owncloud/ocis/pull/1822)

   Without any configuration the ocis runtime will start on `localhost:9250` unless specified
   otherwise. Usage:

   - `OCIS_RUNTIME_PORT=6061 bin/ocis server` - overrides the oCIS runtime and starts on port
   6061 - `OCIS_RUNTIME_PORT=6061 bin/ocis list` - lists running extensions for the runtime on
   `localhost:6061`

   All subcommands are updated and expected to work with the following environment variables:

   ``` OCIS_RUNTIME_HOST OCIS_RUNTIME_PORT ```

   https://github.com/owncloud/ocis/pull/1822

* Enhancement - Generate thumbnails for .gif files: [#1791](https://github.com/owncloud/ocis/pull/1791)

   Added support for gifs to the thumbnails service.

   https://github.com/owncloud/ocis/pull/1791

* Enhancement - Tracing Refactor: [#1819](https://github.com/owncloud/ocis/pull/1819)

   Centralize tracing handling per extension.

   https://github.com/owncloud/ocis/pull/1819

* Enhancement - Update reva to v1.6.1-0.20210326165326-e8a00d9b2368: [#1683](https://github.com/owncloud/ocis/pull/1683)

  * quota querying and tree accounting [cs3org/reva#1405](https://github.com/cs3org/reva/pull/1405)
  * Fix webdav file versions endpoint bugs [cs3org/reva#1526](https://github.com/cs3org/reva/pull/1526)
  * Fix etag changing only once a second [cs3org/reva#1576](https://github.com/cs3org/reva/pull/1576)
  * Trashbin API parity [cs3org/reva#1552](https://github.com/cs3org/reva/pull/1552)
  * Signature authentication for public links [cs3org/reva#1590](https://github.com/cs3org/reva/pull/1590)

   https://github.com/owncloud/ocis/pull/1683
   https://github.com/cs3org/reva/pull/1405
   https://github.com/owncloud/ocis/pull/1861
# Changelog for [1.3.0] (2021-03-09)

The following sections list the changes for 1.3.0.

[1.3.0]: https://github.com/owncloud/ocis/compare/v1.2.0...v1.3.0

## Summary

* Bugfix - Purposely delay accounts service startup: [#1734](https://github.com/owncloud/ocis/pull/1734)
* Bugfix - Add missing gateway config: [#1716](https://github.com/owncloud/ocis/pull/1716)
* Bugfix - Fix accounts initialization: [#1696](https://github.com/owncloud/ocis/pull/1696)
* Bugfix - Fix the ttl of the authentication middleware cache: [#1699](https://github.com/owncloud/ocis/pull/1699)
* Change - Update ownCloud Web to v2.0.1: [#1683](https://github.com/owncloud/ocis/pull/1683)
* Change - Update ownCloud Web to v2.0.2: [#1776](https://github.com/owncloud/ocis/pull/1776)
* Enhancement - Remove the JWT from the log: [#1758](https://github.com/owncloud/ocis/pull/1758)
* Enhancement - Update go-micro to v3.5.1-0.20210217182006-0f0ace1a44a9: [#1670](https://github.com/owncloud/ocis/pull/1670)
* Enhancement - Update reva to v1.6.1-0.20210223065028-53f39499762e: [#1683](https://github.com/owncloud/ocis/pull/1683)
* Enhancement - Add initial nats and kubernetes registry support: [#1697](https://github.com/owncloud/ocis/pull/1697)

## Details

* Bugfix - Purposely delay accounts service startup: [#1734](https://github.com/owncloud/ocis/pull/1734)

   As it turns out the race condition between `accounts <-> storage-metadata` still remains.
   This PR is a hotfix, and it should be followed up with a proper fix. Either:

   - block the accounts' initialization until the storage metadata is ready (using the registry)
   or - allow the accounts service to initialize and use a message broker to signal the accounts the
   metadata storage is ready to receive requests.

   https://github.com/owncloud/ocis/pull/1734

* Bugfix - Add missing gateway config: [#1716](https://github.com/owncloud/ocis/pull/1716)

   The auth provider `ldap` and `oidc` drivers now need to be able talk to the reva gateway. We added
   the `gatewayscv` to the config that is passed to reva.

   https://github.com/owncloud/ocis/pull/1716

* Bugfix - Fix accounts initialization: [#1696](https://github.com/owncloud/ocis/pull/1696)

   Originally the accounts service relies on both the `settings` and `storage-metadata` to be up
   and running at the moment it starts. This is an antipattern as it will cause the entire service to
   panic if the dependants are not present.

   We inverted this dependency and moved the default initialization data (i.e: creating roles,
   permissions, settings bundles) and instead of notifying the settings service that the
   account has to provide with such options, the settings is instead initialized with the options
   the accounts rely on. Essentially saving bandwith as there is no longer a gRPC call to the
   settings service.

   For the `storage-metadata` a retry mechanism was added that retries by default 20 times to
   fetch the `com.owncloud.storage.metadata` from the service registry every `500`
   miliseconds. If this retry expires the accounts panics, as its dependency on the
   `storage-metadata` service cannot be resolved.

   We also introduced a client wrapper that acts as middleware between a client and a server. For
   more information on how it works further read [here](https://github.com/sony/gobreaker)

   https://github.com/owncloud/ocis/pull/1696

* Bugfix - Fix the ttl of the authentication middleware cache: [#1699](https://github.com/owncloud/ocis/pull/1699)

   The authentication cache ttl was multiplied with `time.Second` multiple times. This
   resulted in a ttl that was not intended.

   https://github.com/owncloud/ocis/pull/1699

* Change - Update ownCloud Web to v2.0.1: [#1683](https://github.com/owncloud/ocis/pull/1683)

   Tags: web

   We updated ownCloud Web to v2.0.1. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/1683
   https://github.com/owncloud/web/releases/tag/v2.0.1

* Change - Update ownCloud Web to v2.0.2: [#1776](https://github.com/owncloud/ocis/pull/1776)

   Tags: web

   We updated ownCloud Web to v2.0.2. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/1776
   https://github.com/owncloud/web/releases/tag/v2.0.2

* Enhancement - Remove the JWT from the log: [#1758](https://github.com/owncloud/ocis/pull/1758)

   We were logging the JWT in some places. Secrets should not be exposed in logs so it got removed.

   https://github.com/owncloud/ocis/pull/1758

* Enhancement - Update go-micro to v3.5.1-0.20210217182006-0f0ace1a44a9: [#1670](https://github.com/owncloud/ocis/pull/1670)

   - We updated from go micro v2 (v2.9.1) go-micro v3 (v3.5.1 edge). - oCIS runtime is now aware of
   `MICRO_LOG_LEVEL` and is set to `error` by default. This decision was made because ownCloud,
   as framework builders, want to log everything oCIS related and hide everything unrelated by
   default. It can be re-enabled by setting it to a log level other than `error`. i.e:
   `MICRO_LOG_LEVEL=info`. - Updated `protoc-gen-micro` to the [latest
   version](https://github.com/asim/go-micro/tree/master/cmd/protoc-gen-micro). -
   We're using Prometheus wrappers from go-micro.

   https://github.com/owncloud/ocis/pull/1670
   https://github.com/asim/go-micro/pull/2126

* Enhancement - Update reva to v1.6.1-0.20210223065028-53f39499762e: [#1683](https://github.com/owncloud/ocis/pull/1683)

  * quota querying and tree accounting [cs3org/reva#1405](https://github.com/cs3org/reva/pull/1405)

   https://github.com/owncloud/ocis/pull/1683
   https://github.com/cs3org/reva/pull/1405

* Enhancement - Add initial nats and kubernetes registry support: [#1697](https://github.com/owncloud/ocis/pull/1697)

   We added initial support to use nats and kubernetes as a service registry using
   `MICRO_REGISTRY=nats` and `MICRO_REGISTRY=kubernetes` respectively. Multiple nodes can
   be given with `MICRO_REGISTRY_ADDRESS=1.2.3.4,5.6.7.8,9.10.11.12`.

   https://github.com/owncloud/ocis/pull/1697
# Changelog for [1.2.0] (2021-02-17)

The following sections list the changes for 1.2.0.

[1.2.0]: https://github.com/owncloud/ocis/compare/v1.1.0...v1.2.0

## Summary

* Bugfix - Check if roles are present in user object before looking those up: [#1388](https://github.com/owncloud/ocis/pull/1388)
* Bugfix - Fix etcd address configuration: [#1546](https://github.com/owncloud/ocis/pull/1546)
* Bugfix - Remove unimplemented config file option for oCIS root command: [#1636](https://github.com/owncloud/ocis/pull/1636)
* Bugfix - Fix thumbnail generation when using different idp: [#1624](https://github.com/owncloud/ocis/issues/1624)
* Change - Initial release of graph and graph explorer: [#1594](https://github.com/owncloud/ocis/pull/1594)
* Change - Move runtime code on refs/pman over to owncloud/ocis/ocis: [#1483](https://github.com/owncloud/ocis/pull/1483)
* Change - Update ownCloud Web to v2.0.0: [#1661](https://github.com/owncloud/ocis/pull/1661)
* Enhancement - Make use of new design-system oc-table: [#1597](https://github.com/owncloud/ocis/pull/1597)
* Enhancement - Use a default protocol parameter instead of explicitly disabling tus: [#1331](https://github.com/cs3org/reva/pull/1331)
* Enhancement - Functionality to map home directory to different storage providers: [#1186](https://github.com/owncloud/ocis/pull/1186)
* Enhancement - Introduce ADR: [#1042](https://github.com/owncloud/ocis/pull/1042)
* Enhancement - Switch to opencontainers annotation scheme: [#1381](https://github.com/owncloud/ocis/pull/1381)
* Enhancement - Migrate ocis-graph-explorer to ocis monorepo: [#1596](https://github.com/owncloud/ocis/pull/1596)
* Enhancement - Migrate ocis-graph to ocis monorepo: [#1594](https://github.com/owncloud/ocis/pull/1594)
* Enhancement - Enable group sharing and add config for sharing SQL driver: [#1626](https://github.com/owncloud/ocis/pull/1626)
* Enhancement - Update reva to v1.5.2-0.20210125114636-0c10b333ee69: [#1482](https://github.com/owncloud/ocis/pull/1482)

## Details

* Bugfix - Check if roles are present in user object before looking those up: [#1388](https://github.com/owncloud/ocis/pull/1388)

   https://github.com/owncloud/ocis/pull/1388

* Bugfix - Fix etcd address configuration: [#1546](https://github.com/owncloud/ocis/pull/1546)

   The etcd server address in `MICRO_REGISTRY_ADDRESS` was not picked up when etcd was set as
   service discovery registry `MICRO_REGISTRY=etcd`. Therefore etcd was only working if
   available on localhost / 127.0.0.1.

   https://github.com/owncloud/ocis/pull/1546

* Bugfix - Remove unimplemented config file option for oCIS root command: [#1636](https://github.com/owncloud/ocis/pull/1636)

   https://github.com/owncloud/ocis/pull/1636

* Bugfix - Fix thumbnail generation when using different idp: [#1624](https://github.com/owncloud/ocis/issues/1624)

   The thumbnail service was relying on a konnectd specific field in the access token. This logic
   was now replaced by a service parameter for the username.

   https://github.com/owncloud/ocis/issues/1624
   https://github.com/owncloud/ocis/pull/1628

* Change - Initial release of graph and graph explorer: [#1594](https://github.com/owncloud/ocis/pull/1594)

   Tags: graph, graph-explorer

   We brought initial basic Graph and Graph-Explorer support for the ownCloud Infinite Scale
   project.

   https://github.com/owncloud/ocis/pull/1594
   https://github.com/owncloud/ocis-graph-explorer/pull/3

* Change - Move runtime code on refs/pman over to owncloud/ocis/ocis: [#1483](https://github.com/owncloud/ocis/pull/1483)

   Tags: ocis, runtime

   Currently, the runtime is under the private account of an oCIS developer. For future-proofing
   we don't want oCIS mission critical components to depend on external repositories, so we're
   including refs/pman module as an oCIS package instead.

   https://github.com/owncloud/ocis/pull/1483

* Change - Update ownCloud Web to v2.0.0: [#1661](https://github.com/owncloud/ocis/pull/1661)

   Tags: web

   We updated ownCloud Web to v2.0.0. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/1661
   https://github.com/owncloud/web/releases/tag/v2.0.0

* Enhancement - Make use of new design-system oc-table: [#1597](https://github.com/owncloud/ocis/pull/1597)

   Tags: ui, accounts

   The design-system table component has changed the way it's used. We updated accounts-ui to use
   the new 'oc-table-simple' component.

   https://github.com/owncloud/ocis/pull/1597

* Enhancement - Use a default protocol parameter instead of explicitly disabling tus: [#1331](https://github.com/cs3org/reva/pull/1331)

   https://github.com/cs3org/reva/pull/1331
   https://github.com/owncloud/ocis/pull/1374

* Enhancement - Functionality to map home directory to different storage providers: [#1186](https://github.com/owncloud/ocis/pull/1186)

   We added a parameter in reva that allows us to redirect /home requests to different storage
   providers based on a mapping derived from the user attributes, which was previously not
   possible since we hardcode the /home path for all users. For example, having its value as
   `/home/{{substr 0 1 .Username}}` can be used to redirect home requests for different users to
   different storage providers.

   https://github.com/owncloud/ocis/pull/1186
   https://github.com/cs3org/reva/pull/1142

* Enhancement - Introduce ADR: [#1042](https://github.com/owncloud/ocis/pull/1042)

   We will keep track of [Architectual Decision Records using
   Markdown](https://adr.github.io/madr/) in `/docs/adr`.

   https://github.com/owncloud/ocis/pull/1042

* Enhancement - Switch to opencontainers annotation scheme: [#1381](https://github.com/owncloud/ocis/pull/1381)

   Switch docker image annotation scheme to org.opencontainers standard because
   org.label-schema is depreciated.

   https://github.com/owncloud/ocis/pull/1381

* Enhancement - Migrate ocis-graph-explorer to ocis monorepo: [#1596](https://github.com/owncloud/ocis/pull/1596)

   Tags: ocis, ocis-graph-explorer

   Ocis-graph-explorer was not migrated during the monorepo conversion.

   https://github.com/owncloud/ocis/pull/1596

* Enhancement - Migrate ocis-graph to ocis monorepo: [#1594](https://github.com/owncloud/ocis/pull/1594)

   Tags: ocis, ocis-graph

   Ocis-graph was not migrated during the monorepo conversion.

   https://github.com/owncloud/ocis/pull/1594

* Enhancement - Enable group sharing and add config for sharing SQL driver: [#1626](https://github.com/owncloud/ocis/pull/1626)

   This PR adds config to support sharing with groups. It also introduces a breaking change for the
   CS3APIs definitions since grantees can now refer to both users as well as groups. Since we store
   the grantee information in a json file, `/var/tmp/ocis/storage/shares.json`, its previous
   version needs to be removed as we won't be able to unmarshal data corresponding to the previous
   definitions.

   https://github.com/owncloud/ocis/pull/1626
   https://github.com/cs3org/reva/pull/1453

* Enhancement - Update reva to v1.5.2-0.20210125114636-0c10b333ee69: [#1482](https://github.com/owncloud/ocis/pull/1482)

  * initial checksum support for ocis [cs3org/reva#1400](https://github.com/cs3org/reva/pull/1400)
  * Use updated etag of home directory even if it is cached [cs3org/reva#1416](https://github.com/cs3org/reva/pull/#1416)
  * Indicate in EOS containers that TUS is not supported [cs3org/reva#1415](https://github.com/cs3org/reva/pull/#1415)
  * Get status code from recycle response [cs3org/reva#1408](https://github.com/cs3org/reva/pull/#1408)

   https://github.com/owncloud/ocis/pull/1482
   https://github.com/cs3org/reva/pull/1400
   https://github.com/cs3org/reva/pull/1416
   https://github.com/cs3org/reva/pull/1415
   https://github.com/cs3org/reva/pull/1408
# Changelog for [1.1.0] (2021-01-22)

The following sections list the changes for 1.1.0.

[1.1.0]: https://github.com/owncloud/ocis/compare/v1.0.0...v1.1.0

## Summary

* Change - Disable pretty logging by default: [#1133](https://github.com/owncloud/ocis/pull/1133)
* Change - Add "volume" declaration to docker images: [#1375](https://github.com/owncloud/ocis/pull/1375)
* Change - Add "expose" information to docker images: [#1366](https://github.com/owncloud/ocis/pull/1366)
* Change - Generate cryptographically secure state token: [#1203](https://github.com/owncloud/ocis/pull/1203)
* Change - Move k6 to cdperf: [#1358](https://github.com/owncloud/ocis/pull/1358)
* Change - Update go version: [#1364](https://github.com/owncloud/ocis/pull/1364)
* Change - Update ownCloud Web to v1.0.1: [#1191](https://github.com/owncloud/ocis/pull/1191)
* Enhancement - Add OCIS_URL env var: [#1148](https://github.com/owncloud/ocis/pull/1148)
* Enhancement - Use sync.cache for roles cache: [#1367](https://github.com/owncloud/ocis/pull/1367)
* Enhancement - Add named locks and refactor cache: [#1212](https://github.com/owncloud/ocis/pull/1212)
* Enhancement - Update reva to v1.5.1: [#1372](https://github.com/owncloud/ocis/pull/1372)
* Enhancement - Update reva to v1.4.1-0.20210111080247-f2b63bfd6825: [#1194](https://github.com/owncloud/ocis/pull/1194)

## Details

* Change - Disable pretty logging by default: [#1133](https://github.com/owncloud/ocis/pull/1133)

   Tags: ocis

   Disable pretty logging default for performance reasons.

   https://github.com/owncloud/ocis/pull/1133

* Change - Add "volume" declaration to docker images: [#1375](https://github.com/owncloud/ocis/pull/1375)

   Tags: docker

   Add "volume" declaration to docker images. This makes it easier for Docker users to see where
   oCIS stores data.

   https://github.com/owncloud/ocis/pull/1375

* Change - Add "expose" information to docker images: [#1366](https://github.com/owncloud/ocis/pull/1366)

   Tags: docker

   Add "expose" information to docker images. Docker users will now see that we offer services on
   port 9200.

   https://github.com/owncloud/ocis/pull/1366

* Change - Generate cryptographically secure state token: [#1203](https://github.com/owncloud/ocis/pull/1203)

   Replaced Math.random with a cryptographically secure way to generate the oidc state token
   using the javascript crypto api.

   https://github.com/owncloud/ocis/pull/1203
   https://developer.mozilla.org/en-US/docs/Web/API/Crypto/getRandomValues
   https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Math/random

* Change - Move k6 to cdperf: [#1358](https://github.com/owncloud/ocis/pull/1358)

   Tags: performance, testing, k6

   The ownCloud performance tests can not only be used to test oCIS. This is why we have decided to
   move the k6 tests to https://github.com/owncloud/cdperf

   https://github.com/owncloud/ocis/pull/1358

* Change - Update go version: [#1364](https://github.com/owncloud/ocis/pull/1364)

   Tags: go

   Update go from 1.13 to 1.15

   https://github.com/owncloud/ocis/pull/1364

* Change - Update ownCloud Web to v1.0.1: [#1191](https://github.com/owncloud/ocis/pull/1191)

   Tags: web

   We updated ownCloud Web to v1.0.1. Please refer to the changelog (linked) for details on the web
   release.

   https://github.com/owncloud/ocis/pull/1191
   https://github.com/owncloud/web/releases/tag/v1.0.1

* Enhancement - Add OCIS_URL env var: [#1148](https://github.com/owncloud/ocis/pull/1148)

   Tags: ocis

   We introduced a new environment variable `OCIS_URL` that expects a URL including protocol,
   host and optionally port to simplify configuring all the different services. These existing
   environment variables still take precedence, but will also fall back to `OCIS_URL`:
   `STORAGE_LDAP_IDP`, `STORAGE_OIDC_ISSUER`, `PROXY_OIDC_ISSUER`,
   `STORAGE_FRONTEND_PUBLIC_URL`, `KONNECTD_ISS`, `WEB_OIDC_AUTHORITY`, and
   `WEB_UI_CONFIG_SERVER`.

   Some environment variables are now built dynamically if they are not set: -
   `STORAGE_DATAGATEWAY_PUBLIC_URL` defaults to `<STORAGE_FRONTEND_PUBLIC_URL>/data`,
   also falling back to `OCIS_URL` - `WEB_OIDC_METADATA_URL` defaults to
   `<WEB_OIDC_AUTHORITY>/.well-known/openid-configuration`, also falling back to
   `OCIS_URL`

   Furthermore, the built in konnectd will generate an `identifier-registration.yaml` that
   uses the `KONNECTD_ISS` in the allowed `redirect_uris` and `origins`. It simplifies the
   default `https://localhost:9200` and remote deployment with `OCIS_URL` which is evaluated
   as a fallback if `KONNECTD_ISS` is not set.

   An oCIS server can now be started on a remote machine as easy as
   `OCIS_URL=https://cloud.ocis.test PROXY_HTTP_ADDR=0.0.0.0:443 ocis server`.

   Note that the `OCIS_DOMAIN` environment variable is not used by oCIS, but by the docker
   containers.

   https://github.com/owncloud/ocis/pull/1148

* Enhancement - Use sync.cache for roles cache: [#1367](https://github.com/owncloud/ocis/pull/1367)

   Tags: ocis-pkg

   Update ocis-pkg/roles cache to use ocis-pkg/sync cache

   https://github.com/owncloud/ocis/pull/1367

* Enhancement - Add named locks and refactor cache: [#1212](https://github.com/owncloud/ocis/pull/1212)

   Tags: ocis-pkg, accounts

   We had the case that we needed kind of a named locking mechanism which enables us to lock only
   under certain conditions. It's used in the indexer package where we do not need to lock
   everything, instead just lock the requested parts and differentiate between reads and
   writes.

   This made it possible to entirely remove locks from the accounts service and move them to the
   ocis-pkg indexer. Another part of this refactor was to make the cache atomic and write tests for
   it.

   - remove locking from accounts service - add sync package with named mutex - add named locking to
   indexer - move cache to sync package

   https://github.com/owncloud/ocis/issues/966
   https://github.com/owncloud/ocis/pull/1212

* Enhancement - Update reva to v1.5.1: [#1372](https://github.com/owncloud/ocis/pull/1372)

   Summary -------

  * Fix #1401: Use the user in request for deciding the layout for non-home DAV requests
  * Fix #1413: Re-include the '.git' dir in the Docker images to pass the version tag
  * Fix #1399: Fix ocis trash-bin purge
  * Enh #1397: Bump the Copyright date to 2021
  * Enh #1398: Support site authorization status in Mentix
  * Enh #1393: Allow setting favorites, mtime and a temporary etag
  * Enh #1403: Support remote cloud gathering metrics

   Details -------

  * Bugfix #1401: Use the user in request for deciding the layout for non-home DAV requests

   For the incoming /dav/files/userID requests, we have different namespaces depending on
   whether the request is for the logged-in user's namespace or not. Since in the storage drivers,
   we specify the layout depending only on the user whose resources are to be accessed, this fails
   when a user wants to access another user's namespace when the storage provider depends on the
   logged in user's namespace. This PR fixes that.

   For example, consider the following case. The owncloud fs uses a layout {{substr 0 1
   .Id.OpaqueId}}/{{.Id.OpaqueId}}. The user einstein sends a request to access a resource
   shared with him, say /dav/files/marie/abcd, which should be allowed. However, based on the
   way we applied the layout, there's no way in which this can be translated to /m/marie/.

   Https://github.com/cs3org/reva/pull/1401

  * Bugfix #1413: Re-include the '.git' dir in the Docker images to pass the version tag

   And git SHA to the release tool.

   Https://github.com/cs3org/reva/pull/1413

  * Bugfix #1399: Fix ocis trash-bin purge

   Fixes the empty trash-bin functionality for ocis-storage

   Https://github.com/owncloud/product/issues/254
   https://github.com/cs3org/reva/pull/1399

  * Enhancement #1397: Bump the Copyright date to 2021

   Https://github.com/cs3org/reva/pull/1397

  * Enhancement #1398: Support site authorization status in Mentix

   This enhancement adds support for a site authorization status to Mentix. This way, sites
   registered via a web app can now be excluded until authorized manually by an administrator.

   Furthermore, Mentix now sets the scheme for Prometheus targets. This allows us to also support
   monitoring of sites that do not support the default HTTPS scheme.

   Https://github.com/cs3org/reva/pull/1398

  * Enhancement #1393: Allow setting favorites, mtime and a temporary etag

   We now let the oCIS driver persist favorites, set temporary etags and the mtime as arbitrary
   metadata.

   Https://github.com/owncloud/ocis/issues/567
   https://github.com/cs3org/reva/issues/1394
   https://github.com/cs3org/reva/pull/1393

  * Enhancement #1403: Support remote cloud gathering metrics

   The current metrics package can only gather metrics either from json files. With this feature,
   the metrics can be gathered polling the http endpoints exposed by the owncloud/nextcloud
   sciencemesh apps.

   Https://github.com/cs3org/reva/pull/1403

   https://github.com/owncloud/ocis/pull/1372

* Enhancement - Update reva to v1.4.1-0.20210111080247-f2b63bfd6825: [#1194](https://github.com/owncloud/ocis/pull/1194)

  * Enhancement: calculate and expose actual file permission set [cs3org/reva#1368](https://github.com/cs3org/reva/pull/1368)
  * initial range request support [cs3org/reva#1326](https://github.com/cs3org/reva/pull/1388)

   https://github.com/owncloud/ocis/pull/1194
   https://github.com/cs3org/reva/pull/1368
   https://github.com/cs3org/reva/pull/1388
# Changelog for [1.0.0] (2020-12-17)

The following sections list the changes for 1.0.0.

## Summary

* Bugfix - Enable scrolling in accounts list: [#909](https://github.com/owncloud/ocis/pull/909)
* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)
* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)
* Bugfix - Lower Bound was not working for the cs3 api index implementation: [#741](https://github.com/owncloud/ocis/pull/741)
* Bugfix - Accounts config sometimes being overwritten: [#808](https://github.com/owncloud/ocis/pull/808)
* Bugfix - Make settings service start without go coroutines: [#835](https://github.com/owncloud/ocis/pull/835)
* Bugfix - Fix button layout after phoenix update: [#625](https://github.com/owncloud/ocis/pull/625)
* Bugfix - Fix choose account dialogue: [#846](https://github.com/owncloud/ocis/pull/846)
* Bugfix - Fix id or username query handling: [#745](https://github.com/owncloud/ocis/pull/745)
* Bugfix - Fix konnectd build: [#809](https://github.com/owncloud/ocis/pull/809)
* Bugfix - Fix path of files shared with me in ocs api: [#204](https://github.com/owncloud/product/issues/204)
* Bugfix - Use micro default client: [#718](https://github.com/owncloud/ocis/pull/718)
* Bugfix - Allow consent-prompt with switch-account: [#788](https://github.com/owncloud/ocis/pull/788)
* Bugfix - Mint token with uid and gid: [#737](https://github.com/owncloud/ocis/pull/737)
* Bugfix - Serve index.html for directories: [#912](https://github.com/owncloud/ocis/pull/912)
* Bugfix - Don't create account if id/mail/username already taken: [#709](https://github.com/owncloud/ocis/pull/709)
* Bugfix - Fix director selection in proxy: [#521](https://github.com/owncloud/ocis/pull/521)
* Bugfix - Permission checks for settings write access: [#1092](https://github.com/owncloud/ocis/pull/1092)
* Bugfix - Fix minor ui bugs: [#1043](https://github.com/owncloud/ocis/issues/1043)
* Bugfix - Disable public link expiration by default: [#987](https://github.com/owncloud/ocis/issues/987)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)
* Change - Accounts UI shows message when no permissions: [#656](https://github.com/owncloud/ocis/pull/656)
* Change - Cache password validation: [#958](https://github.com/owncloud/ocis/pull/958)
* Change - Filesystem based index: [#709](https://github.com/owncloud/ocis/pull/709)
* Change - Rebuild index command for accounts: [#748](https://github.com/owncloud/ocis/pull/748)
* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)
* Change - CS3 can be used as accounts-backend: [#1020](https://github.com/owncloud/ocis/pull/1020)
* Change - Use bcrypt to hash the user passwords: [#510](https://github.com/owncloud/ocis/issues/510)
* Change - Replace the library which scales the images: [#910](https://github.com/owncloud/ocis/pull/910)
* Change - Choose disk or cs3 storage for accounts and groups: [#623](https://github.com/owncloud/ocis/pull/623)
* Change - Enable OpenID dynamic client registration: [#811](https://github.com/owncloud/ocis/issues/811)
* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)
* Change - Improve reva service descriptions: [#536](https://github.com/owncloud/ocis/pull/536)
* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)
* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)
* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)
* Change - Properly style konnectd consent page: [#754](https://github.com/owncloud/ocis/pull/754)
* Change - Make all paths configurable and default to a common temp dir: [#1080](https://github.com/owncloud/ocis/pull/1080)
* Change - Move the indexer package from ocis/accounts to ocis/ocis-pkg: [#794](https://github.com/owncloud/ocis/pull/794)
* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Move ocis default config to root level: [#842](https://github.com/owncloud/ocis/pull/842)
* Change - Remove username field in OCS: [#709](https://github.com/owncloud/ocis/pull/709)
* Change - Account management permissions for Admin role: [#124](https://github.com/owncloud/product/issues/124)
* Change - Update phoenix to v0.18.0: [#651](https://github.com/owncloud/ocis/pull/651)
* Change - Default apps in ownCloud Web: [#688](https://github.com/owncloud/ocis/pull/688)
* Change - Proxy allow insecure upstreams: [#1007](https://github.com/owncloud/ocis/pull/1007)
* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)
* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)
* Change - Theme welcome and choose account pages: [#887](https://github.com/owncloud/ocis/pull/887)
* Change - Bring oC theme: [#698](https://github.com/owncloud/ocis/pull/698)
* Change - Unify Configuration Parsing: [#675](https://github.com/owncloud/ocis/pull/675)
* Change - Update phoenix to v0.20.0: [#674](https://github.com/owncloud/ocis/pull/674)
* Change - Update phoenix to v0.21.0: [#728](https://github.com/owncloud/ocis/pull/728)
* Change - Update phoenix to v0.22.0: [#757](https://github.com/owncloud/ocis/pull/757)
* Change - Update phoenix to v0.23.0: [#785](https://github.com/owncloud/ocis/pull/785)
* Change - Update phoenix to v0.24.0: [#817](https://github.com/owncloud/ocis/pull/817)
* Change - Update phoenix to v0.25.0: [#868](https://github.com/owncloud/ocis/pull/868)
* Change - Update phoenix to v0.26.0: [#935](https://github.com/owncloud/ocis/pull/935)
* Change - Update phoenix to v0.27.0: [#943](https://github.com/owncloud/ocis/pull/943)
* Change - Update phoenix to v0.28.0: [#1027](https://github.com/owncloud/ocis/pull/1027)
* Change - Update phoenix to v0.29.0: [#1034](https://github.com/owncloud/ocis/pull/1034)
* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)
* Change - Update reva to v1.4.1-0.20201209113234-e791b5599a89: [#1089](https://github.com/owncloud/ocis/pull/1089)
* Change - Clarify storage driver env vars: [#729](https://github.com/owncloud/ocis/pull/729)
* Change - Update ownCloud Web to v1.0.0-beta3: [#1105](https://github.com/owncloud/ocis/pull/1105)
* Change - Update ownCloud Web to v1.0.0-beta4: [#1110](https://github.com/owncloud/ocis/pull/1110)
* Change - Settings and accounts appear in the user menu: [#656](https://github.com/owncloud/ocis/pull/656)
* Enhancement - Add tracing to the accounts service: [#1016](https://github.com/owncloud/ocis/issues/1016)
* Enhancement - Add the accounts service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add basic auth option: [#627](https://github.com/owncloud/ocis/pull/627)
* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)
* Enhancement - Add the glauth service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add k6: [#941](https://github.com/owncloud/ocis/pull/941)
* Enhancement - Add the konnectd service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocis-phoenix service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocis-pkg package: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the ocs service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the proxy service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the settings service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the storage service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the store service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add the thumbnails service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Add a command to list the versions of running instances: [#226](https://github.com/owncloud/product/issues/226)
* Enhancement - Add the webdav service: [#244](https://github.com/owncloud/product/issues/244)
* Enhancement - Better adopt Go-Micro: [#840](https://github.com/owncloud/ocis/pull/840)
* Enhancement - Add permission check when assigning and removing roles: [#879](https://github.com/owncloud/ocis/issues/879)
* Enhancement - Create OnlyOffice extension: [#857](https://github.com/owncloud/ocis/pull/857)
* Enhancement - Show basic-auth warning only once: [#886](https://github.com/owncloud/ocis/pull/886)
* Enhancement - Add glauth fallback backend: [#649](https://github.com/owncloud/ocis/pull/649)
* Enhancement - Tidy dependencies: [#845](https://github.com/owncloud/ocis/pull/845)
* Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)
* Enhancement - Add a version command to ocis: [#915](https://github.com/owncloud/ocis/pull/915)
* Enhancement - Create a proxy access-log: [#889](https://github.com/owncloud/ocis/pull/889)
* Enhancement - Cache userinfo in proxy: [#877](https://github.com/owncloud/ocis/pull/877)
* Enhancement - Update reva to v1.4.1-0.20201125144025-57da0c27434c: [#1320](https://github.com/cs3org/reva/pull/1320)
* Enhancement - Runtime Cleanup: [#1066](https://github.com/owncloud/ocis/pull/1066)
* Enhancement - Update OCIS Runtime: [#1108](https://github.com/owncloud/ocis/pull/1108)
* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)
* Enhancement - Update glauth to dev fd3ac7e4bbdc93578655d9a08d8e23f105aaa5b2: [#834](https://github.com/owncloud/ocis/pull/834)
* Enhancement - Update glauth to dev 4f029234b2308: [#786](https://github.com/owncloud/ocis/pull/786)
* Enhancement - Update konnectd to v0.33.8: [#744](https://github.com/owncloud/ocis/pull/744)
* Enhancement - Update reva to v1.4.1-0.20201123062044-b2c4af4e897d: [#823](https://github.com/owncloud/ocis/pull/823)
* Enhancement - Update reva to v1.4.1-0.20201130061320-ac85e68e0600: [#980](https://github.com/owncloud/ocis/pull/980)
* Enhancement - Update reva to cdb3d6688da5: [#748](https://github.com/owncloud/ocis/pull/748)
* Enhancement - Update reva to dd3a8c0f38: [#725](https://github.com/owncloud/ocis/pull/725)
* Enhancement - Update reva to v1.4.1-0.20201127111856-e6a6212c1b7b: [#971](https://github.com/owncloud/ocis/pull/971)
* Enhancement - Update reva to 063b3db9162b: [#1091](https://github.com/owncloud/ocis/pull/1091)
* Enhancement - Add www-authenticate based on user agent: [#1009](https://github.com/owncloud/ocis/pull/1009)

## Details

* Bugfix - Enable scrolling in accounts list: [#909](https://github.com/owncloud/ocis/pull/909)

   Tags: accounts

   We've fixed the accounts list to enable scrolling.

   https://github.com/owncloud/ocis/pull/909

* Bugfix - Add missing env vars to docker compose: [#392](https://github.com/owncloud/ocis/pull/392)

   Tags: docker

   Without setting `REVA_FRONTEND_URL` and `REVA_DATAGATEWAY_URL` uploads would default to
   locahost and fail if `OCIS_DOMAIN` was used to run ocis on a remote host.

   https://github.com/owncloud/ocis/pull/392

* Bugfix - Don't enforce empty external apps slice: [#473](https://github.com/owncloud/ocis/pull/473)

   Tags: web

   The command for ocis-phoenix enforced an empty external apps configuration. This was
   removed, as it was blocking a new set of default external apps in ocis-phoenix.

   https://github.com/owncloud/ocis/pull/473

* Bugfix - Lower Bound was not working for the cs3 api index implementation: [#741](https://github.com/owncloud/ocis/pull/741)

   Tags: accounts

   Lower bound working on the cs3 index implementation

   https://github.com/owncloud/ocis/pull/741

* Bugfix - Accounts config sometimes being overwritten: [#808](https://github.com/owncloud/ocis/pull/808)

   Tags: accounts

   Sometimes when running the accounts extensions flags were not being taken into
   consideration.

   https://github.com/owncloud/ocis/pull/808

* Bugfix - Make settings service start without go coroutines: [#835](https://github.com/owncloud/ocis/pull/835)

   The go routines cause a race condition that sometimes causes the tests to fail. The ListRoles
   request would not return all permissions.

   https://github.com/owncloud/ocis/pull/835

* Bugfix - Fix button layout after phoenix update: [#625](https://github.com/owncloud/ocis/pull/625)

   Tags: accounts

   With the phoenix update to v0.17.0 a new ODS version was released which has a breaking change for
   buttons regarding their layouting. We adjusted the button layout in the accounts UI
   accordingly.

   https://github.com/owncloud/ocis/pull/625

* Bugfix - Fix choose account dialogue: [#846](https://github.com/owncloud/ocis/pull/846)

   Tags: konnectd

   We've fixed the choose account dialogue in konnectd bug that the user hasn't been logged in
   after selecting account.

   https://github.com/owncloud/ocis/pull/846

* Bugfix - Fix id or username query handling: [#745](https://github.com/owncloud/ocis/pull/745)

   Tags: accounts

   The code was stopping execution when encountering an error while loading an account by id. But
   for or queries we can continue execution.

   https://github.com/owncloud/ocis/pull/745

* Bugfix - Fix konnectd build: [#809](https://github.com/owncloud/ocis/pull/809)

   Tags: konnectd

   We fixed the default config for konnectd and updated the Makefile to include the `yarn
   install`and `yarn build` steps if the static assets are missing.

   https://github.com/owncloud/ocis/pull/809

* Bugfix - Fix path of files shared with me in ocs api: [#204](https://github.com/owncloud/product/issues/204)

   The path of files shared with me using the ocs api was pointing to an incorrect location.

   https://github.com/owncloud/product/issues/204
   https://github.com/owncloud/ocis/pull/994

* Bugfix - Use micro default client: [#718](https://github.com/owncloud/ocis/pull/718)

   Tags: glauth

   We found a file descriptor leak in the glauth connections to the accounts service. Fixed it by
   using the micro default client.

   https://github.com/owncloud/ocis/pull/718

* Bugfix - Allow consent-prompt with switch-account: [#788](https://github.com/owncloud/ocis/pull/788)

   Multiple prompt values are allowed and this change fixes the check for select_account if it was
   used together with other prompt values. Where select_account previously was ignored, it is
   now processed as required, fixing the use case when a RP wants to trigger select_account first
   while at the same time wants also to request interactive consent.

   https://github.com/owncloud/ocis/pull/788

* Bugfix - Mint token with uid and gid: [#737](https://github.com/owncloud/ocis/pull/737)

   Tags: accounts

   The eos driver expects the uid and gid from the opaque map of a user. While the proxy does mint
   tokens correctly, the accounts service wasn't.

   https://github.com/owncloud/ocis/pull/737

* Bugfix - Serve index.html for directories: [#912](https://github.com/owncloud/ocis/pull/912)

   The static middleware in ocis-pkg now serves index.html instead of returning 404 on paths with
   a trailing `/`.

   https://github.com/owncloud/ocis-pkg/issues/63
   https://github.com/owncloud/ocis/pull/912

* Bugfix - Don't create account if id/mail/username already taken: [#709](https://github.com/owncloud/ocis/pull/709)

   Tags: accounts

   We don't allow anymore to create a new account if the provided id/mail/username is already
   taken.

   https://github.com/owncloud/ocis/pull/709

* Bugfix - Fix director selection in proxy: [#521](https://github.com/owncloud/ocis/pull/521)

   Tags: proxy

   We fixed a bug in ocis-proxy where simultaneous requests could be executed on the wrong
   backend.

   https://github.com/owncloud/ocis/pull/521
   https://github.com/owncloud/ocis-proxy/pull/99

* Bugfix - Permission checks for settings write access: [#1092](https://github.com/owncloud/ocis/pull/1092)

   Tags: settings

   There were several endpoints with write access to the settings service that were not protected
   by permission checks. We introduced a generic settings management permission to fix this for
   now. Will be more fine grained later on.

   https://github.com/owncloud/ocis/pull/1092

* Bugfix - Fix minor ui bugs: [#1043](https://github.com/owncloud/ocis/issues/1043)

   - the ui haven't updated the language of the items in the settings view menu. Now we listen to the
   selected language and update the ui - deduplicate resetMenuItems call

   https://github.com/owncloud/ocis/issues/1043
   https://github.com/owncloud/ocis/pull/1044

* Bugfix - Disable public link expiration by default: [#987](https://github.com/owncloud/ocis/issues/987)

   Tags: storage

   The public link expiration was enabled by default and didn't have a default expiration span by
   default, which resulted in already expired public links coming from the public link quick
   action. We fixed this by disabling the public link expiration by default.

   https://github.com/owncloud/ocis/issues/987
   https://github.com/owncloud/ocis/pull/1035

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#416](https://github.com/owncloud/ocis/pull/416)

   Tags: docker

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis/pull/416

* Change - Accounts UI shows message when no permissions: [#656](https://github.com/owncloud/ocis/pull/656)

   We improved the UX of the accounts UI by showing a message information the user about missing
   permissions when the accounts or roles fail to load. This was showing an indeterminate
   progress bar before.

   https://github.com/owncloud/ocis/pull/656

* Change - Cache password validation: [#958](https://github.com/owncloud/ocis/pull/958)

   Tags: accounts

   The password validity check for requests like `login eq '%s' and password eq '%s'` is now cached
   for 10 minutes. This improves the performance for basic auth requests.

   https://github.com/owncloud/ocis/pull/958

* Change - Filesystem based index: [#709](https://github.com/owncloud/ocis/pull/709)

   Tags: accounts, storage

   We replaced `bleve` with a new filesystem based index implementation. There is an `indexer`
   which is capable of orchestrating different index types to build indices on documents by
   field. You can choose from the index types `unique`, `non-unique` or `autoincrement`.
   Indices can be utilized to run search queries (full matches or globbing) on document fields.
   The accounts service is using this index internally to run the search queries coming in via
   `ListAccounts` and `ListGroups` and to generate UIDs for new accounts as well as GIDs for new
   groups.

   The accounts service can be configured to store the index on the local FS / a NFS (`disk`
   implementation of the index) or to use an arbitrary storage ( `cs3` implementation of the
   index). `cs3` is the new default, which is configured to use the `metadata` storage.

   https://github.com/owncloud/ocis/pull/709

* Change - Rebuild index command for accounts: [#748](https://github.com/owncloud/ocis/pull/748)

   Tags: accounts

   The index for the accounts service can now be rebuilt by running the cli command `./bin/ocis
   accounts rebuild`. It deletes all configured indices and rebuilds them from the documents
   found on storage. For this we also introduced a `LoadAccounts` and `LoadGroups` function on
   storage for loading all existing documents.

   https://github.com/owncloud/ocis/pull/748

* Change - Add the thumbnails command: [#156](https://github.com/owncloud/ocis/issues/156)

   Tags: thumbnails

   Added the thumbnails command so that the thumbnails service can get started via ocis.

   https://github.com/owncloud/ocis/issues/156

* Change - CS3 can be used as accounts-backend: [#1020](https://github.com/owncloud/ocis/pull/1020)

   Tags: proxy

   PROXY_ACCOUNT_BACKEND_TYPE=cs3 PROXY_ACCOUNT_BACKEND_TYPE=accounts (default)

   By using a backend which implements the CS3 user-api (currently provided by reva/storage) it
   is possible to bypass the ocis-accounts service and for example use ldap directly.

   https://github.com/owncloud/ocis/pull/1020

* Change - Use bcrypt to hash the user passwords: [#510](https://github.com/owncloud/ocis/issues/510)

   Change the hashing algorithm from SHA-512 to bcrypt since the latter is better suitable for
   password hashing. This is a breaking change. Existing deployments need to regenerate the
   accounts folder.

   https://github.com/owncloud/ocis/issues/510

* Change - Replace the library which scales the images: [#910](https://github.com/owncloud/ocis/pull/910)

   The library went out of support. Also did some refactoring of the thumbnails service code.

   https://github.com/owncloud/ocis/pull/910

* Change - Choose disk or cs3 storage for accounts and groups: [#623](https://github.com/owncloud/ocis/pull/623)

   Tags: accounts

   The accounts service now has an abstraction layer for the storage. In addition to the local disk
   implementation we implemented a cs3 storage, which is the new default for the accounts
   service.

   https://github.com/owncloud/ocis/pull/623

* Change - Enable OpenID dynamic client registration: [#811](https://github.com/owncloud/ocis/issues/811)

   Enable OpenID dynamic client registration

   https://github.com/owncloud/ocis/issues/811
   https://github.com/owncloud/ocis/pull/813

* Change - Integrate import command from ocis-migration: [#249](https://github.com/owncloud/ocis/pull/249)

   Tags: migration

   https://github.com/owncloud/ocis/pull/249
   https://github.com/owncloud/ocis-migration

* Change - Improve reva service descriptions: [#536](https://github.com/owncloud/ocis/pull/536)

   Tags: docs

   The descriptions make it clearer that the services actually represent a mount point in the
   combined storage. Each mount point can have a different driver.

   https://github.com/owncloud/ocis/pull/536

* Change - Initial release of basic version: [#2](https://github.com/owncloud/ocis/issues/2)

   Just prepared an initial basic version which simply embeds the minimum of required services in
   the context of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis/issues/2

* Change - Add cli-commands to manage accounts: [#115](https://github.com/owncloud/product/issues/115)

   Tags: accounts

   COMMANDS:

  * list, ls        List existing accounts
  * add, create     Create a new account
  * update          Make changes to an existing account
  * remove, rm      Removes an existing account
  * inspect         Show detailed data on an existing account
  * help, h         Shows a list of commands or help for one command

   https://github.com/owncloud/product/issues/115

* Change - Start ocis-accounts with the ocis server command: [#25](https://github.com/owncloud/product/issues/25)

   Tags: accounts

   Starts ocis-accounts in single binary mode (./ocis server). This service stores the
   user-account information.

   https://github.com/owncloud/product/issues/25
   https://github.com/owncloud/ocis/pull/239/files

* Change - Properly style konnectd consent page: [#754](https://github.com/owncloud/ocis/pull/754)

   Tags: konnectd

   After bringing our theme into konnectd, we've had to adjust the styles of the consent page so the
   text is visible and button reflects our theme.

   https://github.com/owncloud/ocis/pull/754

* Change - Make all paths configurable and default to a common temp dir: [#1080](https://github.com/owncloud/ocis/pull/1080)

   Aligned all services to use a dir following`/var/tmp/ocis/<service>/...` by default. Also
   made some missing temp paths configurable via env vars and config flags.

   https://github.com/owncloud/ocis/pull/1080

* Change - Move the indexer package from ocis/accounts to ocis/ocis-pkg: [#794](https://github.com/owncloud/ocis/pull/794)

   We are making that change for semantic reasons. So consumers of any index don't necessarily
   need to know of the accounts service.

   https://github.com/owncloud/ocis/pull/794

* Change - Switch over to a new custom-built runtime: [#287](https://github.com/owncloud/ocis/pull/287)

   We moved away from using the go-micro runtime and are now using [our own
   runtime](https://github.com/refs/pman). This allows us to spawn service processes even
   when they are using different versions of go-micro. On top of that we now have the commands `ocis
   list`, `ocis kill` and `ocis run` available for service runtime management.

   https://github.com/owncloud/ocis/pull/287

* Change - Move ocis default config to root level: [#842](https://github.com/owncloud/ocis/pull/842)

   Tags: ocis

   We moved the tracing config to the `root` flagset so that they are parsed on all commands. We also
   introduced a `JWTSecret` flag in the root flagset, in order to apply a common default JWTSecret
   to all services that have one.

   https://github.com/owncloud/ocis/pull/842
   https://github.com/owncloud/ocis/pull/843

* Change - Remove username field in OCS: [#709](https://github.com/owncloud/ocis/pull/709)

   Tags: ocs

   We use the incoming userid as both the `id` and the `on_premises_sam_account_name` for new
   accounts in the accounts service. The userid in OCS requests is in fact the username, not our
   internal account id. We need to enforce the userid as our internal account id though, because
   the account id is part of various `path` formats.

   https://github.com/owncloud/ocis/pull/709
   https://github.com/owncloud/ocis/pull/816

* Change - Account management permissions for Admin role: [#124](https://github.com/owncloud/product/issues/124)

   Tags: accounts, settings

   We created an `AccountManagement` permission and added it to the default admin role. There are
   permission checks in place to protected http endpoints in ocis-accounts against requests
   without the permission. All existing default users (einstein, marie, richard) have the
   default user role now (doesn't have the `AccountManagement` permission). Additionally,
   there is a new default Admin user with credentials `moss:vista`.

   Known issue: for users without the `AccountManagement` permission, the accounts UI
   extension is still available in the ocis-web app switcher, but the requests for loading the
   users will fail (as expected). We are working on a way to hide the accounts UI extension if the
   user doesn't have the `AccountManagement` permission.

   https://github.com/owncloud/product/issues/124
   https://github.com/owncloud/ocis-settings/pull/59
   https://github.com/owncloud/ocis-settings/pull/66
   https://github.com/owncloud/ocis-settings/pull/67
   https://github.com/owncloud/ocis-settings/pull/69
   https://github.com/owncloud/ocis-proxy/pull/95
   https://github.com/owncloud/ocis-pkg/pull/59
   https://github.com/owncloud/ocis-accounts/pull/95
   https://github.com/owncloud/ocis-accounts/pull/100
   https://github.com/owncloud/ocis-accounts/pull/102

* Change - Update phoenix to v0.18.0: [#651](https://github.com/owncloud/ocis/pull/651)

   Tags: web

   We updated phoenix to v0.18.0. Please refer to the changelog (linked) for details on the
   phoenix release. With the ODS release brought in by phoenix we now have proper oc-checkbox and
   oc-radio components for the settings and accounts UI.

   https://github.com/owncloud/ocis/pull/651
   https://github.com/owncloud/phoenix/releases/tag/v0.18.0
   https://github.com/owncloud/owncloud-design-system/releases/tag/v1.12.1

* Change - Default apps in ownCloud Web: [#688](https://github.com/owncloud/ocis/pull/688)

   Tags: web

   We changed the default apps for ownCloud Web to be only files and media-viewer.
   Markdown-editor and draw-io have been removed as defaults.

   https://github.com/owncloud/ocis/pull/688

* Change - Proxy allow insecure upstreams: [#1007](https://github.com/owncloud/ocis/pull/1007)

   Tags: proxy

   We can now configure the proxy if insecure upstream servers are allowed. This was added since
   you need to disable certificate checks fore some situations like testing.

   https://github.com/owncloud/ocis/pull/1007

* Change - Make ocis-settings available: [#287](https://github.com/owncloud/ocis/pull/287)

   Tags: settings

   This version delivers `settings` as a new service. It is part of the array of services in the
   `server` command.

   https://github.com/owncloud/ocis/pull/287

* Change - Start ocis-proxy with the ocis server command: [#119](https://github.com/owncloud/ocis/issues/119)

   Tags: proxy

   Starts the proxy in single binary mode (./ocis server) on port 9200. The proxy serves as a
   single-entry point for all http-clients.

   https://github.com/owncloud/ocis/issues/119
   https://github.com/owncloud/ocis/issues/136

* Change - Theme welcome and choose account pages: [#887](https://github.com/owncloud/ocis/pull/887)

   Tags: konnectd

   We've themed the konnectd pages Welcome and Choose account. All text has a white color now to be
   easily readable on the dark background.

   https://github.com/owncloud/ocis/pull/887

* Change - Bring oC theme: [#698](https://github.com/owncloud/ocis/pull/698)

   Tags: konnectd

   We've styled our konnectd login page to reflect ownCloud theme.

   https://github.com/owncloud/ocis/pull/698

* Change - Unify Configuration Parsing: [#675](https://github.com/owncloud/ocis/pull/675)

   Tags: ocis

   - responsibility for config parsing should be on the subcommand - if there is a config file in the
   environment location, env var should take precedence - general rule of thumb: the more
   explicit the config file is that would be picked up. Order from less to more explicit: - config
   location (/etc/ocis) - environment variable - cli flag

   https://github.com/owncloud/ocis/pull/675

* Change - Update phoenix to v0.20.0: [#674](https://github.com/owncloud/ocis/pull/674)

   Tags: web

   We updated phoenix to v0.20.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/674
   https://github.com/owncloud/phoenix/releases/tag/v0.20.0

* Change - Update phoenix to v0.21.0: [#728](https://github.com/owncloud/ocis/pull/728)

   Tags: web

   We updated phoenix to v0.21.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/728
   https://github.com/owncloud/phoenix/releases/tag/v0.21.0

* Change - Update phoenix to v0.22.0: [#757](https://github.com/owncloud/ocis/pull/757)

   Tags: web

   We updated phoenix to v0.22.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/757
   https://github.com/owncloud/phoenix/releases/tag/v0.22.0

* Change - Update phoenix to v0.23.0: [#785](https://github.com/owncloud/ocis/pull/785)

   Tags: web

   We updated phoenix to v0.23.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/785
   https://github.com/owncloud/phoenix/releases/tag/v0.23.0

* Change - Update phoenix to v0.24.0: [#817](https://github.com/owncloud/ocis/pull/817)

   Tags: web

   We updated phoenix to v0.24.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/817
   https://github.com/owncloud/phoenix/releases/tag/v0.24.0

* Change - Update phoenix to v0.25.0: [#868](https://github.com/owncloud/ocis/pull/868)

   Tags: web

   We updated phoenix to v0.25.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/868
   https://github.com/owncloud/phoenix/releases/tag/v0.25.0

* Change - Update phoenix to v0.26.0: [#935](https://github.com/owncloud/ocis/pull/935)

   Tags: web

   We updated phoenix to v0.26.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/935
   https://github.com/owncloud/phoenix/releases/tag/v0.26.0

* Change - Update phoenix to v0.27.0: [#943](https://github.com/owncloud/ocis/pull/943)

   Tags: web

   We updated phoenix to v0.27.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/943
   https://github.com/owncloud/phoenix/releases/tag/v0.27.0

* Change - Update phoenix to v0.28.0: [#1027](https://github.com/owncloud/ocis/pull/1027)

   Tags: web

   We updated phoenix to v0.28.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/1027
   https://github.com/owncloud/phoenix/releases/tag/v0.28.0

* Change - Update phoenix to v0.29.0: [#1034](https://github.com/owncloud/ocis/pull/1034)

   Tags: web

   We updated phoenix to v0.29.0. Please refer to the changelog (linked) for details on the
   phoenix release.

   https://github.com/owncloud/ocis/pull/1034
   https://github.com/owncloud/phoenix/releases/tag/v0.29.0

* Change - Update reva config: [#336](https://github.com/owncloud/ocis/pull/336)

  * EOS homes are not configured with an enable-flag anymore, but with a dedicated storage driver.
  * We're using it now and adapted default configs of storages

   https://github.com/owncloud/ocis/pull/336
   https://github.com/owncloud/ocis/pull/337
   https://github.com/owncloud/ocis/pull/338
   https://github.com/owncloud/ocis-reva/pull/891

* Change - Update reva to v1.4.1-0.20201209113234-e791b5599a89: [#1089](https://github.com/owncloud/ocis/pull/1089)

   Updated reva to v1.4.1-0.20201209113234-e791b5599a89

   https://github.com/owncloud/ocis/pull/1089

* Change - Clarify storage driver env vars: [#729](https://github.com/owncloud/ocis/pull/729)

   After renaming ocsi-reva to storage and combining the storage and data providers some env vars
   were confusingly named `STORAGE_STORAGE_...`. We are changing the prefix for driver related
   env vars to `STORAGE_DRIVER_...`. This makes changing the storage driver using eg.:
   `STORAGE_HOME_DRIVER=eos` and setting driver options using
   `STORAGE_DRIVER_EOS_LAYOUT=...` less confusing.

   https://github.com/owncloud/ocis/pull/729

* Change - Update ownCloud Web to v1.0.0-beta3: [#1105](https://github.com/owncloud/ocis/pull/1105)

   Tags: web

   We updated ownCloud Web to v1.0.0-beta3. Please refer to the changelog (linked) for details on
   the web release.

   https://github.com/owncloud/ocis/pull/1105
   https://github.com/owncloud/phoenix/releases/tag/v1.0.0-beta3

* Change - Update ownCloud Web to v1.0.0-beta4: [#1110](https://github.com/owncloud/ocis/pull/1110)

   Tags: web

   We updated ownCloud Web to v1.0.0-beta4. Please refer to the changelog (linked) for details on
   the web release.

   https://github.com/owncloud/ocis/pull/1110
   https://github.com/owncloud/phoenix/releases/tag/v1.0.0-beta4

* Change - Settings and accounts appear in the user menu: [#656](https://github.com/owncloud/ocis/pull/656)

   We moved settings and accounts to the user menu.

   https://github.com/owncloud/ocis/pull/656

* Enhancement - Add tracing to the accounts service: [#1016](https://github.com/owncloud/ocis/issues/1016)

   Added tracing to the accounts service.

   https://github.com/owncloud/ocis/issues/1016

* Enhancement - Add the accounts service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: accounts

  * Bugfix - Initialize roleService client in GRPC server: [#114](https://github.com/owncloud/ocis-accounts/pull/114)
  * Bugfix - Cleanup separated indices in memory: [#224](https://github.com/owncloud/product/issues/224)
  * Change - Set user role on builtin users: [#102](https://github.com/owncloud/ocis-accounts/pull/102)
  * Change - Add new builtin admin user: [#102](https://github.com/owncloud/ocis-accounts/pull/102)
  * Change - We make use of the roles cache to enforce permission checks: [#100](https://github.com/owncloud/ocis-accounts/pull/100)
  * Change - We make use of the roles manager to enforce permission checks: [#108](https://github.com/owncloud/ocis-accounts/pull/108)
  * Enhancement - Add create account form: [#148](https://github.com/owncloud/product/issues/148)
  * Enhancement - Add delete accounts action: [#148](https://github.com/owncloud/product/issues/148)
  * Enhancement - Add enable/disable capabilities to the WebUI: [#118](https://github.com/owncloud/product/issues/118)
  * Enhancement - Improve visual appearance of accounts UI: [#222](https://github.com/owncloud/product/issues/222)
  * Bugfix - Adapting to new settings API for fetching roles: [#96](https://github.com/owncloud/ocis-accounts/pull/96)
  * Change - Create account api-call implicitly adds "default-user" role: [#173](https://github.com/owncloud/product/issues/173)
  * Change - Add role selection to accounts UI: [#103](https://github.com/owncloud/product/issues/103)
  * Bugfix - Atomic Requests: [#82](https://github.com/owncloud/ocis-accounts/pull/82)
  * Bugfix - Unescape value for prefix query: [#76](https://github.com/owncloud/ocis-accounts/pull/76)
  * Change - Adapt to new ocis-settings data model: [#87](https://github.com/owncloud/ocis-accounts/pull/87)
  * Change - Add permissions for language to default roles: [#88](https://github.com/owncloud/ocis-accounts/pull/88)
  * Bugfix - Add write mutexes: [#71](https://github.com/owncloud/ocis-accounts/pull/71)
  * Bugfix - Fix the accountId and groupId mismatch in DeleteGroup Method: [#60](https://github.com/owncloud/ocis-accounts/pull/60)
  * Bugfix - Fix index mapping: [#73](https://github.com/owncloud/ocis-accounts/issues/73)
  * Bugfix - Use NewNumericRangeInclusiveQuery for numeric literals: [#28](https://github.com/owncloud/ocis-glauth/issues/28)
  * Bugfix - Prevent segfault when no password is set: [#65](https://github.com/owncloud/ocis-accounts/pull/65)
  * Bugfix - Update account return value not used: [#70](https://github.com/owncloud/ocis-accounts/pull/70)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#64](https://github.com/owncloud/ocis-accounts/pull/64)
  * Change - Align structure of this extension with other extensions: [#51](https://github.com/owncloud/ocis-accounts/pull/51)
  * Change - Change api errors: [#11](https://github.com/owncloud/ocis-accounts/issues/11)
  * Change - Enable accounts on creation: [#43](https://github.com/owncloud/ocis-accounts/issues/43)
  * Change - Fix index update on create/update: [#57](https://github.com/owncloud/ocis-accounts/issues/57)
  * Change - Pass around the correct logger throughout the code: [#41](https://github.com/owncloud/ocis-accounts/issues/41)
  * Change - Remove timezone setting: [#33](https://github.com/owncloud/ocis-accounts/pull/33)
  * Change - Tighten screws on usernames and email addresses: [#65](https://github.com/owncloud/ocis-accounts/pull/65)
  * Enhancement - Add early version of cli tools for user-management: [#69](https://github.com/owncloud/ocis-accounts/pull/69)
  * Enhancement - Update accounts API: [#30](https://github.com/owncloud/ocis-accounts/pull/30)
  * Enhancement - Add simple user listing UI: [#51](https://github.com/owncloud/ocis-accounts/pull/51)
  * Enhancement - Logging is configurable: [#24](https://github.com/owncloud/ocis-accounts/pull/24)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-accounts/issues/1)
  * Enhancement - Configuration: [#15](https://github.com/owncloud/ocis-accounts/pull/15)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add basic auth option: [#627](https://github.com/owncloud/ocis/pull/627)

   We added a new `enable-basic-auth` option and `PROXY_ENABLE_BASIC_AUTH` environment
   variable that can be set to `true` to make the proxy verify the basic auth header with the
   accounts service. This should only be used for testing and development and is disabled by
   default.

   https://github.com/owncloud/product/issues/198
   https://github.com/owncloud/ocis/pull/627

* Enhancement - Document how to run OCIS on top of EOS: [#172](https://github.com/owncloud/ocis/pull/172)

   Tags: eos

   We have added rules to the Makefile that use the official [eos docker
   images](https://gitlab.cern.ch/eos/eos-docker) to boot an eos cluster and configure OCIS
   to use it.

   https://github.com/owncloud/ocis/pull/172

* Enhancement - Add the glauth service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: glauth

  * Bugfix - Return invalid credentials when user was not found: [#30](https://github.com/owncloud/ocis-glauth/pull/30)
  * Bugfix - Query numeric attribute values without quotes: [#28](https://github.com/owncloud/ocis-glauth/issues/28)
  * Bugfix - Use searchBaseDN if already a user/group name: [#214](https://github.com/owncloud/product/issues/214)
  * Bugfix - Fix LDAP substring startswith filters: [#31](https://github.com/owncloud/ocis-glauth/pull/31)
  * Enhancement - Add build information to the metrics: [#226](https://github.com/owncloud/product/issues/226)
  * Enhancement - Reenable configuring backends: [#600](https://github.com/owncloud/ocis/pull/600)
  * Bugfix - Ignore case when comparing objectclass values: [#26](https://github.com/owncloud/ocis-glauth/pull/26)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#24](https://github.com/owncloud/ocis-glauth/pull/24)
  * Enhancement - Handle ownCloudUUID attribute: [#27](https://github.com/owncloud/ocis-glauth/pull/27)
  * Enhancement - Implement group queries: [#22](https://github.com/owncloud/ocis-glauth/issues/22)
  * Enhancement - Configuration: [#11](https://github.com/owncloud/ocis-glauth/pull/11)
  * Enhancement - Improve default settings: [#12](https://github.com/owncloud/ocis-glauth/pull/12)
  * Enhancement - Generate temporary ldap certificates if LDAPS is enabled: [#12](https://github.com/owncloud/ocis-glauth/pull/12)
  * Enhancement - Provide additional tls-endpoint: [#12](https://github.com/owncloud/ocis-glauth/pull/12)
  * Change - Use physicist demo users: [#5](https://github.com/owncloud/ocis-glauth/issues/5)
  * Change - Default to config based user backend: [#6](https://github.com/owncloud/ocis-glauth/pull/6)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add k6: [#941](https://github.com/owncloud/ocis/pull/941)

   Tags: tests

   Add k6 as a performance testing framework

   https://github.com/owncloud/ocis/pull/941
   https://github.com/owncloud/ocis/pull/983

* Enhancement - Add the konnectd service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: konnectd

  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Add silent redirect url: [#69](https://github.com/owncloud/ocis-konnectd/issues/69)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#71](https://github.com/owncloud/ocis-konnectd/pull/71)
  * Bugfix - Include the assets for #62: [#64](https://github.com/owncloud/ocis-konnectd/pull/64)
  * Bugfix - Redirect to the provided uri: [#26](https://github.com/owncloud/ocis-konnectd/issues/26)
  * Change - Add a trailing slash to trusted redirect uris: [#26](https://github.com/owncloud/ocis-konnectd/issues/26)
  * Change - Improve client identifiers for end users: [#62](https://github.com/owncloud/ocis-konnectd/pull/62)
  * Enhancement - Use upstream version of konnect library: [#14](https://github.com/owncloud/product/issues/14)
  * Enhancement - Change default config for single-binary: [#55](https://github.com/owncloud/ocis-konnectd/pull/55)
  * Bugfix - Generate a random CSP-Nonce in the webapp: [#17](https://github.com/owncloud/ocis-konnectd/issues/17)
  * Change - Dummy index.html is not required anymore by upstream: [#25](https://github.com/owncloud/ocis-konnectd/issues/25)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-konnectd/issues/1)
  * Change - Use glauth as ldap backend, default to running behind ocis-proxy: [#52](https://github.com/owncloud/ocis-konnectd/pull/52)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the ocis-phoenix service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: web

  * Bugfix - Fix external app URLs: [#218](https://github.com/owncloud/product/issues/218)
  * Change - Remove pdf-viewer from default apps: [#85](https://github.com/owncloud/ocis-phoenix/pull/85)
  * Change - Enable Settings and Accounts apps by default: [#80](https://github.com/owncloud/ocis-phoenix/pull/80)
  * Bugfix - Exit when assets or config are not found: [#76](https://github.com/owncloud/ocis-phoenix/pull/76)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#73](https://github.com/owncloud/ocis-phoenix/pull/73)
  * Change - Hide searchbar by default: [#116](https://github.com/owncloud/product/issues/116)
  * Bugfix - Allow silent refresh of access token: [#69](https://github.com/owncloud/ocis-konnectd/issues/69)
  * Change - Update Phoenix: [#60](https://github.com/owncloud/ocis-phoenix/pull/60)
  * Enhancement - Configuration: [#57](https://github.com/owncloud/ocis-phoenix/pull/57)
  * Bugfix - Config file value not being read: [#45](https://github.com/owncloud/ocis-phoenix/pull/45)
  * Change - Default to running behind ocis-proxy: [#55](https://github.com/owncloud/ocis-phoenix/pull/55)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the ocis-pkg package: [#244](https://github.com/owncloud/product/issues/244)

   Tags: ocis-pkg

  * Change - Unwrap roleIDs from access-token into metadata context: [#59](https://github.com/owncloud/ocis-pkg/pull/59)
  * Change - Provide cache for roles: [#59](https://github.com/owncloud/ocis-pkg/pull/59)
  * Change - Roles manager: [#60](https://github.com/owncloud/ocis-pkg/pull/60)
  * Change - Use go-micro's metadata context for account id: [#56](https://github.com/owncloud/ocis-pkg/pull/56)
  * Bugfix - Remove redigo 2.0.0+incompatible dependency: [#33](https://github.com/owncloud/ocis-graph/pull/33)
  * Change - Add middleware for x-access-token distmantling: [#46](https://github.com/owncloud/ocis-pkg/pull/46)
  * Enhancement - Add `ocis.id` and numeric id claims: [#50](https://github.com/owncloud/ocis-pkg/pull/50)
  * Bugfix - Pass flags to micro service: [#44](https://github.com/owncloud/ocis-pkg/pull/44)
  * Change - Add header to cors handler: [#41](https://github.com/owncloud/ocis-pkg/issues/41)
  * Enhancement - Tracing middleware: [#35](https://github.com/owncloud/ocis-pkg/pull/35/)
  * Enhancement - Allow http services to register handlers: [#33](https://github.com/owncloud/ocis-pkg/pull/33)
  * Change - Upgrade the micro libraries: [#22](https://github.com/owncloud/ocis-pkg/pull/22)
  * Bugfix - Fix Module Path: [#25](https://github.com/owncloud/ocis-pkg/pull/25)
  * Bugfix - Change import paths to ocis-pkg/v2: [#27](https://github.com/owncloud/ocis-pkg/pull/27)
  * Bugfix - Fix serving static assets: [#14](https://github.com/owncloud/ocis-pkg/pull/14)
  * Change - Add TLS support for http services: [#19](https://github.com/owncloud/ocis-pkg/issues/19)
  * Enhancement - Introduce OpenID Connect middleware: [#8](https://github.com/owncloud/ocis-pkg/issues/8)
  * Change - Add root path to static middleware: [#9](https://github.com/owncloud/ocis-pkg/issues/9)
  * Change - Better log level handling within micro: [#2](https://github.com/owncloud/ocis-pkg/issues/2)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the ocs service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: ocs

  * Bugfix - Match the user response to the OC10 format: [#181](https://github.com/owncloud/product/issues/181)
  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Add the top level response structure to json responses: [#181](https://github.com/owncloud/product/issues/181)
  * Enhancement - Update ocis-accounts: [#42](https://github.com/owncloud/ocis-ocs/pull/42)
  * Bugfix - Mimic oc10 user enabled as string in provisioning api: [#39](https://github.com/owncloud/ocis-ocs/pull/39)
  * Bugfix - Use opaque ID of a user for signing keys: [#436](https://github.com/owncloud/ocis/issues/436)
  * Enhancement - Add option to create user with uidnumber and gidnumber: [#34](https://github.com/owncloud/ocis-ocs/pull/34)
  * Bugfix - Fix file descriptor leak: [#79](https://github.com/owncloud/ocis-accounts/issues/79)
  * Enhancement - Add Group management for OCS Povisioning API: [#25](https://github.com/owncloud/ocis-ocs/pull/25)
  * Enhancement - Basic Support for the User Provisioning API: [#23](https://github.com/owncloud/ocis-ocs/pull/23)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#20](https://github.com/owncloud/ocis-ocs/pull/20)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-ocs/issues/1)
  * Change - Upgrade micro libraries: [#11](https://github.com/owncloud/ocis-ocs/issues/11)
  * Enhancement - Configuration: [#14](https://github.com/owncloud/ocis-ocs/pull/14)
  * Enhancement - Support signing key: [#18](https://github.com/owncloud/ocis-ocs/pull/18)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the proxy service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: proxy

  * Bugfix - Fix director selection: [#99](https://github.com/owncloud/ocis-proxy/pull/99)
  * Bugfix - Add settings API and app endpoints to example config: [#93](https://github.com/owncloud/ocis-proxy/pull/93)
  * Change - Remove accounts caching: [#100](https://github.com/owncloud/ocis-proxy/pull/100)
  * Enhancement - Add autoprovision accounts flag: [#219](https://github.com/owncloud/product/issues/219)
  * Enhancement - Add hello API and app endpoints to example config and builtin config: [#96](https://github.com/owncloud/ocis-proxy/pull/96)
  * Enhancement - Add roleIDs to the access token: [#95](https://github.com/owncloud/ocis-proxy/pull/95)
  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Enhancement - Add numeric uid and gid to the access token: [#89](https://github.com/owncloud/ocis-proxy/pull/89)
  * Enhancement - Add configuration options for the pre-signed url middleware: [#91](https://github.com/owncloud/ocis-proxy/issues/91)
  * Bugfix - Enable new accounts by default: [#79](https://github.com/owncloud/ocis-proxy/pull/79)
  * Bugfix - Lookup user by id for presigned URLs: [#85](https://github.com/owncloud/ocis-proxy/pull/85)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#78](https://github.com/owncloud/ocis-proxy/pull/78)
  * Change - Add settings and ocs group routes: [#81](https://github.com/owncloud/ocis-proxy/pull/81)
  * Change - Add route for user provisioning API in ocis-ocs: [#80](https://github.com/owncloud/ocis-proxy/pull/80)
  * Bugfix - Provide token configuration from config: [#69](https://github.com/owncloud/ocis-proxy/pull/69)
  * Bugfix - Provide token configuration from config: [#76](https://github.com/owncloud/ocis-proxy/pull/76)
  * Change - Add OIDC config flags: [#66](https://github.com/owncloud/ocis-proxy/pull/66)
  * Change - Mint new username property in the reva token: [#62](https://github.com/owncloud/ocis-proxy/pull/62)
  * Enhancement - Add Accounts UI routes: [#65](https://github.com/owncloud/ocis-proxy/pull/65)
  * Enhancement - Add option to disable TLS: [#71](https://github.com/owncloud/ocis-proxy/issues/71)
  * Enhancement - Only send create home request if an account has been migrated: [#52](https://github.com/owncloud/ocis-proxy/issues/52)
  * Enhancement - Create a root span on proxy that propagates down to consumers: [#64](https://github.com/owncloud/ocis-proxy/pull/64)
  * Enhancement - Support signed URLs: [#73](https://github.com/owncloud/ocis-proxy/issues/73)
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
  * Change - Update ocis-pkg: [#30](https://github.com/owncloud/ocis-proxy/pull/30)
  * Change - Insecure http-requests are now redirected to https: [#29](https://github.com/owncloud/ocis-proxy/pull/29)
  * Enhancement - Configurable OpenID Connect client: [#27](https://github.com/owncloud/ocis-proxy/pull/27)
  * Enhancement - Add policy selectors: [#4](https://github.com/owncloud/ocis-proxy/issues/4)
  * Bugfix - Set TLS-Certificate correctly: [#25](https://github.com/owncloud/ocis-proxy/pull/25)
  * Change - Route requests based on regex or query parameters: [#21](https://github.com/owncloud/ocis-proxy/issues/21)
  * Enhancement - Proxy client urls in default configuration: [#19](https://github.com/owncloud/ocis-proxy/issues/19)
  * Enhancement - Make TLS-Cert configurable: [#14](https://github.com/owncloud/ocis-proxy/pull/14)
  * Enhancement - Load Proxy Policies at Runtime: [#17](https://github.com/owncloud/ocis-proxy/issues/17)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the settings service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: settings

  * Bugfix - Fix loading and saving system scoped values: [#66](https://github.com/owncloud/ocis-settings/pull/66)
  * Bugfix - Complete input validation: [#66](https://github.com/owncloud/ocis-settings/pull/66)
  * Change - Add filter option for bundle ids in ListBundles and ListRoles: [#59](https://github.com/owncloud/ocis-settings/pull/59)
  * Change - Reuse roleIDs from the metadata context: [#69](https://github.com/owncloud/ocis-settings/pull/69)
  * Change - Update ocis-pkg/v2: [#72](https://github.com/owncloud/ocis-settings/pull/72)
  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Fix fetching bundles in settings UI: [#61](https://github.com/owncloud/ocis-settings/pull/61)
  * Change - Filter settings by permissions: [#99](https://github.com/owncloud/product/issues/99)
  * Change - Add role service: [#110](https://github.com/owncloud/product/issues/110)
  * Change - Rename endpoints and message types: [#36](https://github.com/owncloud/ocis-settings/issues/36)
  * Change - Use UUIDs instead of alphanumeric identifiers: [#46](https://github.com/owncloud/ocis-settings/pull/46)
  * Bugfix - Adjust UUID validation to be more tolerant: [#41](https://github.com/owncloud/ocis-settings/issues/41)
  * Bugfix - Fix runtime error when type asserting on nil value: [#38](https://github.com/owncloud/ocis-settings/pull/38)
  * Bugfix - Fix multiple submits on string and number form elements: [#745](https://github.com/owncloud/owncloud-design-system/issues/745)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#39](https://github.com/owncloud/ocis-settings/pull/39)
  * Change - Dynamically add navItems for extensions with settings bundles: [#25](https://github.com/owncloud/ocis-settings/pull/25)
  * Change - Introduce input validation: [#22](https://github.com/owncloud/ocis-settings/pull/22)
  * Change - Use account uuid from x-access-token: [#14](https://github.com/owncloud/ocis-settings/pull/14)
  * Change - Use server config variable from ocis-web: [#34](https://github.com/owncloud/ocis-settings/pull/34)
  * Enhancement - Remove paths from Makefile: [#33](https://github.com/owncloud/ocis-settings/pull/33)
  * Enhancement - Extend the docs: [#11](https://github.com/owncloud/ocis-settings/issues/11)
  * Enhancement - Update ocis-pkg/v2: [#42](https://github.com/owncloud/ocis-settings/pull/42)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the storage service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: storage, reva

  * Enhancement - Enable ocis driver treetime accounting: [#620](https://github.com/owncloud/ocis/pull/620)
  * Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)

   In the future accounts, settings etc. should be stored in a dedicated metadata storage. The
   services should talk to this storage directly, bypassing reva-gateway.

   Https://github.com/owncloud/ocis/pull/602

  * Enhancement - Update reva to v1.2.2-0.20200924071957-e6676516e61e: [#601](https://github.com/owncloud/ocis/pull/601)

   - Update reva to v1.2.2-0.20200924071957-e6676516e61e - eos client: Handle eos EPERM as
   permission denied [(reva/#1183)](https://github.com/cs3org/reva/pull/1183) - ocis
   driver: synctime based etag propagation
   [(reva/#1180)](https://github.com/cs3org/reva/pull/1180) - ocis driver: fix litmus
   [(reva/#1179)](https://github.com/cs3org/reva/pull/1179) - ocis driver: fix move
   [(reva/#1177)](https://github.com/cs3org/reva/pull/1177) - ocs service: cache
   displaynames [(reva/#1161)](https://github.com/cs3org/reva/pull/1161)

   Https://github.com/owncloud/ocis-reva/issues/262
   https://github.com/owncloud/ocis-reva/issues/357
   https://github.com/owncloud/ocis-reva/issues/301
   https://github.com/owncloud/ocis-reva/issues/302
   https://github.com/owncloud/ocis/pull/601

  * Bugfix - Fix default configuration for accessing shares: [#205](https://github.com/owncloud/product/issues/205)

   The storage provider mounted at `/home` should always have EnableHome set to `true`. The other
   storage providers should have it set to `false`.

   Https://github.com/owncloud/product/issues/205
   https://github.com/owncloud/ocis-reva/pull/461

  * Enhancement - Allow configuring arbitrary storage registry rules: [#193](https://github.com/owncloud/product/issues/193)

   We added a new config flag `storage-registry-rule` that can be given multiple times for the
   gateway to specify arbitrary storage registry rules. You can also use a comma separated list of
   rules in the `REVA_STORAGE_REGISTRY_RULES` environment variable.

   Https://github.com/owncloud/product/issues/193
   https://github.com/owncloud/ocis-reva/pull/461

  * Enhancement - Update reva to v1.2.1-0.20200826162318-c0f54e1f37ea: [#454](https://github.com/owncloud/ocis-reva/pull/454)

   - Update reva to v1.2.1-0.20200826162318-c0f54e1f37ea - Do not swallow 'not found' errors in
   Stat [(reva/#1124)](https://github.com/cs3org/reva/pull/1124) - Rewire dav files to the
   home storage [(reva/#1125)](https://github.com/cs3org/reva/pull/1125) - Do not restore
   recycle entry on purge [(reva/#1099)](https://github.com/cs3org/reva/pull/1099) -
   Allow listing the trashbin [(reva/#1091)](https://github.com/cs3org/reva/pull/1091) -
   Restore and delete trash items via ocs
   [(reva/#1103)](https://github.com/cs3org/reva/pull/1103) - Ensure ignoring public
   stray shares [(reva/#1090)](https://github.com/cs3org/reva/pull/1090) - Ensure
   ignoring stray shares [(reva/#1064)](https://github.com/cs3org/reva/pull/1064) -
   Minor fixes in reva cmd, gateway uploads and smtpclient
   [(reva/#1082)](https://github.com/cs3org/reva/pull/1082) - Owncloud driver -
   propagate mtime on RemoveGrant
   [(reva/#1115)](https://github.com/cs3org/reva/pull/1115) - Handle redirection
   prefixes when extracting destination from URL
   [(reva/#1111)](https://github.com/cs3org/reva/pull/1111) - Add UID and GID in ldap auth
   driver [(reva/#1101)](https://github.com/cs3org/reva/pull/1101) - Add calens check to
   verify changelog entries in CI
   [(reva/#1077)](https://github.com/cs3org/reva/pull/1077) - Refactor Reva CLI with
   prompts [(reva/#1072)](https://github.com/cs3org/reva/pull/1072j) - Get file info
   using fxids from EOS [(reva/#1079)](https://github.com/cs3org/reva/pull/1079) - Update
   LDAP user driver [(reva/#1088)](https://github.com/cs3org/reva/pull/1088) - System
   information metrics cleanup
   [(reva/#1114)](https://github.com/cs3org/reva/pull/1114) - System information
   included in Prometheus metrics
   [(reva/#1071)](https://github.com/cs3org/reva/pull/1071) - Add logic for resolving
   storage references over webdav
   [(reva/#1094)](https://github.com/cs3org/reva/pull/1094)

   Https://github.com/owncloud/ocis-reva/pull/454

  * Enhancement - Update reva to v1.2.1-0.20200911111727-51649e37df2d: [#466](https://github.com/owncloud/ocis-reva/pull/466)

   - Update reva to v1.2.1-0.20200911111727-51649e37df2d - Added new OCIS storage driver ocis
   [(reva/#1155)](https://github.com/cs3org/reva/pull/1155) - App provider: fallback to
   env. variable if 'iopsecret' unset
   [(reva/#1146)](https://github.com/cs3org/reva/pull/1146) - Add switch to database
   [(reva/#1135)](https://github.com/cs3org/reva/pull/1135) - Add the ocdav HTTP svc to the
   standalone config [(reva/#1128)](https://github.com/cs3org/reva/pull/1128)

   Https://github.com/owncloud/ocis-reva/pull/466

  * Enhancement - Separate user and auth providers, add config for rest user: [#412](https://github.com/owncloud/ocis-reva/pull/412)

   Previously, the auth and user provider services used to have the same driver, which restricted
   using separate drivers and configs for both. This PR separates the two and adds the config for
   the rest user driver and the gatewaysvc parameter to EOS fs.

   Https://github.com/owncloud/ocis-reva/pull/412
   https://github.com/cs3org/reva/pull/995

  * Enhancement - Update reva to v1.1.1-0.20200819100654-dcbf0c8ea187: [#447](https://github.com/owncloud/ocis-reva/pull/447)

   - Update reva to v1.1.1-0.20200819100654-dcbf0c8ea187 - fix restoring and deleting trash
   items via ocs [(reva/#1103)](https://github.com/cs3org/reva/pull/1103) - Add UID and GID
   in ldap auth driver [(reva/#1101)](https://github.com/cs3org/reva/pull/1101) - Allow
   listing the trashbin [(reva/#1091)](https://github.com/cs3org/reva/pull/1091) -
   Ignore Stray Public Shares [(reva/#1090)](https://github.com/cs3org/reva/pull/1090) -
   Implement GetUserByClaim for LDAP user driver
   [(reva/#1088)](https://github.com/cs3org/reva/pull/1088) - eosclient: get file info by
   fxid [(reva/#1079)](https://github.com/cs3org/reva/pull/1079) - Ensure stray shares
   get ignored [(reva/#1064)](https://github.com/cs3org/reva/pull/1064) - Improve
   timestamp precision while logging
   [(reva/#1059)](https://github.com/cs3org/reva/pull/1059) - Ocfs lookup userid
   (update) [(reva/#1052)](https://github.com/cs3org/reva/pull/1052) - Disallow sharing
   the shares directory [(reva/#1051)](https://github.com/cs3org/reva/pull/1051) - Local
   storage provider: Fixed resolution of fileid
   [(reva/#1046)](https://github.com/cs3org/reva/pull/1046) - List public shares only
   created by the current user [(reva/#1042)](https://github.com/cs3org/reva/pull/1042)

   Https://github.com/owncloud/ocis-reva/pull/447

  * Bugfix - Update LDAP filters: [#399](https://github.com/owncloud/ocis-reva/pull/399)

   With the separation of use and find filters we can now use a filter that taken into account a users
   uuid as well as his username. This is necessary to make sharing work with the new account service
   which assigns accounts an immutable account id that is different from the username.
   Furthermore, the separate find filters now allows searching users by their displayname or
   email as well.

   ``` userfilter =
   "(&(objectclass=posixAccount)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))"
   findfilter =
   "(&(objectclass=posixAccount)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))"
   ```

   Https://github.com/owncloud/ocis-reva/pull/399
   https://github.com/cs3org/reva/pull/996

  * Change - Environment updates for the username userid split: [#420](https://github.com/owncloud/ocis-reva/pull/420)

   We updated the owncloud storage driver in reva to properly look up users by userid or username
   using the userprovider instead of taking the path segment as is. This requires the user service
   address as well as changing the default layout to the userid instead of the username. The latter
   is not considered a stable and persistent identifier.

   Https://github.com/owncloud/ocis-reva/pull/420
   https://github.com/cs3org/reva/pull/1033

  * Enhancement - Update storage documentation: [#384](https://github.com/owncloud/ocis-reva/pull/384)

   We added details to the documentation about storage requirements known from ownCloud 10, the
   local storage driver and the ownCloud storage driver.

   Https://github.com/owncloud/ocis-reva/pull/384
   https://github.com/owncloud/ocis-reva/pull/390

  * Enhancement - Update reva to v0.1.1-0.20200724135750-b46288b375d6: [#399](https://github.com/owncloud/ocis-reva/pull/399)

   - Update reva to v0.1.1-0.20200724135750-b46288b375d6 - Split LDAP user filters
   (reva/#996) - meshdirectory: Add invite forward API to provider links (reva/#1000) - OCM:
   Pass the link to the meshdirectory service in token mail (reva/#1002) - Update
   github.com/go-ldap/ldap to v3 (reva/#1004)

   Https://github.com/owncloud/ocis-reva/pull/399
   https://github.com/cs3org/reva/pull/996 https://github.com/cs3org/reva/pull/1000
   https://github.com/cs3org/reva/pull/1002 https://github.com/cs3org/reva/pull/1004

  * Enhancement - Update reva to v0.1.1-0.20200728071211-c948977dd3a0: [#407](https://github.com/owncloud/ocis-reva/pull/407)

   - Update reva to v0.1.1-0.20200728071211-c948977dd3a0 - Use proper logging for ldap auth
   requests (reva/#1008) - Update github.com/eventials/go-tus to
   v0.0.0-20200718001131-45c7ec8f5d59 (reva/#1007) - Check if SMTP credentials are nil
   (reva/#1006)

   Https://github.com/owncloud/ocis-reva/pull/407
   https://github.com/cs3org/reva/pull/1008 https://github.com/cs3org/reva/pull/1007
   https://github.com/cs3org/reva/pull/1006

  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#393](https://github.com/owncloud/ocis-reva/pull/393)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   Https://github.com/owncloud/ocis-reva/pull/393

  * Enhancement - Update reva to v0.1.1-0.20200710143425-cf38a45220c5: [#371](https://github.com/owncloud/ocis-reva/pull/371)

   - Update reva to v0.1.1-0.20200710143425-cf38a45220c5 (#371) - Add wopi open (reva/#920) -
   Added a CS3API compliant data exporter to Mentix (reva/#955) - Read SMTP password from env if
   not set in config (reva/#953) - OCS share fix including file info after update (reva/#958) - Add
   flag to smtpclient for for unauthenticated SMTP (reva/#963)

   Https://github.com/owncloud/ocis-reva/pull/371
   https://github.com/cs3org/reva/pull/920 https://github.com/cs3org/reva/pull/953
   https://github.com/cs3org/reva/pull/955 https://github.com/cs3org/reva/pull/958
   https://github.com/cs3org/reva/pull/963

  * Enhancement - Update reva to v0.1.1-0.20200722125752-6dea7936f9d1: [#392](https://github.com/owncloud/ocis-reva/pull/392)

   - Update reva to v0.1.1-0.20200722125752-6dea7936f9d1 - Added signing key capability
   (reva/#986) - Add functionality to create webdav references for OCM shares (reva/#974) -
   Added a site locations exporter to Mentix (reva/#972) - Add option to config to allow requests
   to hosts with unverified certificates (reva/#969)

   Https://github.com/owncloud/ocis-reva/pull/392
   https://github.com/cs3org/reva/pull/986 https://github.com/cs3org/reva/pull/974
   https://github.com/cs3org/reva/pull/972 https://github.com/cs3org/reva/pull/969

  * Enhancement - Make frontend prefixes configurable: [#363](https://github.com/owncloud/ocis-reva/pull/363)

   We introduce three new environment variables and preconfigure them the following way:

  * `REVA_FRONTEND_DATAGATEWAY_PREFIX="data"`
  * `REVA_FRONTEND_OCDAV_PREFIX=""`
  * `REVA_FRONTEND_OCS_PREFIX="ocs"`

   This restores the reva defaults that were changed upstream.

   Https://github.com/owncloud/ocis-reva/pull/363
   https://github.com/cs3org/reva/pull/936/files#diff-51bf4fb310f7362f5c4306581132fc3bR63

  * Enhancement - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66: [#341](https://github.com/owncloud/ocis-reva/pull/341)

   - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66 (#341) - Added country information
   to Mentix (reva/#924) - Refactor metrics package to implement reader interface (reva/#934) -
   Fix OCS public link share update values logic (#252, #288, reva/#930)

   Https://github.com/owncloud/ocis-reva/issues/252
   https://github.com/owncloud/ocis-reva/issues/288
   https://github.com/owncloud/ocis-reva/pull/341
   https://github.com/cs3org/reva/pull/924 https://github.com/cs3org/reva/pull/934
   https://github.com/cs3org/reva/pull/930

  * Enhancement - Update reva to v0.1.1-0.20200709064551-91eed007038f: [#362](https://github.com/owncloud/ocis-reva/pull/362)

   - Update reva to v0.1.1-0.20200709064551-91eed007038f (#362) - Fix config for uploads when
   data server is not exposed (reva/#936) - Update OCM partners endpoints (reva/#937) - Update
   Ailleron endpoint (reva/#938) - OCS: Fix initialization of shares json file (reva/#940) -
   OCS: Fix returned public link URL (#336, reva/#945) - OCS: Share wrap resource id correctly
   (#344, reva/#951) - OCS: Implement share handling for accepting and listing shares (#11,
   reva/#929) - ocm: dynamically lookup IPs for provider check (reva/#946) - ocm: add
   functionality to mail OCM invite tokens (reva/#944) - Change percentagused to
   percentageused (reva/#903) - Fix file-descriptor leak (reva/#954)

   Https://github.com/owncloud/ocis-reva/issues/344
   https://github.com/owncloud/ocis-reva/issues/336
   https://github.com/owncloud/ocis-reva/issues/11
   https://github.com/owncloud/ocis-reva/pull/362
   https://github.com/cs3org/reva/pull/936 https://github.com/cs3org/reva/pull/937
   https://github.com/cs3org/reva/pull/938 https://github.com/cs3org/reva/pull/940
   https://github.com/cs3org/reva/pull/951 https://github.com/cs3org/reva/pull/945
   https://github.com/cs3org/reva/pull/929 https://github.com/cs3org/reva/pull/946
   https://github.com/cs3org/reva/pull/944 https://github.com/cs3org/reva/pull/903
   https://github.com/cs3org/reva/pull/954

  * Enhancement - Add new config options for the http client: [#330](https://github.com/owncloud/ocis-reva/pull/330)

   The internal certificates are checked for validity after
   https://github.com/cs3org/reva/pull/914, which causes the acceptance tests to fail. This
   change sets new hardcoded defaults.

   Https://github.com/owncloud/ocis-reva/pull/330

  * Enhancement - Allow datagateway transfers to take 24h: [#323](https://github.com/owncloud/ocis-reva/pull/323)

   - Increase transfer token life time to 24h (PR #323)

   Https://github.com/owncloud/ocis-reva/pull/323

  * Enhancement - Update reva to v0.1.1-0.20200630075923-39a90d431566: [#320](https://github.com/owncloud/ocis-reva/pull/320)

   - Update reva to v0.1.1-0.20200630075923-39a90d431566 (#320) - Return special value for
   public link password (#294, reva/#904) - Fix public stat and listcontainer response to
   contain the correct prefix (#310, reva/#902)

   Https://github.com/owncloud/ocis-reva/issues/310
   https://github.com/owncloud/ocis-reva/issues/294
   https://github.com/owncloud/ocis-reva/pull/320
   https://github.com/cs3org/reva/pull/902 https://github.com/cs3org/reva/pull/904

  * Enhancement - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66: [#328](https://github.com/owncloud/ocis-reva/pull/328)

   - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66 (#328) - Use sync.Map on pool package
   (reva/#909) - Use mutex instead of sync.Map (reva/#915) - Use gatewayProviders instead of
   storageProviders on conn pool (reva/#916) - Add logic to ls and stat to process arbitrary
   metadata keys (reva/#905) - Preliminary implementation of Set/UnsetArbitraryMetadata
   (reva/#912) - Make datagateway forward headers (reva/#913, reva/#926) - Add option to cmd
   upload to disable tus (reva/#911) - OCS Share Allow date-only expiration for public shares
   (#288, reva/#918) - OCS Share Remove array from OCS Share update response (#252, reva/#919) -
   OCS Share Implement GET request for single shares (#249, reva/#921)

   Https://github.com/owncloud/ocis-reva/issues/288
   https://github.com/owncloud/ocis-reva/issues/252
   https://github.com/owncloud/ocis-reva/issues/249
   https://github.com/owncloud/ocis-reva/pull/328
   https://github.com/cs3org/reva/pull/909 https://github.com/cs3org/reva/pull/915
   https://github.com/cs3org/reva/pull/916 https://github.com/cs3org/reva/pull/905
   https://github.com/cs3org/reva/pull/912 https://github.com/cs3org/reva/pull/913
   https://github.com/cs3org/reva/pull/926 https://github.com/cs3org/reva/pull/911
   https://github.com/cs3org/reva/pull/918 https://github.com/cs3org/reva/pull/919
   https://github.com/cs3org/reva/pull/921

  * Enhancement - Update reva to v0.1.1-0.20200629131207-04298ea1c088: [#309](https://github.com/owncloud/ocis-reva/pull/309)

   - Update reva to v0.1.1-0.20200629094927-e33d65230abc (#309) - Fix public link file share
   (#278, reva/#895, reva/#900) - Delete public share (reva/#899) - Updated reva to
   v0.1.1-0.20200629131207-04298ea1c088 (#313)

   Https://github.com/owncloud/ocis-reva/issues/278
   https://github.com/owncloud/ocis-reva/pull/309
   https://github.com/cs3org/reva/pull/895 https://github.com/cs3org/reva/pull/899
   https://github.com/cs3org/reva/pull/900
   https://github.com/owncloud/ocis-reva/pull/313

  * Enhancement - Update reva to v0.1.1-0.20200626111234-e21c32db9614: [#261](https://github.com/owncloud/ocis-reva/pull/261)

   - Updated reva to v0.1.1-0.20200626111234-e21c32db9614 (#304) - TUS upload support through
   datagateway (#261, reva/#878, reva/#888) - Added support for differing metrics path for
   Prometheus to Mentix (reva/#875) - More data exported by Mentix (reva/#881) - Implementation
   of file operations in public folder shares (#49, #293, reva/#877) - Make httpclient trust
   local certificates for now (reva/#880) - EOS homes are not configured with an enable-flag
   anymore, but with a dedicated storage driver. We're using it now and adapted default configs of
   storages (reva/#891, #304)

   Https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/issues/293
   https://github.com/owncloud/ocis-reva/issues/261
   https://github.com/owncloud/ocis-reva/pull/261
   https://github.com/cs3org/reva/pull/875 https://github.com/cs3org/reva/pull/877
   https://github.com/cs3org/reva/pull/878 https://github.com/cs3org/reva/pull/881
   https://github.com/cs3org/reva/pull/880 https://github.com/cs3org/reva/pull/888
   https://github.com/owncloud/ocis-reva/pull/304
   https://github.com/cs3org/reva/pull/891

  * Enhancement - Update reva to v0.1.1-0.20200624063447-db5e6635d5f0: [#279](https://github.com/owncloud/ocis-reva/pull/279)

   - Updated reva to v0.1.1-0.20200624063447-db5e6635d5f0 (#279) - Local storage: URL-encode
   file ids to ease integration with other microservices like WOPI (reva/#799) - Mentix fixes
   (reva/#803, reva/#817) - OCDAV: fix returned timestamp format (#116, reva/#805) - OCM: add
   default prefix (#814) - add the content-length header to the responses (reva/#816) - Deps:
   clean (reva/#818) - Fix trashbin listing (#112, #253, #254, reva/#819) - Make the json
   publicshare driver configurable (reva/#820) - TUS: Return metadata headers after direct
   upload (ocis/#216, reva/#813) - Set mtime to storage after simple upload (#174, reva/#823,
   reva/#841) - Configure grpc client to allow for insecure conns and skip server certificate
   verification (reva/#825) - Deployment: simplify config with more default values
   (reva/#826, reva/#837, reva/#843, reva/#848, reva/#842) - Separate local fs into home and
   with home disabled (reva/#829) - Register reflection after other services (reva/#831) -
   Refactor EOS fs (reva/#830) - Add ocs-share-permissions to the propfind response (#47,
   reva/#836) - OCS: Properly read permissions when creating public link (reva/#852) - localfs:
   make normalize return associated error (reva/#850) - EOS grpc driver (reva/#664) - OCS: Add
   support for legacy public link arg publicUpload (reva/#853) - Add cache layer to user REST
   package (reva/#849) - Meshdirectory: pass query params to selected provider (reva/#863) -
   Pass etag in quotes from the fs layer (#269, reva/#866, reva/#867) - OCM: use refactored
   cs3apis provider definition (reva/#864)

   Https://github.com/owncloud/ocis-reva/issues/116
   https://github.com/owncloud/ocis-reva/issues/112
   https://github.com/owncloud/ocis-reva/issues/253
   https://github.com/owncloud/ocis-reva/issues/254
   https://github.com/owncloud/ocis/issues/216
   https://github.com/owncloud/ocis-reva/issues/174
   https://github.com/owncloud/ocis-reva/issues/47
   https://github.com/owncloud/ocis-reva/issues/269
   https://github.com/owncloud/ocis-reva/pull/279
   https://github.com/owncloud/cs3org/reva/pull/799
   https://github.com/owncloud/cs3org/reva/pull/803
   https://github.com/owncloud/cs3org/reva/pull/817
   https://github.com/owncloud/cs3org/reva/pull/805
   https://github.com/owncloud/cs3org/reva/pull/814
   https://github.com/owncloud/cs3org/reva/pull/816
   https://github.com/owncloud/cs3org/reva/pull/818
   https://github.com/owncloud/cs3org/reva/pull/819
   https://github.com/owncloud/cs3org/reva/pull/820
   https://github.com/owncloud/cs3org/reva/pull/823
   https://github.com/owncloud/cs3org/reva/pull/841
   https://github.com/owncloud/cs3org/reva/pull/813
   https://github.com/owncloud/cs3org/reva/pull/825
   https://github.com/owncloud/cs3org/reva/pull/826
   https://github.com/owncloud/cs3org/reva/pull/837
   https://github.com/owncloud/cs3org/reva/pull/843
   https://github.com/owncloud/cs3org/reva/pull/848
   https://github.com/owncloud/cs3org/reva/pull/842
   https://github.com/owncloud/cs3org/reva/pull/829
   https://github.com/owncloud/cs3org/reva/pull/831
   https://github.com/owncloud/cs3org/reva/pull/830
   https://github.com/owncloud/cs3org/reva/pull/836
   https://github.com/owncloud/cs3org/reva/pull/852
   https://github.com/owncloud/cs3org/reva/pull/850
   https://github.com/owncloud/cs3org/reva/pull/664
   https://github.com/owncloud/cs3org/reva/pull/853
   https://github.com/owncloud/cs3org/reva/pull/849
   https://github.com/owncloud/cs3org/reva/pull/863
   https://github.com/owncloud/cs3org/reva/pull/866
   https://github.com/owncloud/cs3org/reva/pull/867
   https://github.com/owncloud/cs3org/reva/pull/864

  * Enhancement - Add TUS global capability: [#177](https://github.com/owncloud/ocis-reva/issues/177)

   The TUS global capabilities from Reva are now exposed.

   The advertised max chunk size can be configured using the "--upload-max-chunk-size" CLI
   switch or "REVA_FRONTEND_UPLOAD_MAX_CHUNK_SIZE" environment variable. The advertised
   http method override can be configured using the "--upload-http-method-override" CLI
   switch or "REVA_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE" environment variable.

   Https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/owncloud/ocis-reva/pull/228

  * Enhancement - Update reva to v0.1.1-0.20200603071553-e05a87521618: [#244](https://github.com/owncloud/ocis-reva/issues/244)

   - Updated reva to v0.1.1-0.20200603071553-e05a87521618 (#244) - Add option to disable TUS on
   OC layer (#177, reva/#791) - Dataprovider now supports method override (#177, reva/#792) -
   OCS fixes for create public link (reva/#798)

   Https://github.com/owncloud/ocis-reva/issues/244
   https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/cs3org/reva/pull/791 https://github.com/cs3org/reva/pull/792
   https://github.com/cs3org/reva/pull/798

  * Enhancement - Add public shares service: [#49](https://github.com/owncloud/ocis-reva/issues/49)

   Added Public Shares service with CRUD operations and File Public Shares Manager

   Https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/pull/232

  * Enhancement - Update reva to v0.1.1-0.20200529120551-4f2d9c85d3c9: [#49](https://github.com/owncloud/ocis-reva/issues/49)

   - Updated reva to v0.1.1-0.20200529120551 (#232) - Public Shares CRUD, File Public Shares
   Manager (#49, #232, reva/#681, reva/#788) - Disable HTTP-KeepAlives to reduce fd count
   (ocis/#268, reva/#787) - Fix trashbin listing (#229, reva/#782) - Create PUT wrapper for TUS
   uploads (reva/#770) - Add security access headers for ocdav requests (#66, reva/#780) - Add
   option to revad cmd to specify logging level (reva/#772) - New metrics package (reva/#740) -
   Remove implicit data member from memory store (reva/#774) - Added TUS global capabilities
   (#177, reva/#775) - Fix PROPFIND with Depth 1 for cross-storage operations (reva/#779)

   Https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/issues/229
   https://github.com/owncloud/ocis-reva/issues/66
   https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/owncloud/ocis/issues/268
   https://github.com/owncloud/ocis-reva/pull/232
   https://github.com/cs3org/reva/pull/787 https://github.com/cs3org/reva/pull/681
   https://github.com/cs3org/reva/pull/788 https://github.com/cs3org/reva/pull/782
   https://github.com/cs3org/reva/pull/770 https://github.com/cs3org/reva/pull/780
   https://github.com/cs3org/reva/pull/772 https://github.com/cs3org/reva/pull/740
   https://github.com/cs3org/reva/pull/774 https://github.com/cs3org/reva/pull/775
   https://github.com/cs3org/reva/pull/779

  * Enhancement - Update reva to v0.1.1-0.20200520150229: [#161](https://github.com/owncloud/ocis-reva/pull/161)

   - Update reva to v0.1.1-0.20200520150229 (#161, #180, #192, #207, #221) - Return arbitrary
   metadata with stat, upload without TUS (reva/#766) - Stat file before returning datagateway
   URL when initiating download (reva/#765) - REST driver for user package (reva/#747) - Sharing
   behavior now consistent with the old backend (#20, #26, #43, #44, #46, #94 ,reva/#748) - Mentix
   service (reva/#755) - meshdirectory: add mentix driver for gocdb sites integration
   (reva/#754) - Add functionality to commit to storage for OCM shares (reva/#760) - Add option in
   config to disable tus (reva/#759) - ocdav: fix custom property XML parsing in PROPPATCH
   handler (#203, reva/#743) - ocdav: fix PROPPATCH response for removed properties (#186,
   reva/#742) - ocdav: implement PROPFIND infinity depth (#212, reva/#758) - Local fs: Allow
   setting of arbitrary metadata, minor bug fixes (reva/#764) - Local fs: metadata handling and
   share persistence (reva/#732) - Local fs: return file owner info in stat (reva/#750) - Fixed
   regression when uploading empty files to OCFS or EOS with PUT and TUS (#188, reva/#734) - On
   delete move the file versions to the trashbin (#94, reva/#731) - Fix OCFS move operation (#182,
   reva/#729) - Fix OCFS custom property / xattr removal (reva/#728) - Retry trashbin in case of
   timestamp collision (reva/#730) - Disable chunking v1 by default (reva/#678) - Implement ocs
   to http status code mapping (#26, reva/#696, reva/#707, reva/#711) - Handle the case if
   directory already exists (reva/#695) - Added TUS upload support (reva/#674, reva/#725,
   reva/#717) - Always return file sizes in Webdav PROPFIND (reva/#712) - Use default mime type
   when none was detected (reva/#713) - Fixed Webdav shallow COPY (reva/#714) - Fixed arbitrary
   namespace usage for custom properties in PROPFIND (#57, reva/#720) - Implement returning
   Webdav custom properties from xattr (#57, reva/#721) - Minor fix in OCM share pkg (reva/#718)

   Https://github.com/owncloud/ocis-reva/issues/20
   https://github.com/owncloud/ocis-reva/issues/26
   https://github.com/owncloud/ocis-reva/issues/43
   https://github.com/owncloud/ocis-reva/issues/44
   https://github.com/owncloud/ocis-reva/issues/46
   https://github.com/owncloud/ocis-reva/issues/94
   https://github.com/owncloud/ocis-reva/issues/26
   https://github.com/owncloud/ocis-reva/issues/67
   https://github.com/owncloud/ocis-reva/issues/57
   https://github.com/owncloud/ocis-reva/issues/94
   https://github.com/owncloud/ocis-reva/issues/188
   https://github.com/owncloud/ocis-reva/issues/182
   https://github.com/owncloud/ocis-reva/issues/212
   https://github.com/owncloud/ocis-reva/issues/186
   https://github.com/owncloud/ocis-reva/issues/203
   https://github.com/owncloud/ocis-reva/pull/161
   https://github.com/owncloud/ocis-reva/pull/180
   https://github.com/owncloud/ocis-reva/pull/192
   https://github.com/owncloud/ocis-reva/pull/207
   https://github.com/owncloud/ocis-reva/pull/221
   https://github.com/cs3org/reva/pull/766 https://github.com/cs3org/reva/pull/765
   https://github.com/cs3org/reva/pull/755 https://github.com/cs3org/reva/pull/754
   https://github.com/cs3org/reva/pull/747 https://github.com/cs3org/reva/pull/748
   https://github.com/cs3org/reva/pull/760 https://github.com/cs3org/reva/pull/759
   https://github.com/cs3org/reva/pull/678 https://github.com/cs3org/reva/pull/696
   https://github.com/cs3org/reva/pull/707 https://github.com/cs3org/reva/pull/711
   https://github.com/cs3org/reva/pull/695 https://github.com/cs3org/reva/pull/674
   https://github.com/cs3org/reva/pull/725 https://github.com/cs3org/reva/pull/717
   https://github.com/cs3org/reva/pull/712 https://github.com/cs3org/reva/pull/713
   https://github.com/cs3org/reva/pull/720 https://github.com/cs3org/reva/pull/718
   https://github.com/cs3org/reva/pull/731 https://github.com/cs3org/reva/pull/734
   https://github.com/cs3org/reva/pull/729 https://github.com/cs3org/reva/pull/728
   https://github.com/cs3org/reva/pull/730 https://github.com/cs3org/reva/pull/758
   https://github.com/cs3org/reva/pull/742 https://github.com/cs3org/reva/pull/764
   https://github.com/cs3org/reva/pull/743 https://github.com/cs3org/reva/pull/732
   https://github.com/cs3org/reva/pull/750

  * Bugfix - Stop advertising unsupported chunking v2: [#145](https://github.com/owncloud/ocis-reva/pull/145)

   Removed "chunking" attribute in the DAV capabilities. Please note that chunking v2 is
   advertised as "chunking 1.0" while chunking v1 is the attribute "bigfilechunking" which is
   already false.

   Https://github.com/owncloud/ocis-reva/pull/145

  * Enhancement - Allow configuring the gateway for dataproviders: [#136](https://github.com/owncloud/ocis-reva/pull/136)

   This allows using basic or bearer auth when directly talking to dataproviders.

   Https://github.com/owncloud/ocis-reva/pull/136

  * Enhancement - Use a configured logger on reva runtime: [#153](https://github.com/owncloud/ocis-reva/pull/153)

   For consistency reasons we need a configured logger that is inline with an ocis logger, so the
   log cascade can be easily parsed by a human.

   Https://github.com/owncloud/ocis-reva/pull/153

  * Bugfix - Fix eos user sharing config: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   We have added missing config options for the user sharing manager and added a dedicated eos
   storage command with pre configured settings for the eos-docker container. It configures a
   `Shares` folder in a users home when using eos as the storage driver.

   Https://github.com/owncloud/ocis-reva/pull/127

  * Enhancement - Update reva to v1.1.0-20200414133413: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   Adds initial public sharing and ocm implementation.

   Https://github.com/owncloud/ocis-reva/pull/127

  * Bugfix - Fix eos config: [#125](https://github.com/owncloud/ocis-reva/pull/125)

   We have added missing config options for the home layout to the config struct that is passed to
   eos.

   Https://github.com/owncloud/ocis-reva/pull/125

  * Bugfix - Set correct flag type in the flagsets: [#75](https://github.com/owncloud/ocis-reva/issues/75)

   While upgrading to the micro/cli version 2 there where two instances of `StringFlag` which had
   not been changed to `StringSliceFlag`. This caused `ocis-reva users` and `ocis-reva
   storage-root` to fail on startup.

   Https://github.com/owncloud/ocis-reva/issues/75
   https://github.com/owncloud/ocis-reva/pull/76

  * Bugfix - We fixed a typo in the `REVA_LDAP_SCHEMA_MAIL` environment variable: [#113](https://github.com/owncloud/ocis-reva/pull/113)

   It was misspelled as `REVA_LDAP_SCHEMA_Mail`.

   Https://github.com/owncloud/ocis-reva/pull/113

  * Bugfix - Allow different namespaces for /webdav and /dav/files: [#68](https://github.com/owncloud/ocis-reva/pull/68)

   After fbf131c the path for the "new" webdav path does not contain a username
   `/remote.php/dav/files/textfile0.txt`. It used to be
   `/remote.php/dav/files/oc/einstein/textfile0.txt` So it lost `oc/einstein`.

   This PR allows setting up different namespaces for `/webav` and `/dav/files`:

   `/webdav` is jailed into `/home` - which uses the home storage driver and uses the logged in user
   to construct the path `/dav/files` is jailed into `/oc` - which uses the owncloud storage
   driver and expects a username as the first path segment

   This mimics oc10

   The `WEBDAV_NAMESPACE_JAIL` environment variable is split into - `WEBDAV_NAMESPACE` and -
   `DAV_FILES_NAMESPACE` accordingly.

   Https://github.com/owncloud/ocis-reva/pull/68 related:

  * Change - Use /home as default namespace: [#68](https://github.com/owncloud/ocis-reva/pull/68)

   Currently, cross storage etag propagation is not yet implemented, which prevents the desktop
   client from detecting changes via the PROPFIND to /. / is managed by the root storage provider
   which is independend of the home and oc storage providers. If a file changes in /home/foo, the
   etag change will only be propagated to the root of the home storage provider.

   This change jails users into the `/home` namespace, and allows configuring the namespace to
   use for the two webdav endpoints using the new environment variable `WEBDAV_NAMESPACE_JAIL`
   which affects both endpoints `/dav/files` and `/webdav`.

   This will allow us to focus on getting a single storage driver like eos or owncloud tested and
   better resembles what owncloud 10 does.

   To get back the global namespace, which ultimately is the goal, just set the above environment
   variable to `/`.

   Https://github.com/owncloud/ocis-reva/pull/68

  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-reva/issues/1)

   Just prepared an initial basic version to start a reva server and start integrating with the
   go-micro base dextension framework of ownCloud Infinite Scale.

   Https://github.com/owncloud/ocis-reva/issues/1

  * Change - Start multiple services with dedicated commands: [#6](https://github.com/owncloud/ocis-reva/issues/6)

   The initial version would only allow us to use a set of reva configurations to start multiple
   services. We use a more opinionated set of commands to start dedicated services that allows us
   to configure them individually. It allows us to switch eg. the user backend to LDAP and fully use
   it on the cli.

   Https://github.com/owncloud/ocis-reva/issues/6

  * Change - Storage providers now default to exposing data servers: [#89](https://github.com/owncloud/ocis-reva/issues/89)

   The flags that let reva storage providers announce that they expose a data server now defaults
   to true:

   `REVA_STORAGE_HOME_EXPOSE_DATA_SERVER=1` `REVA_STORAGE_OC_EXPOSE_DATA_SERVER=1`

   Https://github.com/owncloud/ocis-reva/issues/89

  * Change - Default to running behind ocis-proxy: [#113](https://github.com/owncloud/ocis-reva/pull/113)

   We changed the default configuration to integrate better with ocis.

   - We use ocis-glauth as the default ldap server on port 9125 with base `dc=example,dc=org`. - We
   use a dedicated technical `reva` user to make ldap binds - Clients are supposed to use the
   ocis-proxy endpoint `https://localhost:9200` - We removed unneeded ocis configuration
   from the frontend which no longer serves an oidc provider. - We changed the default user
   OpaqueID attribute from `sub` to `preferred_username`. The latter is a claim populated by
   konnectd that can also be used by the reva ldap user manager to look up users by their OpaqueId

   Https://github.com/owncloud/ocis-reva/pull/113

  * Enhancement - Expose owncloud storage driver config in flagset: [#87](https://github.com/owncloud/ocis-reva/issues/87)

   Three new flags are now available:

   - scan files on startup to generate missing fileids default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_SCAN` cli option: `--storage-owncloud-scan`

   - autocreate home path for new users default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_AUTOCREATE` cli option: `--storage-owncloud-autocreate`

   - the address of the redis server default: `:6379` env var:
   `REVA_STORAGE_OWNCLOUD_REDIS_ADDR` cli option: `--storage-owncloud-redis`

   Https://github.com/owncloud/ocis-reva/issues/87

  * Enhancement - Update reva to v0.0.2-0.20200212114015-0dbce24f7e8b: [#91](https://github.com/owncloud/ocis-reva/pull/91)

   Reva has seen a lot of changes that allow us to - reduce the configuration overhead - use the
   autocreato home folder option - use the home folder path layout option - no longer start the root
   storage

   Https://github.com/owncloud/ocis-reva/pull/91 related:

  * Enhancement - Allow configuring user sharing driver: [#115](https://github.com/owncloud/ocis-reva/pull/115)

   We now default to `json` which persists shares in the sharing manager in a json file instead of an
   in memory db.

   Https://github.com/owncloud/ocis-reva/pull/115

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the store service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: store

  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Removed code from other service: [#7](https://github.com/owncloud/ocis-store/pull/7)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#5](https://github.com/owncloud/ocis-store/pull/5)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-store/pull/1)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add the thumbnails service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: thumbnails

  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#35](https://github.com/owncloud/ocis-thumbnails/pull/35)
  * Enhancement - Serve the metrics endpoint: [#37](https://github.com/owncloud/ocis-thumbnails/issues/37)
  * Change - Add more default resolutions: [#23](https://github.com/owncloud/ocis-thumbnails/issues/23)
  * Change - Refactor code to remove code smells: [#21](https://github.com/owncloud/ocis-thumbnails/issues/21)
  * Change - Use micro service error api: [#31](https://github.com/owncloud/ocis-thumbnails/issues/31)
  * Enhancement - Limit users to access own thumbnails: [#5](https://github.com/owncloud/ocis-thumbnails/issues/5)
  * Bugfix - Fix usage of context.Context: [#18](https://github.com/owncloud/ocis-thumbnails/issues/18)
  * Bugfix - Fix execution when passing program flags: [#15](https://github.com/owncloud/ocis-thumbnails/issues/15)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-thumbnails/issues/1)
  * Change - Use predefined resolutions for thumbnail generation: [#7](https://github.com/owncloud/ocis-thumbnails/issues/7)
  * Change - Implement the first working version: [#3](https://github.com/owncloud/ocis-thumbnails/pull/3)

   https://github.com/owncloud/product/issues/244

* Enhancement - Add a command to list the versions of running instances: [#226](https://github.com/owncloud/product/issues/226)

   Tags: accounts

   Added a micro command to list the versions of running accounts services.

   https://github.com/owncloud/product/issues/226

* Enhancement - Add the webdav service: [#244](https://github.com/owncloud/product/issues/244)

   Tags: webdav

  * Enhancement - Add version command: [#226](https://github.com/owncloud/product/issues/226)
  * Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#22](https://github.com/owncloud/ocis-webdav/pull/22)
  * Change Change status not found on missing thumbnail: [#20](https://github.com/owncloud/ocis-webdav/issues/20)
  * Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-webdav/issues/1)
  * Change - Update ocis-pkg to version 2.2.0: [#16](https://github.com/owncloud/ocis-webdav/issues/16)
  * Enhancement - Configuration: [#14](https://github.com/owncloud/ocis-webdav/pull/14)
  * Enhancement - Implement preview API: [#13](https://github.com/owncloud/ocis-webdav/pull/13)

   https://github.com/owncloud/product/issues/244

* Enhancement - Better adopt Go-Micro: [#840](https://github.com/owncloud/ocis/pull/840)

   Tags: ocis

   There are a few building blocks that we were relying on default behavior, such as
   `micro.Registry` and the go-micro client. In order for oCIS to work in any environment and not
   relying in black magic configuration or running daemons we need to be able to:

   - Provide with a configurable go-micro registry. - Use our own go-micro client adjusted to our
   own needs (i.e: custom timeout, custom dial timeout, custom transport...)

   This PR is relying on 2 env variables from Micro: `MICRO_REGISTRY` and
   `MICRO_REGISTRY_ADDRESS`. The latter does not make sense to provide if the registry is not
   `etcd`.

   The current implementation only accounts for `mdns` and `etcd` registries, defaulting to
   `mdns` when not explicitly defined to use `etcd`.

   https://github.com/owncloud/ocis/pull/840

* Enhancement - Add permission check when assigning and removing roles: [#879](https://github.com/owncloud/ocis/issues/879)

   Everyone could add and remove roles from users. Added a new permission and a check so that only
   users with the role management permissions can assign and unassign roles.

   https://github.com/owncloud/ocis/issues/879

* Enhancement - Create OnlyOffice extension: [#857](https://github.com/owncloud/ocis/pull/857)

   Tags: OnlyOffice

   We've created an OnlyOffice extension which enables users to create and edit docx documents
   and open spreadsheets and presentations.

   https://github.com/owncloud/ocis/pull/857

* Enhancement - Show basic-auth warning only once: [#886](https://github.com/owncloud/ocis/pull/886)

   Show basic-auth warning only on startup instead on every request.

   https://github.com/owncloud/ocis/pull/886

* Enhancement - Add glauth fallback backend: [#649](https://github.com/owncloud/ocis/pull/649)

   We introduced the `fallback-datastore` config option and the corresponding options to allow
   configuring a simple chain of two handlers.

   Simple, because it is intended for bind and single result search queries. Merging large sets of
   results is currently out of scope. For now, the implementation will only search the fallback
   backend if the default backend returns an error or the number of results is 0. This is sufficient
   to allow an IdP to authenticate users from ocis as well as owncloud 10 as described in the [bridge
   scenario](https://owncloud.github.io/ocis/deployment/bridge/).

   https://github.com/owncloud/ocis-glauth/issues/18
   https://github.com/owncloud/ocis/pull/649

* Enhancement - Tidy dependencies: [#845](https://github.com/owncloud/ocis/pull/845)

   Methodology:

   ``` go-modules() { find . \( -name vendor -o -name '[._].*' -o -name node_modules \) -prune -o
   -name go.mod -print | sed 's:/go.mod$::' } ```

   ``` for m in $(go-modules); do (cd $m && go mod tidy); done ```

   https://github.com/owncloud/ocis/pull/845

* Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)

   Tags: metadata, accounts, settings

   In the future accounts, settings etc. should be stored in a dedicated metadata storage. The
   services should talk to this storage directly, bypassing reva-gateway.

   https://github.com/owncloud/ocis/pull/602

* Enhancement - Add a version command to ocis: [#915](https://github.com/owncloud/ocis/pull/915)

   The version command was only implemented in the extensions. This adds the version command to
   ocis to list all services in the ocis namespace.

   https://github.com/owncloud/ocis/pull/915

* Enhancement - Create a proxy access-log: [#889](https://github.com/owncloud/ocis/pull/889)

   Logs client access at the proxy

   https://github.com/owncloud/ocis/pull/889

* Enhancement - Cache userinfo in proxy: [#877](https://github.com/owncloud/ocis/pull/877)

   Tags: proxy

   We introduced caching for the userinfo response. The token expiration is used for cache
   invalidation if available. Otherwise we fall back to a preconfigured TTL (default 10
   seconds).

   https://github.com/owncloud/ocis/pull/877

* Enhancement - Update reva to v1.4.1-0.20201125144025-57da0c27434c: [#1320](https://github.com/cs3org/reva/pull/1320)

   Mostly to bring fixes to pressing changes.

   https://github.com/cs3org/reva/pull/1320
   https://github.com/cs3org/reva/pull/1338

* Enhancement - Runtime Cleanup: [#1066](https://github.com/owncloud/ocis/pull/1066)

   Small runtime cleanup prior to Tech Preview release

   https://github.com/owncloud/ocis/pull/1066

* Enhancement - Update OCIS Runtime: [#1108](https://github.com/owncloud/ocis/pull/1108)

   - enhances the overall behavior of our runtime - runtime `db` file configurable - two new env
   variables to deal with the runtime - `RUNTIME_DB_FILE` and `RUNTIME_KEEP_ALIVE` -
   `RUNTIME_KEEP_ALIVE` defaults to `false` to provide backwards compatibility - if
   `RUNTIME_KEEP_ALIVE` is set to `true`, if a supervised process terminates the runtime will
   attempt to start with the same environment provided.

   https://github.com/owncloud/ocis/pull/1108

* Enhancement - Simplify tracing config: [#92](https://github.com/owncloud/product/issues/92)

   We now apply the oCIS tracing config to all services which have tracing. With this it is possible
   to set one tracing config for all services at the same time.

   https://github.com/owncloud/product/issues/92
   https://github.com/owncloud/ocis/pull/329
   https://github.com/owncloud/ocis/pull/409

* Enhancement - Update glauth to dev fd3ac7e4bbdc93578655d9a08d8e23f105aaa5b2: [#834](https://github.com/owncloud/ocis/pull/834)

   We updated glauth to dev commit fd3ac7e4bbdc93578655d9a08d8e23f105aaa5b2, which allows to
   skip certificate checks for the owncloud backend.

   https://github.com/owncloud/ocis/pull/834

* Enhancement - Update glauth to dev 4f029234b2308: [#786](https://github.com/owncloud/ocis/pull/786)

   Includes a bugfix, don't mix graph and provisioning api.

   https://github.com/owncloud/ocis/pull/786

* Enhancement - Update konnectd to v0.33.8: [#744](https://github.com/owncloud/ocis/pull/744)

   This update adds options which allow the configuration of oidc-token expiration parameters:
   KONNECTD_ACCESS_TOKEN_EXPIRATION, KONNECTD_ID_TOKEN_EXPIRATION and
   KONNECTD_REFRESH_TOKEN_EXPIRATION.

   Other changes from upstream:

   - Generate random endsession state for external authority - Update dependencies in
   Dockerfile - Set prompt=None to avoid loops with external authority - Update Jenkins
   reporting plugin from checkstyle to recordIssues - Remove extra kty key from JWKS top level
   document - Fix regression which encodes URL fragments twice - Avoid generating fragmet/query
   URLs with wrong order - Return state for oidc endsession response redirects - Use server
   provided username to avoid case mismatch - Use signed-out-uri if set as fallback for goodbye
   redirect on saml slo - Add checks to ensure post_logout_redirect_uri is not empty - Fix SAML2
   logout request parsing - Cure panic when no state is found in saml esr - Use SAML IdP Issuer value
   from meta data entityID - Allow configuration of expiration of oidc access, id and refresh
   tokens - Implement trampolin for external OIDC authority end session - Update
   ca-certificates version

   https://github.com/owncloud/ocis/pull/744

* Enhancement - Update reva to v1.4.1-0.20201123062044-b2c4af4e897d: [#823](https://github.com/owncloud/ocis/pull/823)

  * Refactor the uploading files workflow from various clients [cs3org/reva#1285](https://github.com/cs3org/reva/pull/1285), [cs3org/reva#1314](https://github.com/cs3org/reva/pull/1314)
  * [OCS] filter share with me requests [cs3org/reva#1302](https://github.com/cs3org/reva/pull/1302)
  * Fix listing shares for nonexisting path [cs3org/reva#1316](https://github.com/cs3org/reva/pull/1316)
  * prevent nil pointer when listing shares [cs3org/reva#1317](https://github.com/cs3org/reva/pull/1317)
  * Sharee retrieves the information about a share -but gets response containing all the shares [owncloud/ocis-reva#260](https://github.com/owncloud/ocis-reva/issues/260)
  * Deleting a public link after renaming a file [owncloud/ocis-reva#311](https://github.com/owncloud/ocis-reva/issues/311)
  * Avoid log spam [cs3org/reva#1323](https://github.com/cs3org/reva/pull/1323), [cs3org/reva#1324](https://github.com/cs3org/reva/pull/1324)
  * Fix trashbin [cs3org/reva#1326](https://github.com/cs3org/reva/pull/1326)

   https://github.com/owncloud/ocis-reva/issues/260
   https://github.com/owncloud/ocis-reva/issues/311
   https://github.com/owncloud/ocis/pull/823
   https://github.com/cs3org/reva/pull/1285
   https://github.com/cs3org/reva/pull/1302
   https://github.com/cs3org/reva/pull/1314
   https://github.com/cs3org/reva/pull/1316
   https://github.com/cs3org/reva/pull/1317
   https://github.com/cs3org/reva/pull/1323
   https://github.com/cs3org/reva/pull/1324
   https://github.com/cs3org/reva/pull/1326

* Enhancement - Update reva to v1.4.1-0.20201130061320-ac85e68e0600: [#980](https://github.com/owncloud/ocis/pull/980)

  * Fix move operation in ocis storage driver [csorg/reva#1343](https://github.com/cs3org/reva/pull/1343)

   https://github.com/owncloud/ocis/issues/975
   https://github.com/owncloud/ocis/pull/980
   https://github.com/cs3org/reva/pull/1343

* Enhancement - Update reva to cdb3d6688da5: [#748](https://github.com/owncloud/ocis/pull/748)

  * let the gateway filter invalid references

   https://github.com/owncloud/ocis/pull/748
   https://github.com/cs3org/reva/pull/1274

* Enhancement - Update reva to dd3a8c0f38: [#725](https://github.com/owncloud/ocis/pull/725)

  * fixes etag propagation in the ocis driver

   https://github.com/owncloud/ocis/pull/725
   https://github.com/cs3org/reva/pull/1264

* Enhancement - Update reva to v1.4.1-0.20201127111856-e6a6212c1b7b: [#971](https://github.com/owncloud/ocis/pull/971)

   Tags: reva

  * Fix capabilities response for multiple client versions #1331 [cs3org/reva#1331](https://github.com/cs3org/reva/pull/1331)
  * Fix home storage redirect for remote.php/dav/files [cs3org/reva#1342](https://github.com/cs3org/reva/pull/1342)

   https://github.com/owncloud/ocis/pull/971
   https://github.com/cs3org/reva/pull/1331
   https://github.com/cs3org/reva/pull/1342

* Enhancement - Update reva to 063b3db9162b: [#1091](https://github.com/owncloud/ocis/pull/1091)

   - bring public link removal changes to OCIS. - fix subcommand name collision from renaming
   phoenix -> web.

   https://github.com/owncloud/ocis/issues/1098
   https://github.com/owncloud/ocis/pull/1091

* Enhancement - Add www-authenticate based on user agent: [#1009](https://github.com/owncloud/ocis/pull/1009)

   Tags: reva, proxy

   We now comply with HTTP spec by adding Www-Authenticate headers on every `401` request.
   Furthermore, we not only take care of such a thing at the Proxy but also Reva will take care of it.
   In addition, we now are able to lock-in a set of User-Agent to specific challenges.

   Admins can use this feature by configuring oCIS + Reva following this approach:

   ``` STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT="mirall:basic,
   Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:83.0) Gecko/20100101
   Firefox/83.0:bearer" \
   PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT="mirall:basic, Mozilla/5.0
   (Macintosh; Intel Mac OS X 10.15; rv:83.0) Gecko/20100101 Firefox/83.0:bearer" \
   PROXY_ENABLE_BASIC_AUTH=true \ go run cmd/ocis/main.go server ```

   We introduced two new environment variables:

   `STORAGE_FRONTEND_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT` as well as
   `PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT`, The reason they have the same value
   is not to rely on the os env on a distributed environment, so in redundancy we trust. They both
   configure the same on the backend storage and oCIS Proxy.

   https://github.com/owncloud/ocis/pull/1009

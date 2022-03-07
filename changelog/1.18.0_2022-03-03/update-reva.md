Enhancement: Update reva to v2.0.0

We updated reva to the version 2.0.0.

* Fix [cs3org/reva#2457](https://github.com/cs3org/reva/pull/2457) :  Do not swallow error
* Fix [cs3org/reva#2422](https://github.com/cs3org/reva/pull/2422) :  Handle non existing spaces correctly
* Fix [cs3org/reva#2327](https://github.com/cs3org/reva/pull/2327) :  Enable changelog on edge branch
* Fix [cs3org/reva#2370](https://github.com/cs3org/reva/pull/2370) :  Fixes for apps in public shares, project spaces for EOS driver
* Fix [cs3org/reva#2464](https://github.com/cs3org/reva/pull/2464) :  Pass spacegrants when adding member to space
* Fix [cs3org/reva#2430](https://github.com/cs3org/reva/pull/2430) :  Fix aggregated child folder id
* Fix [cs3org/reva#2348](https://github.com/cs3org/reva/pull/2348) :  Make archiver handle spaces protocol
* Fix [cs3org/reva#2452](https://github.com/cs3org/reva/pull/2452) :  Fix create space error message
* Fix [cs3org/reva#2445](https://github.com/cs3org/reva/pull/2445) :  Don't handle ids containing "/" in decomposedfs
* Fix [cs3org/reva#2285](https://github.com/cs3org/reva/pull/2285) :  Accept new userid idp format
* Fix [cs3org/reva#2503](https://github.com/cs3org/reva/pull/2503) :  Remove the protection from /v?.php/config endpoints
* Fix [cs3org/reva#2462](https://github.com/cs3org/reva/pull/2462) :  Public shares path needs to be set
* Fix [cs3org/reva#2427](https://github.com/cs3org/reva/pull/2427) :  Fix registry caching
* Fix [cs3org/reva#2298](https://github.com/cs3org/reva/pull/2298) :  Remove share refs from trashbin
* Fix [cs3org/reva#2433](https://github.com/cs3org/reva/pull/2433) :  Fix shares provider filter
* Fix [cs3org/reva#2351](https://github.com/cs3org/reva/pull/2351) :  Fix Statcache removing
* Fix [cs3org/reva#2374](https://github.com/cs3org/reva/pull/2374) :  Fix webdav copy of zero byte files
* Fix [cs3org/reva#2336](https://github.com/cs3org/reva/pull/2336) :  Handle sending all permissions when creating public links
* Fix [cs3org/reva#2440](https://github.com/cs3org/reva/pull/2440) :  Add ArbitraryMetadataKeys to statcache key
* Fix [cs3org/reva#2582](https://github.com/cs3org/reva/pull/2582) :  Keep lock structs in a local map protected by a mutex
* Fix [cs3org/reva#2372](https://github.com/cs3org/reva/pull/2372) :  Make owncloudsql work with the spaces registry
* Fix [cs3org/reva#2416](https://github.com/cs3org/reva/pull/2416) :  The registry now returns complete space structs
* Fix [cs3org/reva#3066](https://github.com/cs3org/reva/pull/3066) :  Fix propfind listing for files
* Fix [cs3org/reva#2428](https://github.com/cs3org/reva/pull/2428) :  Remove unused home provider from config
* Fix [cs3org/reva#2334](https://github.com/cs3org/reva/pull/2334) :  Revert fix decomposedfs upload
* Fix [cs3org/reva#2415](https://github.com/cs3org/reva/pull/2415) :  Services should never return transport level errors
* Fix [cs3org/reva#2419](https://github.com/cs3org/reva/pull/2419) :  List project spaces for share recipients
* Fix [cs3org/reva#2501](https://github.com/cs3org/reva/pull/2501) :  Fix spaces stat
* Fix [cs3org/reva#2432](https://github.com/cs3org/reva/pull/2432) :  Use space reference when listing containers
* Fix [cs3org/reva#2572](https://github.com/cs3org/reva/pull/2572) :  Wait for nats server on middleware start
* Fix [cs3org/reva#2454](https://github.com/cs3org/reva/pull/2454) :  Fix webdav paths in PROPFINDS
* Chg [cs3org/reva#2329](https://github.com/cs3org/reva/pull/2329) :  Activate the statcache
* Chg [cs3org/reva#2596](https://github.com/cs3org/reva/pull/2596) :  Remove hash from public link urls
* Chg [cs3org/reva#2495](https://github.com/cs3org/reva/pull/2495) :  Remove the ownCloud storage driver
* Chg [cs3org/reva#2527](https://github.com/cs3org/reva/pull/2527) :  Store space attributes in decomposedFS
* Chg [cs3org/reva#2581](https://github.com/cs3org/reva/pull/2581) :  Update hard-coded status values
* Chg [cs3org/reva#2524](https://github.com/cs3org/reva/pull/2524) :  Use description during space creation
* Chg [cs3org/reva#2554](https://github.com/cs3org/reva/pull/2554) :  Shard nodes per space in decomposedfs
* Chg [cs3org/reva#2576](https://github.com/cs3org/reva/pull/2576) :  Harden xattrs errors
* Chg [cs3org/reva#2436](https://github.com/cs3org/reva/pull/2436) :  Replace template in GroupFilter for UserProvider with a simple string
* Chg [cs3org/reva#2429](https://github.com/cs3org/reva/pull/2429) :  Make archiver id based
* Chg [cs3org/reva#2340](https://github.com/cs3org/reva/pull/2340) :  Allow multiple space configurations per provider
* Chg [cs3org/reva#2396](https://github.com/cs3org/reva/pull/2396) :  The ocdav handler is now spaces aware
* Chg [cs3org/reva#2349](https://github.com/cs3org/reva/pull/2349) :  Require `ListRecycle` when listing trashbin
* Chg [cs3org/reva#2353](https://github.com/cs3org/reva/pull/2353) :  Reduce log output
* Chg [cs3org/reva#2542](https://github.com/cs3org/reva/pull/2542) :  Do not encode webDAV ids to base64
* Chg [cs3org/reva#2519](https://github.com/cs3org/reva/pull/2519) :  Remove the auto creation of the .space folder
* Chg [cs3org/reva#2394](https://github.com/cs3org/reva/pull/2394) :  Remove logic from gateway
* Chg [cs3org/reva#2023](https://github.com/cs3org/reva/pull/2023) :  Add a sharestorageprovider
* Chg [cs3org/reva#2234](https://github.com/cs3org/reva/pull/2234) :  Add a spaces registry
* Chg [cs3org/reva#2339](https://github.com/cs3org/reva/pull/2339) :  Fix static registry regressions
* Chg [cs3org/reva#2370](https://github.com/cs3org/reva/pull/2370) :  Fix static registry regressions
* Chg [cs3org/reva#2354](https://github.com/cs3org/reva/pull/2354) :  Return not found when updating non existent space
* Chg [cs3org/reva#2589](https://github.com/cs3org/reva/pull/2589) :  Remove deprecated linter modules
* Chg [cs3org/reva#2016](https://github.com/cs3org/reva/pull/2016) :  Move wrapping and unwrapping of paths to the storage gateway
* Enh [cs3org/reva#2591](https://github.com/cs3org/reva/pull/2591) :  Set up App Locks with basic locks
* Enh [cs3org/reva#1209](https://github.com/cs3org/reva/pull/1209) :  Reva CephFS module v0.2.1
* Enh [cs3org/reva#2511](https://github.com/cs3org/reva/pull/2511) :  Error handling cleanup in decomposed FS
* Enh [cs3org/reva#2516](https://github.com/cs3org/reva/pull/2516) :  Cleaned up some code
* Enh [cs3org/reva#2512](https://github.com/cs3org/reva/pull/2512) :  Consolidate xattr setter and getter
* Enh [cs3org/reva#2341](https://github.com/cs3org/reva/pull/2341) :  Use CS3 permissions API
* Enh [cs3org/reva#2343](https://github.com/cs3org/reva/pull/2343) :  Allow multiple space type fileters on decomposedfs
* Enh [cs3org/reva#2460](https://github.com/cs3org/reva/pull/2460) :  Add locking support to decomposedfs
* Enh [cs3org/reva#2540](https://github.com/cs3org/reva/pull/2540) :  Refactored the xattrs package in the decomposedfs
* Enh [cs3org/reva#2463](https://github.com/cs3org/reva/pull/2463) :  Do not log whole nodes
* Enh [cs3org/reva#2350](https://github.com/cs3org/reva/pull/2350) :  Add file locking methods to the storage and filesystem interfaces
* Enh [cs3org/reva#2379](https://github.com/cs3org/reva/pull/2379) :  Add new file url of the app provider to the ocs capabilities
* Enh [cs3org/reva#2369](https://github.com/cs3org/reva/pull/2369) :  Implement TouchFile from the CS3apis
* Enh [cs3org/reva#2385](https://github.com/cs3org/reva/pull/2385) :  Allow to create new files with the app provider on public links
* Enh [cs3org/reva#2397](https://github.com/cs3org/reva/pull/2397) :  Product field in OCS version
* Enh [cs3org/reva#2393](https://github.com/cs3org/reva/pull/2393) :  Update tus/tusd to version 1.8.0
* Enh [cs3org/reva#2522](https://github.com/cs3org/reva/pull/2522) :  Introduce events
* Enh [cs3org/reva#2528](https://github.com/cs3org/reva/pull/2528) :  Use an exclcusive write lock when writing multiple attributes
* Enh [cs3org/reva#2595](https://github.com/cs3org/reva/pull/2595) :  Add integration test for the groupprovider
* Enh [cs3org/reva#2439](https://github.com/cs3org/reva/pull/2439) :  Ignore handled errors when creating spaces
* Enh [cs3org/reva#2500](https://github.com/cs3org/reva/pull/2500) :  Invalidate listproviders cache
* Enh [cs3org/reva#2345](https://github.com/cs3org/reva/pull/2345) :  Don't assume that the LDAP groupid in reva matches the name
* Enh [cs3org/reva#2525](https://github.com/cs3org/reva/pull/2525) :  Allow using AD UUID as userId values
* Enh [cs3org/reva#2584](https://github.com/cs3org/reva/pull/2584) :  Allow running userprovider integration tests for the LDAP driver
* Enh [cs3org/reva#2585](https://github.com/cs3org/reva/pull/2585) :  Add metadata storage layer and indexer
* Enh [cs3org/reva#2163](https://github.com/cs3org/reva/pull/2163) :  Nextcloud-based share manager for pkg/ocm/share
* Enh [cs3org/reva#2278](https://github.com/cs3org/reva/pull/2278) :  OIDC driver changes for lightweight users
* Enh [cs3org/reva#2315](https://github.com/cs3org/reva/pull/2315) :  Add new attributes to public link propfinds
* Enh [cs3org/reva#2431](https://github.com/cs3org/reva/pull/2431) :  Delete shares when purging spaces
* Enh [cs3org/reva#2434](https://github.com/cs3org/reva/pull/2434) :  Refactor ocdav into smaller chunks
* Enh [cs3org/reva#2524](https://github.com/cs3org/reva/pull/2524) :  Add checks when removing space members
* Enh [cs3org/reva#2457](https://github.com/cs3org/reva/pull/2457) :  Restore spaces that were previously deleted
* Enh [cs3org/reva#2498](https://github.com/cs3org/reva/pull/2498) :  Include grants in list storage spaces response
* Enh [cs3org/reva#2344](https://github.com/cs3org/reva/pull/2344) :  Allow listing all storage spaces
* Enh [cs3org/reva#2547](https://github.com/cs3org/reva/pull/2547) :  Add an if-match check to the storage provider
* Enh [cs3org/reva#2486](https://github.com/cs3org/reva/pull/2486) :  Update cs3apis to include lock api changes
* Enh [cs3org/reva#2526](https://github.com/cs3org/reva/pull/2526) :  Upgrade ginkgo to v2

https://github.com/owncloud/ocis/pull/3231
https://github.com/owncloud/ocis/pull/3258

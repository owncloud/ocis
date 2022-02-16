Enhancement: Update REVA to xxx

Updated REVA to xxx
This update includes:

* Enh [cs3org/reva#2524](https://github.com/cs3org/reva/pull/2524): Remove space members 
* Fix [cs3org/reva#2541](https://github.com/cs3org/reva/pull/2541): fix xattr error types, remove error wrapper 
* Chg [cs3org/reva#2540](https://github.com/cs3org/reva/pull/2540): decomposedfs: refactor xattrs package errors 
* Enh [cs3org/reva#2533](https://github.com/cs3org/reva/pull/2533): Use space description on creation 
* Enh [cs3org/reva#2527](https://github.com/cs3org/reva/pull/2527): Add space props 
* Enh [cs3org/reva#2522](https://github.com/cs3org/reva/pull/2522): Events 
* Chg [cs3org/reva#2512](https://github.com/cs3org/reva/pull/2512): Consolidate all metadata Get's and Set's to central functions. 
* Chg [cs3org/reva#2511](https://github.com/cs3org/reva/pull/2511): Some error cleanup steps in the decomposed FS 
* Enh [cs3org/reva#2460](https://github.com/cs3org/reva/pull/2460): decomposedfs: add locking support 
* Chg [cs3org/reva#2519](https://github.com/cs3org/reva/pull/2519): remove creation of .space folder 
* Fix [cs3org/reva#2506](https://github.com/cs3org/reva/pull/2506): fix propfind listing for files 
* Chg [cs3org/reva#2503](https://github.com/cs3org/reva/pull/2503): unprotected ocs config endpoint 
* Enh [cs3org/reva#2458](https://github.com/cs3org/reva/pull/2458): Restoring Spaces 
* Enh [cs3org/reva#2498](https://github.com/cs3org/reva/pull/2498): add grants to list-spaces 
* Fix [cs3org/reva#2500](https://github.com/cs3org/reva/pull/2500): invalidate cache when modifying or deleting a space 
* Fix [cs3org/reva#2501](https://github.com/cs3org/reva/pull/2501): fix spaces stat requests 
* Enh [cs3org/reva#2472](https://github.com/cs3org/reva/pull/2472): Make owncloudsql spaces aware 
* Enh [cs3org/reva#2464](https://github.com/cs3org/reva/pull/2464): Space grants 
* Fix [cs3org/reva#2463](https://github.com/cs3org/reva/pull/2463): Do not log nodes 
* Enh [cs3org/reva#2437](https://github.com/cs3org/reva/pull/2437): Make gateway dumb again 
* Enh [cs3org/reva#2459](https://github.com/cs3org/reva/pull/2459): prevent purging of enabled spaces 
* Fix [cs3org/reva#2457](https://github.com/cs3org/reva/pull/2457): decomposedfs: do not swallow errors when creating nodes 
* Fix [cs3org/reva#2454](https://github.com/cs3org/reva/pull/2454): fix path construction in webdav propfind 
* Fix [cs3org/reva#2452](https://github.com/cs3org/reva/pull/2452): fix create space error message 
* Enh [cs3org/reva#2431](https://github.com/cs3org/reva/pull/2431): Purge spaces 
* Fix [cs3org/reva#2445](https://github.com/cs3org/reva/pull/2445): Fix publiclinks and decomposedfs 
* Chg [cs3org/reva#2439](https://github.com/cs3org/reva/pull/2439): ignore handled errors when creating spaces 
* Enh [cs3org/reva#2436](https://github.com/cs3org/reva/pull/2436): Adjust "groupfilter" to be able to search by member name 
* Fix [cs3org/reva#2434](https://github.com/cs3org/reva/pull/2434): Start splitting up ocdav 
* Fix [cs3org/reva#2433](https://github.com/cs3org/reva/pull/2433): fix shares provider filter 
* Chg [cs3org/reva#2432](https://github.com/cs3org/reva/pull/2432): use space reference when listing containers 
* Fix [cs3org/reva#2430](https://github.com/cs3org/reva/pull/2430): fix aggregated child folder id 
* Enh [cs3org/reva#2429](https://github.com/cs3org/reva/pull/2429): make archiver id based 
* Fix [cs3org/reva#2427](https://github.com/cs3org/reva/pull/2427): fix registry caching 
* Fix [cs3org/reva#2422](https://github.com/cs3org/reva/pull/2422): handle space does not exist 
* Fix [cs3org/reva#2419](https://github.com/cs3org/reva/pull/2419): Spaces fixes 
* Chg [cs3org/reva#2415](https://github.com/cs3org/reva/pull/2415): services should never return transport level errors 
* Chg [cs3org/reva#2396](https://github.com/cs3org/reva/pull/2396): Ocdav spaces aware 
* Fix [cs3org/reva#2348](https://github.com/cs3org/reva/pull/2348): fix-archiver 
* Chg [cs3org/reva#2344](https://github.com/cs3org/reva/pull/2344): allow listing all storage spaces 
* Chg [cs3org/reva#2345](https://github.com/cs3org/reva/pull/2345): Switch LDAP test to use entryUUID as unique id for groups 
* Chg [cs3org/reva#2343](https://github.com/cs3org/reva/pull/2343): allow multiple space type filters on decomposedfs 
* Enh [cs3org/reva#2329](https://github.com/cs3org/reva/pull/2329): Activate Statcache 
* Enh [cs3org/reva#2340](https://github.com/cs3org/reva/pull/2340): Space registry multiple spaces per provider 
* Chg [cs3org/reva#2336](https://github.com/cs3org/reva/pull/2336): handle sending all permissions when creating public links 
* Fix [cs3org/reva#2330](https://github.com/cs3org/reva/pull/2330): fix decomposedfs upload 
* Enh [cs3org/reva#2234](https://github.com/cs3org/reva/pull/2234): Spaces registry 
* Enh [cs3org/reva#2217](https://github.com/cs3org/reva/pull/2217): New OIDC ESCAPE auth driver. 
* Enh [cs3org/reva#2250](https://github.com/cs3org/reva/pull/2250): Implement space membership endpoints 
* Fix [cs3org/reva#1941](https://github.com/cs3org/reva/pull/1941): fix tus with transfer token only 
* Fix [cs3org/reva#2309](https://github.com/cs3org/reva/pull/2309): Bugfix: Remove early finish for zero byte file uploads 
* Fix [cs3org/reva#2303](https://github.com/cs3org/reva/pull/2303): Fix content disposition 
* Fix [cs3org/reva#2314](https://github.com/cs3org/reva/pull/2314): OIDC: fallback to "email" if IDP doesn't provide "preferred_username" claim 
* Enh [cs3org/reva#2256](https://github.com/cs3org/reva/pull/2256): Return user type in the response of the ocs GET user call 
* Enh [cs3org/reva#2310](https://github.com/cs3org/reva/pull/2310): Implement setting arbitrary metadata for the public storage provider 
* Fix [cs3org/reva#2305](https://github.com/cs3org/reva/pull/2305): Make sure /app/new takes target as absolute path 
* Fix [cs3org/reva#2297](https://github.com/cs3org/reva/pull/2297): Fix public link paths for file shares 

https://github.com/owncloud/ocis/pull/2878
https://github.com/owncloud/ocis/pull/2901
https://github.com/owncloud/ocis/pull/2997
https://github.com/owncloud/ocis/pull/3116
https://github.com/owncloud/ocis/pull/3130
https://github.com/owncloud/ocis/pull/3175
https://github.com/owncloud/ocis/pull/3182

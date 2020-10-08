Enhancement: Add the storage service

Tags: storage, reva

* Enhancement - Enable ocis driver treetime accounting: [#620](https://github.com/owncloud/ocis/pull/620)
* Enhancement - Launch a storage to store ocis-metadata: [#602](https://github.com/owncloud/ocis/pull/602)

   In the future accounts, settings etc. should be stored in a dedicated metadata storage. The
   services should talk to this storage directly, bypassing reva-gateway.

   https://github.com/owncloud/ocis/pull/602

* Enhancement - Update reva to v1.2.2-0.20200924071957-e6676516e61e: [#601](https://github.com/owncloud/ocis/pull/601)

   - Update reva to v1.2.2-0.20200924071957-e6676516e61e - eos client: Handle eos EPERM as
   permission denied [(reva/#1183)](https://github.com/cs3org/reva/pull/1183) - ocis
   driver: synctime based etag propagation
   [(reva/#1180)](https://github.com/cs3org/reva/pull/1180) - ocis driver: fix litmus
   [(reva/#1179)](https://github.com/cs3org/reva/pull/1179) - ocis driver: fix move
   [(reva/#1177)](https://github.com/cs3org/reva/pull/1177) - ocs service: cache
   displaynames [(reva/#1161)](https://github.com/cs3org/reva/pull/1161)

   https://github.com/owncloud/ocis-reva/issues/262
   https://github.com/owncloud/ocis-reva/issues/357
   https://github.com/owncloud/ocis-reva/issues/301
   https://github.com/owncloud/ocis-reva/issues/302
   https://github.com/owncloud/ocis/pull/601

* Bugfix - Fix default configuration for accessing shares: [#205](https://github.com/owncloud/product/issues/205)

   The storage provider mounted at `/home` should always have EnableHome set to `true`. The other
   storage providers should have it set to `false`.

   https://github.com/owncloud/product/issues/205
   https://github.com/owncloud/ocis-reva/pull/461

* Enhancement - Allow configuring arbitrary storage registry rules: [#193](https://github.com/owncloud/product/issues/193)

   We added a new config flag `storage-registry-rule` that can be given multiple times for the
   gateway to specify arbitrary storage registry rules. You can also use a comma separated list of
   rules in the `REVA_STORAGE_REGISTRY_RULES` environment variable.

   https://github.com/owncloud/product/issues/193
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

   https://github.com/owncloud/ocis-reva/pull/454

* Enhancement - Update reva to v1.2.1-0.20200911111727-51649e37df2d: [#466](https://github.com/owncloud/ocis-reva/pull/466)

   - Update reva to v1.2.1-0.20200911111727-51649e37df2d - Added new OCIS storage driver ocis
   [(reva/#1155)](https://github.com/cs3org/reva/pull/1155) - App provider: fallback to
   env. variable if 'iopsecret' unset
   [(reva/#1146)](https://github.com/cs3org/reva/pull/1146) - Add switch to database
   [(reva/#1135)](https://github.com/cs3org/reva/pull/1135) - Add the ocdav HTTP svc to the
   standalone config [(reva/#1128)](https://github.com/cs3org/reva/pull/1128)

   https://github.com/owncloud/ocis-reva/pull/466

* Enhancement - Separate user and auth providers, add config for rest user: [#412](https://github.com/owncloud/ocis-reva/pull/412)

   Previously, the auth and user provider services used to have the same driver, which restricted
   using separate drivers and configs for both. This PR separates the two and adds the config for
   the rest user driver and the gatewaysvc parameter to EOS fs.

   https://github.com/owncloud/ocis-reva/pull/412
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

   https://github.com/owncloud/ocis-reva/pull/447

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

   https://github.com/owncloud/ocis-reva/pull/399
   https://github.com/cs3org/reva/pull/996

* Change - Environment updates for the username userid split: [#420](https://github.com/owncloud/ocis-reva/pull/420)

   We updated the owncloud storage driver in reva to properly look up users by userid or username
   using the userprovider instead of taking the path segment as is. This requires the user service
   address as well as changing the default layout to the userid instead of the username. The latter
   is not considered a stable and persistent identifier.

   https://github.com/owncloud/ocis-reva/pull/420
   https://github.com/cs3org/reva/pull/1033

* Enhancement - Update storage documentation: [#384](https://github.com/owncloud/ocis-reva/pull/384)

   We added details to the documentation about storage requirements known from ownCloud 10, the
   local storage driver and the ownCloud storage driver.

   https://github.com/owncloud/ocis-reva/pull/384
   https://github.com/owncloud/ocis-reva/pull/390

* Enhancement - Update reva to v0.1.1-0.20200724135750-b46288b375d6: [#399](https://github.com/owncloud/ocis-reva/pull/399)

   - Update reva to v0.1.1-0.20200724135750-b46288b375d6 - Split LDAP user filters
   (reva/#996) - meshdirectory: Add invite forward API to provider links (reva/#1000) - OCM:
   Pass the link to the meshdirectory service in token mail (reva/#1002) - Update
   github.com/go-ldap/ldap to v3 (reva/#1004)

   https://github.com/owncloud/ocis-reva/pull/399
   https://github.com/cs3org/reva/pull/996
   https://github.com/cs3org/reva/pull/1000
   https://github.com/cs3org/reva/pull/1002
   https://github.com/cs3org/reva/pull/1004

* Enhancement - Update reva to v0.1.1-0.20200728071211-c948977dd3a0: [#407](https://github.com/owncloud/ocis-reva/pull/407)

   - Update reva to v0.1.1-0.20200728071211-c948977dd3a0 - Use proper logging for ldap auth
   requests (reva/#1008) - Update github.com/eventials/go-tus to
   v0.0.0-20200718001131-45c7ec8f5d59 (reva/#1007) - Check if SMTP credentials are nil
   (reva/#1006)

   https://github.com/owncloud/ocis-reva/pull/407
   https://github.com/cs3org/reva/pull/1008
   https://github.com/cs3org/reva/pull/1007
   https://github.com/cs3org/reva/pull/1006

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#393](https://github.com/owncloud/ocis-reva/pull/393)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-reva/pull/393

* Enhancement - Update reva to v0.1.1-0.20200710143425-cf38a45220c5: [#371](https://github.com/owncloud/ocis-reva/pull/371)

   - Update reva to v0.1.1-0.20200710143425-cf38a45220c5 (#371) - Add wopi open (reva/#920) -
   Added a CS3API compliant data exporter to Mentix (reva/#955) - Read SMTP password from env if
   not set in config (reva/#953) - OCS share fix including file info after update (reva/#958) - Add
   flag to smtpclient for for unauthenticated SMTP (reva/#963)

   https://github.com/owncloud/ocis-reva/pull/371
   https://github.com/cs3org/reva/pull/920
   https://github.com/cs3org/reva/pull/953
   https://github.com/cs3org/reva/pull/955
   https://github.com/cs3org/reva/pull/958
   https://github.com/cs3org/reva/pull/963

* Enhancement - Update reva to v0.1.1-0.20200722125752-6dea7936f9d1: [#392](https://github.com/owncloud/ocis-reva/pull/392)

   - Update reva to v0.1.1-0.20200722125752-6dea7936f9d1 - Added signing key capability
   (reva/#986) - Add functionality to create webdav references for OCM shares (reva/#974) -
   Added a site locations exporter to Mentix (reva/#972) - Add option to config to allow requests
   to hosts with unverified certificates (reva/#969)

   https://github.com/owncloud/ocis-reva/pull/392
   https://github.com/cs3org/reva/pull/986
   https://github.com/cs3org/reva/pull/974
   https://github.com/cs3org/reva/pull/972
   https://github.com/cs3org/reva/pull/969

* Enhancement - Make frontend prefixes configurable: [#363](https://github.com/owncloud/ocis-reva/pull/363)

   We introduce three new environment variables and preconfigure them the following way:

   * `REVA_FRONTEND_DATAGATEWAY_PREFIX="data"`
   * `REVA_FRONTEND_OCDAV_PREFIX=""`
   * `REVA_FRONTEND_OCS_PREFIX="ocs"`

   This restores the reva defaults that were changed upstream.

   https://github.com/owncloud/ocis-reva/pull/363
   https://github.com/cs3org/reva/pull/936/files#diff-51bf4fb310f7362f5c4306581132fc3bR63


* Enhancement - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66: [#341](https://github.com/owncloud/ocis-reva/pull/341)

   - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66 (#341) - Added country information
   to Mentix (reva/#924) - Refactor metrics package to implement reader interface (reva/#934) -
   Fix OCS public link share update values logic (#252, #288, reva/#930)

   https://github.com/owncloud/ocis-reva/issues/252
   https://github.com/owncloud/ocis-reva/issues/288
   https://github.com/owncloud/ocis-reva/pull/341
   https://github.com/cs3org/reva/pull/924
   https://github.com/cs3org/reva/pull/934
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

   https://github.com/owncloud/ocis-reva/issues/344
   https://github.com/owncloud/ocis-reva/issues/336
   https://github.com/owncloud/ocis-reva/issues/11
   https://github.com/owncloud/ocis-reva/pull/362
   https://github.com/cs3org/reva/pull/936
   https://github.com/cs3org/reva/pull/937
   https://github.com/cs3org/reva/pull/938
   https://github.com/cs3org/reva/pull/940
   https://github.com/cs3org/reva/pull/951
   https://github.com/cs3org/reva/pull/945
   https://github.com/cs3org/reva/pull/929
   https://github.com/cs3org/reva/pull/946
   https://github.com/cs3org/reva/pull/944
   https://github.com/cs3org/reva/pull/903
   https://github.com/cs3org/reva/pull/954

* Enhancement - Add new config options for the http client: [#330](https://github.com/owncloud/ocis-reva/pull/330)

   The internal certificates are checked for validity after
   https://github.com/cs3org/reva/pull/914, which causes the acceptance tests to fail. This
   change sets new hardcoded defaults.

   https://github.com/owncloud/ocis-reva/pull/330

* Enhancement - Allow datagateway transfers to take 24h: [#323](https://github.com/owncloud/ocis-reva/pull/323)

   - Increase transfer token life time to 24h (PR #323)

   https://github.com/owncloud/ocis-reva/pull/323

* Enhancement - Update reva to v0.1.1-0.20200630075923-39a90d431566: [#320](https://github.com/owncloud/ocis-reva/pull/320)

   - Update reva to v0.1.1-0.20200630075923-39a90d431566 (#320) - Return special value for
   public link password (#294, reva/#904) - Fix public stat and listcontainer response to
   contain the correct prefix (#310, reva/#902)

   https://github.com/owncloud/ocis-reva/issues/310
   https://github.com/owncloud/ocis-reva/issues/294
   https://github.com/owncloud/ocis-reva/pull/320
   https://github.com/cs3org/reva/pull/902
   https://github.com/cs3org/reva/pull/904

* Enhancement - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66: [#328](https://github.com/owncloud/ocis-reva/pull/328)

   - Update reva to v0.1.1-0.20200701152626-2f6cc60e2f66 (#328) - Use sync.Map on pool package
   (reva/#909) - Use mutex instead of sync.Map (reva/#915) - Use gatewayProviders instead of
   storageProviders on conn pool (reva/#916) - Add logic to ls and stat to process arbitrary
   metadata keys (reva/#905) - Preliminary implementation of Set/UnsetArbitraryMetadata
   (reva/#912) - Make datagateway forward headers (reva/#913, reva/#926) - Add option to cmd
   upload to disable tus (reva/#911) - OCS Share Allow date-only expiration for public shares
   (#288, reva/#918) - OCS Share Remove array from OCS Share update response (#252, reva/#919) -
   OCS Share Implement GET request for single shares (#249, reva/#921)

   https://github.com/owncloud/ocis-reva/issues/288
   https://github.com/owncloud/ocis-reva/issues/252
   https://github.com/owncloud/ocis-reva/issues/249
   https://github.com/owncloud/ocis-reva/pull/328
   https://github.com/cs3org/reva/pull/909
   https://github.com/cs3org/reva/pull/915
   https://github.com/cs3org/reva/pull/916
   https://github.com/cs3org/reva/pull/905
   https://github.com/cs3org/reva/pull/912
   https://github.com/cs3org/reva/pull/913
   https://github.com/cs3org/reva/pull/926
   https://github.com/cs3org/reva/pull/911
   https://github.com/cs3org/reva/pull/918
   https://github.com/cs3org/reva/pull/919
   https://github.com/cs3org/reva/pull/921

* Enhancement - Update reva to v0.1.1-0.20200629131207-04298ea1c088: [#309](https://github.com/owncloud/ocis-reva/pull/309)

   - Update reva to v0.1.1-0.20200629094927-e33d65230abc (#309) - Fix public link file share
   (#278, reva/#895, reva/#900) - Delete public share (reva/#899) - Updated reva to
   v0.1.1-0.20200629131207-04298ea1c088 (#313)

   https://github.com/owncloud/ocis-reva/issues/278
   https://github.com/owncloud/ocis-reva/pull/309
   https://github.com/cs3org/reva/pull/895
   https://github.com/cs3org/reva/pull/899
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

   https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/issues/293
   https://github.com/owncloud/ocis-reva/issues/261
   https://github.com/owncloud/ocis-reva/pull/261
   https://github.com/cs3org/reva/pull/875
   https://github.com/cs3org/reva/pull/877
   https://github.com/cs3org/reva/pull/878
   https://github.com/cs3org/reva/pull/881
   https://github.com/cs3org/reva/pull/880
   https://github.com/cs3org/reva/pull/888
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

   https://github.com/owncloud/ocis-reva/issues/116
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

   https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/owncloud/ocis-reva/pull/228

* Enhancement - Update reva to v0.1.1-0.20200603071553-e05a87521618: [#244](https://github.com/owncloud/ocis-reva/issues/244)

   - Updated reva to v0.1.1-0.20200603071553-e05a87521618 (#244) - Add option to disable TUS on
   OC layer (#177, reva/#791) - Dataprovider now supports method override (#177, reva/#792) -
   OCS fixes for create public link (reva/#798)

   https://github.com/owncloud/ocis-reva/issues/244
   https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/cs3org/reva/pull/791
   https://github.com/cs3org/reva/pull/792
   https://github.com/cs3org/reva/pull/798

* Enhancement - Add public shares service: [#49](https://github.com/owncloud/ocis-reva/issues/49)

   Added Public Shares service with CRUD operations and File Public Shares Manager

   https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/pull/232

* Enhancement - Update reva to v0.1.1-0.20200529120551-4f2d9c85d3c9: [#49](https://github.com/owncloud/ocis-reva/issues/49)

   - Updated reva to v0.1.1-0.20200529120551 (#232) - Public Shares CRUD, File Public Shares
   Manager (#49, #232, reva/#681, reva/#788) - Disable HTTP-KeepAlives to reduce fd count
   (ocis/#268, reva/#787) - Fix trashbin listing (#229, reva/#782) - Create PUT wrapper for TUS
   uploads (reva/#770) - Add security access headers for ocdav requests (#66, reva/#780) - Add
   option to revad cmd to specify logging level (reva/#772) - New metrics package (reva/#740) -
   Remove implicit data member from memory store (reva/#774) - Added TUS global capabilities
   (#177, reva/#775) - Fix PROPFIND with Depth 1 for cross-storage operations (reva/#779)

   https://github.com/owncloud/ocis-reva/issues/49
   https://github.com/owncloud/ocis-reva/issues/229
   https://github.com/owncloud/ocis-reva/issues/66
   https://github.com/owncloud/ocis-reva/issues/177
   https://github.com/owncloud/ocis/issues/268
   https://github.com/owncloud/ocis-reva/pull/232
   https://github.com/cs3org/reva/pull/787
   https://github.com/cs3org/reva/pull/681
   https://github.com/cs3org/reva/pull/788
   https://github.com/cs3org/reva/pull/782
   https://github.com/cs3org/reva/pull/770
   https://github.com/cs3org/reva/pull/780
   https://github.com/cs3org/reva/pull/772
   https://github.com/cs3org/reva/pull/740
   https://github.com/cs3org/reva/pull/774
   https://github.com/cs3org/reva/pull/775
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

   https://github.com/owncloud/ocis-reva/issues/20
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
   https://github.com/cs3org/reva/pull/766
   https://github.com/cs3org/reva/pull/765
   https://github.com/cs3org/reva/pull/755
   https://github.com/cs3org/reva/pull/754
   https://github.com/cs3org/reva/pull/747
   https://github.com/cs3org/reva/pull/748
   https://github.com/cs3org/reva/pull/760
   https://github.com/cs3org/reva/pull/759
   https://github.com/cs3org/reva/pull/678
   https://github.com/cs3org/reva/pull/696
   https://github.com/cs3org/reva/pull/707
   https://github.com/cs3org/reva/pull/711
   https://github.com/cs3org/reva/pull/695
   https://github.com/cs3org/reva/pull/674
   https://github.com/cs3org/reva/pull/725
   https://github.com/cs3org/reva/pull/717
   https://github.com/cs3org/reva/pull/712
   https://github.com/cs3org/reva/pull/713
   https://github.com/cs3org/reva/pull/720
   https://github.com/cs3org/reva/pull/718
   https://github.com/cs3org/reva/pull/731
   https://github.com/cs3org/reva/pull/734
   https://github.com/cs3org/reva/pull/729
   https://github.com/cs3org/reva/pull/728
   https://github.com/cs3org/reva/pull/730
   https://github.com/cs3org/reva/pull/758
   https://github.com/cs3org/reva/pull/742
   https://github.com/cs3org/reva/pull/764
   https://github.com/cs3org/reva/pull/743
   https://github.com/cs3org/reva/pull/732
   https://github.com/cs3org/reva/pull/750

* Bugfix - Stop advertising unsupported chunking v2: [#145](https://github.com/owncloud/ocis-reva/pull/145)

   Removed "chunking" attribute in the DAV capabilities. Please note that chunking v2 is
   advertised as "chunking 1.0" while chunking v1 is the attribute "bigfilechunking" which is
   already false.

   https://github.com/owncloud/ocis-reva/pull/145

* Enhancement - Allow configuring the gateway for dataproviders: [#136](https://github.com/owncloud/ocis-reva/pull/136)

   This allows using basic or bearer auth when directly talking to dataproviders.

   https://github.com/owncloud/ocis-reva/pull/136

* Enhancement - Use a configured logger on reva runtime: [#153](https://github.com/owncloud/ocis-reva/pull/153)

   For consistency reasons we need a configured logger that is inline with an ocis logger, so the
   log cascade can be easily parsed by a human.

   https://github.com/owncloud/ocis-reva/pull/153

* Bugfix - Fix eos user sharing config: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   We have added missing config options for the user sharing manager and added a dedicated eos
   storage command with pre configured settings for the eos-docker container. It configures a
   `Shares` folder in a users home when using eos as the storage driver.

   https://github.com/owncloud/ocis-reva/pull/127

* Enhancement - Update reva to v1.1.0-20200414133413: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   Adds initial public sharing and ocm implementation.

   https://github.com/owncloud/ocis-reva/pull/127

* Bugfix - Fix eos config: [#125](https://github.com/owncloud/ocis-reva/pull/125)

   We have added missing config options for the home layout to the config struct that is passed to
   eos.

   https://github.com/owncloud/ocis-reva/pull/125

* Bugfix - Set correct flag type in the flagsets: [#75](https://github.com/owncloud/ocis-reva/issues/75)

   While upgrading to the micro/cli version 2 there where two instances of `StringFlag` which had
   not been changed to `StringSliceFlag`. This caused `ocis-reva users` and `ocis-reva
   storage-root` to fail on startup.

   https://github.com/owncloud/ocis-reva/issues/75
   https://github.com/owncloud/ocis-reva/pull/76

* Bugfix - We fixed a typo in the `REVA_LDAP_SCHEMA_MAIL` environment variable: [#113](https://github.com/owncloud/ocis-reva/pull/113)

   It was misspelled as `REVA_LDAP_SCHEMA_Mail`.

   https://github.com/owncloud/ocis-reva/pull/113

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

   https://github.com/owncloud/ocis-reva/pull/68
   related:

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

   https://github.com/owncloud/ocis-reva/pull/68

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-reva/issues/1)

   Just prepared an initial basic version to start a reva server and start integrating with the
   go-micro base dextension framework of ownCloud Infinite Scale.

   https://github.com/owncloud/ocis-reva/issues/1

* Change - Start multiple services with dedicated commands: [#6](https://github.com/owncloud/ocis-reva/issues/6)

   The initial version would only allow us to use a set of reva configurations to start multiple
   services. We use a more opinionated set of commands to start dedicated services that allows us
   to configure them individually. It allows us to switch eg. the user backend to LDAP and fully use
   it on the cli.

   https://github.com/owncloud/ocis-reva/issues/6

* Change - Storage providers now default to exposing data servers: [#89](https://github.com/owncloud/ocis-reva/issues/89)

   The flags that let reva storage providers announce that they expose a data server now defaults
   to true:

   `REVA_STORAGE_HOME_EXPOSE_DATA_SERVER=1` `REVA_STORAGE_OC_EXPOSE_DATA_SERVER=1`

   https://github.com/owncloud/ocis-reva/issues/89

* Change - Default to running behind ocis-proxy: [#113](https://github.com/owncloud/ocis-reva/pull/113)

   We changed the default configuration to integrate better with ocis.

   - We use ocis-glauth as the default ldap server on port 9125 with base `dc=example,dc=org`. - We
   use a dedicated technical `reva` user to make ldap binds - Clients are supposed to use the
   ocis-proxy endpoint `https://localhost:9200` - We removed unneeded ocis configuration
   from the frontend which no longer serves an oidc provider. - We changed the default user
   OpaqueID attribute from `sub` to `preferred_username`. The latter is a claim populated by
   konnectd that can also be used by the reva ldap user manager to look up users by their OpaqueId

   https://github.com/owncloud/ocis-reva/pull/113

* Enhancement - Expose owncloud storage driver config in flagset: [#87](https://github.com/owncloud/ocis-reva/issues/87)

   Three new flags are now available:

   - scan files on startup to generate missing fileids default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_SCAN` cli option: `--storage-owncloud-scan`

   - autocreate home path for new users default: `true` env var:
   `REVA_STORAGE_OWNCLOUD_AUTOCREATE` cli option: `--storage-owncloud-autocreate`

   - the address of the redis server default: `:6379` env var:
   `REVA_STORAGE_OWNCLOUD_REDIS_ADDR` cli option: `--storage-owncloud-redis`

   https://github.com/owncloud/ocis-reva/issues/87

* Enhancement - Update reva to v0.0.2-0.20200212114015-0dbce24f7e8b: [#91](https://github.com/owncloud/ocis-reva/pull/91)

   Reva has seen a lot of changes that allow us to - reduce the configuration overhead - use the
   autocreato home folder option - use the home folder path layout option - no longer start the root
   storage

   https://github.com/owncloud/ocis-reva/pull/91
   related:

* Enhancement - Allow configuring user sharing driver: [#115](https://github.com/owncloud/ocis-reva/pull/115)

   We now default to `json` which persists shares in the sharing manager in a json file instead of an
   in memory db.

   https://github.com/owncloud/ocis-reva/pull/115

https://github.com/owncloud/product/issues/244

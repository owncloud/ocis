# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-reva unreleased.

[unreleased]: https://github.com/owncloud/ocis-reva/compare/v0.5.0...master

## Summary

* Enhancement - Update reva to v0.1.1-0.20200624063447-db5e6635d5f0: [#279](https://github.com/owncloud/ocis-reva/pull/279)

## Details

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

# Changelog for [0.5.0] (2020-06-04)

The following sections list the changes in ocis-reva 0.5.0.

[0.5.0]: https://github.com/owncloud/ocis-reva/compare/v0.4.0...v0.5.0

## Summary

* Enhancement - Add TUS global capability: [#177](https://github.com/owncloud/ocis-reva/issues/177)
* Enhancement - Update reva to v0.1.1-0.20200603071553-e05a87521618: [#244](https://github.com/owncloud/ocis-reva/issues/244)

## Details

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

# Changelog for [0.4.0] (2020-05-29)

The following sections list the changes in ocis-reva 0.4.0.

[0.4.0]: https://github.com/owncloud/ocis-reva/compare/v0.3.0...v0.4.0

## Summary

* Enhancement - Add public shares service: [#49](https://github.com/owncloud/ocis-reva/issues/49)
* Enhancement - Update reva to v0.1.1-0.20200529120551-4f2d9c85d3c9: [#49](https://github.com/owncloud/ocis-reva/issues/49)

## Details

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

# Changelog for [0.3.0] (2020-05-20)

The following sections list the changes in ocis-reva 0.3.0.

[0.3.0]: https://github.com/owncloud/ocis-reva/compare/v0.2.1...v0.3.0

## Summary

* Enhancement - Update reva to v0.1.1-0.20200520150229: [#161](https://github.com/owncloud/ocis-reva/pull/161)

## Details

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

# Changelog for [0.2.1] (2020-04-28)

The following sections list the changes in ocis-reva 0.2.1.

[0.2.1]: https://github.com/owncloud/ocis-reva/compare/v0.2.0...v0.2.1

## Summary

* Bugfix - Stop advertising unsupported chunking v2: [#145](https://github.com/owncloud/ocis-reva/pull/145)
* Enhancement - Allow configuring the gateway for dataproviders: [#136](https://github.com/owncloud/ocis-reva/pull/136)
* Enhancement - Use a configured logger on reva runtime: [#153](https://github.com/owncloud/ocis-reva/pull/153)

## Details

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

# Changelog for [0.2.0] (2020-04-15)

The following sections list the changes in ocis-reva 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis-reva/compare/v0.1.1...v0.2.0

## Summary

* Bugfix - Fix eos user sharing config: [#127](https://github.com/owncloud/ocis-reva/pull/127)
* Enhancement - Update reva to v1.1.0-20200414133413: [#127](https://github.com/owncloud/ocis-reva/pull/127)

## Details

* Bugfix - Fix eos user sharing config: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   We have added missing config options for the user sharing manager and added a dedicated eos
   storage command with pre configured settings for the eos-docker container. It configures a
   `Shares` folder in a users home when using eos as the storage driver.

   https://github.com/owncloud/ocis-reva/pull/127


* Enhancement - Update reva to v1.1.0-20200414133413: [#127](https://github.com/owncloud/ocis-reva/pull/127)

   Adds initial public sharing and ocm implementation.

   https://github.com/owncloud/ocis-reva/pull/127

# Changelog for [0.1.1] (2020-03-31)

The following sections list the changes in ocis-reva 0.1.1.

[0.1.1]: https://github.com/owncloud/ocis-reva/compare/v0.1.0...v0.1.1

## Summary

* Bugfix - Fix eos config: [#125](https://github.com/owncloud/ocis-reva/pull/125)

## Details

* Bugfix - Fix eos config: [#125](https://github.com/owncloud/ocis-reva/pull/125)

   We have added missing config options for the home layout to the config struct that is passed to
   eos.

   https://github.com/owncloud/ocis-reva/pull/125

# Changelog for [0.1.0] (2020-03-23)

The following sections list the changes in ocis-reva 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-reva/compare/6702be7f9045a382d40691a9bcd04f572203e9ed...v0.1.0

## Summary

* Bugfix - Set correct flag type in the flagsets: [#75](https://github.com/owncloud/ocis-reva/issues/75)
* Bugfix - We fixed a typo in the `REVA_LDAP_SCHEMA_MAIL` environment variable: [#113](https://github.com/owncloud/ocis-reva/pull/113)
* Bugfix - Allow different namespaces for /webdav and /dav/files: [#68](https://github.com/owncloud/ocis-reva/pull/68)
* Change - Use /home as default namespace: [#68](https://github.com/owncloud/ocis-reva/pull/68)
* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-reva/issues/1)
* Change - Start multiple services with dedicated commands: [#6](https://github.com/owncloud/ocis-reva/issues/6)
* Change - Storage providers now default to exposing data servers: [#89](https://github.com/owncloud/ocis-reva/issues/89)
* Change - Default to running behind ocis-proxy: [#113](https://github.com/owncloud/ocis-reva/pull/113)
* Enhancement - Expose owncloud storage driver config in flagset: [#87](https://github.com/owncloud/ocis-reva/issues/87)
* Enhancement - Update reva to v0.0.2-0.20200212114015-0dbce24f7e8b: [#91](https://github.com/owncloud/ocis-reva/pull/91)
* Enhancement - Allow configuring user sharing driver: [#115](https://github.com/owncloud/ocis-reva/pull/115)

## Details

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


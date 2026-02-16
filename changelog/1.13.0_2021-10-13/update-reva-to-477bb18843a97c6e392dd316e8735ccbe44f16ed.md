Enhancement: Update reva to v1.14.0

This update includes:


 * Bugfix [cs3org/reva#2103](https://github.com/cs3org/reva/pull/2103): AppProvider: propagate back errors reported by WOPI
 * Bugfix [cs3org/reva#2149](https://github.com/cs3org/reva/pull/2149): Remove excess info from the http list app providers endpoint
 * Bugfix [cs3org/reva#2114](https://github.com/cs3org/reva/pull/2114): Add as default app while registering and skip unset mimetypes
 * Bugfix [cs3org/reva#2095](https://github.com/cs3org/reva/pull/2095): Fix app open when multiple app providers are present
 * Bugfix [cs3org/reva#2135](https://github.com/cs3org/reva/pull/2135): Make TUS capabilities configurable
 * Bugfix [cs3org/reva#2076](https://github.com/cs3org/reva/pull/2076): Fix chi routing
 * Bugfix [cs3org/reva#2077](https://github.com/cs3org/reva/pull/2077): Fix concurrent registration of mimetypes
 * Bugfix [cs3org/reva#2154](https://github.com/cs3org/reva/pull/2154): Return OK when trying to delete a non existing reference
 * Bugfix [cs3org/reva#2078](https://github.com/cs3org/reva/pull/2078): Fix nil pointer exception in stat
 * Bugfix [cs3org/reva#2073](https://github.com/cs3org/reva/pull/2073): Fix opening a readonly filetype with WOPI
 * Bugfix [cs3org/reva#2140](https://github.com/cs3org/reva/pull/2140): Map GRPC error codes to REVA errors
 * Bugfix [cs3org/reva#2147](https://github.com/cs3org/reva/pull/2147): Follow up of #2138: this is the new expected format
 * Bugfix [cs3org/reva#2116](https://github.com/cs3org/reva/pull/2116): Differentiate share types when retrieving received shares in sql driver
 * Bugfix [cs3org/reva#2074](https://github.com/cs3org/reva/pull/2074): Fix Stat() for EOS storage provider
 * Bugfix [cs3org/reva#2151](https://github.com/cs3org/reva/pull/2151): Fix return code for webdav uploads when the token expired
 * Change [cs3org/reva#2121](https://github.com/cs3org/reva/pull/2121): Sharemanager API change
 * Enhancement [cs3org/reva#2090](https://github.com/cs3org/reva/pull/2090): Return space name during list storage spaces
 * Enhancement [cs3org/reva#2138](https://github.com/cs3org/reva/pull/2138): Default AppProvider on top of the providers list
 * Enhancement [cs3org/reva#2137](https://github.com/cs3org/reva/pull/2137): Revamp app registry and add parameter to control file creation
 * Enhancement [cs3org/reva#145](https://github.com/cs3org/reva/pull/2137): UI improvements for the AppProviders
 * Enhancement [cs3org/reva#2088](https://github.com/cs3org/reva/pull/2088): Add archiver and app provider to ocs capabilities
 * Enhancement [cs3org/reva#2537](https://github.com/cs3org/reva/pull/2537): Add maximum files and size to archiver capabilities
 * Enhancement [cs3org/reva#2100](https://github.com/cs3org/reva/pull/2100): Add support for resource id to the archiver
 * Enhancement [cs3org/reva#2158](https://github.com/cs3org/reva/pull/2158): Augment the Id of new spaces
 * Enhancement [cs3org/reva#2085](https://github.com/cs3org/reva/pull/2085): Make encoding user groups in access tokens configurable
 * Enhancement [cs3org/reva#146](https://github.com/cs3org/reva/pull/146): Filter the denial shares (permission = 0) out of
 * Enhancement [cs3org/reva#2141](https://github.com/cs3org/reva/pull/2141): Use golang v1.17
 * Enhancement [cs3org/reva#2053](https://github.com/cs3org/reva/pull/2053): Safer defaults for TLS verification on LDAP connections
 * Enhancement [cs3org/reva#2115](https://github.com/cs3org/reva/pull/2115): Reduce code duplication in LDAP related drivers
 * Enhancement [cs3org/reva#1989](https://github.com/cs3org/reva/pull/1989): Add redirects from OC10 URL formats
 * Enhancement [cs3org/reva#2479](https://github.com/cs3org/reva/pull/2479): Limit publicshare and resourceinfo scope content
 * Enhancement [cs3org/reva#2071](https://github.com/cs3org/reva/pull/2071): Implement listing favorites via the dav report API
 * Enhancement [cs3org/reva#2091](https://github.com/cs3org/reva/pull/2091): Nextcloud share managers
 * Enhancement [cs3org/reva#2070](https://github.com/cs3org/reva/pull/2070): More unit tests for the Nextcloud storage provider
 * Enhancement [cs3org/reva#2087](https://github.com/cs3org/reva/pull/2087): More unit tests for the Nextcloud auth and user managers
 * Enhancement [cs3org/reva#2075](https://github.com/cs3org/reva/pull/2075): Make owncloudsql leverage existing filecache index
 * Enhancement [cs3org/reva#2050](https://github.com/cs3org/reva/pull/2050): Add a share types filter to the OCS API
 * Enhancement [cs3org/reva#2134](https://github.com/cs3org/reva/pull/2134): Use space Type from request
 * Enhancement [cs3org/reva#2132](https://github.com/cs3org/reva/pull/2132): Align local tests with drone setup
 * Enhancement [cs3org/reva#2095](https://github.com/cs3org/reva/pull/2095): Whitelisting for apps
 * Enhancement [cs3org/reva#2155](https://github.com/cs3org/reva/pull/2155): Pass an extra query parameter to WOPI /openinapp with a

https://github.com/owncloud/ocis/pull/2615
https://github.com/owncloud/ocis/pull/2566
https://github.com/owncloud/ocis/pull/2520

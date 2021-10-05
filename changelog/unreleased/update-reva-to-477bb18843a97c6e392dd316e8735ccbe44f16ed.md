Enhancement: Update reva to v1.13.1-0.20211001063718-477bb18843a9

This update includes:

* Bugfix [cs3org/reva#2076](https://github.com/cs3org/reva/pull/2076): Fix chi routing
* Bugfix [cs3org/reva#2077](https://github.com/cs3org/reva/pull/2077): Fix concurrent registration of mimetypes
* Bugfix [cs3org/reva#2074](https://github.com/cs3org/reva/pull/2074): Fix Stat() for eos storage provider
* Bugfix [cs3org/reva#2072](https://github.com/cs3org/reva/pull/2072): Fix denial shares being visible on Shared-with-me page
* Bugfix [cs3org/reva#2073](https://github.com/cs3org/reva/pull/2073): Fix opening a readonly filetype with WOPI
* Bugfix [cs3org/reva#2114](https://github.com/cs3org/reva/pull/2114): Fix apps as default while registering and skip unset mimetypes
* Security [cs3org/reva#2093](https://github.com/cs3org/reva/pull/2093): Limit the data exposed to resourceinfo and publicshare scopes
* Security [cs3org/reva#2053](https://github.com/cs3org/reva/pull/2053): Use safer defaults for TLS verification on LDAP connections
* Enhancement [cs3org/reva#1989](https://github.com/cs3org/reva/pull/1989): Implement url translation for legacy urls
* Enhancement [cs3org/reva#2075](https://github.com/cs3org/reva/pull/2075): Make owncloudsql leverage existing filecache index
* Enhancement [cs3org/reva#2090](https://github.com/cs3org/reva/pull/2090): Add space name during listStorageSpaces on decomposedfs
* Enhancement [cs3org/reva#2088](https://github.com/cs3org/reva/pull/2088): Add archiver and app provider capabilities
* Enhancement [cs3org/reva#2106](https://github.com/cs3org/reva/pull/2106): Add max num files and max size to the archiver capabilities
* Enhancement [cs3org/reva#2067](https://github.com/cs3org/reva/pull/2067): Extend AppRegistry and AppProvider
* Enhancement [cs3org/reva#2095](https://github.com/cs3org/reva/pull/2095): Whitelist apps via AppRegistry and AppProvider
* Enhancement [cs3org/reva#2115](https://github.com/cs3org/reva/pull/2115): Reduce code duplication in LDAP related drivers
* Enhancement [cs3org/reva#2100](https://github.com/cs3org/reva/pull/2100): Resource id based archiver for zip/tar downloads

https://github.com/owncloud/ocis/pull/2566

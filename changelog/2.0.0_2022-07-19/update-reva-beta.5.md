Enhancement: Update reva to version 2.7.2

Changelog for reva 2.7.2 (2022-07-18)
=======================================

* Bugfix [cs3org/reva#3079](https://github.com/cs3org/reva/pull/3079): Allow empty permissions
* Bugfix [cs3org/reva#3084](https://github.com/cs3org/reva/pull/3084): Spaces related permissions and providerID cleanup
* Bugfix [cs3org/reva#3083](https://github.com/cs3org/reva/pull/3083): Add space id to ItemTrashed event

Changelog for reva 2.7.1 (2022-07-15)
=======================================

* Bugfix [cs3org/reva#3080](https://github.com/cs3org/reva/pull/3080): Make dataproviders return more headers
* Enhancement [cs3org/reva#3046](https://github.com/cs3org/reva/pull/3046): Add user filter

Changelog for reva 2.7.0 (2022-07-15)
=======================================

* Bugfix [cs3org/reva#3075](https://github.com/cs3org/reva/pull/3075): Check permissions of the move operation destination
* Bugfix [cs3org/reva#3036](https://github.com/cs3org/reva/pull/3036): * Bugfix revad with EOS docker image
* Bugfix [cs3org/reva#3037](https://github.com/cs3org/reva/pull/3037): Add uid- and gidNumber to LDAP queries
* Bugfix [cs3org/reva#4061](https://github.com/cs3org/reva/pull/4061): Forbid resharing with higher permissions
* Bugfix [cs3org/reva#3017](https://github.com/cs3org/reva/pull/3017): Removed unused gateway config "commit_share_to_storage_ref"
* Bugfix [cs3org/reva#3031](https://github.com/cs3org/reva/pull/3031): Return proper response code when detecting recursive copy/move operations
* Bugfix [cs3org/reva#3071](https://github.com/cs3org/reva/pull/3071): Make CS3 sharing drivers parse legacy resource id
* Bugfix [cs3org/reva#3035](https://github.com/cs3org/reva/pull/3035): Prevent cross space move
* Bugfix [cs3org/reva#3074](https://github.com/cs3org/reva/pull/3074): Send storage provider and space id to wopi server
* Bugfix [cs3org/reva#3022](https://github.com/cs3org/reva/pull/3022): Improve the sharing internals
* Bugfix [cs3org/reva#2977](https://github.com/cs3org/reva/pull/2977): Test valid filename on spaces tus upload
* Change [cs3org/reva#3006](https://github.com/cs3org/reva/pull/3006): Use spaceID on the cs3api
* Enhancement [cs3org/reva#3043](https://github.com/cs3org/reva/pull/3043): Introduce LookupCtx for index interface
* Enhancement [cs3org/reva#3009](https://github.com/cs3org/reva/pull/3009): Prevent recursive copy/move operations
* Enhancement [cs3org/reva#2977](https://github.com/cs3org/reva/pull/2977): Skip space lookup on space propfind

https://github.com/owncloud/ocis/pull/4115
https://github.com/owncloud/ocis/pull/4201
https://github.com/owncloud/ocis/pull/4203
https://github.com/owncloud/ocis/pull/4025
https://github.com/owncloud/ocis/pull/4211

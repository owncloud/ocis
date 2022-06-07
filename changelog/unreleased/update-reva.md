Enhancement: update reva to version 2.5.0

Changelog for reva 2.5.0 (2022-06-07)
=======================================

The following sections list the changes in reva 2.5.0 relevant to
reva users. The changes are ordered by importance.

Summary
-------

* Bugfix [cs3org/reva#2909](https://github.com/cs3org/reva/pull/2909): The decomposedfs now checks the GetPath permission
* Bugfix [cs3org/reva#2899](https://github.com/cs3org/reva/pull/2899): Empty meta requests should return body
* Bugfix [cs3org/reva#2928](https://github.com/cs3org/reva/pull/2928): Fix mkcol response code
* Bugfix [cs3org/reva#2907](https://github.com/cs3org/reva/pull/2907): Correct share jail child aggregation
* Bugfix [cs3org/reva#3810](https://github.com/cs3org/reva/pull/3810): Fix unlimited quota in spaces
* Bugfix [cs3org/reva#3498](https://github.com/cs3org/reva/pull/3498): Check user permissions before updating/removing public shares
* Bugfix [cs3org/reva#2904](https://github.com/cs3org/reva/pull/2904): Share jail now works properly when accessed as a space
* Bugfix [cs3org/reva#2903](https://github.com/cs3org/reva/pull/2903): User owncloudsql now uses the correct userid
* Change [cs3org/reva#2920](https://github.com/cs3org/reva/pull/2920): Clean up the propfind code
* Change [cs3org/reva#2913](https://github.com/cs3org/reva/pull/2913): Rename ocs parameter "space_ref"
* Enhancement [cs3org/reva#2919](https://github.com/cs3org/reva/pull/2919): EOS Spaces implementation
* Enhancement [cs3org/reva#2888](https://github.com/cs3org/reva/pull/2888): Introduce spaces field mask
* Enhancement [cs3org/reva#2922](https://github.com/cs3org/reva/pull/2922): Refactor webdav error handling

https://github.com/owncloud/ocis/pull/3922
https://github.com/owncloud/ocis/pull/3928

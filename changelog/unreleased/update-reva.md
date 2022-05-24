Enhancement: update reva to version 2.4.0

Changelog for reva 2.4.0 (2022-05-24)
=======================================

The following sections list the changes in reva 2.4.0 relevant to
reva users. The changes are ordered by importance.

Summary
-------

* Bugfix [cs3org/reva#2854](https://github.com/cs3org/reva/pull/2854): Handle non uuid space and nodeid in decomposedfs
* Bugfix [cs3org/reva#2853](https://github.com/cs3org/reva/pull/2853): Filter CS3 share manager listing
* Bugfix [cs3org/reva#2868](https://github.com/cs3org/reva/pull/2868): Actually remove blobs when purging
* Bugfix [cs3org/reva#2882](https://github.com/cs3org/reva/pull/2882): Fix FileUploaded event being emitted too early
* Bugfix [cs3org/reva#2848](https://github.com/cs3org/reva/pull/2848): Fix storage id in the references in the ItemTrashed events
* Bugfix [cs3org/reva#2852](https://github.com/cs3org/reva/pull/2852): Fix rcbox dependency on reva 1.18
* Bugfix [cs3org/reva#3505](https://github.com/cs3org/reva/pull/3505): Fix creating a new file with wopi
* Bugfix [cs3org/reva#2885](https://github.com/cs3org/reva/pull/2885): Move stat out of usershareprovider
* Bugfix [cs3org/reva#2883](https://github.com/cs3org/reva/pull/2883): Fix role consideration when updating a share
* Bugfix [cs3org/reva#2864](https://github.com/cs3org/reva/pull/2864): Fix Grant Space IDs
* Bugfix [cs3org/reva#2870](https://github.com/cs3org/reva/pull/2870): Update quota calculation
* Bugfix [cs3org/reva#2876](https://github.com/cs3org/reva/pull/2876): Fix version number in status page
* Bugfix [cs3org/reva#2829](https://github.com/cs3org/reva/pull/2829): Don't include versions in quota
* Change [cs3org/reva#2856](https://github.com/cs3org/reva/pull/2856): Do not allow to edit disabled spaces
* Enhancement [cs3org/reva#3741](https://github.com/cs3org/reva/pull/3741): Add download endpoint to ocdav versions API
* Enhancement [cs3org/reva#2884](https://github.com/cs3org/reva/pull/2884): Show mounted shares in virtual share jail root
* Enhancement [cs3org/reva#2792](https://github.com/cs3org/reva/pull/2792): Use storageproviderid for spaces routing

https://github.com/owncloud/ocis/pull/3746
https://github.com/owncloud/ocis/pull/3771
https://github.com/owncloud/ocis/pull/3778
https://github.com/owncloud/ocis/pull/3842
https://github.com/owncloud/ocis/pull/3854
https://github.com/owncloud/ocis/pull/3858

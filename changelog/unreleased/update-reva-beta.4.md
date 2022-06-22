Enhancement: Update reva

Changelog for reva 2.6.0 (2022-06-21)
=======================================

The following sections list the changes in reva 2.6.0 relevant to
reva users. The changes are ordered by importance.

* Bugfix [cs3org/reva#2985](https://github.com/cs3org/reva/pull/2985): Make stat requests route based on storage providerid
* Bugfix [cs3org/reva#2987](https://github.com/cs3org/reva/pull/2987): Let archiver handle all error codes
* Bugfix [cs3org/reva#2994](https://github.com/cs3org/reva/pull/2994): Bugfix errors when loading shares
* Bugfix [cs3org/reva#2996](https://github.com/cs3org/reva/pull/2996): Do not close share dump channels
* Bugfix [cs3org/reva#2993](https://github.com/cs3org/reva/pull/2993): Remove unused configuration
* Bugfix [cs3org/reva#2950](https://github.com/cs3org/reva/pull/2950): Bugfix sharing with space ref
* Bugfix [cs3org/reva#2991](https://github.com/cs3org/reva/pull/2991): Make sharesstorageprovider get accepted share
* Change [cs3org/reva#2877](https://github.com/cs3org/reva/pull/2877): Enable resharing
* Change [cs3org/reva#2984](https://github.com/cs3org/reva/pull/2984): Update CS3Apis
* Enhancement [cs3org/reva#3753](https://github.com/cs3org/reva/pull/3753): Add executant to the events
* Enhancement [cs3org/reva#2820](https://github.com/cs3org/reva/pull/2820): Instrument GRPC and HTTP requests with OTel
* Enhancement [cs3org/reva#2975](https://github.com/cs3org/reva/pull/2975): Leverage shares space storageid and type when listing shares
* Enhancement [cs3org/reva#3882](https://github.com/cs3org/reva/pull/3882): Explicitly return on ocdav move requests with body
* Enhancement [cs3org/reva#2932](https://github.com/cs3org/reva/pull/2932): Stat accepted shares mountpoints, configure existing share updates
* Enhancement [cs3org/reva#2944](https://github.com/cs3org/reva/pull/2944): Improve owncloudsql connection management
* Enhancement [cs3org/reva#2962](https://github.com/cs3org/reva/pull/2962): Per service TracerProvider
* Enhancement [cs3org/reva#2911](https://github.com/cs3org/reva/pull/2911): Allow for dumping and loading shares
* Enhancement [cs3org/reva#2938](https://github.com/cs3org/reva/pull/2938): Sharpen tooling

https://github.com/owncloud/ocis/pull/3944
https://github.com/owncloud/ocis/pull/3975
https://github.com/owncloud/ocis/pull/3982
https://github.com/owncloud/ocis/pull/4000
https://github.com/owncloud/ocis/pull/4006

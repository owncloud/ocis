# Changelog for [0.1.1] (2020-08-17)

The following sections list the changes in ocis-store 0.1.1.

[0.1.1]: https://github.com/owncloud/ocis-store/compare/v0.1.0...v0.1.1

## Summary

* Bugfix - Removed code from other service: [#7](https://github.com/owncloud/ocis-store/pull/7)

## Details

* Bugfix - Removed code from other service: [#7](https://github.com/owncloud/ocis-store/pull/7)

   We had code copied over from another repository which doesn't belong in here and now removed it
   again.

   https://github.com/owncloud/ocis-store/pull/7

# Changelog for [0.1.0] (2020-07-23)

The following sections list the changes in ocis-store 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-store/compare/a57983ac33e65bcd73702eb1d74f56d36a94ef6c...v0.1.0

## Summary

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#5](https://github.com/owncloud/ocis-store/pull/5)
* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-store/pull/1)

## Details

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#5](https://github.com/owncloud/ocis-store/pull/5)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-store/pull/5


* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-store/pull/1)

   We have built a new service which implements go micro's [store
   interface](https://github.com/micro/development/blob/master/design/framework/store.md),
   so that we can provide a key-value-store through GRPC for our services. The underlying
   implementation stores data as json files on disk and maintains an index using
   [bleve](https://github.com/blevesearch/bleve).

   https://github.com/owncloud/ocis-store/pull/1
   https://github.com/owncloud/ocis-store/pull/2

-}}

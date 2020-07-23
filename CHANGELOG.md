# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-thumbnails unreleased.

[unreleased]: https://github.com/owncloud/ocis-thumbnails/compare/v0.2.0...master

## Summary

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#35](https://github.com/owncloud/ocis-thumbnails/pull/35)

## Details

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#35](https://github.com/owncloud/ocis-thumbnails/pull/35)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-thumbnails/pull/35

# Changelog for [0.2.0] (2020-07-10)

The following sections list the changes in ocis-thumbnails 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis-thumbnails/compare/v0.1.2...v0.2.0

## Summary

* Change - Add more default resolutions: [#23](https://github.com/owncloud/ocis-thumbnails/issues/23)
* Change - Refactor code to remove code smells: [#21](https://github.com/owncloud/ocis-thumbnails/issues/21)
* Change - Use micro service error api: [#31](https://github.com/owncloud/ocis-thumbnails/issues/31)
* Enhancement - Limit users to access own thumbnails: [#5](https://github.com/owncloud/ocis-thumbnails/issues/5)

## Details

* Change - Add more default resolutions: [#23](https://github.com/owncloud/ocis-thumbnails/issues/23)

   The `ocis-thumbnails` service was also used by the mediaviewer but the returned images were
   too small since the highest resolution was 128x128 pixels.

   https://github.com/owncloud/ocis-thumbnails/issues/23


* Change - Refactor code to remove code smells: [#21](https://github.com/owncloud/ocis-thumbnails/issues/21)

   Scanning the code using a static code analyzer showed some code smells. This change fixes them.

   https://github.com/owncloud/ocis-thumbnails/issues/21


* Change - Use micro service error api: [#31](https://github.com/owncloud/ocis-thumbnails/issues/31)

   The service now also returns a status code when an error occurs.

   https://github.com/owncloud/ocis-thumbnails/issues/31


* Enhancement - Limit users to access own thumbnails: [#5](https://github.com/owncloud/ocis-thumbnails/issues/5)

   Users of the service can no longer request thumbnails of another users by guessing the etag. The
   thumbnails are now only accessible by the users who created the thumbnail.

   https://github.com/owncloud/ocis-thumbnails/issues/5

# Changelog for [0.1.2] (2020-05-11)

The following sections list the changes in ocis-thumbnails 0.1.2.

[0.1.2]: https://github.com/owncloud/ocis-thumbnails/compare/v0.1.1...v0.1.2

## Summary

* Bugfix - Fix usage of context.Context: [#18](https://github.com/owncloud/ocis-thumbnails/issues/18)

## Details

* Bugfix - Fix usage of context.Context: [#18](https://github.com/owncloud/ocis-thumbnails/issues/18)

   The context was filled with a key defined in the package service but read with a key from the
   package imgsource. Since `service.key` and `imgsource.key` are different types imgsource
   could not read the value provided by service.

   https://github.com/owncloud/ocis-thumbnails/issues/18

# Changelog for [0.1.1] (2020-04-21)

The following sections list the changes in ocis-thumbnails 0.1.1.

[0.1.1]: https://github.com/owncloud/ocis-thumbnails/compare/v0.1.0...v0.1.1

## Summary

* Bugfix - Fix execution when passing program flags: [#15](https://github.com/owncloud/ocis-thumbnails/issues/15)

## Details

* Bugfix - Fix execution when passing program flags: [#15](https://github.com/owncloud/ocis-thumbnails/issues/15)

   The program flags where not correctly recognized because we didn't pass them to the micro
   framework when initializing a grpc service.

   https://github.com/owncloud/ocis-thumbnails/issues/15

# Changelog for [0.1.0] (2020-03-31)

The following sections list the changes in ocis-thumbnails 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-thumbnails/compare/c43f3a33cb0b57d7e25ebc88c138d22e95f88cfe...v0.1.0

## Summary

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-thumbnails/issues/1)
* Change - Use predefined resolutions for thumbnail generation: [#7](https://github.com/owncloud/ocis-thumbnails/issues/7)
* Change - Implement the first working version: [#3](https://github.com/owncloud/ocis-thumbnails/pull/3)

## Details

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-thumbnails/issues/1)

   Just prepare an initial basic version to embed a thumbnailer into our microservice
   infrastructure in the scope of the ownCloud Infinite Scale project.

   https://github.com/owncloud/ocis-thumbnails/issues/1


* Change - Use predefined resolutions for thumbnail generation: [#7](https://github.com/owncloud/ocis-thumbnails/issues/7)

   We implemented predefined resolutions to prevent attacker from flooding the service with a
   large number of thumbnails. The requested resolution gets mapped to the closest matching
   predefined resolution.

   https://github.com/owncloud/ocis-thumbnails/issues/7


* Change - Implement the first working version: [#3](https://github.com/owncloud/ocis-thumbnails/pull/3)

   We implemented the first simple version. It can load images via webdav and store them locally in
   the filesystem.

   https://github.com/owncloud/ocis-thumbnails/pull/3


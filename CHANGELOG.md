# Changelog for [2.1.0] (2020-03-16)

The following sections list the changes in ocis-pkg 2.1.0.

[2.1.0]: https://github.com/owncloud/ocis-pkg/compare/v2.0.2...v2.1.0

## Summary

* Enhancement - Tracing middleware: [#35](https://github.com/owncloud/ocis-pkg/pull/35/)

## Details

* Enhancement - Tracing middleware: [#35](https://github.com/owncloud/ocis-pkg/pull/35/)

   A new tracing middleware has been added to unpack context propagation

   https://github.com/owncloud/ocis-pkg/pull/35/

# Changelog for [2.0.2] (2020-03-04)

The following sections list the changes in ocis-pkg 2.0.2.

[2.0.2]: https://github.com/owncloud/ocis-pkg/compare/v2.0.1...v2.0.2

## Summary

* Enhancement - Allow http services to register handlers: [#33](https://github.com/owncloud/ocis-pkg/pull/33)

## Details

* Enhancement - Allow http services to register handlers: [#33](https://github.com/owncloud/ocis-pkg/pull/33)

   Added a handler option on http services

   https://github.com/owncloud/ocis-pkg/pull/33

# Changelog for [2.0.1] (2020-02-05)

The following sections list the changes in ocis-pkg 2.0.1.

[2.0.1]: https://github.com/owncloud/ocis-pkg/compare/v2.0.0...v2.0.1

## Summary

* Bugfix - Fix Module Path: [#25](https://github.com/owncloud/ocis-pkg/pull/25)
* Bugfix - Change import paths to ocis-pkg/v2: [#27](https://github.com/owncloud/ocis-pkg/pull/27)

## Details

* Bugfix - Fix Module Path: [#25](https://github.com/owncloud/ocis-pkg/pull/25)

   The module version must be in the path. See
   https://github.com/golang/go/wiki/Modules#semantic-import-versioning for more
   information. > If the module is version v2 or higher, the major version of the module must be
   included as a /vN at the end of the module paths used in go.mod files (e.g., module
   github.com/my/mod/v2, require github.com/my/mod/v2 v2.0.1) and in the package import path
   (e.g., import "github.com/my/mod/v2/mypkg"). This includes the paths used in go get
   commands (e.g., go get github.com/my/mod/v2@v2.0.1. Note there is both a /v2 and a @v2.0.1 in
   that example. One way to think about it is that the module name now includes the /v2, so include
   /v2 whenever you are using the module name).

   https://github.com/owncloud/ocis-pkg/pull/25


* Bugfix - Change import paths to ocis-pkg/v2: [#27](https://github.com/owncloud/ocis-pkg/pull/27)

   Changed the import paths to the current version

   https://github.com/owncloud/ocis-pkg/pull/27

# Changelog for [2.0.0] (2020-02-04)

The following sections list the changes in ocis-pkg 2.0.0.

[2.0.0]: https://github.com/owncloud/ocis-pkg/compare/v1.3.0...v2.0.0

## Summary

* Change - Upgrade the micro libraries: [#22](https://github.com/owncloud/ocis-pkg/pull/22)

## Details

* Change - Upgrade the micro libraries: [#22](https://github.com/owncloud/ocis-pkg/pull/22)

   Upgraded the go-micro libraries to v2.

   https://github.com/owncloud/ocis-pkg/pull/22

# Changelog for [1.3.0] (2020-01-20)

The following sections list the changes in ocis-pkg 1.3.0.

[1.3.0]: https://github.com/owncloud/ocis-pkg/compare/v1.2.0...v1.3.0

## Summary

* Bugfix - Fix serving static assets: [#14](https://github.com/owncloud/ocis-pkg/pull/14)
* Change - Add TLS support for http services: [#19](https://github.com/owncloud/ocis-pkg/issues/19)
* Enhancement - Introduce OpenID Connect middleware: [#8](https://github.com/owncloud/ocis-pkg/issues/8)

## Details

* Bugfix - Fix serving static assets: [#14](https://github.com/owncloud/ocis-pkg/pull/14)

   Ocis-hello used "/" as root. adding another / caused the static middleware to always fail
   stripping that prefix. All requests will return 404. Setting root to `""` in the `ocis-hello`
   flag does not work because Chi reports that routes need to start with `/`.
   `path.Clean(root+"/")` would yield `""` for `root="/"`.

   https://github.com/owncloud/ocis-pkg/pull/14


* Change - Add TLS support for http services: [#19](https://github.com/owncloud/ocis-pkg/issues/19)

   `ocis-pkg` http services support TLS. The idea behind is setting the issuer on phoenix's
   `config.json` to `https`. Or in other words, use https to access the Kopano extension, and
   authenticate using an SSL certificate.

   https://github.com/owncloud/ocis-pkg/issues/19


* Enhancement - Introduce OpenID Connect middleware: [#8](https://github.com/owncloud/ocis-pkg/issues/8)

   Added an openid connect middleware that will try to authenticate users using OpenID Connect.
   The claims will be added to the context of the request.

   https://github.com/owncloud/ocis-pkg/issues/8

# Changelog for [1.2.0] (2019-12-09)

The following sections list the changes in ocis-pkg 1.2.0.

[1.2.0]: https://github.com/owncloud/ocis-pkg/compare/v1.1.0...v1.2.0

## Summary

* Change - Add root path to static middleware: [#9](https://github.com/owncloud/ocis-pkg/issues/9)

## Details

* Change - Add root path to static middleware: [#9](https://github.com/owncloud/ocis-pkg/issues/9)

   Currently the `Static` middleware always serves from the root path, but all our HTTP handlers
   accept a custom root path which also got to be applied to the static file handling.

   https://github.com/owncloud/ocis-pkg/issues/9

# Changelog for [1.1.0] (2019-12-06)

The following sections list the changes in ocis-pkg 1.1.0.

[1.1.0]: https://github.com/owncloud/ocis-pkg/compare/v1.0.0...v1.1.0

## Summary

* Change - Better log level handling within micro: [#2](https://github.com/owncloud/ocis-pkg/issues/2)

## Details

* Change - Better log level handling within micro: [#2](https://github.com/owncloud/ocis-pkg/issues/2)

   Currently every log message from the micro internals are logged with the info level, we really
   need to respect the proper defined log level within our log wrapper package.

   https://github.com/owncloud/ocis-pkg/issues/2

# Changelog for [1.0.0] (2019-12-05)

The following sections list the changes in ocis-pkg 1.0.0.

[1.0.0]: https://github.com/owncloud/ocis-pkg/compare/63fa90a673cbc3238a503ea5e6304f1db7fdf47b...v1.0.0

## Summary

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-pkg/issues/1)

## Details

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-pkg/issues/1)

   Just prepared an initial basic version to have some shared functionality published which can
   be used by all other ownCloud Infinite Scale extensions.

   https://github.com/owncloud/ocis-pkg/issues/1


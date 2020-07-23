# Changes in 0.1.0

## Summary

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#20](https://github.com/owncloud/ocis-ocs/pull/20)
* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-ocs/issues/1)
* Change - Upgrade micro libraries: [#11](https://github.com/owncloud/ocis-ocs/issues/11)
* Enhancement - Configuration: [#14](https://github.com/owncloud/ocis-ocs/pull/14)
* Enhancement - Support signing key: [#18](https://github.com/owncloud/ocis-ocs/pull/18)

## Details

* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#20](https://github.com/owncloud/ocis-ocs/pull/20)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-ocs/pull/20


* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-ocs/issues/1)

   Just prepared an initial basic version to serve OCS for the ownCloud Infinite Scale project. It
   just provides a minimal viable product to demonstrate the microservice pattern.

   https://github.com/owncloud/ocis-ocs/issues/1


* Change - Upgrade micro libraries: [#11](https://github.com/owncloud/ocis-ocs/issues/11)

   Updated the micro and ocis-pkg libraries to version 2.

   https://github.com/owncloud/ocis-ocs/issues/11


* Enhancement - Configuration: [#14](https://github.com/owncloud/ocis-ocs/pull/14)

   Extensions should be responsible of configuring themselves. We use Viper for config loading
   from default paths. Environment variables **WILL** take precedence over config files.

   https://github.com/owncloud/ocis-ocs/pull/14


* Enhancement - Support signing key: [#18](https://github.com/owncloud/ocis-ocs/pull/18)

   We added support for the `/v[12].php/cloud/user/signing-key` endpoint that is used by the
   owncloud-sdk to generate signed URLs. This allows directly downloading large files with
   browsers instead of using `blob://` urls, which eats memory ...

   https://github.com/owncloud/ocis-ocs/pull/18
   https://github.com/owncloud/ocis-proxy/pull/75
   https://github.com/owncloud/owncloud-sdk/pull/504


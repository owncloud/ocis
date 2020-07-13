# Changes in unreleased

## Summary

* Change - Initial release of basic version: [#1](https://github.com/owncloud/ocis-ocs/issues/1)
* Change - Upgrade micro libraries: [#11](https://github.com/owncloud/ocis-ocs/issues/11)
* Enhancement - Configuration: [#14](https://github.com/owncloud/ocis-ocs/pull/14)

## Details

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


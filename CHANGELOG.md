# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-ocs unreleased.

[unreleased]: https://github.com/owncloud/ocis-ocs/compare/v0.2.0...master

## Summary

* Bugfix - Mimic oc10 user enabled as string in provisioning api: [#39](https://github.com/owncloud/ocis-ocs/pull/39)
* Bugfix - Use opaque ID of a user for signing keys: [#436](https://github.com/owncloud/ocis/issues/436)
* Enhancement - Add option to create user with uidnumber and gidnumber: [#34](https://github.com/owncloud/ocis-ocs/pull/34)

## Details

* Bugfix - Mimic oc10 user enabled as string in provisioning api: [#39](https://github.com/owncloud/ocis-ocs/pull/39)

   The oc10 user provisioning API uses a string for the boolean `enabled` flag. 😭

   https://github.com/owncloud/ocis-ocs/pull/39


* Bugfix - Use opaque ID of a user for signing keys: [#436](https://github.com/owncloud/ocis/issues/436)

   OCIS switched from user the user's opaque ID (UUID) everywhere, so to keep compatible we have
   adjusted the signing keys endpoint to also use the UUID when storing and generating the keys.

   https://github.com/owncloud/ocis/issues/436
   https://github.com/owncloud/ocis-ocs/pull/32


* Enhancement - Add option to create user with uidnumber and gidnumber: [#34](https://github.com/owncloud/ocis-ocs/pull/34)

   We have added an option to pass uidnumber and gidnumber to the ocis api while creating a new user

   https://github.com/owncloud/ocis-ocs/pull/34

# Changelog for [0.2.0] (2020-08-17)

The following sections list the changes in ocis-ocs 0.2.0.

[0.2.0]: https://github.com/owncloud/ocis-ocs/compare/v0.1.0...v0.2.0

## Summary

* Bugfix - Fix file descriptor leak: [#79](https://github.com/owncloud/ocis-accounts/issues/79)
* Enhancement - Add Group management for OCS Povisioning API: [#25](https://github.com/owncloud/ocis-ocs/pull/25)
* Enhancement - Basic Support for the User Provisioning API: [#23](https://github.com/owncloud/ocis-ocs/pull/23)

## Details

* Bugfix - Fix file descriptor leak: [#79](https://github.com/owncloud/ocis-accounts/issues/79)

   Only use a single instance of go-micro's GRPC client as it already does connection pooling.
   This prevents connection and file descriptor leaks.

   https://github.com/owncloud/ocis-accounts/issues/79
   https://github.com/owncloud/ocis-ocs/pull/29


* Enhancement - Add Group management for OCS Povisioning API: [#25](https://github.com/owncloud/ocis-ocs/pull/25)

   We added support for the group management related set of API calls of the provisioning API.
   [Reference](https://doc.owncloud.com/server/admin_manual/configuration/user/user_provisioning_api.html)

   https://github.com/owncloud/ocis-ocs/pull/25


* Enhancement - Basic Support for the User Provisioning API: [#23](https://github.com/owncloud/ocis-ocs/pull/23)

   We added support for a basic set of API calls for the user provisioning API.
   [Reference](https://doc.owncloud.com/server/admin_manual/configuration/user/user_provisioning_api.html)

   https://github.com/owncloud/ocis-ocs/pull/23

# Changelog for [0.1.0] (2020-07-23)

The following sections list the changes in ocis-ocs 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-ocs/compare/acd6d6e7f59d1a44bcedb4dd60564910b474c38a...v0.1.0

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


# Changelog for [unreleased] (UNRELEASED)

The following sections list the changes in ocis-settings unreleased.

[unreleased]: https://github.com/owncloud/ocis-settings/compare/v0.1.0...master

## Summary

* Change - Add role service: [#110](https://github.com/owncloud/product/issues/110)
* Change - Rename endpoints and message types: [#36](https://github.com/owncloud/ocis-settings/issues/36)
* Change - Use UUIDs instead of alphanumeric identifiers: [#46](https://github.com/owncloud/ocis-settings/pull/46)

## Details

* Change - Add role service: [#110](https://github.com/owncloud/product/issues/110)

   We added service endpoints for registering roles and maintaining permissions.

   https://github.com/owncloud/product/issues/110
   https://github.com/owncloud/ocis-settings/issues/10
   https://github.com/owncloud/ocis-settings/pull/47


* Change - Rename endpoints and message types: [#36](https://github.com/owncloud/ocis-settings/issues/36)

   We decided to rename endpoints and message types to be less verbose. Specifically,
   `SettingsBundle` became `Bundle`, `Setting` (inside a bundle) kept its name and
   `SettingsValue` became `Value`.

   https://github.com/owncloud/ocis-settings/issues/36
   https://github.com/owncloud/ocis-settings/issues/32
   https://github.com/owncloud/ocis-settings/pull/46


* Change - Use UUIDs instead of alphanumeric identifiers: [#46](https://github.com/owncloud/ocis-settings/pull/46)

   `Bundles`, `Settings` and `Values` were identified by a set of alphanumeric identifiers so
   far. We switched to UUIDs in order to achieve a flat file hierarchy on disk. Referencing the
   respective entities by their alphanumeric identifiers (as used in UI code) is still
   supported.

   https://github.com/owncloud/ocis-settings/pull/46

# Changelog for [0.1.0] (2020-08-17)

The following sections list the changes in ocis-settings 0.1.0.

[0.1.0]: https://github.com/owncloud/ocis-settings/compare/6fdbbd7e8eb39f18ada1a8a3c45a1c925e239889...v0.1.0

## Summary

* Bugfix - Adjust UUID validation to be more tolerant: [#41](https://github.com/owncloud/ocis-settings/issues/41)
* Bugfix - Fix runtime error when type asserting on nil value: [#38](https://github.com/owncloud/ocis-settings/pull/38)
* Bugfix - Fix multiple submits on string and number form elements: [#745](https://github.com/owncloud/owncloud-design-system/issues/745)
* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#39](https://github.com/owncloud/ocis-settings/pull/39)
* Change - Dynamically add navItems for extensions with settings bundles: [#25](https://github.com/owncloud/ocis-settings/pull/25)
* Change - Introduce input validation: [#22](https://github.com/owncloud/ocis-settings/pull/22)
* Change - Use account uuid from x-access-token: [#14](https://github.com/owncloud/ocis-settings/pull/14)
* Change - Use server config variable from ocis-web: [#34](https://github.com/owncloud/ocis-settings/pull/34)
* Enhancement - Remove paths from Makefile: [#33](https://github.com/owncloud/ocis-settings/pull/33)
* Enhancement - Extend the docs: [#11](https://github.com/owncloud/ocis-settings/issues/11)
* Enhancement - Update ocis-pkg/v2: [#42](https://github.com/owncloud/ocis-settings/pull/42)

## Details

* Bugfix - Adjust UUID validation to be more tolerant: [#41](https://github.com/owncloud/ocis-settings/issues/41)

   The UUID now allows any alphanumeric character and "-", "_", ".", "+" and "@" which can also
   allow regular user names.

   https://github.com/owncloud/ocis-settings/issues/41


* Bugfix - Fix runtime error when type asserting on nil value: [#38](https://github.com/owncloud/ocis-settings/pull/38)

   Fixed the case where an account UUID present in the context is nil, and type asserting it as a
   string would produce a runtime error.

   https://github.com/owncloud/ocis-settings/issues/37
   https://github.com/owncloud/ocis-settings/pull/38


* Bugfix - Fix multiple submits on string and number form elements: [#745](https://github.com/owncloud/owncloud-design-system/issues/745)

   We had a bug with keyboard event listeners triggering multiple submits on input fields. This
   was recently fixed in the ownCloud design system (ODS). We rolled out that bugfix to the
   settings ui as well.

   https://github.com/owncloud/owncloud-design-system/issues/745
   https://github.com/owncloud/owncloud-design-system/pull/768
   https://github.com/owncloud/ocis-settings/pulls/31


* Bugfix - Build docker images with alpine:latest instead of alpine:edge: [#39](https://github.com/owncloud/ocis-settings/pull/39)

   ARM builds were failing when built on alpine:edge, so we switched to alpine:latest instead.

   https://github.com/owncloud/ocis-settings/pull/39


* Change - Dynamically add navItems for extensions with settings bundles: [#25](https://github.com/owncloud/ocis-settings/pull/25)

   We now make use of a new feature in ocis-web-core, allowing us to add navItems not only through
   configuration, but also after app initialization. With this we now have navItems available
   for all extensions within the settings ui, that have at least one settings bundle registered.

   https://github.com/owncloud/ocis-settings/pull/25


* Change - Introduce input validation: [#22](https://github.com/owncloud/ocis-settings/pull/22)

   We set up input validation, starting with enforcing alphanumeric identifier values and UUID
   format on account uuids. As a result, traversal into parent folders is not possible anymore. We
   also made sure that get and list requests are side effect free, i.e. not creating any folders.

   https://github.com/owncloud/ocis-settings/issues/15
   https://github.com/owncloud/ocis-settings/issues/16
   https://github.com/owncloud/ocis-settings/issues/19
   https://github.com/owncloud/ocis-settings/pull/22


* Change - Use account uuid from x-access-token: [#14](https://github.com/owncloud/ocis-settings/pull/14)

   We are now using an ocis-pkg middleware for extracting the account uuid of the authenticated
   user from the `x-access-token` of the http request header and inject it into the Identifier
   protobuf messages wherever possible. This allows us to use `me` instead of an actual account
   uuid, when the request comes through the proxy.

   https://github.com/owncloud/ocis-settings/pull/14


* Change - Use server config variable from ocis-web: [#34](https://github.com/owncloud/ocis-settings/pull/34)

   We are not providing an api url anymore but use the server url from ocis-web config instead. This
   still - as before - requires that ocis-proxy is in place for routing requests to ocis-settings.

   https://github.com/owncloud/ocis-settings/pull/34


* Enhancement - Remove paths from Makefile: [#33](https://github.com/owncloud/ocis-settings/pull/33)

   We have a variable for the proto files path in our Makefile, but were not using it. Changed the
   Makefile to use the PROTO_SRC variable where possible.

   https://github.com/owncloud/ocis-settings/pull/33


* Enhancement - Extend the docs: [#11](https://github.com/owncloud/ocis-settings/issues/11)

   We have extended the documentation by adding a chapter about settings values.

   https://github.com/owncloud/ocis-settings/issues/11
   https://github.com/owncloud/ocis-settings/pulls/28


* Enhancement - Update ocis-pkg/v2: [#42](https://github.com/owncloud/ocis-settings/pull/42)

   Update to ocis-pkg/v2 v2.2.2-0.20200812103920-db41b5a3d14d

   https://github.com/owncloud/ocis-settings/pull/42


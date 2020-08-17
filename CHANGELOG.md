# Changes in 0.1.0

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


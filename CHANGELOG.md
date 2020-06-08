# Changes in unreleased

## Summary

* Bugfix - Fix multiple submits on string and number form elements: [#745](https://github.com/owncloud/owncloud-design-system/issues/745)
* Change - Introduce input validation: [#22](https://github.com/owncloud/ocis-settings/pull/22)
* Change - Use account uuid from x-access-token: [#14](https://github.com/owncloud/ocis-settings/pull/14)

## Details

* Bugfix - Fix multiple submits on string and number form elements: [#745](https://github.com/owncloud/owncloud-design-system/issues/745)

   We had a bug with keyboard event listeners triggering multiple submits on input fields. This
   was recently fixed in the ownCloud design system (ODS). We rolled out that bugfix to the
   settings ui as well.

   https://github.com/owncloud/owncloud-design-system/issues/745
   https://github.com/owncloud/owncloud-design-system/pull/768
   https://github.com/owncloud/ocis-settings/pulls/31


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


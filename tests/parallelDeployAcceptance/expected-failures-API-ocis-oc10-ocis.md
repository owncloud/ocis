## Scenarios from ownCloud10 core API tests that are expected to fail with parallel deployment matrix: ocis-oc10-ocis

The expected failures in this file are from features in the owncloud/core repo.

#### [Ability to return error messages in Webdav response bodies](https://github.com/owncloud/ocis/issues/1293)

- [apiAuthOcs/ocsPOSTAuth.feature:10](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPOSTAuth.feature#L10) Scenario: send POST requests to OCS endpoints as normal user with wrong password
- [apiAuthOcs/ocsPUTAuth.feature:10](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPUTAuth.feature#L10) Scenario: send PUT request to OCS endpoints as admin with wrong password

#### [various sharing settings cannot be set](https://github.com/owncloud/ocis/issues/1328)

- [apiCapabilities/capabilities.feature:8](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L8)
- [apiCapabilities/capabilities.feature:15](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L15)
- [apiCapabilities/capabilities.feature:41](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L41)
- [apiCapabilities/capabilities.feature:85](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L85)
- [apiCapabilities/capabilities.feature:100](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L100)
- [apiCapabilities/capabilities.feature:116](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L116)
- [apiCapabilities/capabilities.feature:127](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L127)
- [apiCapabilities/capabilities.feature:139](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L139)
- [apiCapabilities/capabilities.feature:282](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L282)
- [apiCapabilities/capabilities.feature:959](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiCapabilities/capabilities.feature#L959)
- [apiMain/caldav.feature:31](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/caldav.feature#L31)
- [apiMain/carddav.feature:31](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiMain/carddav.feature#L31)

#### [Support for favorites](https://github.com/owncloud/ocis/issues/1228)

- [apiFavorites/favorites.feature:115](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L115)
- [apiFavorites/favorites.feature:116](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiFavorites/favorites.feature#L116)

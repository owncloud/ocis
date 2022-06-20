## Scenarios from ownCloud10 core API tests that are expected to fail with parallel deployment matrix: oc10-ocis-oc10

The expected failures in this file are from features in the owncloud/core repo.

#### [Ability to return error messages in Webdav response bodies](https://github.com/owncloud/ocis/issues/1293)

- [apiAuthOcs/ocsPOSTAuth.feature:10](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPOSTAuth.feature#L10) Scenario: send POST requests to OCS endpoints as normal user with wrong password
- [apiAuthOcs/ocsPUTAuth.feature:10](https://github.com/owncloud/core/blob/master/tests/acceptance/features/apiAuthOcs/ocsPUTAuth.feature#L10) Scenario: send PUT request to OCS endpoints as admin with wrong password

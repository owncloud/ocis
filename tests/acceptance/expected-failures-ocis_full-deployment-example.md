## Scenarios that are expected to fail against the ocis_full deployment example

The expected failures in this file are from features in the owncloud/ocis repo.

These are differences specific to running the acceptance suite against the `ocis_full` docker-compose
deployment example (`.github/workflows/deployment-test.yml`), not against the from-source instance the
shared `expected-failures-localAPI-on-OCIS-storage.md` baseline is tuned for.

#### [ocis_full enables Collabora by default, which adds an extra "Secure Viewer" permission role - the role list has 8 entries instead of the vanilla 7](https://github.com/owncloud/ocis/blob/master/deployments/examples/ocis_full/collabora.yml)

- [apiGraph/roleManagementEndpoint.feature:10](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/roleManagementEndpoint.feature#L10)

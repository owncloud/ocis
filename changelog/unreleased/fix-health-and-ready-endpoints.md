Bugfix: Fix health and ready endpoints

We added new checks to the `/readyz` and `/healthz` endpoints to ensure that the services are ready and healthy.
This change ensures that the endpoints return the correct status codes, which is needed to stabilize the k8s deployments.

https://github.com/owncloud/ocis/pull/10163
https://github.com/owncloud/ocis/pull/10301
https://github.com/owncloud/ocis/pull/10302
https://github.com/owncloud/ocis/pull/10303
https://github.com/owncloud/ocis/pull/10308
https://github.com/owncloud/ocis/pull/10323
https://github.com/owncloud/ocis/pull/10163
https://github.com/owncloud/ocis/pull/10333
https://github.com/owncloud/ocis/issues/10316
https://github.com/owncloud/ocis/issues/10281

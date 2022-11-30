Bugfix: Remove static ocs user backend config

We've remove the `OCS_ACCOUNT_BACKEND_TYPE` configuration option.
It was intended to allow configuration of different user backends for the ocs service.
Right now the ocs service only has a "cs3" backend. Therefor it's a static entry and not configurable.

https://github.com/owncloud/ocis/pull/4077

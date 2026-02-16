Enhancement: disable the color logging in docker compose examples

Disabled the color logging in the example docker compose deployments.
Although colored logs are helpful during the development process they may be undesired in other situations like production deployments, where the logs aren't consumed by humans directly but instead by a log aggregator.

https://github.com/owncloud/ocis/issues/871
https://github.com/owncloud/ocis/pull/3935

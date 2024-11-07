Bugfix: fix gateway nats checks

We now only check if nats is available when the gateway actually uses it. Furthermore, we added a backoff for checking the readys endpoint.

https://github.com/owncloud/ocis/pull/10502

Bugfix: Fix panic when stopping the nats

The nats server itself runs signal handling that the Shutdown() call in the ocis code is redundant and led to a panic.

https://github.com/owncloud/ocis/pull/10363
https://github.com/owncloud/ocis/issues/10360

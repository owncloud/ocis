Bugfix: Panic when service fails to start

Tags: runtime

When attempting to run a service through the runtime that is currently running and fails to start, a race condition still redirect os Interrupt signals to a closed channel.

https://github.com/owncloud/ocis/pull/2252

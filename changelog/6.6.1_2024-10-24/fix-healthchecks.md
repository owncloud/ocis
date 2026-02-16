Bugfix: Fix healthchecks

We needed to replace 0.0.0.0 bind addresses by outbound IP addresses in the healthcheck routine.

https://github.com/owncloud/ocis/pull/10405

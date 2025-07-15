Bugfix: Fix logging levels

We've fixed the configuration of logging levels. Previously it was not possible
to configure a service with a more or less verbose log level then all other services
when running in the supervised / runtime mode `ocis server`.

For example `OCIS_LOG_LEVEL=error PROXY_LOG_LEVEL=debug ocis server` did not configure
error logging for all services except the proxy, which should be on debug logging. This is now fixed
and working properly.

Also we fixed the format of go-micro logs to always default to error level.
Previously this was only ensured in the supervised / runtime mode.

https://github.com/owncloud/ocis/pull/4102
https://github.com/owncloud/ocis/issues/4089

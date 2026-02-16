Enhancement: File Logging

When running supervised, support for configuring all logs to a single log file:
`OCIS_LOG_FILE=/Users/foo/bar/ocis.log MICRO_REGISTRY=etcd bin/ocis server`

Supports directing log from single extensions to a log file:
`PROXY_LOG_FILE=/Users/foo/bar/proxy.log MICRO_REGISTRY=etcd bin/ocis proxy`

https://github.com/owncloud/ocis/pull/1816

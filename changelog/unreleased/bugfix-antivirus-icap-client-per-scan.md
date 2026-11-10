Bugfix: Use a dedicated ICAP client per scan

The antivirus service now creates a dedicated ICAP client for every scan
instead of sharing a single client between its workers. This only affects
deployments using `ANTIVIRUS_SCANNER_TYPE=icap`.

https://github.com/owncloud/ocis/pull/13387

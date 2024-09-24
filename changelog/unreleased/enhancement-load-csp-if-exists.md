Enhancement: Load CSP configuration file if it exists

The Content Security Policy (CSP) configuration file is now loaded by default if it exists.
The configuration file looked for should be located at `$OCIS_BASE_DATA_PATH/proxy/csp.yaml`.
If the file does not exist, the default CSP configuration is used.

https://github.com/owncloud/ocis/pull/10139
https://github.com/owncloud/ocis/issues/10021

Enhancement: make frontend prefixes configurable

We introduce three new environment variables and preconfigure them the following way:

```
REVA_FRONTEND_DATAGATEWAY_PREFIX="data"
REVA_FRONTEND_OCDAV_PREFIX=""
REVA_FRONTEND_OCS_PREFIX="ocs"
```

This restores the reva defaults that were changed upstream.

https://github.com/owncloud/ocis/ocis-revapull/363
https://github.com/cs3org/reva/pull/936/files#diff-51bf4fb310f7362f5c4306581132fc3bR63

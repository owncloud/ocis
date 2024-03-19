Enhancement: Change Cors default settings

We have changed the default CORS settings to set `Access-Control-Allow-Origin` to the `OCIS_URL` if not explicitely set
and `Access-Control-Allow-Credentials` to `false` if not explicitely set.

https://github.com/owncloud/ocis/pull/8518
https://github.com/owncloud/ocis/issues/8514

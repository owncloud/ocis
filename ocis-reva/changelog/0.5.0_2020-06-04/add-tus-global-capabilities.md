Enhancement: Add TUS global capability

The TUS global capabilities from Reva are now exposed.

The advertised max chunk size can be configured using the "--upload-max-chunk-size" CLI switch or "REVA_FRONTEND_UPLOAD_MAX_CHUNK_SIZE" environment variable.
The advertised http method override can be configured using the "--upload-http-method-override" CLI switch or "REVA_FRONTEND_UPLOAD_HTTP_METHOD_OVERRIDE" environment variable.

https://github.com/owncloud/ocis/ocis-revaissues/177
https://github.com/owncloud/ocis/ocis-revapull/228

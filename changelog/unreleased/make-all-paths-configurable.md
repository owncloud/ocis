Change: Make all paths configurable and default to a common temp dir

Aligned all services to use a dir following`/var/tmp/ocis/<service>/...` by default. Also made some missing temp paths configurable via env vars and config flags.

https://github.com/owncloud/ocis/pull/1080

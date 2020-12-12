Change: Make all paths configurable and default to `/var/tmp/ocis/<service>/...`

Aligned all services to use a subdir of `/var/tmp/ocis/` by default. Also made some missing temp paths configurable via env vars and config flags.

https://github.com/owncloud/ocis/pulls/1080

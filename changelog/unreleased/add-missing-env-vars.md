Bugfix: add missing env vars to docker compose

Tags: docker

Without setting `REVA_FRONTEND_URL` and `REVA_DATAGATEWAY_URL` uploads would default to locahost and fail if `OCIS_DOMAIN` was used to run ocis on a remote host.

https://github.com/owncloud/ocis/pull/392

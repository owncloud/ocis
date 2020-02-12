Enhancement: expose owncloud storage driver config in flagset

Three new flags are now available:

- scan files on startup to generate missing fileids
  default: `true`
  env var: `REVA_STORAGE_OWNCLOUD_SCAN`
  cli option: `--storage-owncloud-scan`

- autocreate home path for new users
  default: `true`
  env var: `REVA_STORAGE_OWNCLOUD_AUTOCREATE`
  cli option: `--storage-owncloud-autocreate`

- the address of the redis server
  default: `:6379`
  env var: `REVA_STORAGE_OWNCLOUD_REDIS_ADDR`
  cli option: `--storage-owncloud-redis`

https://github.com/owncloud/ocis-reva/issues/87
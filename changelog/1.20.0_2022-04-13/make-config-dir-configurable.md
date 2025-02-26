Enhancement: make config dir configurable

We have added an `OCIS_CONFIG_DIR` environment variable the will take precedence over the default `/etc/ocis`, `~/.ocis` and `.config` locations. When it is set the default locations will be ignored and only the configuration files in that directory will be read.

https://github.com/owncloud/ocis/pull/3440

Change: Load configuration files just from one directory

We've changed the configuration file loading behavior and are now only loading
configuration files from ONE single directory. This directory can be set on
compile time or via an environment variable on startup (`OCIS_CONFIG_DIR`).

We are using following configuration default paths:

- Docker images: `/etc/ocis/`
- Binary releases: `$HOME/.ocis/config/`

https://github.com/owncloud/ocis/pull/3587

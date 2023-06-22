Bugfix: Thumbnailer can be disabled globally now

We added an env var `OCIS_DISABLE_PREVIEWS` to disable the thumbnailer globally.
The web-only local env-var `WEB_OPTION_DISABLE_PREVIEWS` will be deprecated.

https://github.com/owncloud/ocis/pull/6577
https://github.com/owncloud/ocis/issues/192
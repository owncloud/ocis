Change: Rename `APP_PROVIDER_BASIC_*` environment variables

We've renamed the `APP_PROVIDER_BASIC_*` to `APP_PROVIDER_*` since
the `_BASIC_` part is a copy and paste error. Now all app provider
environment variables are consistently starting with `APP_PROVIDER_*`.

https://github.com/owncloud/ocis/pull/2812
https://github.com/owncloud/ocis/pull/2811

Change: Remove accessDeniedHelpUrl from the config

We've removed the option accessDeniedHelpUrl from the config, since other clients weren't able to consume it.
In order to be accessible by other clients, not just Web, it should be configured via the theme.json file.

https://github.com/owncloud/ocis/pull/7970

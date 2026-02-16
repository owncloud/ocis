Change: remove privacyURL and imprintURL from the config

We've removed the option privacyURL and imprintURL from the config, since other clients weren't able to consume these.
In order to be accessible by other clients, not just Web, those should be configured via the theme.json file.

https://github.com/owncloud/ocis/pull/7938/

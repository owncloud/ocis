Bugfix: Removed outdated and unused dependency from idp package

We've removed the outdated and apparently unused dependency `cldr` from the `kpop` dependency inside the idp web ui. This resolves a security issue around an oudated `xmldom` package version, originating from said `kpop` library.

https://github.com/owncloud/ocis/issues/7957
https://github.com/owncloud/ocis/pull/7988

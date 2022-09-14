Bugfix: Remove unused transfer secret from app provider

We've fixed the startup of the app provider by removing the startup dependency
on a configured transfer secret, which was not used. This only happend if you
start the app provider without runtime (eg. `ocis app-provider server`) and didn't
have configured all oCIS secrets.

https://github.com/owncloud/ocis/pull/3798

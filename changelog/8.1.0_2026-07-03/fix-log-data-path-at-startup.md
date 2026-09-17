Enhancement: Log effective data and config paths at startup

oCIS now logs the effective data path and config path at startup so
operators can immediately verify that data is written to the expected
location. This helps catch misconfigured Docker volume mounts where
data silently falls back to an ephemeral container path instead of
the intended persistent mount.

https://github.com/owncloud/ocis/pull/12117
https://github.com/owncloud/ocis/issues/12044

Bugfix: Fix configuration validation for extensions' server commands

We've fixed the configuration validation for the extensions' server commands.
Before this fix error messages have occurred when trying to start individual services
without certain oCIS fullstack configuration values.

We now no longer do the common oCIS configuration validation for extensions' server
commands and now rely only on the extensions' validation function.

https://github.com/owncloud/ocis/pull/3911

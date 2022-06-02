Bugfix: Fix configuration validation for extensions' server commands

We've fixed the configuration validation for the extensions' server commands.
Before that fix error messages have occurred when started services which don't need
certain configuration values, that are needed for the oCIS fullstack command.

We now no longer do the common oCIS configuration validation for extensions' server
commands and now rely only on the extensions' validation function.


https://github.com/owncloud/ocis/pull/3911

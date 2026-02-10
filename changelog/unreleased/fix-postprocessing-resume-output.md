Bugfix: Fix postprocessing resume command --restart flag

The `--restart` / `-r` flag for `ocis postprocessing resume` was broken due to a flag
name mismatch (`retrigger` vs `restart`) and silently did nothing. This has been fixed
and the command now prints a confirmation message on success.

https://github.com/owncloud/ocis/issues/11692
https://github.com/owncloud/ocis/pull/12002

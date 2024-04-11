Bugfix: Fix restarting of postprocessing

When an upload is not found, the logic to restart postprocessing was bunked. Additionally we extended the upload sessions
command to be able to restart the uploads without using a second command.

NOTE: This also includes a breaking fix for the deprecated `ocis storage-users uploads list` command

https://github.com/owncloud/ocis/pull/8782

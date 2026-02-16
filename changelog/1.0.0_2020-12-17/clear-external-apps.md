Bugfix: Don't enforce empty external apps slice

Tags: web

The command for ocis-phoenix enforced an empty external apps configuration. This was removed, as it was blocking a new set of default external apps in ocis-phoenix.

https://github.com/owncloud/ocis/pull/473

Bugfix: Avoid superfluous GetPublicShare call when deleting space permissions

We fixed `DeletePermission` to recognise space permission IDs (prefixed with
`u:` or `g:`) by their format before making any gateway calls. Previously,
deleting a space member always triggered a `GetPublicShare` lookup that was
guaranteed to fail, producing a confusing error log.

https://github.com/owncloud/ocis/pull/12122
https://github.com/owncloud/ocis/issues/12012

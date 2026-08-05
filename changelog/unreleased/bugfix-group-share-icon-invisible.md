Bugfix: Fix invisible group icon in share collaborator list

The icon shown next to group entries in the share collaborator autocomplete
list used a color variable that no longer resolved to a valid value, making
the icon effectively invisible against the background. The icon now uses the
correct default text color variable so it renders properly again.

https://github.com/owncloud/ocis/pull/TODO

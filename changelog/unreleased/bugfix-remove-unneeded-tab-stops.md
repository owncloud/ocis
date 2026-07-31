Bugfix: Remove unneeded keyboard focus stops

Some icons and date columns in tables (spaces list, file list) could
unexpectedly receive keyboard focus even though they don't trigger any
action, adding unnecessary stops when tabbing through the page. The
contextual helper icon also carried a label that duplicated the label
already present on its surrounding button. These redundant focus stops
and duplicate labels have been removed.

https://github.com/owncloud/ocis/pull/12645

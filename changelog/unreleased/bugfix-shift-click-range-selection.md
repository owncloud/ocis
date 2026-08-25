Bugfix: Fix shift-click range selection in the files list

We have fixed an issue where shift-click did not select resources between
the anchor and the clicked one, and ignored ctrl/cmd. It now selects the
full range, replacing the current selection unless ctrl/cmd is held, and no
longer highlights the text of the affected rows. Checkboxes now also keep their
checked state in sync with the selection when a click does not change it.

https://github.com/owncloud/ocis/pull/12804

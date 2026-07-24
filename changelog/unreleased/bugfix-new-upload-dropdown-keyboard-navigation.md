Bugfix: Fix keyboard navigation in "New" and "Upload" dropdown menus

Opening the "New" or "Upload" dropdown menu in the Files app toolbar did
not move focus into the menu, so pressing the arrow keys had no effect
on the highlighted item, unlike the existing right-click context menu.
Both dropdowns now focus their first item on open, restoring arrow key
navigation.

https://github.com/owncloud/ocis/pull/12646

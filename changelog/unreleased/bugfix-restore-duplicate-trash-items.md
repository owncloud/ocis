Bugfix: Restore multiple trash items with the same name

Restoring several deleted files with the same name at once showed a conflict
dialog claiming the file already exists, even when it didn't. Choosing to
replace it then showed both a success and a failure notification, because
all files were restored to the same location at the same time and only one
of them could succeed.

Trash items that only collide with each other are now restored side by side
with unique names, without asking. When they collide with an existing file
and the user chooses to replace it, only one of them replaces the file and
the others are restored with unique names.

https://github.com/owncloud/ocis/pull/13238

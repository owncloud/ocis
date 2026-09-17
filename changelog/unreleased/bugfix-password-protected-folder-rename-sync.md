Bugfix: Sync the `.psec` file name when a password-protected folder is renamed

Renaming a password-protected folder did not update its `.psec` pointer file,
so the main file listing kept showing the old name. This happened both when
renaming the folder via the hidden `.PasswordProtectedFolders` location and
when renaming it from inside the folder-view modal.

Renaming the real folder now also renames its `.psec` counterpart. Inside the
folder-view modal, the framed public-link session notifies the parent window
at the moment of rename, so the pointer file stays in sync even on refresh or
another device.

https://github.com/owncloud/ocis/pull/12684

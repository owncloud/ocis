Bugfix: Show live OnlyOffice edits to read-only users

Read-only OnlyOffice sessions now use the edit discovery action so they join
the active collaborative session and see changes before the editor closes the
document. CheckFileInfo continues to advertise that these users cannot write.

https://github.com/owncloud/ocis/issues/12605
https://github.com/owncloud/ocis/pull/12617
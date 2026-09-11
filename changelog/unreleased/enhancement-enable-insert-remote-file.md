Enhancement: Enable EnableInsertRemoteFile and EnableInsertRemoteImage WOPI flags

Enable the EnableInsertRemoteFile and EnableInsertRemoteImage flags in the Collabora CheckFileInfo response. This activates the multimedia insertion and document comparison menus in Collabora Online via the UI_InsertFile and UI_InsertGraphic postMessages.

Requires Collabora Online >= 24.04.10 for multimedia insertion, >= 25.04.9.1 for document comparison. Also requires a companion frontend change in owncloud/web to handle the postMessages.

https://github.com/owncloud/ocis/pull/12192

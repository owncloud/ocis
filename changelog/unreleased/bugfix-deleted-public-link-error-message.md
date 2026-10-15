Bugfix: Show a clear error for deleted public links and files in app views

Opening a file via a public link in an app such as the preview, after the link
or the file had been deleted, showed the message "Unknown error". The server
answers with an empty 404, which the WebDAV client turns into a generic message,
and app views that failed while loading the folder passed only that message on.
These views now show "The resource could not be located, it may not exist
anymore.", as the other loading paths already did.

https://github.com/owncloud/ocis/issues/13138
https://github.com/owncloud/ocis/pull/13156

Bugfix: Only open supported file types in web viewer apps

Viewer apps in the web UI, like the PDF viewer, could be opened for any file
by entering its address directly, even when the app does not support that file
type. The file was then loaded and shown anyway.

Viewer apps now check that the file type is supported before loading the file
and show an error message otherwise.

https://github.com/owncloud/ocis/pull/13439

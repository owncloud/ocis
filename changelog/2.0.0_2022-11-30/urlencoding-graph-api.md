Bugfix: URL encode the webdav url in the graph API

Fixed the webdav URL in the drives responses. Without encoding the URL could be broken by files with spaces in the file name.

https://github.com/owncloud/ocis/pull/3597
https://github.com/owncloud/ocis/issues/3538

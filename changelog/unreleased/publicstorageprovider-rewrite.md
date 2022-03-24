Bugfix: replace public mountpoint fileid with grant fileid

We now show the same resoucre id for resources when accessing them via a public links as when using a logged in user. This allows the web ui to start a WOPI session with the correct resource id.

https://github.com/owncloud/ocis/pull/3349

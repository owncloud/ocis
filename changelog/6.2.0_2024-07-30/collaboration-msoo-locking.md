Enhancement: Add locking support for MS Office Online Server

We added support for the special kind of lock tokens that MS Office Online Server uses to lock files via the Wopi protocol.
It will only be active if you set the `COLLABORATION_APP_NAME` environment variable to `MicrosoftOfficeOnline`.

https://github.com/owncloud/ocis/pull/9685

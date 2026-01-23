Bugfix: Remove leading dot before checking disabled extension

We have fixed a bug where the leading dot was not removed before checking if an extension is disabled.
The original behavior would have caused the `COLLABORATION_WOPI_DISABLED_EXTENSIONS` config to be ignored.

https://github.com/owncloud/ocis/pull/11814

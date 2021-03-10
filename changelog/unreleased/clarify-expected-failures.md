Enhancement: clarify expected failures

Some features, while covered by the ownCloud 10 acceptance tests, will not be implmented for now:
- blacklisted / ignored files, because ocis does not need to blacklist `.htaccess` files
- `OC-LazyOps` support was [removed from the clients](https://github.com/owncloud/client/pull/8398). We are thinking about [a state machine for uploads to properly solve that scenario and also list the state of files in progress in the web ui](https://github.com/owncloud/ocis/issues/214).
The expected failures files now have a dedicated _Won't fix_ section for these items.

https://github.com/owncloud/ocis/pull/1790
https://github.com/owncloud/client/pull/8398
https://github.com/owncloud/ocis/issues/214
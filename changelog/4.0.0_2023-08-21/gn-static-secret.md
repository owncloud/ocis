Enhancement: Add static secret to gn endpoints

The global notifications POST and DELETE endpoints (used only for deprovision notifications at the moment) can now be called by adding a static secret to the header. Admins can still call this endpoint without knowing the secret

https://github.com/owncloud/ocis/pull/6946

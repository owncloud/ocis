Bugfix: Fix thumbnail generation when using different idp

The thumbnail service was relying on a konnectd specific field in the access token.
This logic was now replaced by a service parameter for the username.

https://github.com/owncloud/ocis/issues/1624
https://github.com/owncloud/ocis/pull/1628

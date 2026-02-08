Bugfix: Fix auth-app token validation performance with many tokens

The auth-app JSON manager performed bcrypt comparison against every stored
token including expired ones. With many accumulated impersonation tokens this
caused authentication to take tens of seconds. Expired tokens are now skipped
before bcrypt comparison, purged on service startup (persisted to disk), and
cleaned per-user during token generation. Additionally, read-only operations
use a read lock to allow concurrent access.

https://github.com/owncloud/ocis/issues/11692
https://github.com/owncloud/ocis/pull/11998

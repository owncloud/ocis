Bugfix: Fix runtime error when type asserting on nil value

Fixed the case where an account UUID present in the context is nil, and type asserting it as a string would produce a runtime error.

https://github.com/owncloud/ocis/settings/pull/38
https://github.com/owncloud/ocis/settings/issues/37

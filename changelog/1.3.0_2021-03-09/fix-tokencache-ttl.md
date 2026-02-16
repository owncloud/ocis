Bugfix: Fix the ttl of the authentication middleware cache 

The authentication cache ttl was multiplied with `time.Second` multiple times. This resulted in a ttl that was not intended.

https://github.com/owncloud/ocis/pull/1699

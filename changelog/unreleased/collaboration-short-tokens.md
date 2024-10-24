Bugfix: Generate short tokens to be used as access tokens for WOPI

Currently, the access tokens being used might be too long.
In particular, Microsoft Office Online complains about the URL (which contains the access token)
is too long and refuses to work. 

https://github.com/owncloud/ocis/pull/10391

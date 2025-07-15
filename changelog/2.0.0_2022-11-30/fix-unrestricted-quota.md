Bugfix: Fix unrestricted quota on the graphAPI

Unrestricted quota needs to show 0 on the API. It is not good for clients when the property is missing.

https://github.com/owncloud/ocis/pull/4363

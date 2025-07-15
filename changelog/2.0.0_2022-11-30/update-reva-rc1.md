Enhancement: Update reva to version 2.11.0

Changelog for reva 2.11.0 (2022-11-03)
=======================================

*   Bugfix  [cs3org/reva#3282](https://github.com/cs3org/reva/pull/3282):  Use Displayname in wopi apps
*   Bugfix  [cs3org/reva#3430](https://github.com/cs3org/reva/pull/3430):  Add missing error check in decomposedfs
*   Bugfix  [cs3org/reva#3298](https://github.com/cs3org/reva/pull/3298):  Make date only expiry dates valid for the whole day
*   Bugfix  [cs3org/reva#3394](https://github.com/cs3org/reva/pull/3394):  Avoid AppProvider panic
*   Bugfix  [cs3org/reva#3267](https://github.com/cs3org/reva/pull/3267):  Reduced default cache sizes for smaller memory footprint
*   Bugfix  [cs3org/reva#3338](https://github.com/cs3org/reva/pull/3338):  Fix malformed uid string in cache
*   Bugfix  [cs3org/reva#3255](https://github.com/cs3org/reva/pull/3255):  Properly escape oc:name in propfind response
*   Bugfix  [cs3org/reva#3324](https://github.com/cs3org/reva/pull/3324):  Correct base URL for download URL and href when listing file public links
*   Bugfix  [cs3org/reva#3278](https://github.com/cs3org/reva/pull/3278):  Fix public share view mode during app open
*   Bugfix  [cs3org/reva#3377](https://github.com/cs3org/reva/pull/3377):  Fix possible race conditions
*   Bugfix  [cs3org/reva#3274](https://github.com/cs3org/reva/pull/3274):  Fix "uploader" role permissions
*   Bugfix  [cs3org/reva#3241](https://github.com/cs3org/reva/pull/3241):  Fix uploading empty files into shares
*   Bugfix  [cs3org/reva#3251](https://github.com/cs3org/reva/pull/3251):  Make listing xattrs more robust
*   Bugfix  [cs3org/reva#3287](https://github.com/cs3org/reva/pull/3287):  Return OCS forbidden error when a share already exists
*   Bugfix  [cs3org/reva#3218](https://github.com/cs3org/reva/pull/3218):  Improve performance when listing received shares
*   Bugfix  [cs3org/reva#3251](https://github.com/cs3org/reva/pull/3251):  Lock source on move
*   Bugfix  [cs3org/reva#3238](https://github.com/cs3org/reva/pull/3238):  Return relative used quota amount as a percent value
*   Bugfix  [cs3org/reva#3279](https://github.com/cs3org/reva/pull/3279):  Polish OCS error responses
*   Bugfix  [cs3org/reva#3307](https://github.com/cs3org/reva/pull/3307):  Refresh lock in decomposedFS needs to overwrite
*   Bugfix  [cs3org/reva#3368](https://github.com/cs3org/reva/pull/3368):  Return 404 when no permission to space
*   Bugfix  [cs3org/reva#3341](https://github.com/cs3org/reva/pull/3341):  Validate s3ng downloads
*   Bugfix  [cs3org/reva#3284](https://github.com/cs3org/reva/pull/3284):  Prevent nil pointer when requesting user
*   Bugfix  [cs3org/reva#3257](https://github.com/cs3org/reva/pull/3257):  Fix wopi access to publicly shared files
*   Change  [cs3org/reva#3267](https://github.com/cs3org/reva/pull/3267):  Decomposedfs no longer stores the idp
*   Change  [cs3org/reva#3381](https://github.com/cs3org/reva/pull/3381):  Changed Name of the Shares Jail
*   Enhancement  [cs3org/reva#3381](https://github.com/cs3org/reva/pull/3381):  Add capability for sharing by role
*   Enhancement  [cs3org/reva#3320](https://github.com/cs3org/reva/pull/3320):  Add the parentID to the ocs and dav responses
*   Enhancement  [cs3org/reva#3239](https://github.com/cs3org/reva/pull/3239):  Add privatelink to PROPFIND response
*   Enhancement  [cs3org/reva#3340](https://github.com/cs3org/reva/pull/3340):  Add SpaceOwner to some event
*   Enhancement  [cs3org/reva#3252](https://github.com/cs3org/reva/pull/3252):  Add SpaceShared event
*   Enhancement  [cs3org/reva#3297](https://github.com/cs3org/reva/pull/3297):  Update dependencies
*   Enhancement  [cs3org/reva#3429](https://github.com/cs3org/reva/pull/3429):  Make max lock cycles configurable
*   Enhancement  [cs3org/reva#3011](https://github.com/cs3org/reva/pull/3011):  Expose capability to deny access in OCS API
*   Enhancement  [cs3org/reva#3224](https://github.com/cs3org/reva/pull/3224):  Make the jsoncs3 share manager cache ttl configurable
*   Enhancement  [cs3org/reva#3290](https://github.com/cs3org/reva/pull/3290):  Harden file system accesses
*   Enhancement  [cs3org/reva#3332](https://github.com/cs3org/reva/pull/3332):  Allow to enable TLS for grpc service
*   Enhancement  [cs3org/reva#3223](https://github.com/cs3org/reva/pull/3223):  Improve CreateShare grpc error reporting
*   Enhancement  [cs3org/reva#3376](https://github.com/cs3org/reva/pull/3376):  Improve logging
*   Enhancement  [cs3org/reva#3250](https://github.com/cs3org/reva/pull/3250):  Allow sharing the gateway caches
*   Enhancement  [cs3org/reva#3240](https://github.com/cs3org/reva/pull/3240):  We now only encode &, < and > in PROPFIND PCDATA
*   Enhancement  [cs3org/reva#3334](https://github.com/cs3org/reva/pull/3334):  Secure the nats connection with TLS
*   Enhancement  [cs3org/reva#3300](https://github.com/cs3org/reva/pull/3300):  Do not leak existence of resources
*   Enhancement  [cs3org/reva#3233](https://github.com/cs3org/reva/pull/3233):  Allow to override default broker for go-micro base ocdav service
*   Enhancement  [cs3org/reva#3258](https://github.com/cs3org/reva/pull/3258):  Allow ocdav to share the registry instance with other services
*   Enhancement  [cs3org/reva#3225](https://github.com/cs3org/reva/pull/3225):  Render file parent id for ocs shares
*   Enhancement  [cs3org/reva#3222](https://github.com/cs3org/reva/pull/3222):  Support Prefer: return=minimal in PROPFIND
*   Enhancement  [cs3org/reva#3395](https://github.com/cs3org/reva/pull/3395):  Reduce lock contention issues
*   Enhancement  [cs3org/reva#3286](https://github.com/cs3org/reva/pull/3286):  Make Refresh Lock operation WOPI compliant
*   Enhancement  [cs3org/reva#3229](https://github.com/cs3org/reva/pull/3229):  Request counting middleware
*   Enhancement  [cs3org/reva#3312](https://github.com/cs3org/reva/pull/3312):  Implemented new share filters
*   Enhancement  [cs3org/reva#3308](https://github.com/cs3org/reva/pull/3308):  Update the ttlcache library
*   Enhancement  [cs3org/reva#3291](https://github.com/cs3org/reva/pull/3291):  The wopi app driver supports more options

https://github.com/owncloud/ocis/pull/4588
https://github.com/owncloud/ocis/pull/4716
https://github.com/owncloud/ocis/pull/4719
https://github.com/owncloud/ocis/pull/4750
https://github.com/owncloud/ocis/pull/4833
https://github.com/owncloud/ocis/pull/4867
https://github.com/owncloud/ocis/pull/4903
https://github.com/owncloud/ocis/pull/4908
https://github.com/owncloud/ocis/pull/4915
https://github.com/owncloud/ocis/pull/4964

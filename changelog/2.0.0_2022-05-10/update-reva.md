Enhancement: update reva to v2.3.1

Updated reva to version 2.3.1. This update includes

* Bugfix [cs3org/reva#2827](https://github.com/cs3org/reva/pull/2827): Check permissions when deleting spaces
* Bugfix [cs3org/reva#2830](https://github.com/cs3org/reva/pull/2830): Correctly render response when accepting merged shares
* Bugfix [cs3org/reva#2831](https://github.com/cs3org/reva/pull/2831): Fix uploads to owncloudsql storage when no mtime is provided
* Enhancement [cs3org/reva#2833](https://github.com/cs3org/reva/pull/2833): Make status.php values configurable
* Enhancement [cs3org/reva#2832](https://github.com/cs3org/reva/pull/2832): Add version option for ocdav go-micro service

Updated reva to version 2.3.0. This update includes:

* Bugfix [cs3org/reva#2693](https://github.com/cs3org/reva/pull/2693): Support editnew actions from MS Office
* Bugfix [cs3org/reva#2588](https://github.com/cs3org/reva/pull/2588): Dockerfile.revad-ceph to use the right base image
* Bugfix [cs3org/reva#2499](https://github.com/cs3org/reva/pull/2499): Removed check DenyGrant in resource permission
* Bugfix [cs3org/reva#2285](https://github.com/cs3org/reva/pull/2285): Accept new userid idp format
* Bugfix [cs3org/reva#2802](https://github.com/cs3org/reva/pull/2802): Bugfix the resource id handling for space shares
* Bugfix [cs3org/reva#2800](https://github.com/cs3org/reva/pull/2800): Bugfix spaceid parsing in spaces trashbin API
* Bugfix [cs3org/reva#2608](https://github.com/cs3org/reva/pull/2608): Respect the tracing_service_name config variable
* Bugfix [cs3org/reva#2742](https://github.com/cs3org/reva/pull/2742): Use exact match in login filter
* Bugfix [cs3org/reva#2759](https://github.com/cs3org/reva/pull/2759): Made uid, gid claims parsing more robust in OIDC auth provider
* Bugfix [cs3org/reva#2788](https://github.com/cs3org/reva/pull/2788): Return the correct file IDs on public link resources
* Bugfix [cs3org/reva#2322](https://github.com/cs3org/reva/pull/2322): Use RFC3339 for parsing dates
* Bugfix [cs3org/reva#2784](https://github.com/cs3org/reva/pull/2784): Disable storageprovider cache for the share jail
* Bugfix [cs3org/reva#2555](https://github.com/cs3org/reva/pull/2555): Bugfix site accounts endpoints
* Bugfix [cs3org/reva#2675](https://github.com/cs3org/reva/pull/2675): Updates Makefile according to latest go standards
* Bugfix [cs3org/reva#2572](https://github.com/cs3org/reva/pull/2572): Wait for nats server on middleware start
* Change [cs3org/reva#2735](https://github.com/cs3org/reva/pull/2735): Avoid user enumeration
* Change [cs3org/reva#2737](https://github.com/cs3org/reva/pull/2737): Bump go-cs3api
* Change [cs3org/reva#2763](https://github.com/cs3org/reva/pull/2763): Change the oCIS and S3NG  storage driver blob store layout
* Change [cs3org/reva#2596](https://github.com/cs3org/reva/pull/2596): Remove hash from public link urls
* Change [cs3org/reva#2785](https://github.com/cs3org/reva/pull/2785): Implement workaround for chi.RegisterMethod
* Change [cs3org/reva#2559](https://github.com/cs3org/reva/pull/2559): Do not encode webDAV ids to base64
* Change [cs3org/reva#2740](https://github.com/cs3org/reva/pull/2740): Rename oc10 share manager driver
* Change [cs3org/reva#2561](https://github.com/cs3org/reva/pull/2561): Merge oidcmapping auth manager into oidc
* Enhancement [cs3org/reva#2698](https://github.com/cs3org/reva/pull/2698): Make capabilities endpoint public, authenticate users is present
* Enhancement [cs3org/reva#2515](https://github.com/cs3org/reva/pull/2515): Enabling tracing by default if not explicitly disabled
* Enhancement [cs3org/reva#2686](https://github.com/cs3org/reva/pull/2686): Features for favorites xattrs in EOS, cache for scope expansion
* Enhancement [cs3org/reva#2494](https://github.com/cs3org/reva/pull/2494): Use sys ACLs for file permissions
* Enhancement [cs3org/reva#2522](https://github.com/cs3org/reva/pull/2522): Introduce events
* Enhancement [cs3org/reva#2811](https://github.com/cs3org/reva/pull/2811): Add event for created directories
* Enhancement [cs3org/reva#2798](https://github.com/cs3org/reva/pull/2798): Add additional fields to events to enable search
* Enhancement [cs3org/reva#2790](https://github.com/cs3org/reva/pull/2790): Fake providerids so API stays stable after beta
* Enhancement [cs3org/reva#2685](https://github.com/cs3org/reva/pull/2685): Enable federated account access
* Enhancement [cs3org/reva#1787](https://github.com/cs3org/reva/pull/1787): Add support for HTTP TPC
* Enhancement [cs3org/reva#2799](https://github.com/cs3org/reva/pull/2799): Add flag to enable unrestriced listing of spaces
* Enhancement [cs3org/reva#2560](https://github.com/cs3org/reva/pull/2560): Mentix PromSD extensions
* Enhancement [cs3org/reva#2741](https://github.com/cs3org/reva/pull/2741): Meta path for user
* Enhancement [cs3org/reva#2613](https://github.com/cs3org/reva/pull/2613): Externalize custom mime types configuration for storage providers
* Enhancement [cs3org/reva#2163](https://github.com/cs3org/reva/pull/2163): Nextcloud-based share manager for pkg/ocm/share
* Enhancement [cs3org/reva#2696](https://github.com/cs3org/reva/pull/2696): Preferences driver refactor and cbox sql implementation
* Enhancement [cs3org/reva#2052](https://github.com/cs3org/reva/pull/2052): New CS3API datatx methods
* Enhancement [cs3org/reva#2743](https://github.com/cs3org/reva/pull/2743): Add capability for public link single file edit
* Enhancement [cs3org/reva#2738](https://github.com/cs3org/reva/pull/2738): Site accounts site-global settings
* Enhancement [cs3org/reva#2672](https://github.com/cs3org/reva/pull/2672): Further Site Accounts improvements
* Enhancement [cs3org/reva#2549](https://github.com/cs3org/reva/pull/2549): Site accounts improvements
* Enhancement [cs3org/reva#2795](https://github.com/cs3org/reva/pull/2795): Add feature flags "projects" and "share_jail" to spaces capability
* Enhancement [cs3org/reva#2514](https://github.com/cs3org/reva/pull/2514): Reuse ocs role objects in other drivers
* Enhancement [cs3org/reva#2781](https://github.com/cs3org/reva/pull/2781): In memory user provider
* Enhancement [cs3org/reva#2752](https://github.com/cs3org/reva/pull/2752): Refactor the rest user and group provider drivers

https://github.com/owncloud/ocis/pull/3552
https://github.com/owncloud/ocis/pull/3570
https://github.com/owncloud/ocis/pull/3601
https://github.com/owncloud/ocis/pull/3602
https://github.com/owncloud/ocis/pull/3605
https://github.com/owncloud/ocis/pull/3611
https://github.com/owncloud/ocis/issues/3621
https://github.com/owncloud/ocis/pull/3637
https://github.com/owncloud/ocis/pull/3652
https://github.com/owncloud/ocis/pull/3681

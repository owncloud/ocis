Enhancement: Update Reva to version 2.14.0

Changelog for reva 2.14.0 (2023-06-05)
=======================================

*   Bugfix [cs3org/reva#3919](https://github.com/cs3org/reva/pull/3919): We added missing timestamps to events
*   Bugfix [cs3org/reva#3911](https://github.com/cs3org/reva/pull/3911): Clean IDCache properly
*   Bugfix [cs3org/reva#3896](https://github.com/cs3org/reva/pull/3896): Do not lose old revisions when overwriting a file during copy
*   Bugfix [cs3org/reva#3918](https://github.com/cs3org/reva/pull/3918): Dont enumerate users
*   Bugfix [cs3org/reva#3902](https://github.com/cs3org/reva/pull/3902): Do not try to use the cache for empty node
*   Bugfix [cs3org/reva#3877](https://github.com/cs3org/reva/pull/3877): Empty exact list while searching for a sharee
*   Bugfix [cs3org/reva#3906](https://github.com/cs3org/reva/pull/3906): Fix preflight requests
*   Bugfix [cs3org/reva#3934](https://github.com/cs3org/reva/pull/3934): Fix the space editor permissions
*   Bugfix [cs3org/reva#3899](https://github.com/cs3org/reva/pull/3899): Harden uploads
*   Bugfix [cs3org/reva#3917](https://github.com/cs3org/reva/pull/3917): Prevent last space manager from leaving
*   Bugfix [cs3org/reva#3866](https://github.com/cs3org/reva/pull/3866): Fix public link lookup performance
*   Bugfix [cs3org/reva#3904](https://github.com/cs3org/reva/pull/3904): Improve performance of directory listings
*   Enhancement [cs3org/reva#3893](https://github.com/cs3org/reva/pull/3893): Cleanup Space Delete permissions
*   Enhancement [cs3org/reva#3894](https://github.com/cs3org/reva/pull/3894): Fix err when the user share the locked file
*   Enhancement [cs3org/reva#3913](https://github.com/cs3org/reva/pull/3913): Introduce FullTextSearch Capability
*   Enhancement [cs3org/reva#3898](https://github.com/cs3org/reva/pull/3898): Add Graph User capabilities
*   Enhancement [cs3org/reva#3496](https://github.com/cs3org/reva/pull/3496): Add otlp tracing exporter
*   Enhancement [cs3org/reva#3922](https://github.com/cs3org/reva/pull/3922): Rename permissions

Changelog for reva 2.13.3 (2023-05-17)
=======================================

*   Bugfix [cs3org/reva#3890](https://github.com/cs3org/reva/pull/3890): Bring back public link sharing of project space roots
*   Bugfix [cs3org/reva#3888](https://github.com/cs3org/reva/pull/3888): We fixed a bug that unnecessarily fetched all members of a group
*   Bugfix [cs3org/reva#3886](https://github.com/cs3org/reva/pull/3886): Decomposedfs no longer deadlocks when cache is disabled
*   Bugfix [cs3org/reva#3892](https://github.com/cs3org/reva/pull/3892): Fix public links
*   Bugfix [cs3org/reva#3876](https://github.com/cs3org/reva/pull/3876): Remove go-micro/store/redis specific workaround
*   Bugfix [cs3org/reva#3889](https://github.com/cs3org/reva/pull/3889): Update space root mtime when changing space metadata
*   Bugfix [cs3org/reva#3836](https://github.com/cs3org/reva/pull/3836): Fix spaceID in the decomposedFS
*   Bugfix [cs3org/reva#3867](https://github.com/cs3org/reva/pull/3867): Restore last version after positive result
*   Bugfix [cs3org/reva#3849](https://github.com/cs3org/reva/pull/3849): Prevent sharing space roots and personal spaces
*   Enhancement [cs3org/reva#3865](https://github.com/cs3org/reva/pull/3865): Remove unneccessary code from gateway
*   Enhancement [cs3org/reva#3895](https://github.com/cs3org/reva/pull/3895): Add missing expiry date to shares

Changelog for reva 2.13.2 (2023-05-08)
=======================================

*   Bugfix [cs3org/reva#3845](https://github.com/cs3org/reva/pull/3845): Fix propagation
*   Bugfix [cs3org/reva#3856](https://github.com/cs3org/reva/pull/3856): Fix response code
*   Bugfix [cs3org/reva#3857](https://github.com/cs3org/reva/pull/3857): Fix trashbin purge

Changelog for reva 2.13.1 (2023-05-03)
=======================================

*   Bugfix [cs3org/reva#3843](https://github.com/cs3org/reva/pull/3843): Allow scope check to impersonate space owners

Changelog for reva 2.13.0 (2023-05-02)
=======================================

*   Bugfix [cs3org/reva#3570](https://github.com/cs3org/reva/pull/3570): Return 425 on HEAD
*   Bugfix [cs3org/reva#3830](https://github.com/cs3org/reva/pull/3830): Be more robust when logging errors
*   Bugfix [cs3org/reva#3815](https://github.com/cs3org/reva/pull/3815): Bump micro redis store
*   Bugfix [cs3org/reva#3596](https://github.com/cs3org/reva/pull/3596): Cache CreateHome calls
*   Bugfix [cs3org/reva#3823](https://github.com/cs3org/reva/pull/3823): Deny correctly in decomposedfs
*   Bugfix [cs3org/reva#3826](https://github.com/cs3org/reva/pull/3826): Add by group index to decomposedfs
*   Bugfix [cs3org/reva#3618](https://github.com/cs3org/reva/pull/3618): Drain body on failed put
*   Bugfix [cs3org/reva#3685](https://github.com/cs3org/reva/pull/3685): Send fileid on copy
*   Bugfix [cs3org/reva#3688](https://github.com/cs3org/reva/pull/3688): Return 425 on GET
*   Bugfix [cs3org/reva#3755](https://github.com/cs3org/reva/pull/3755): Fix app provider language validation
*   Bugfix [cs3org/reva#3800](https://github.com/cs3org/reva/pull/3800): Fix building for freebsd
*   Bugfix [cs3org/reva#3700](https://github.com/cs3org/reva/pull/3700): Fix caching
*   Bugfix [cs3org/reva#3535](https://github.com/cs3org/reva/pull/3535): Fix ceph driver storage fs implementation
*   Bugfix [cs3org/reva#3764](https://github.com/cs3org/reva/pull/3764): Fix missing CORS config in ocdav service
*   Bugfix [cs3org/reva#3710](https://github.com/cs3org/reva/pull/3710): Fix error when try to delete space without permission
*   Bugfix [cs3org/reva#3822](https://github.com/cs3org/reva/pull/3822): Fix deleting spaces
*   Bugfix [cs3org/reva#3718](https://github.com/cs3org/reva/pull/3718): Fix revad-eos docker image which was failing to build
*   Bugfix [cs3org/reva#3559](https://github.com/cs3org/reva/pull/3559): Fix build on freebsd
*   Bugfix [cs3org/reva#3696](https://github.com/cs3org/reva/pull/3696): Fix ldap filters when checking for enabled users
*   Bugfix [cs3org/reva#3767](https://github.com/cs3org/reva/pull/3767): Decode binary UUID when looking up a users group memberships
*   Bugfix [cs3org/reva#3741](https://github.com/cs3org/reva/pull/3741): Fix listing shares to multiple groups
*   Bugfix [cs3org/reva#3834](https://github.com/cs3org/reva/pull/3834): Return correct error during MKCOL
*   Bugfix [cs3org/reva#3841](https://github.com/cs3org/reva/pull/3841): Fix nil pointer and improve logging
*   Bugfix [cs3org/reva#3831](https://github.com/cs3org/reva/pull/3831): Ignore 'null' mtime on tus upload
*   Bugfix [cs3org/reva#3758](https://github.com/cs3org/reva/pull/3758): Fix public links with enforced password
*   Bugfix [cs3org/reva#3814](https://github.com/cs3org/reva/pull/3814): Fix stat cache access
*   Bugfix [cs3org/reva#3650](https://github.com/cs3org/reva/pull/3650): FreeBSD xattr support
*   Bugfix [cs3org/reva#3827](https://github.com/cs3org/reva/pull/3827): Initialize user cache for decomposedfs
*   Bugfix [cs3org/reva#3818](https://github.com/cs3org/reva/pull/3818): Invalidate cache when deleting space
*   Bugfix [cs3org/reva#3812](https://github.com/cs3org/reva/pull/3812): Filemetadata Cache now deletes keys without listing them first
*   Bugfix [cs3org/reva#3817](https://github.com/cs3org/reva/pull/3817): Pipeline cache deletes
*   Bugfix [cs3org/reva#3711](https://github.com/cs3org/reva/pull/3711): Replace ini metadata backend by messagepack backend
*   Bugfix [cs3org/reva#3828](https://github.com/cs3org/reva/pull/3828): Send quota when listing spaces in decomposedfs
*   Bugfix [cs3org/reva#3681](https://github.com/cs3org/reva/pull/3681): Fix etag of "empty" shares jail
*   Bugfix [cs3org/reva#3748](https://github.com/cs3org/reva/pull/3748): Prevent service from panicking
*   Bugfix [cs3org/reva#3816](https://github.com/cs3org/reva/pull/3816): Write Metadata once
*   Change [cs3org/reva#3641](https://github.com/cs3org/reva/pull/3641): Hide file versions for share receivers
*   Change [cs3org/reva#3820](https://github.com/cs3org/reva/pull/3820): Streamline stores
*   Enhancement [cs3org/reva#3732](https://github.com/cs3org/reva/pull/3732): Make method for detecting the metadata backend public
*   Enhancement [cs3org/reva#3789](https://github.com/cs3org/reva/pull/3789): Add capabilities indicating if user attributes are read-only
*   Enhancement [cs3org/reva#3792](https://github.com/cs3org/reva/pull/3792): Add a prometheus gauge to keep track of active uploads and downloads
*   Enhancement [cs3org/reva#3637](https://github.com/cs3org/reva/pull/3637): Add an ID to each events
*   Enhancement [cs3org/reva#3704](https://github.com/cs3org/reva/pull/3704): Add more information to events
*   Enhancement [cs3org/reva#3744](https://github.com/cs3org/reva/pull/3744): Add LDAP user type attribute
*   Enhancement [cs3org/reva#3806](https://github.com/cs3org/reva/pull/3806): Decomposedfs now supports filtering spaces by owner
*   Enhancement [cs3org/reva#3730](https://github.com/cs3org/reva/pull/3730): Antivirus
*   Enhancement [cs3org/reva#3531](https://github.com/cs3org/reva/pull/3531): Async Postprocessing
*   Enhancement [cs3org/reva#3571](https://github.com/cs3org/reva/pull/3571): Async Upload Improvements
*   Enhancement [cs3org/reva#3801](https://github.com/cs3org/reva/pull/3801): Cache node ids
*   Enhancement [cs3org/reva#3690](https://github.com/cs3org/reva/pull/3690): Check set project space quota permission
*   Enhancement [cs3org/reva#3686](https://github.com/cs3org/reva/pull/3686): User disabling functionality
*   Enhancement [cs3org/reva#3505](https://github.com/cs3org/reva/pull/3505): Fix eosgrpc package
*   Enhancement [cs3org/reva#3575](https://github.com/cs3org/reva/pull/3575): Fix skip group grant index cleanup
*   Enhancement [cs3org/reva#3564](https://github.com/cs3org/reva/pull/3564): Fix tag pkg
*   Enhancement [cs3org/reva#3756](https://github.com/cs3org/reva/pull/3756): Prepare for GDPR export
*   Enhancement [cs3org/reva#3612](https://github.com/cs3org/reva/pull/3612): Group feature changed event added
*   Enhancement [cs3org/reva#3729](https://github.com/cs3org/reva/pull/3729): Improve decomposedfs performance, esp. with network fs/cache
*   Enhancement [cs3org/reva#3697](https://github.com/cs3org/reva/pull/3697): Improve the ini file metadata backend
*   Enhancement [cs3org/reva#3819](https://github.com/cs3org/reva/pull/3819): Allow creating internal links without permission
*   Enhancement [cs3org/reva#3740](https://github.com/cs3org/reva/pull/3740): Limit concurrency in decomposedfs
*   Enhancement [cs3org/reva#3569](https://github.com/cs3org/reva/pull/3569): Always list shares jail when listing spaces
*   Enhancement [cs3org/reva#3788](https://github.com/cs3org/reva/pull/3788): Make resharing configurable
*   Enhancement [cs3org/reva#3674](https://github.com/cs3org/reva/pull/3674): Introduce ini file based metadata backend
*   Enhancement [cs3org/reva#3728](https://github.com/cs3org/reva/pull/3728): Automatically migrate file metadata from xattrs to messagepack
*   Enhancement [cs3org/reva#3807](https://github.com/cs3org/reva/pull/3807): Name Validation
*   Enhancement [cs3org/reva#3574](https://github.com/cs3org/reva/pull/3574): Opaque space group
*   Enhancement [cs3org/reva#3598](https://github.com/cs3org/reva/pull/3598): Pass estream to Storage Providers
*   Enhancement [cs3org/reva#3763](https://github.com/cs3org/reva/pull/3763): Add a capability for personal data export
*   Enhancement [cs3org/reva#3577](https://github.com/cs3org/reva/pull/3577): Prepare for SSE
*   Enhancement [cs3org/reva#3731](https://github.com/cs3org/reva/pull/3731): Add config option to enforce passwords on public links
*   Enhancement [cs3org/reva#3693](https://github.com/cs3org/reva/pull/3693): Enforce the PublicLink.Write permission
*   Enhancement [cs3org/reva#3497](https://github.com/cs3org/reva/pull/3497): Introduce owncloud 10 publiclink manager
*   Enhancement [cs3org/reva#3714](https://github.com/cs3org/reva/pull/3714): Add global max quota option and quota for CreateHome
*   Enhancement [cs3org/reva#3759](https://github.com/cs3org/reva/pull/3759): Set correct share type when listing shares
*   Enhancement [cs3org/reva#3594](https://github.com/cs3org/reva/pull/3594): Add expiration to user and group shares
*   Enhancement [cs3org/reva#3580](https://github.com/cs3org/reva/pull/3580): Share expired event
*   Enhancement [cs3org/reva#3620](https://github.com/cs3org/reva/pull/3620): Allow a new ShareType `SpaceMembershipGroup`
*   Enhancement [cs3org/reva#3609](https://github.com/cs3org/reva/pull/3609): Space Management Permissions
*   Enhancement [cs3org/reva#3655](https://github.com/cs3org/reva/pull/3655): Add expiration date to space memberships
*   Enhancement [cs3org/reva#3697](https://github.com/cs3org/reva/pull/3697): Add support for redis sentinel caches
*   Enhancement [cs3org/reva#3552](https://github.com/cs3org/reva/pull/3552): Suppress tusd logs
*   Enhancement [cs3org/reva#3555](https://github.com/cs3org/reva/pull/3555): Tags
*   Enhancement [cs3org/reva#3785](https://github.com/cs3org/reva/pull/3785): Increase unit test coverage in the ocdav service
*   Enhancement [cs3org/reva#3739](https://github.com/cs3org/reva/pull/3739): Try to rename uploaded files to their final position
*   Enhancement [cs3org/reva#3610](https://github.com/cs3org/reva/pull/3610): Walk and log chi routes


https://github.com/owncloud/ocis/pull/6448
https://github.com/owncloud/ocis/pull/6447
https://github.com/owncloud/ocis/pull/6381
https://github.com/owncloud/ocis/pull/6305
https://github.com/owncloud/ocis/pull/6339
https://github.com/owncloud/ocis/pull/6205
https://github.com/owncloud/ocis/pull/6186

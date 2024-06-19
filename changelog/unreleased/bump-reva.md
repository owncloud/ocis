Enhancement: Bump Reva to v2.20.0

*   Bugfix [cs3org/reva#4623](https://github.com/cs3org/reva/pull/4623): Consistently use spaceid and nodeid in logs
*   Bugfix [cs3org/reva#4584](https://github.com/cs3org/reva/pull/4584): Prevent copying a file to a parent folder
*   Bugfix [cs3org/reva#4700](https://github.com/cs3org/reva/pull/4700): Clean empty trash node path on delete
*   Bugfix [cs3org/reva#4567](https://github.com/cs3org/reva/pull/4567): Fix error message in authprovider if user is not found
*   Bugfix [cs3org/reva#4615](https://github.com/cs3org/reva/pull/4615): Write blob based on session id
*   Bugfix [cs3org/reva#4557](https://github.com/cs3org/reva/pull/4557): Fix ceph build
*   Bugfix [cs3org/reva#4711](https://github.com/cs3org/reva/pull/4711): Duplicate headers in DAV responses
*   Bugfix [cs3org/reva#4568](https://github.com/cs3org/reva/pull/4568): Fix sharing invite on virtual drive
*   Bugfix [cs3org/reva#4559](https://github.com/cs3org/reva/pull/4559): Fix graph drive invite
*   Bugfix [cs3org/reva#4593](https://github.com/cs3org/reva/pull/4593): Make initiatorIDs also work on uploads
*   Bugfix [cs3org/reva#4608](https://github.com/cs3org/reva/pull/4608): Use gateway selector in jsoncs3
*   Bugfix [cs3org/reva#4546](https://github.com/cs3org/reva/pull/4546): Fix the mount points naming
*   Bugfix [cs3org/reva#4678](https://github.com/cs3org/reva/pull/4678): Fix nats encoding
*   Bugfix [cs3org/reva#4630](https://github.com/cs3org/reva/pull/4630): Fix ocm-share-id
*   Bugfix [cs3org/reva#4518](https://github.com/cs3org/reva/pull/4518): Fix an error when lock/unlock a file
*   Bugfix [cs3org/reva#4622](https://github.com/cs3org/reva/pull/4622): Fix public share update
*   Bugfix [cs3org/reva#4566](https://github.com/cs3org/reva/pull/4566): Fix public link previews
*   Bugfix [cs3org/reva#4589](https://github.com/cs3org/reva/pull/4589): Fix uploading via a public link
*   Bugfix [cs3org/reva#4660](https://github.com/cs3org/reva/pull/4660): Fix creating documents in nested folders of public shares
*   Bugfix [cs3org/reva#4635](https://github.com/cs3org/reva/pull/4635): Fix nil pointer when removing groups from space
*   Bugfix [cs3org/reva#4709](https://github.com/cs3org/reva/pull/4709): Fix share update
*   Bugfix [cs3org/reva#4661](https://github.com/cs3org/reva/pull/4661): Fix space share update for ocs
*   Bugfix [cs3org/reva#4656](https://github.com/cs3org/reva/pull/4656): Fix space share update
*   Bugfix [cs3org/reva#4561](https://github.com/cs3org/reva/pull/4561): Fix Stat() by Path on re-created resource
*   Bugfix [cs3org/reva#4710](https://github.com/cs3org/reva/pull/4710): Tolerate missing user space index
*   Bugfix [cs3org/reva#4632](https://github.com/cs3org/reva/pull/4632): Fix access to files withing a public link targeting a space root
*   Bugfix [cs3org/reva#4603](https://github.com/cs3org/reva/pull/4603): Mask user email in output
*   Change [cs3org/reva#4542](https://github.com/cs3org/reva/pull/4542): Drop unused service spanning stat cache
*   Enhancement [cs3org/reva#4712](https://github.com/cs3org/reva/pull/4712): Add the error translation to the utils
*   Enhancement [cs3org/reva#4696](https://github.com/cs3org/reva/pull/4696): Add List method to ocis and s3ng blobstore
*   Enhancement [cs3org/reva#4693](https://github.com/cs3org/reva/pull/4693): Add mimetype for sb3 files
*   Enhancement [cs3org/reva#4699](https://github.com/cs3org/reva/pull/4699): Add a Path method to blobstore
*   Enhancement [cs3org/reva#4695](https://github.com/cs3org/reva/pull/4695): Add photo and image props
*   Enhancement [cs3org/reva#4706](https://github.com/cs3org/reva/pull/4706): Add secureview flag when listing apps via http
*   Enhancement [cs3org/reva#4585](https://github.com/cs3org/reva/pull/4585): Move more consistency checks to the usershare API
*   Enhancement [cs3org/reva#4702](https://github.com/cs3org/reva/pull/4702): Added theme capability
*   Enhancement [cs3org/reva#4672](https://github.com/cs3org/reva/pull/4672): Add virus filter to list uploads sessions
*   Enhancement [cs3org/reva#4614](https://github.com/cs3org/reva/pull/4614): Bump mockery to v2.40.2
*   Enhancement [cs3org/reva#4621](https://github.com/cs3org/reva/pull/4621): Use a memory cache for the personal space creation cache
*   Enhancement [cs3org/reva#4556](https://github.com/cs3org/reva/pull/4556): Allow tracing requests by giving util functions a context
*   Enhancement [cs3org/reva#4694](https://github.com/cs3org/reva/pull/4694): Expose SecureView in WebDAV permissions
*   Enhancement [cs3org/reva#4652](https://github.com/cs3org/reva/pull/4652): Better error codes when removing a space member
*   Enhancement [cs3org/reva#4725](https://github.com/cs3org/reva/pull/4725): Unique share mountpoint name
*   Enhancement [cs3org/reva#4689](https://github.com/cs3org/reva/pull/4689): Extend service account permissions
*   Enhancement [cs3org/reva#4545](https://github.com/cs3org/reva/pull/4545): Extend service account permissions
*   Enhancement [cs3org/reva#4581](https://github.com/cs3org/reva/pull/4581): Make decomposedfs more extensible
*   Enhancement [cs3org/reva#4564](https://github.com/cs3org/reva/pull/4564): Send file locked/unlocked events
*   Enhancement [cs3org/reva#4730](https://github.com/cs3org/reva/pull/4730): Improve posixfs storage driver
*   Enhancement [cs3org/reva#4587](https://github.com/cs3org/reva/pull/4587): Allow passing a initiator id
*   Enhancement [cs3org/reva#4645](https://github.com/cs3org/reva/pull/4645): Add ItemID to LinkRemoved
*   Enhancement [cs3org/reva#4686](https://github.com/cs3org/reva/pull/4686): Mint view only token for open in app requests
*   Enhancement [cs3org/reva#4606](https://github.com/cs3org/reva/pull/4606): Remove resharing
*   Enhancement [cs3org/reva#4643](https://github.com/cs3org/reva/pull/4643): Secure viewer share role
*   Enhancement [cs3org/reva#4631](https://github.com/cs3org/reva/pull/4631): Add space-share-updated event
*   Enhancement [cs3org/reva#4685](https://github.com/cs3org/reva/pull/4685): Support t and x in ACEs
*   Enhancement [cs3org/reva#4625](https://github.com/cs3org/reva/pull/4625): Test async processing cornercases
*   Enhancement [cs3org/reva#4653](https://github.com/cs3org/reva/pull/4653): Allow to resolve public shares without the ocs tokeninfo endpoint
*   Enhancement [cs3org/reva#4657](https://github.com/cs3org/reva/pull/4657): Add ScanData to Uploadsession

https://github.com/owncloud/ocis/pull/9415
https://github.com/owncloud/ocis/pull/9377
https://github.com/owncloud/ocis/pull/9330
https://github.com/owncloud/ocis/pull/9318
https://github.com/owncloud/ocis/pull/9269
https://github.com/owncloud/ocis/pull/9236
https://github.com/owncloud/ocis/pull/9188
https://github.com/owncloud/ocis/pull/9132
https://github.com/owncloud/ocis/pull/9041
https://github.com/owncloud/ocis/pull/9002
https://github.com/owncloud/ocis/pull/8917
https://github.com/owncloud/ocis/pull/8795
https://github.com/owncloud/ocis/pull/8701
https://github.com/owncloud/ocis/pull/8606
https://github.com/owncloud/ocis/pull/8937

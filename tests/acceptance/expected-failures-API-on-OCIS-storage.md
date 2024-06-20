## Scenarios from ownCloud10 core API tests that are expected to fail with OCIS storage while running with the Graph API

The expected failures in this file are from features in the owncloud/ocis repo.

### File

Basic file management like up and download, move, copy, properties, trash, versions and chunking.

#### [copy personal space file to shared folder root result share in decline state](https://github.com/owncloud/ocis/issues/6999)

- [coreApiWebdavProperties/copyFile.feature:297](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L297)
- [coreApiWebdavProperties/copyFile.feature:298](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L298)

#### [Custom dav properties with namespaces are rendered incorrectly](https://github.com/owncloud/ocis/issues/2140)

_ocdav: double-check the webdav property parsing when custom namespaces are used_

- [coreApiWebdavProperties/setFileProperties.feature:36](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L36)
- [coreApiWebdavProperties/setFileProperties.feature:37](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L37)
- [coreApiWebdavProperties/setFileProperties.feature:42](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L42)
- [coreApiWebdavProperties/setFileProperties.feature:78](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L78)
- [coreApiWebdavProperties/setFileProperties.feature:79](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L79)
- [coreApiWebdavProperties/setFileProperties.feature:84](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L84)

#### [Cannot set custom webDav properties](https://github.com/owncloud/product/issues/264)

- [coreApiWebdavProperties/getFileProperties.feature:348](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L348)
- [coreApiWebdavProperties/getFileProperties.feature:349](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L349)
- [coreApiWebdavProperties/getFileProperties.feature:354](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L354)
- [coreApiWebdavProperties/getFileProperties.feature:384](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L384)
- [coreApiWebdavProperties/getFileProperties.feature:385](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L385)
- [coreApiWebdavProperties/getFileProperties.feature:390](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L390)

#### [file versions do not report the version author](https://github.com/owncloud/ocis/issues/2914)

- [coreApiVersions/fileVersionAuthor.feature:15](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L15)
- [coreApiVersions/fileVersionAuthor.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L58)
- [coreApiVersions/fileVersionAuthor.feature:88](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L88)
- [coreApiVersions/fileVersionAuthor.feature:117](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L117)
- [coreApiVersions/fileVersionAuthor.feature:153](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L153)
- [coreApiVersions/fileVersionAuthor.feature:183](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L183)
- [coreApiVersions/fileVersionAuthor.feature:217](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L217)
- [coreApiVersions/fileVersionAuthor.feature:258](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L258)
- [coreApiVersions/fileVersionAuthor.feature:334](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L334)
- [coreApiVersions/fileVersionAuthor.feature:407](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L407)
- [coreApiVersions/fileVersionAuthor.feature:436](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L436)

### Sync

Synchronization features like etag propagation, setting mtime and locking files

#### [Uploading an old method chunked file with checksum should fail using new DAV path](https://github.com/owncloud/ocis/issues/2323)

- [coreApiMain/checksums.feature:269](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiMain/checksums.feature#L269)
- [coreApiMain/checksums.feature:274](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiMain/checksums.feature#L274)

### Share

#### [copying a folder within a public link folder to folder with same name as an already existing file overwrites the parent file](https://github.com/owncloud/ocis/issues/1232)

- [coreApiSharePublicLink2/copyFromPublicLink.feature:75](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L75)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:105](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L105)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:199](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L199)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:200](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L200)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:217](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L217)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:218](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L218)

#### [d:quota-available-bytes in dprop of PROPFIND give wrong response value](https://github.com/owncloud/ocis/issues/8197)

- [coreApiWebdavProperties/getQuota.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L56)
- [coreApiWebdavProperties/getQuota.feature:57](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L57)
- [coreApiWebdavProperties/getQuota.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L58)
- [coreApiWebdavProperties/getQuota.feature:72](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L72)
- [coreApiWebdavProperties/getQuota.feature:73](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L73)
- [coreApiWebdavProperties/getQuota.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L74)

#### [deleting a file inside a received shared folder is moved to the trash-bin of the sharer not the receiver](https://github.com/owncloud/ocis/issues/1124)

- [coreApiTrashbin/trashbinSharingToShares.feature:34](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L34)
- [coreApiTrashbin/trashbinSharingToShares.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L55)
- [coreApiTrashbin/trashbinSharingToShares.feature:60](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L60)
- [coreApiTrashbin/trashbinSharingToShares.feature:85](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L85)
- [coreApiTrashbin/trashbinSharingToShares.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L90)
- [coreApiTrashbin/trashbinSharingToShares.feature:146](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L146)
- [coreApiTrashbin/trashbinSharingToShares.feature:151](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L151)
- [coreApiTrashbin/trashbinSharingToShares.feature:209](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L209)
- [coreApiTrashbin/trashbinSharingToShares.feature:214](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L214)
- [coreApiTrashbin/trashbinSharingToShares.feature:241](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L241)
- [coreApiTrashbin/trashbinSharingToShares.feature:269](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L269)

### Other

API, search, favorites, config, capabilities, not existing endpoints, CORS and others

#### [sending MKCOL requests to another or non-existing user's webDav endpoints as normal user should return 404](https://github.com/owncloud/ocis/issues/5049)

_ocdav: api compatibility, return correct status code_

- [coreApiAuth/webDavMKCOLAuth.feature:42](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavMKCOLAuth.feature#L42)
- [coreApiAuth/webDavMKCOLAuth.feature:53](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavMKCOLAuth.feature#L53)

#### [trying to lock file of another user gives http 200](https://github.com/owncloud/ocis/issues/2176)

- [coreApiAuth/webDavLOCKAuth.feature:46](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavLOCKAuth.feature#L46)
- [coreApiAuth/webDavLOCKAuth.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavLOCKAuth.feature#L58)

#### [send POST requests to another user's webDav endpoints as normal user](https://github.com/owncloud/ocis/issues/1287)

_ocdav: api compatibility, return correct status code_

- [coreApiAuth/webDavPOSTAuth.feature:46](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavPOSTAuth.feature#L46)
- [coreApiAuth/webDavPOSTAuth.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavPOSTAuth.feature#L55)

#### Another users space literally does not exist because it is not listed as a space for him, 404 seems correct, expects 403

- [coreApiAuth/webDavPUTAuth.feature:46](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavPUTAuth.feature#L46)
- [coreApiAuth/webDavPUTAuth.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavPUTAuth.feature#L58)

#### [[old/new/spaces] In ocis and oc10, REPORT request response differently](https://github.com/owncloud/ocis/issues/4712)

- [coreApiWebdavOperations/search.feature:208](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L208)
- [coreApiWebdavOperations/search.feature:209](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L209)
- [coreApiWebdavOperations/search.feature:214](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L214)
- [coreApiWebdavOperations/search.feature:240](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L240)
- [coreApiWebdavOperations/search.feature:241](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L241)
- [coreApiWebdavOperations/search.feature:246](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L246)

#### [Support for favorites](https://github.com/owncloud/ocis/issues/1228)

- [coreApiFavorites/favorites.feature:117](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L117)
- [coreApiFavorites/favorites.feature:118](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L118)
- [coreApiFavorites/favorites.feature:144](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L144)
- [coreApiFavorites/favorites.feature:145](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L145)
- [coreApiFavorites/favorites.feature:221](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L221)
- [coreApiFavorites/favorites.feature:222](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L222)

And other missing implementation of favorites

- [coreApiFavorites/favorites.feature:169](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L169)
- [coreApiFavorites/favorites.feature:170](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L170)
- [coreApiFavorites/favorites.feature:175](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L175)
- [coreApiFavorites/favorites.feature:202](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L202)
- [coreApiFavorites/favorites.feature:203](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L203)
- [coreApiFavorites/favorites.feature:208](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L208)
- [coreApiFavorites/favoritesSharingToShares.feature:84](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favoritesSharingToShares.feature#L84)
- [coreApiFavorites/favoritesSharingToShares.feature:85](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favoritesSharingToShares.feature#L85)

#### [WWW-Authenticate header for unauthenticated requests is not clear](https://github.com/owncloud/ocis/issues/2285)

- [coreApiWebdavOperations/refuseAccess.feature:21](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L21)
- [coreApiWebdavOperations/refuseAccess.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L22)

#### [Sharing a same file twice to the same group](https://github.com/owncloud/ocis/issues/1710)

- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:650](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L650)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:651](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L651)

#### [PATCH request for TUS upload with wrong checksum gives incorrect response](https://github.com/owncloud/ocis/issues/1755)

- [coreApiWebdavUploadTUS/checksums.feature:86](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L86)
- [coreApiWebdavUploadTUS/checksums.feature:87](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L87)
- [coreApiWebdavUploadTUS/checksums.feature:88](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L88)
- [coreApiWebdavUploadTUS/checksums.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L89)
- [coreApiWebdavUploadTUS/checksums.feature:94](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L94)
- [coreApiWebdavUploadTUS/checksums.feature:95](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L95)
- [coreApiWebdavUploadTUS/checksums.feature:175](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L175)
- [coreApiWebdavUploadTUS/checksums.feature:176](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L176)
- [coreApiWebdavUploadTUS/checksums.feature:181](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L181)
- [coreApiWebdavUploadTUS/checksums.feature:228](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L228)
- [coreApiWebdavUploadTUS/checksums.feature:229](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L229)
- [coreApiWebdavUploadTUS/checksums.feature:230](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L230)
- [coreApiWebdavUploadTUS/checksums.feature:231](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L231)
- [coreApiWebdavUploadTUS/checksums.feature:236](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L236)
- [coreApiWebdavUploadTUS/checksums.feature:237](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L237)
- [coreApiWebdavUploadTUS/checksums.feature:284](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L284)
- [coreApiWebdavUploadTUS/checksums.feature:285](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L285)
- [coreApiWebdavUploadTUS/checksums.feature:286](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L286)
- [coreApiWebdavUploadTUS/checksums.feature:287](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L287)
- [coreApiWebdavUploadTUS/checksums.feature:292](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L292)
- [coreApiWebdavUploadTUS/checksums.feature:293](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L293)
- [coreApiWebdavUploadTUS/optionsRequest.feature:10](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L10)
- [coreApiWebdavUploadTUS/optionsRequest.feature:25](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L25)
- [coreApiWebdavUploadTUS/uploadToShare.feature:216](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L216)
- [coreApiWebdavUploadTUS/uploadToShare.feature:217](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L217)
- [coreApiWebdavUploadTUS/uploadToShare.feature:239](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L239)
- [coreApiWebdavUploadTUS/uploadToShare.feature:240](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L240)
- [coreApiWebdavUploadTUS/uploadToShare.feature:262](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L262)
- [coreApiWebdavUploadTUS/uploadToShare.feature:263](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L263)
- [coreApiWebdavUploadTUS/uploadToShare.feature:304](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L304)
- [coreApiWebdavUploadTUS/uploadToShare.feature:305](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L305)
- [coreApiWebdavUploadTUS/uploadToShare.feature:354](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L354)
- [coreApiWebdavUploadTUS/uploadToShare.feature:355](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L355)

#### [TUS OPTIONS requests do not reply with TUS headers when invalid password](https://github.com/owncloud/ocis/issues/1012)

- [coreApiWebdavUploadTUS/optionsRequest.feature:40](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L40)
- [coreApiWebdavUploadTUS/optionsRequest.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L55)

#### [Shares to deleted group listed in the response](https://github.com/owncloud/ocis/issues/2441)

- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:502](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L502)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:503](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L503)

#### [copying the file inside Shares folder returns 412](https://github.com/owncloud/ocis/issues/3874)

- [coreApiWebdavProperties/copyFile.feature:439](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L439)
- [coreApiWebdavProperties/copyFile.feature:440](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L440)
- [coreApiWebdavProperties/copyFile.feature:445](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L445)
- [coreApiWebdavProperties/copyFile.feature:469](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L469)
- [coreApiWebdavProperties/copyFile.feature:470](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L470)
- [coreApiWebdavProperties/copyFile.feature:475](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L475)
- [coreApiWebdavProperties/copyFile.feature:275](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L275)
- [coreApiWebdavProperties/copyFile.feature:276](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L276)

### Won't fix

Not everything needs to be implemented for ocis. While the oc10 testsuite covers these things we are not looking at them right now.

- _The `OC-LazyOps` header is [no longer supported by the client](https://github.com/owncloud/client/pull/8398), implementing this is not necessary for a first production release. We plan to have an upload state machine to visualize the state of a file, see https://github.com/owncloud/ocis/issues/214_
- _Blacklisted ignored files are no longer required because ocis can handle `.htaccess` files without security implications introduced by serving user provided files with apache._

#### [Blacklist files extensions](https://github.com/owncloud/ocis/issues/2177)

- [coreApiWebdavProperties/copyFile.feature:117](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L117)
- [coreApiWebdavProperties/copyFile.feature:118](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L118)
- [coreApiWebdavProperties/copyFile.feature:123](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L123)
- [coreApiWebdavProperties/createFileFolder.feature:106](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/createFileFolder.feature#L106)
- [coreApiWebdavProperties/createFileFolder.feature:107](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/createFileFolder.feature#L107)
- [coreApiWebdavProperties/createFileFolder.feature:112](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/createFileFolder.feature#L112)
- [coreApiWebdavUpload/uploadFile.feature:181](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUpload/uploadFile.feature#L181)
- [coreApiWebdavUpload/uploadFile.feature:180](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUpload/uploadFile.feature#L180)
- [coreApiWebdavUpload/uploadFile.feature:186](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUpload/uploadFile.feature#L186)
- [coreApiWebdavMove2/moveFile.feature:217](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L217)
- [coreApiWebdavMove2/moveFile.feature:218](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L218)
- [coreApiWebdavMove2/moveFile.feature:223](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L223)

#### [cannot set blacklisted file names](https://github.com/owncloud/product/issues/260)

- [coreApiWebdavMove1/moveFolderToBlacklistedName.feature:20](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolderToBlacklistedName.feature#L20)
- [coreApiWebdavMove1/moveFolderToBlacklistedName.feature:21](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolderToBlacklistedName.feature#L21)
- [coreApiWebdavMove1/moveFolderToBlacklistedName.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolderToBlacklistedName.feature#L26)
- [coreApiWebdavMove2/moveFileToBlacklistedName.feature:18](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFileToBlacklistedName.feature#L18)
- [coreApiWebdavMove2/moveFileToBlacklistedName.feature:19](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFileToBlacklistedName.feature#L19)

#### [Share path in the response is different between share states](https://github.com/owncloud/ocis/issues/2540)

- [coreApiShareManagementToShares/acceptShares.feature:28](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L28)
- [coreApiShareManagementToShares/acceptShares.feature:64](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L64)
- [coreApiShareManagementToShares/acceptShares.feature:154](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L154)
- [coreApiShareManagementToShares/acceptShares.feature:186](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L186)
- [coreApiShareManagementToShares/acceptShares.feature:235](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L235)
- [coreApiShareManagementToShares/acceptShares.feature:303](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L303)
- [coreApiShareManagementToShares/acceptShares.feature:534](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L534)

#### [Renaming resource to banned name is allowed in spaces webdav](https://github.com/owncloud/ocis/issues/3099)

- [coreApiWebdavMove1/moveFolder.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L44)
- [coreApiWebdavMove1/moveFolder.feature:62](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L62)
- [coreApiWebdavMove1/moveFolder.feature:80](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L80)
- [coreApiWebdavMove2/moveFile.feature:179](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L179)
- [coreApiWebdavMove2/moveFileToBlacklistedName.feature:24](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFileToBlacklistedName.feature#L24)

#### [REPORT method on spaces returns an incorrect d:href response](https://github.com/owncloud/ocis/issues/3111)

- [coreApiFavorites/favorites.feature:123](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L123)
- [coreApiFavorites/favorites.feature:150](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L150)
- [coreApiFavorites/favorites.feature:227](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L227)

#### [HTTP status code differ while deleting file of another user's trash bin](https://github.com/owncloud/ocis/issues/3544)

- [coreApiTrashbin/trashbinDelete.feature:105](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L105)


#### [sharing the shares folder to users exits with different status code than in oc10 backend](https://github.com/owncloud/ocis/issues/2215)

- [coreApiShareCreateSpecialToShares2/createShareDefaultFolderForReceivedShares.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareDefaultFolderForReceivedShares.feature#L22)
- [coreApiShareCreateSpecialToShares2/createShareDefaultFolderForReceivedShares.feature:23](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareDefaultFolderForReceivedShares.feature#L23)

### [MOVE a file into same folder with same name returns 404 instead of 403](https://github.com/owncloud/ocis/issues/1976)

- [coreApiWebdavMove2/moveFile.feature:120](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L120)
- [coreApiWebdavMove2/moveFile.feature:121](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L121)
- [coreApiWebdavMove2/moveFile.feature:126](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L126)
- [coreApiWebdavMove1/moveFolder.feature:253](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L253)
- [coreApiWebdavMove1/moveFolder.feature:254](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L254)
- [coreApiWebdavMove1/moveFolder.feature:259](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L259)
- [coreApiWebdavMove2/moveShareOnOcis.feature:296](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L296)
- [coreApiWebdavMove2/moveShareOnOcis.feature:299](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L299)

Note: always have an empty line at the end of this file.
The bash script that processes this file requires that the last line has a newline on the end.

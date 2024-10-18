## Scenarios from ownCloud10 core API tests that are expected to fail with OCIS storage while running with the Graph API

The expected failures in this file are from features in the owncloud/ocis repo.

### File

Basic file management like up and download, move, copy, properties, trash, versions and chunking.

#### [Custom dav properties with namespaces are rendered incorrectly](https://github.com/owncloud/ocis/issues/2140)

_ocdav: double-check the webdav property parsing when custom namespaces are used_

- [coreApiWebdavProperties/setFileProperties.feature:128](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L128)
- [coreApiWebdavProperties/setFileProperties.feature:129](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L129)
- [coreApiWebdavProperties/setFileProperties.feature:130](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L130)

#### [file versions do not report the version author](https://github.com/owncloud/ocis/issues/2914)

- [coreApiVersions/fileVersionAuthor.feature:15](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L15)
- [coreApiVersions/fileVersionAuthor.feature:61](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L61)
- [coreApiVersions/fileVersionAuthor.feature:93](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L93)
- [coreApiVersions/fileVersionAuthor.feature:124](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L124)
- [coreApiVersions/fileVersionAuthor.feature:162](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L162)
- [coreApiVersions/fileVersionAuthor.feature:194](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L194)
- [coreApiVersions/fileVersionAuthor.feature:230](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L230)
- [coreApiVersions/fileVersionAuthor.feature:275](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L275)
- [coreApiVersions/fileVersionAuthor.feature:357](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L357)
- [coreApiVersions/fileVersionAuthor.feature:436](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L436)
- [coreApiVersions/fileVersionAuthor.feature:467](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L467)

### Sync

Synchronization features like etag propagation, setting mtime and locking files

#### [Uploading an old method chunked file with checksum should fail using new DAV path](https://github.com/owncloud/ocis/issues/2323)

- [coreApiMain/checksums.feature:219](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiMain/checksums.feature#L219)
- [coreApiMain/checksums.feature:220](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiMain/checksums.feature#L220)

### Share

#### [d:quota-available-bytes in dprop of PROPFIND give wrong response value](https://github.com/owncloud/ocis/issues/8197)

- [coreApiWebdavProperties/getQuota.feature:57](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L57)
- [coreApiWebdavProperties/getQuota.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L58)
- [coreApiWebdavProperties/getQuota.feature:59](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L59)
- [coreApiWebdavProperties/getQuota.feature:73](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L73)
- [coreApiWebdavProperties/getQuota.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L74)
- [coreApiWebdavProperties/getQuota.feature:75](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L75)

#### [deleting a file inside a received shared folder is moved to the trash-bin of the sharer not the receiver](https://github.com/owncloud/ocis/issues/1124)

- [coreApiTrashbin/trashbinSharingToShares.feature:53](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L53)
- [coreApiTrashbin/trashbinSharingToShares.feature:54](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L54)
- [coreApiTrashbin/trashbinSharingToShares.feature:81](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L81)
- [coreApiTrashbin/trashbinSharingToShares.feature:82](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L82)
- [coreApiTrashbin/trashbinSharingToShares.feature:138](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L138)
- [coreApiTrashbin/trashbinSharingToShares.feature:139](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L139)
- [coreApiTrashbin/trashbinSharingToShares.feature:196](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L196)
- [coreApiTrashbin/trashbinSharingToShares.feature:197](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L197)

### Other

API, search, favorites, config, capabilities, not existing endpoints, CORS and others

#### [sending MKCOL requests to another or non-existing user's webDav endpoints as normal user should return 404](https://github.com/owncloud/ocis/issues/5049)

_ocdav: api compatibility, return correct status code_

- [coreApiAuth/webDavMKCOLAuth.feature:42](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavMKCOLAuth.feature#L42)
- [coreApiAuth/webDavMKCOLAuth.feature:53](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavMKCOLAuth.feature#L53)

#### [trying to lock file of another user gives http 500](https://github.com/owncloud/ocis/issues/2176)

- [coreApiAuth/webDavLOCKAuth.feature:46](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavLOCKAuth.feature#L46)
- [coreApiAuth/webDavLOCKAuth.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavLOCKAuth.feature#L58)

#### [Missing properties in REPORT response](https://github.com/owncloud/ocis/issues/9780), [REPORT response has different date format for d:getlastmodified](https://github.com/owncloud/ocis/issues/9781), [d:getetag property has empty value in REPORT response](https://github.com/owncloud/ocis/issues/9783)

- [coreApiWebdavOperations/search.feature:180](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L180)
- [coreApiWebdavOperations/search.feature:181](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L181)
- [coreApiWebdavOperations/search.feature:182](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L182)
- [coreApiWebdavOperations/search.feature:208](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L208)
- [coreApiWebdavOperations/search.feature:209](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L209)
- [coreApiWebdavOperations/search.feature:210](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L210)

#### [Support for favorites](https://github.com/owncloud/ocis/issues/1228)

- [coreApiFavorites/favorites.feature:101](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L101)
- [coreApiFavorites/favorites.feature:102](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L102)
- [coreApiFavorites/favorites.feature:103](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L103)
- [coreApiFavorites/favorites.feature:124](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L124)
- [coreApiFavorites/favorites.feature:125](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L125)
- [coreApiFavorites/favorites.feature:126](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L126)
- [coreApiFavorites/favorites.feature:189](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L189)
- [coreApiFavorites/favorites.feature:190](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L190)
- [coreApiFavorites/favorites.feature:191](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L191)
- [coreApiFavorites/favorites.feature:145](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L145)
- [coreApiFavorites/favorites.feature:146](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L146)
- [coreApiFavorites/favorites.feature:147](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L147)
- [coreApiFavorites/favorites.feature:174](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L174)
- [coreApiFavorites/favorites.feature:175](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L175)
- [coreApiFavorites/favorites.feature:176](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L176)
- [coreApiFavorites/favoritesSharingToShares.feature:88](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favoritesSharingToShares.feature#L88)
- [coreApiFavorites/favoritesSharingToShares.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favoritesSharingToShares.feature#L89)

#### [WWW-Authenticate header for unauthenticated requests is not clear](https://github.com/owncloud/ocis/issues/2285)

- [coreApiWebdavOperations/refuseAccess.feature:21](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L21)
- [coreApiWebdavOperations/refuseAccess.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L22)

#### [PATCH request for TUS upload with wrong checksum gives incorrect response](https://github.com/owncloud/ocis/issues/1755)

- [coreApiWebdavUploadTUS/checksums.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L74)
- [coreApiWebdavUploadTUS/checksums.feature:75](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L75)
- [coreApiWebdavUploadTUS/checksums.feature:76](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L76)
- [coreApiWebdavUploadTUS/checksums.feature:77](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L77)
- [coreApiWebdavUploadTUS/checksums.feature:79](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L79)
- [coreApiWebdavUploadTUS/checksums.feature:78](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L78)
- [coreApiWebdavUploadTUS/checksums.feature:147](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L147)
- [coreApiWebdavUploadTUS/checksums.feature:148](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L148)
- [coreApiWebdavUploadTUS/checksums.feature:149](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L149)
- [coreApiWebdavUploadTUS/checksums.feature:192](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L192)
- [coreApiWebdavUploadTUS/checksums.feature:193](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L193)
- [coreApiWebdavUploadTUS/checksums.feature:194](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L194)
- [coreApiWebdavUploadTUS/checksums.feature:195](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L195)
- [coreApiWebdavUploadTUS/checksums.feature:196](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L196)
- [coreApiWebdavUploadTUS/checksums.feature:197](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L197)
- [coreApiWebdavUploadTUS/checksums.feature:240](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L240)
- [coreApiWebdavUploadTUS/checksums.feature:241](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L241)
- [coreApiWebdavUploadTUS/checksums.feature:242](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L242)
- [coreApiWebdavUploadTUS/checksums.feature:243](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L243)
- [coreApiWebdavUploadTUS/checksums.feature:244](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L244)
- [coreApiWebdavUploadTUS/checksums.feature:245](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L245)
- [coreApiWebdavUploadTUS/uploadToShare.feature:250](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L250)
- [coreApiWebdavUploadTUS/uploadToShare.feature:251](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L251)
- [coreApiWebdavUploadTUS/uploadToShare.feature:274](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L274)
- [coreApiWebdavUploadTUS/uploadToShare.feature:275](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L275)
- [coreApiWebdavUploadTUS/uploadToShare.feature:369](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L369)
- [coreApiWebdavUploadTUS/uploadToShare.feature:370](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L370)

#### [Renaming resource to banned name is allowed in spaces webdav](https://github.com/owncloud/ocis/issues/3099)

- [coreApiWebdavMove2/moveFile.feature:143](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L143)
- [coreApiWebdavMove1/moveFolder.feature:36](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L36)
- [coreApiWebdavMove1/moveFolder.feature:50](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L50)
- [coreApiWebdavMove1/moveFolder.feature:64](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L64)

#### [Trying to delete other user's trashbin item returns 409 for spaces path instead of 404](https://github.com/owncloud/ocis/issues/9791)

- [coreApiTrashbin/trashbinDelete.feature:92](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L92)

#### [MOVE a file into same folder with same name returns 404 instead of 403](https://github.com/owncloud/ocis/issues/1976)

- [coreApiWebdavMove2/moveFile.feature:100](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L100)
- [coreApiWebdavMove2/moveFile.feature:101](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L101)
- [coreApiWebdavMove2/moveFile.feature:102](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L102)
- [coreApiWebdavMove1/moveFolder.feature:217](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L217)
- [coreApiWebdavMove1/moveFolder.feature:218](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L218)
- [coreApiWebdavMove1/moveFolder.feature:219](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L219)
- [coreApiWebdavMove2/moveShareOnOcis.feature:307](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L307)
- [coreApiWebdavMove2/moveShareOnOcis.feature:310](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L310)

#### [COPY file/folder to same name is possible (but 500 code error for folder with spaces path)](https://github.com/owncloud/ocis/issues/8711)

- [coreApiSharePublicLink2/copyFromPublicLink.feature:198](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L198)
- [coreApiWebdavProperties/copyFile.feature:1070](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L1070)
- [coreApiWebdavProperties/copyFile.feature:1071](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L1071)
- [coreApiWebdavProperties/copyFile.feature:1072](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L1072)

#### [same href in REPORT request for all dav-path-version](https://github.com/owncloud/ocis/issues/7060)

- [coreApiWebdavOperations/search.feature:42](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L42)
- [coreApiWebdavOperations/search.feature:43](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L43)
- [coreApiWebdavOperations/search.feature:60](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L60)
- [coreApiWebdavOperations/search.feature:61](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L61)
- [coreApiWebdavOperations/search.feature:79](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L79)
- [coreApiWebdavOperations/search.feature:80](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L80)
- [coreApiWebdavOperations/search.feature:110](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L110)
- [coreApiWebdavOperations/search.feature:111](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L111)
- [coreApiWebdavOperations/search.feature:130](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L130)
- [coreApiWebdavOperations/search.feature:131](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L131)
- [coreApiWebdavOperations/search.feature:151](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L151)
- [coreApiWebdavOperations/search.feature:152](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L152)
- [coreApiWebdavOperations/search.feature:229](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L229)
- [coreApiWebdavOperations/search.feature:230](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L230)

#### [Trying to restore personal file to file of share received folder returns 403 but the share file is deleted (new dav path)](https://github.com/owncloud/ocis/issues/10356)

- [coreApiTrashbin/trashbinSharingToShares.feature:271](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L271)

### Won't fix

Not everything needs to be implemented for ocis. While the oc10 testsuite covers these things we are not looking at them right now.

- _The `OC-LazyOps` header is [no longer supported by the client](https://github.com/owncloud/client/pull/8398), implementing this is not necessary for a first production release. We plan to have an upload state machine to visualize the state of a file, see https://github.com/owncloud/ocis/issues/214_
- _Blacklisted ignored files are no longer required because ocis can handle `.htaccess` files without security implications introduced by serving user provided files with apache._

Note: always have an empty line at the end of this file.
The bash script that processes this file requires that the last line has a newline on the end.

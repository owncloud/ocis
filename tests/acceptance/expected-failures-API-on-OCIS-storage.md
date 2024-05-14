## Scenarios from ownCloud10 core API tests that are expected to fail with OCIS storage while running with the Graph API

The expected failures in this file are from features in the owncloud/ocis repo.

### File

Basic file management like up and download, move, copy, properties, trash, versions and chunking.

#### [copy personal space file to shared folder root result share in decline state](https://github.com/owncloud/ocis/issues/6999)

- [coreApiWebdavProperties/copyFile.feature:285](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L285)
- [coreApiWebdavProperties/copyFile.feature:286](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L286)

#### [Custom dav properties with namespaces are rendered incorrectly](https://github.com/owncloud/ocis/issues/2140)

_ocdav: double-check the webdav property parsing when custom namespaces are used_

- [coreApiWebdavProperties/setFileProperties.feature:36](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L36)
- [coreApiWebdavProperties/setFileProperties.feature:37](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L37)
- [coreApiWebdavProperties/setFileProperties.feature:42](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L42)
- [coreApiWebdavProperties/setFileProperties.feature:78](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L78)
- [coreApiWebdavProperties/setFileProperties.feature:77](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L77)
- [coreApiWebdavProperties/setFileProperties.feature:83](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/setFileProperties.feature#L83)

#### [Cannot set custom webDav properties](https://github.com/owncloud/product/issues/264)

- [coreApiWebdavProperties/getFileProperties.feature:342](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L342)
- [coreApiWebdavProperties/getFileProperties.feature:343](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L343)
- [coreApiWebdavProperties/getFileProperties.feature:348](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L348)
- [coreApiWebdavProperties/getFileProperties.feature:378](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L378)
- [coreApiWebdavProperties/getFileProperties.feature:379](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L379)
- [coreApiWebdavProperties/getFileProperties.feature:384](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getFileProperties.feature#L384)

#### [file versions do not report the version author](https://github.com/owncloud/ocis/issues/2914)

- [coreApiVersions/fileVersionAuthor.feature:15](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L15)
- [coreApiVersions/fileVersionAuthor.feature:43](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L43)
- [coreApiVersions/fileVersionAuthor.feature:68](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L68)
- [coreApiVersions/fileVersionAuthor.feature:92](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L92)
- [coreApiVersions/fileVersionAuthor.feature:123](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L123)
- [coreApiVersions/fileVersionAuthor.feature:148](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L148)
- [coreApiVersions/fileVersionAuthor.feature:177](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L177)
- [coreApiVersions/fileVersionAuthor.feature:208](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L208)
- [coreApiVersions/fileVersionAuthor.feature:254](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L254)
- [coreApiVersions/fileVersionAuthor.feature:297](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L297)
- [coreApiVersions/fileVersionAuthor.feature:316](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L316)

### Sync

Synchronization features like etag propagation, setting mtime and locking files

#### [Uploading an old method chunked file with checksum should fail using new DAV path](https://github.com/owncloud/ocis/issues/2323)

- [coreApiMain/checksums.feature:258](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiMain/checksums.feature#L258)
- [coreApiMain/checksums.feature:263](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiMain/checksums.feature#L263)

### Share

File and sync features in a shared scenario

#### [accepting matching name shared resources from different users/groups sets no serial identifiers on the resource name for the receiver](https://github.com/owncloud/ocis/issues/4289)

- [coreApiShareManagementToShares/acceptShares.feature:249](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L249)
- [coreApiShareManagementToShares/acceptShares.feature:269](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L269)
- [coreApiShareManagementToShares/acceptShares.feature:480](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L480)
- [coreApiShareManagementToShares/acceptShares.feature:546](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L546)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:39](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L39)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:40](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L40)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:126](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L126)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:127](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L127)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:159](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L159)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:160](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L160)

#### [file_target of an auto-renamed file is not correct directly after sharing](https://github.com/owncloud/core/issues/32322)

- [coreApiShareManagementToShares/mergeShare.feature:112](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/mergeShare.feature#L112)
- [coreApiShareManagementToShares/mergeShare.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/mergeShare.feature#L132)

#### [File deletion using dav gives unique string in filename in the trashbin](https://github.com/owncloud/product/issues/178)

- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:75](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L75)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:93](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L93)

cannot share a folder with create permission

#### [Resource with share permission create is readable for sharee](https://github.com/owncloud/ocis/issues/4524)

- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:155](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L155)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:172](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L172)

#### [Public link enforce permissions](https://github.com/owncloud/ocis/issues/1269)

- [coreApiSharePublicLink1/createPublicLinkShare.feature:290](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink1/createPublicLinkShare.feature#L290)

#### [copying a folder within a public link folder to folder with same name as an already existing file overwrites the parent file](https://github.com/owncloud/ocis/issues/1232)

- [coreApiSharePublicLink2/copyFromPublicLink.feature:67](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L67)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:93](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L93)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:177](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L177)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:178](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L178)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:193](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L193)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:194](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L194)

#### [d:quota-available-bytes in dprop of PROPFIND give wrong response value](https://github.com/owncloud/ocis/issues/8197)

- [coreApiWebdavProperties/getQuota.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L55)
- [coreApiWebdavProperties/getQuota.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L56)
- [coreApiWebdavProperties/getQuota.feature:57](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L57)
- [coreApiWebdavProperties/getQuota.feature:71](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L71)
- [coreApiWebdavProperties/getQuota.feature:72](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L72)
- [coreApiWebdavProperties/getQuota.feature:73](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/getQuota.feature#L73)

#### [deleting a file inside a received shared folder is moved to the trash-bin of the sharer not the receiver](https://github.com/owncloud/ocis/issues/1124)

- [coreApiTrashbin/trashbinSharingToShares.feature:29](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L29)
- [coreApiTrashbin/trashbinSharingToShares.feature:45](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L45)
- [coreApiTrashbin/trashbinSharingToShares.feature:50](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L50)
- [coreApiTrashbin/trashbinSharingToShares.feature:70](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L70)
- [coreApiTrashbin/trashbinSharingToShares.feature:75](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L75)
- [coreApiTrashbin/trashbinSharingToShares.feature:121](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L121)
- [coreApiTrashbin/trashbinSharingToShares.feature:126](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L126)
- [coreApiTrashbin/trashbinSharingToShares.feature:174](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L174)
- [coreApiTrashbin/trashbinSharingToShares.feature:179](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L179)
- [coreApiTrashbin/trashbinSharingToShares.feature:201](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L201)
- [coreApiTrashbin/trashbinSharingToShares.feature:224](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L224)

#### [Sharing folder and sub-folder with same user but different permission,the permission of sub-folder is not obeyed ](https://github.com/owncloud/ocis/issues/2440)

- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:213](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L213)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:239](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L239)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:292](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L292)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:318](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L318)

#### [Empty OCS response for a share create request using a disabled user](https://github.com/owncloud/ocis/issues/2212)

- [coreApiShareCreateSpecialToShares2/createShareWithDisabledUser.feature:21](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareWithDisabledUser.feature#L21)
- [coreApiShareCreateSpecialToShares2/createShareWithDisabledUser.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareWithDisabledUser.feature#L22)

#### [Edit user share response has a "name" field](https://github.com/owncloud/ocis/issues/1225)

- [coreApiShareUpdateToShares/updateShare.feature:161](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareUpdateToShares/updateShare.feature#L161)
- [coreApiShareUpdateToShares/updateShare.feature:162](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareUpdateToShares/updateShare.feature#L162)


#### [deleting a share with wrong authentication returns OCS status 996 / HTTP 500](https://github.com/owncloud/ocis/issues/1229)

- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:277](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L277)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:278](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L278)

### Other

API, search, favorites, config, capabilities, not existing endpoints, CORS and others

#### [Ability to return error messages in Webdav response bodies](https://github.com/owncloud/ocis/issues/1293)

- [coreApiAuth/ocsDELETEAuth.feature:7](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsDELETEAuth.feature#L7)
- [coreApiAuth/ocsGETAuth.feature:10](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsGETAuth.feature#L10)
- [coreApiAuth/ocsGETAuth.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsGETAuth.feature#L44)
- [coreApiAuth/ocsGETAuth.feature:75](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsGETAuth.feature#L75)
- [coreApiAuth/ocsGETAuth.feature:106](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsGETAuth.feature#L106)
- [coreApiAuth/ocsGETAuth.feature:123](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsGETAuth.feature#L123)
- [coreApiAuth/ocsPOSTAuth.feature:10](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsPOSTAuth.feature#L10)
- [coreApiAuth/ocsPUTAuth.feature:7](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsPUTAuth.feature#L7)
- [coreApiSharePublicLink1/createPublicLinkShare.feature:277](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink1/createPublicLinkShare.feature#L277)

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

#### [Using double slash in URL to access a folder gives 501 and other status codes](https://github.com/owncloud/ocis/issues/1667)

- [coreApiAuth/webDavSpecialURLs.feature:15](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L15)
- [coreApiAuth/webDavSpecialURLs.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L26)
- [coreApiAuth/webDavSpecialURLs.feature:78](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L78)
- [coreApiAuth/webDavSpecialURLs.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L90)
- [coreApiAuth/webDavSpecialURLs.feature:102](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L102)
- [coreApiAuth/webDavSpecialURLs.feature:113](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L113)
- [coreApiAuth/webDavSpecialURLs.feature:123](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L123)
- [coreApiAuth/webDavSpecialURLs.feature:134](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L134)
- [coreApiAuth/webDavSpecialURLs.feature:144](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L144)
- [coreApiAuth/webDavSpecialURLs.feature:155](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L155)
- [coreApiAuth/webDavSpecialURLs.feature:165](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L165)
- [coreApiAuth/webDavSpecialURLs.feature:176](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L176)
- [coreApiAuth/webDavSpecialURLs.feature:186](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L186)
- [coreApiAuth/webDavSpecialURLs.feature:197](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/webDavSpecialURLs.feature#L197)

#### [Difference in response content of status.php and default capabilities](https://github.com/owncloud/ocis/issues/1286)

- [coreApiCapabilities/capabilitiesWithNormalUser.feature:13](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiCapabilities/capabilitiesWithNormalUser.feature#L13)

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

#### [Request to edit non-existing user by authorized admin gets unauthorized in http response](https://github.com/owncloud/core/issues/38423)

- [coreApiAuth/ocsPUTAuth.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuth/ocsPUTAuth.feature#L22)

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
- [coreApiWebdavUploadTUS/uploadToShare.feature:166](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L166)
- [coreApiWebdavUploadTUS/uploadToShare.feature:167](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L167)
- [coreApiWebdavUploadTUS/uploadToShare.feature:184](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L184)
- [coreApiWebdavUploadTUS/uploadToShare.feature:185](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L185)
- [coreApiWebdavUploadTUS/uploadToShare.feature:202](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L202)
- [coreApiWebdavUploadTUS/uploadToShare.feature:203](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L203)
- [coreApiWebdavUploadTUS/uploadToShare.feature:239](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L239)
- [coreApiWebdavUploadTUS/uploadToShare.feature:240](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L240)
- [coreApiWebdavUploadTUS/uploadToShare.feature:279](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L279)
- [coreApiWebdavUploadTUS/uploadToShare.feature:280](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L280)

#### [TUS OPTIONS requests do not reply with TUS headers when invalid password](https://github.com/owncloud/ocis/issues/1012)

- [coreApiWebdavUploadTUS/optionsRequest.feature:40](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L40)
- [coreApiWebdavUploadTUS/optionsRequest.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L55)

#### [Trying to accept a share with invalid ID gives incorrect OCS and HTTP status](https://github.com/owncloud/ocis/issues/2111)

- [coreApiShareOperationsToShares2/shareAccessByID.feature:83](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L83)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:84](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L84)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:85](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L85)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:86](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L86)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:87](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L87)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:88](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L88)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L89)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L90)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:101](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L101)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:102](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L102)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:131](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L131)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L132)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:133](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L133)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:134](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L134)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L135)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:136](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L136)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:137](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L137)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:138](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L138)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:149](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L149)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:150](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L150)

#### [Shares to deleted group listed in the response](https://github.com/owncloud/ocis/issues/2441)

- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:502](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L502)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:503](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L503)

#### [copying the file inside Shares folder returns 404](https://github.com/owncloud/ocis/issues/3874)

- [coreApiWebdavProperties/copyFile.feature:398](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L398)
- [coreApiWebdavProperties/copyFile.feature:399](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L399)
- [coreApiWebdavProperties/copyFile.feature:404](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L404)
- [coreApiWebdavProperties/copyFile.feature:423](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L423)
- [coreApiWebdavProperties/copyFile.feature:424](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L424)
- [coreApiWebdavProperties/copyFile.feature:429](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L429)
- [coreApiWebdavProperties/copyFile.feature:268](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L268)
- [coreApiWebdavProperties/copyFile.feature:269](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L269)

### Won't fix

Not everything needs to be implemented for ocis. While the oc10 testsuite covers these things we are not looking at them right now.

- _The `OC-LazyOps` header is [no longer supported by the client](https://github.com/owncloud/client/pull/8398), implementing this is not necessary for a first production release. We plan to have an upload state machine to visualize the state of a file, see https://github.com/owncloud/ocis/issues/214_
- _Blacklisted ignored files are no longer required because ocis can handle `.htaccess` files without security implications introduced by serving user provided files with apache._

#### [Blacklist files extensions](https://github.com/owncloud/ocis/issues/2177)

- [coreApiWebdavProperties/copyFile.feature:115](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L115)
- [coreApiWebdavProperties/copyFile.feature:116](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L116)
- [coreApiWebdavProperties/copyFile.feature:121](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/copyFile.feature#L121)
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
- [coreApiShareManagementToShares/acceptShares.feature:139](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L139)
- [coreApiShareManagementToShares/acceptShares.feature:161](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L161)
- [coreApiShareManagementToShares/acceptShares.feature:190](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L190)
- [coreApiShareManagementToShares/acceptShares.feature:238](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L238)
- [coreApiShareManagementToShares/acceptShares.feature:457](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L457)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:119](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L119)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:120](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L120)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:213](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L213)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:214](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L214)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:215](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L215)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:216](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L216)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:237](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L237)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:238](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L238)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:239](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L239)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:240](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L240)

#### [Content-type is not multipart/byteranges when downloading file with Range Header](https://github.com/owncloud/ocis/issues/2677)

- [coreApiWebdavOperations/downloadFile.feature:183](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/downloadFile.feature#L183)
- [coreApiWebdavOperations/downloadFile.feature:184](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/downloadFile.feature#L184)
- [coreApiWebdavOperations/downloadFile.feature:189](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/downloadFile.feature#L189)

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

#### [Default capabilities for normal user and admin user not same as in oC-core](https://github.com/owncloud/ocis/issues/1285)

- [coreApiCapabilities/capabilities.feature:10](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiCapabilities/capabilities.feature#L10)
- [coreApiCapabilities/capabilities.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiCapabilities/capabilities.feature#L135)
- [coreApiCapabilities/capabilities.feature:174](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiCapabilities/capabilities.feature#L174)

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
- [coreApiWebdavMove2/moveShareOnOcis.feature:283](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L283)
- [coreApiWebdavMove2/moveShareOnOcis.feature:284](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L284)
- [coreApiWebdavMove2/moveShareOnOcis.feature:286](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L286)
- [coreApiWebdavMove2/moveShareOnOcis.feature:287](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L287)

Note: always have an empty line at the end of this file.
The bash script that processes this file requires that the last line has a newline on the end.

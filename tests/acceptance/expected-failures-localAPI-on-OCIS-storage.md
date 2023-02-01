## Scenarios from OCIS API tests that are expected to fail with OCIS storage
The expected failures in this file are from features in the owncloud/ocis repo.

#### [Downloading the archive of the resource (files | folder) using resource path is not possible](https://github.com/owncloud/ocis/issues/4637)
- [apiArchiver/downloadByPath.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L26)
- [apiArchiver/downloadByPath.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L27)
- [apiArchiver/downloadByPath.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L44)
- [apiArchiver/downloadByPath.feature:45](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L45)
- [apiArchiver/downloadByPath.feature:48](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L48)
- [apiArchiver/downloadByPath.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L74)
- [apiArchiver/downloadByPath.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L132)
- [apiArchiver/downloadByPath.feature:133](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L133)

### [Downloaded /Shares tar contains resource (files|folder) with leading / in Response](https://github.com/owncloud/ocis/issues/4636)
- [apiArchiver/downloadById.feature:134](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L134)
- [apiArchiver/downloadById.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L135)

### [create request for already existing user exits with status code 500 ](https://github.com/owncloud/ocis/issues/3516)
- [apiGraph/createGroupCaseSensitive.feature:17](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L17)
- [apiGraph/createGroupCaseSensitive.feature:18](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L18)
- [apiGraph/createGroupCaseSensitive.feature:19](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L19)
- [apiGraph/createGroupCaseSensitive.feature:20](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L20)
- [apiGraph/createGroupCaseSensitive.feature:21](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L21)
- [apiGraph/createGroupCaseSensitive.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L22)
- [apiGraph/createGroup.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroup.feature#L26)
- [apiGraph/createUser.feature:28](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createUser.feature#L28)

### [PROPFIND on accepted shares with identical names containing brackets exit with 404](https://github.com/owncloud/ocis/issues/4421)
- [apiSpacesShares/changingFilesShare.feature:12](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/changingFilesShare.feature#L12)

### [copy to overwrite (file and folder) from Personal to Shares Jail behaves differently](https://github.com/owncloud/ocis/issues/4393)
- [apiSpacesShares/copySpaces.feature:487](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/copySpaces.feature#L487)
- [apiSpacesShares/copySpaces.feature:501](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/copySpaces.feature#L501)

#### [PATCH request for TUS upload with wrong checksum gives incorrect response](https://github.com/owncloud/ocis/issues/1755)
- [apiSpacesShares/shareUploadTUS.feature:204](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L204)
- [apiSpacesShares/shareUploadTUS.feature:219](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L219)
- [apiSpacesShares/shareUploadTUS.feature:284](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L284)

### [Copy or move on an existing resource doesn't create a new version but deletes instead](https://github.com/owncloud/ocis/issues/4797)
- [apiSpacesShares/moveSpaces.feature:306](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/moveSpaces.feature#L306)
- [apiSpacesShares/copySpaces.feature:710](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/copySpaces.feature#L710)
- [apiSpacesShares/copySpaces.feature:751](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/copySpaces.feature#L751)

### [Creating group with empty name returns status code 200](https://github.com/owncloud/ocis/issues/5050)
- [apiGraph/createGroup.feature:40](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroup.feature#L40)

### [Settings service user can list other peoples assignments](https://github.com/owncloud/ocis/issues/5032)
- [apiAccountsHashDifficulty/assignRole.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAccountsHashDifficulty/assignRole.feature#L27)
- [apiAccountsHashDifficulty/assignRole.feature:28](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAccountsHashDifficulty/assignRole.feature#L28)

### [Settings service user can see roles list](https://github.com/owncloud/ocis/issues/5079)
- [apiAccountsHashDifficulty/assignRole.feature:15](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAccountsHashDifficulty/assignRole.feature#L15)
- [apiAccountsHashDifficulty/assignRole.feature:16](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAccountsHashDifficulty/assignRole.feature#L16)

### [Group having percentage (%) can be created but cannot be GET](https://github.com/owncloud/ocis/issues/5083)
- [apiGraph/deleteGroup.feature:49](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/deleteGroup.feature#L49)
- [apiGraph/deleteGroup.feature:50](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/deleteGroup.feature#L50)
- [apiGraph/deleteGroup.feature:51](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/deleteGroup.feature#L51)

#### [Share lists deleted user as 'user'](https://github.com/owncloud/ocis/issues/903)
- [apiGraph/deleteGroup.feature:62](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/deleteGroup.feature#L62)

#### [Updating group displayName request seems OK but group is not being renamed](https://github.com/owncloud/ocis/issues/5099)
- [apiGraph/editGroup.feature:20](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L20)
- [apiGraph/editGroup.feature:21](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L21)
- [apiGraph/editGroup.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L22)
- [apiGraph/editGroup.feature:23](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L23)
- [apiGraph/editGroup.feature:24](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L24)
- [apiGraph/editGroup.feature:25](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L25)

#### [CORS headers are not identical with oC10 headers](https://github.com/owncloud/ocis/issues/5195)
- [apiCors/cors.feature:25](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L25)
- [apiCors/cors.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L26)
- [apiCors/cors.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L27)
- [apiCors/cors.feature:28](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L28)

#### [Requests with invalid credentials do not return CORS headers](https://github.com/owncloud/ocis/issues/5194)
- [apiCors/cors.feature:67](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L67)
- [apiCors/cors.feature:68](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L68)

#### [POST response does not return correct path when creating public link](https://github.com/owncloud/ocis/issues/5139)
- [apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature:62](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature#L62)
- [apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature:63](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature#L63)
- [apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature:64](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature#L64)
- [apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature#L90)
- [apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature:160](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature#L160)
- [apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature:161](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature#L161)
- [apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature:162](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareSubItemOfSpaceViaPublicLink.feature#L162)

#### [Public cannot download folder via the public link of the folder inside the project space](https://github.com/owncloud/ocis/issues/5229)
- [apiSpacesShares/publicLinkDownload.feature:30](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/publicLinkDownload.feature#L30)

#### [A User can get information of another user with Graph API](https://github.com/owncloud/ocis/issues/5125)
- [apiGraph/getUser.feature:23](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L23)
- [apiGraph/getUser.feature:92](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L92)

#### [GET a file while it's in processing doesn't return 425 code (async uploads)](https://github.com/owncloud/ocis/issues/5326)
- [apiAsyncUpload/delayPostprocessing.feature:14](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAsyncUpload/delayPostprocessing.feature#L14)
- [apiAsyncUpload/delayPostprocessing.feature:15](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAsyncUpload/delayPostprocessing.feature#L15)
- [apiAsyncUpload/delayPostprocessing.feature:16](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAsyncUpload/delayPostprocessing.feature#L16)

#### [Sharing to a group with an expiration date does not work #5442](https://github.com/owncloud/ocis/issues/5442)
- [apiSpacesShares/shareSubItemOfSpace.feature:99](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareSubItemOfSpace.feature#L99)

Note: always have an empty line at the end of this file.
The bash script that processes this file requires that the last line has a newline on the end.

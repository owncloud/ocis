## Scenarios from OCIS API tests that are expected to fail with OCIS storage

The expected failures in this file are from features in the owncloud/ocis repo.

#### [Downloading the archive of the resource (files | folder) using resource path is not possible](https://github.com/owncloud/ocis/issues/4637)

- [apiArchiver/downloadByPath.feature:25](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L25)
- [apiArchiver/downloadByPath.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L26)
- [apiArchiver/downloadByPath.feature:43](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L43)
- [apiArchiver/downloadByPath.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L44)
- [apiArchiver/downloadByPath.feature:47](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L47)
- [apiArchiver/downloadByPath.feature:73](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L73)
- [apiArchiver/downloadByPath.feature:131](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L131)
- [apiArchiver/downloadByPath.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L132)

### [Downloaded /Shares tar contains resource (files|folder) with leading / in Response](https://github.com/owncloud/ocis/issues/4636)

- [apiArchiver/downloadById.feature:133](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L133)
- [apiArchiver/downloadById.feature:134](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L134)

### [PROPFIND on accepted shares with identical names containing brackets exit with 404](https://github.com/owncloud/ocis/issues/4421)

- [apiSpacesShares/changingFilesShare.feature:14](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/changingFilesShare.feature#L14)

### [Shared mount folder gets deleted when overwritten by a file from personal space](https://github.com/owncloud/ocis/issues/7208)

- [apiSpacesShares/copySpaces.feature:528](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/copySpaces.feature#L528)
- [apiSpacesShares/copySpaces.feature:542](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/copySpaces.feature#L542)

#### [PATCH request for TUS upload with wrong checksum gives incorrect response](https://github.com/owncloud/ocis/issues/1755)

- [apiSpacesShares/shareUploadTUS.feature:203](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L203)
- [apiSpacesShares/shareUploadTUS.feature:218](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L218)
- [apiSpacesShares/shareUploadTUS.feature:283](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L283)

### [Settings service user can list other peoples assignments](https://github.com/owncloud/ocis/issues/5032)

- [apiAccountsHashDifficulty/assignRole.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAccountsHashDifficulty/assignRole.feature#L27)
- [apiAccountsHashDifficulty/assignRole.feature:28](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiAccountsHashDifficulty/assignRole.feature#L28)
- [apiGraph/assignRole.feature:30](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/assignRole.feature#L30)
- [apiGraph/assignRole.feature:31](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/assignRole.feature#L31)
- [apiGraph/assignRole.feature:32](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/assignRole.feature#L32)

#### [Share lists deleted user as 'user'](https://github.com/owncloud/ocis/issues/903)

- [apiGraph/deleteGroup.feature:67](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/deleteGroup.feature#L67)

#### [CORS headers are not identical with oC10 headers](https://github.com/owncloud/ocis/issues/5195)

- [apiCors/cors.feature:28](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L28)
- [apiCors/cors.feature:29](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L29)
- [apiCors/cors.feature:30](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L30)
- [apiCors/cors.feature:31](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L31)

#### [Requests with invalid credentials do not return CORS headers](https://github.com/owncloud/ocis/issues/5194)

- [apiCors/cors.feature:70](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L70)
- [apiCors/cors.feature:71](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiCors/cors.feature#L71)

#### [A User can get information of another user with Graph API](https://github.com/owncloud/ocis/issues/5125)

- [apiGraph/getUser.feature:82](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L82)
- [apiGraph/getUser.feature:83](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L83)
- [apiGraph/getUser.feature:84](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L84)
- [apiGraph/getUser.feature:85](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L85)
- [apiGraph/getUser.feature:86](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L86)
- [apiGraph/getUser.feature:87](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L87)
- [apiGraph/getUser.feature:88](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L88)
- [apiGraph/getUser.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L89)
- [apiGraph/getUser.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L90)
- [apiGraph/getUser.feature:91](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L91)
- [apiGraph/getUser.feature:92](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L92)
- [apiGraph/getUser.feature:93](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L93)
- [apiGraph/getUser.feature:607](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L607)
- [apiGraph/getUser.feature:608](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L608)
- [apiGraph/getUser.feature:609](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L609)
- [apiGraph/getUser.feature:610](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L610)
- [apiGraph/getUser.feature:611](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L611)
- [apiGraph/getUser.feature:612](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L612)
- [apiGraph/getUser.feature:613](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L613)
- [apiGraph/getUser.feature:614](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L614)
- [apiGraph/getUser.feature:615](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L615)
- [apiGraph/getUser.feature:616](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L616)
- [apiGraph/getUser.feature:617](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L617)
- [apiGraph/getUser.feature:618](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getUser.feature#L618)

#### [Normal user can get expanded members information of a group](https://github.com/owncloud/ocis/issues/5604)

- [apiGraph/getGroup.feature:381](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L381)
- [apiGraph/getGroup.feature:382](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L382)
- [apiGraph/getGroup.feature:383](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L383)

#### [Changing user with an uppercase name gives 404 error](https://github.com/owncloud/ocis/issues/7044)

- [apiGraph/editUser.feature:67](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editUser.feature#L67)

#### [Same users can be added in a group multiple time](https://github.com/owncloud/ocis/issues/5702)

- [apiGraph/addUserToGroup.feature:285](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L285)

#### [API requests from an unauthorized user should return 403](https://github.com/owncloud/ocis/issues/5938)

- [apiGraph/addUserToGroup.feature:150](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L150)
- [apiGraph/addUserToGroup.feature:151](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L151)
- [apiGraph/addUserToGroup.feature:152](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L152)
- [apiGraph/addUserToGroup.feature:184](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L184)
- [apiGraph/addUserToGroup.feature:185](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L185)
- [apiGraph/addUserToGroup.feature:186](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L186)
- [apiGraph/createGroup.feature:42](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroup.feature#L42)
- [apiGraph/createGroup.feature:43](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroup.feature#L43)
- [apiGraph/createGroup.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroup.feature#L44)
- [apiGraph/deleteGroup.feature:63](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/deleteGroup.feature#L63)
- [apiGraph/deleteGroup.feature:62](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/deleteGroup.feature#L62)
- [apiGraph/deleteGroup.feature:64](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/deleteGroup.feature#L64)
- [apiGraph/editGroup.feature:35](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L35)
- [apiGraph/editGroup.feature:34](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L34)
- [apiGraph/editGroup.feature:36](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/editGroup.feature#L36)
- [apiGraph/getGroup.feature:54](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L54)
- [apiGraph/getGroup.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L55)
- [apiGraph/getGroup.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L56)
- [apiGraph/getGroup.feature:103](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L103)
- [apiGraph/getGroup.feature:104](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L104)
- [apiGraph/getGroup.feature:105](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L105)
- [apiGraph/getGroup.feature:267](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L267)
- [apiGraph/getGroup.feature:268](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L268)
- [apiGraph/getGroup.feature:269](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/getGroup.feature#L269)
- [apiGraph/removeUserFromGroup.feature:191](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/removeUserFromGroup.feature#L191)
- [apiGraph/removeUserFromGroup.feature:192](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/removeUserFromGroup.feature#L192)
- [apiGraph/removeUserFromGroup.feature:193](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/removeUserFromGroup.feature#L193)

#### [API requests for a non-existent resources should return 404](https://github.com/owncloud/ocis/issues/5939)

- [apiGraph/addUserToGroup.feature:201](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L201)
- [apiGraph/addUserToGroup.feature:202](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L202)
- [apiGraph/addUserToGroup.feature:203](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L203)

### [Users are added in a group with wrong host in host-part of user](https://github.com/owncloud/ocis/issues/5871)

- [apiGraph/addUserToGroup.feature:369](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L369)
- [apiGraph/addUserToGroup.feature:383](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L383)

### [Adding the same user as multiple members in a single request results in listing the same user twice in the group](https://github.com/owncloud/ocis/issues/5855)

- [apiGraph/addUserToGroup.feature:420](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/addUserToGroup.feature#L420)

Note: always have an empty line at the end of this file.
The bash script that processes this file requires that the last line has a newline on the end.

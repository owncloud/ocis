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

### [search doesn't find the project space by name](https://github.com/owncloud/ocis/issues/4506)
- [apiSpaces/search.feature:96](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/search.feature#L96)

#### [PATCH request for TUS upload with wrong checksum gives incorrect response](https://github.com/owncloud/ocis/issues/1755)
- [apiSpacesShares/shareUploadTUS.feature:204](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L204)
- [apiSpacesShares/shareUploadTUS.feature:219](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L219)
- [apiSpacesShares/shareUploadTUS.feature:284](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/shareUploadTUS.feature#L284)

### [Copy or move on an existing resource doesn't create a new version but deletes instead](https://github.com/owncloud/ocis/issues/4797)
- [apiSpacesShares/moveSpaces.feature:306](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/moveSpaces.feature#L306)
- [apiSpacesShares/copySpaces.feature:710](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/copySpaces.feature#L710)
- [apiSpacesShares/copySpaces.feature:748](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpacesShares/copySpaces.feature#L748)

### [Creating group with empty name returns status code 200](https://github.com/owncloud/ocis/issues/5050)
- [apiGraph/createGroup.feature:40](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroup.feature#L40)

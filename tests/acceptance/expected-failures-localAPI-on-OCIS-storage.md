## Scenarios from OCIS API tests that are expected to fail with OCIS storage
The expected failures in this file are from features in the owncloud/ocis repo.

#### [downloading an archive with invalid path returns HTTP/500](https://github.com/owncloud/ocis/issues/2768)
- [apiArchiver/downloadByPath.feature:69](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L69)

#### [Hardcoded call to /home/..., but /home no longer exists](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature)
- [apiArchiver/downloadByPath.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L26)
- [apiArchiver/downloadByPath.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L27)
- [apiArchiver/downloadByPath.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L44)
- [apiArchiver/downloadByPath.feature:45](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L45)
- [apiArchiver/downloadByPath.feature:48](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L48)
- [apiArchiver/downloadByPath.feature:69](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L69)
- [apiArchiver/downloadByPath.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L74)
- [apiArchiver/downloadByPath.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L132)
- [apiArchiver/downloadByPath.feature:133](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L133)

### Tries to download /Shares/ folder but it cannot be downloaded any more directly
- [apiArchiver/downloadById.feature:134](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L134)
- [apiArchiver/downloadById.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L135)

### [Search by shares jail works incorrect](https://github.com/owncloud/ocis/issues/4014)
- [apiSpaces/search.feature:43](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/search.feature#L43)

### [create request for already existing user exits with status code 500 ](https://github.com/owncloud/ocis/issues/3516)
- [apiGraph/createGroupCaseSensitive.feature:16](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L16)
- [apiGraph/createGroupCaseSensitive.feature:17](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L17)
- [apiGraph/createGroupCaseSensitive.feature:18](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L18)
- [apiGraph/createGroupCaseSensitive.feature:19](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L19)
- [apiGraph/createGroupCaseSensitive.feature:20](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L20)
- [apiGraph/createGroupCaseSensitive.feature:21](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiGraph/createGroupCaseSensitive.feature#L21)

### [PROPFIND on accepted shares with identical names containing brackets exit with 404](https://github.com/owncloud/ocis/issues/4421)

- [apiSpaces/changingFilesShare.feature:12](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/changingFilesShare.feature#L12)

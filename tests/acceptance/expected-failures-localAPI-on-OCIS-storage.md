## Scenarios from OCIS API tests that are expected to fail with OCIS storage

#### [downloading an archive with invalid path returns HTTP/500](https://github.com/owncloud/ocis/issues/2768)
-   [apiArchiver/downloadByPath.feature:69]

#### [Hardcoded call to /home/..., but /home no longer exists](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/apiArchiver/downloadByPath.feature#L26)
-   [apiArchiver/downloadByPath.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L26)
-   [apiArchiver/downloadByPath.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L27)
-   [apiArchiver/downloadByPath.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L44)
-   [apiArchiver/downloadByPath.feature:45](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L45)
-   [apiArchiver/downloadByPath.feature:48](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L48)
-   [apiArchiver/downloadByPath.feature:69](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L69)
-   [apiArchiver/downloadByPath.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L74)
-   [apiArchiver/downloadByPath.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L132)
-   [apiArchiver/downloadByPath.feature:133](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L133)

### Tries to download /Shares/ folder but it cannot be downloaded any more directly
- [apiArchiver/downloadById.feature:134](https://github.com/owncloud/web/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L134)
- [apiArchiver/downloadById.feature:135](https://github.com/owncloud/web/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L135)

#### [Overwriting a file in the space within the allowed quota does not work](https://github.com/owncloud/ocis/issues/2829)
- [apiSpaces/quota.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/quota.feature#L56)


## Scenarios from OCIS API tests that are expected to fail with OCIS storage

#### [downloading the /Shares folder using the archiver endpoint does not work](https://github.com/owncloud/ocis/issues/2751)
-   [apiArchiver/downloadById.feature:134](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L134)
-   [apiArchiver/downloadById.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadById.feature#L135)

#### [downloading an archive with invalid path returns HTTP/500](https://github.com/owncloud/ocis/issues/2768)
-   [apiArchiver/downloadByPath.feature:69](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/apiArchiver/downloadByPath.feature#L69)

#### [downloading an archive with non existing / accessible id returns HTTP/500](https://github.com/owncloud/ocis/issues/2795)
- [apiArchiver/downloadById.feature:69](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiArchiver/downloadByPath.feature#L69)

#### [Overwriting a file in the space within the allowed quota does not work](https://github.com/owncloud/ocis/issues/2829)
- [apiSpaces/quota.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/apiSpaces/quota.feature#L56)

## Scenarios that are expected to fail when remote.php is not used

#### [Trying to create .. resource with /webdav root (old dav path) without remote.php returns html](https://github.com/owncloud/ocis/issues/10339)

- [coreApiWebdavProperties/createFileFolder.feature:176](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/createFileFolder.feature#L176)
- [coreApiWebdavProperties/createFileFolder.feature:177](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/createFileFolder.feature#L177)
- [coreApiWebdavProperties/createFileFolder.feature:196](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/createFileFolder.feature#L196)
- [coreApiWebdavProperties/createFileFolder.feature:197](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties/createFileFolder.feature#L197)
- [coreApiWebdavUploadTUS/uploadFile.feature:177](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadFile.feature#L177)

Note: always have an empty line at the end of this file.
The bash script that processes this file requires that the last line has a newline on the end.

## Scenarios from ownCloud10 core API tests that are expected to fail with OCIS storage while running with the Graph API

The expected failures in this file are from features in the owncloud/ocis repo.

### File

Basic file management like up and download, move, copy, properties, trash, versions and chunking.

#### [Getting information about a folder overwritten by a file gives 500 error instead of 404](https://github.com/owncloud/ocis/issues/1239)

- [coreApiWebdavProperties1/copyFile.feature:273](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L273)
- [coreApiWebdavProperties1/copyFile.feature:274](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L274)
- [coreApiWebdavProperties1/copyFile.feature:291](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L291)
- [coreApiWebdavProperties1/copyFile.feature:292](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L292)

#### [Custom dav properties with namespaces are rendered incorrectly](https://github.com/owncloud/ocis/issues/2140)

_ocdav: double-check the webdav property parsing when custom namespaces are used_

- [coreApiWebdavProperties1/setFileProperties.feature:37](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/setFileProperties.feature#L37)
- [coreApiWebdavProperties1/setFileProperties.feature:38](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/setFileProperties.feature#L38)
- [coreApiWebdavProperties1/setFileProperties.feature:43](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/setFileProperties.feature#L43)
- [coreApiWebdavProperties1/setFileProperties.feature:78](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/setFileProperties.feature#L78)
- [coreApiWebdavProperties1/setFileProperties.feature:79](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/setFileProperties.feature#L79)
- [coreApiWebdavProperties1/setFileProperties.feature:84](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/setFileProperties.feature#L84)

#### [Cannot set custom webDav properties](https://github.com/owncloud/product/issues/264)

- [coreApiWebdavProperties2/getFileProperties.feature:341](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties2/getFileProperties.feature#L341)
- [coreApiWebdavProperties2/getFileProperties.feature:346](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties2/getFileProperties.feature#L346)
- [coreApiWebdavProperties2/getFileProperties.feature:351](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties2/getFileProperties.feature#L351)
- [coreApiWebdavProperties2/getFileProperties.feature:381](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties2/getFileProperties.feature#L381)
- [coreApiWebdavProperties2/getFileProperties.feature:386](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties2/getFileProperties.feature#L386)
- [coreApiWebdavProperties2/getFileProperties.feature:391](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties2/getFileProperties.feature#L391)

#### [file versions do not report the version author](https://github.com/owncloud/ocis/issues/2914)

- [coreApiVersions/fileVersionAuthor.feature:13](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L13)
- [coreApiVersions/fileVersionAuthor.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L44)
- [coreApiVersions/fileVersionAuthor.feature:71](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L71)
- [coreApiVersions/fileVersionAuthor.feature:97](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L97)
- [coreApiVersions/fileVersionAuthor.feature:130](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L130)
- [coreApiVersions/fileVersionAuthor.feature:157](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L157)
- [coreApiVersions/fileVersionAuthor.feature:188](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L188)
- [coreApiVersions/fileVersionAuthor.feature:223](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L223)
- [coreApiVersions/fileVersionAuthor.feature:275](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L275)
- [coreApiVersions/fileVersionAuthor.feature:324](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L324)
- [coreApiVersions/fileVersionAuthor.feature:345](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L345)
- [coreApiVersions/fileVersionAuthor.feature:365](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L365)
- [coreApiVersions/fileVersionAuthor.feature:390](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersionAuthor.feature#L390)

### Sync

Synchronization features like etag propagation, setting mtime and locking files

#### [Uploading an old method chunked file with checksum should fail using new DAV path](https://github.com/owncloud/ocis/issues/2323)

- [coreApiMain/checksums.feature:266](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiMain/checksums.feature#L266)
- [coreApiMain/checksums.feature:271](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiMain/checksums.feature#L271)

#### [Webdav LOCK operations](https://github.com/owncloud/ocis/issues/1284)

- [coreApiWebdavLocks/exclusiveLocks.feature:46](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L46)
- [coreApiWebdavLocks/exclusiveLocks.feature:47](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L47)
- [coreApiWebdavLocks/exclusiveLocks.feature:48](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L48)
- [coreApiWebdavLocks/exclusiveLocks.feature:49](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L49)
- [coreApiWebdavLocks/exclusiveLocks.feature:54](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L54)
- [coreApiWebdavLocks/exclusiveLocks.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L55)
- [coreApiWebdavLocks/exclusiveLocks.feature:73](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L73)
- [coreApiWebdavLocks/exclusiveLocks.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L74)
- [coreApiWebdavLocks/exclusiveLocks.feature:75](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L75)
- [coreApiWebdavLocks/exclusiveLocks.feature:76](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L76)
- [coreApiWebdavLocks/exclusiveLocks.feature:81](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L81)
- [coreApiWebdavLocks/exclusiveLocks.feature:82](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L82)
- [coreApiWebdavLocks/exclusiveLocks.feature:100](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L100)
- [coreApiWebdavLocks/exclusiveLocks.feature:101](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L101)
- [coreApiWebdavLocks/exclusiveLocks.feature:102](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L102)
- [coreApiWebdavLocks/exclusiveLocks.feature:103](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L103)
- [coreApiWebdavLocks/exclusiveLocks.feature:108](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L108)
- [coreApiWebdavLocks/exclusiveLocks.feature:109](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/exclusiveLocks.feature#L109)
- [coreApiWebdavLocks/requestsWithToken.feature:29](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/requestsWithToken.feature#L29)
- [coreApiWebdavLocks/requestsWithToken.feature:30](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/requestsWithToken.feature#L30)
- [coreApiWebdavLocks/requestsWithToken.feature:35](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks/requestsWithToken.feature#L35)
- [coreApiWebdavLocks2/independentLocks.feature:23](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L23)
- [coreApiWebdavLocks2/independentLocks.feature:24](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L24)
- [coreApiWebdavLocks2/independentLocks.feature:25](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L25)
- [coreApiWebdavLocks2/independentLocks.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L26)
- [coreApiWebdavLocks2/independentLocks.feature:31](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L31)
- [coreApiWebdavLocks2/independentLocks.feature:32](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L32)
- [coreApiWebdavLocks2/independentLocks.feature:51](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L51)
- [coreApiWebdavLocks2/independentLocks.feature:52](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L52)
- [coreApiWebdavLocks2/independentLocks.feature:53](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L53)
- [coreApiWebdavLocks2/independentLocks.feature:54](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L54)
- [coreApiWebdavLocks2/independentLocks.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L55)
- [coreApiWebdavLocks2/independentLocks.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L56)
- [coreApiWebdavLocks2/independentLocks.feature:57](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L57)
- [coreApiWebdavLocks2/independentLocks.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L58)
- [coreApiWebdavLocks2/independentLocks.feature:63](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L63)
- [coreApiWebdavLocks2/independentLocks.feature:64](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L64)
- [coreApiWebdavLocks2/independentLocks.feature:65](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L65)
- [coreApiWebdavLocks2/independentLocks.feature:66](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocks.feature#L66)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:29](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L29)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:30](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L30)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:31](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L31)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:32](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L32)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:37](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L37)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:38](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L38)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L58)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:59](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L59)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:60](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L60)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:61](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L61)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:66](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L66)
- [coreApiWebdavLocks2/independentLocksShareToShares.feature:67](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocks2/independentLocksShareToShares.feature#L67)
- [coreApiWebdavLocksUnlock/unlock.feature:20](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L20)
- [coreApiWebdavLocksUnlock/unlock.feature:21](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L21)
- [coreApiWebdavLocksUnlock/unlock.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L26)
- [coreApiWebdavLocksUnlock/unlock.feature:40](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L40)
- [coreApiWebdavLocksUnlock/unlock.feature:41](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L41)
- [coreApiWebdavLocksUnlock/unlock.feature:63](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L63)
- [coreApiWebdavLocksUnlock/unlock.feature:64](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L64)
- [coreApiWebdavLocksUnlock/unlock.feature:65](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L65)
- [coreApiWebdavLocksUnlock/unlock.feature:66](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L66)
- [coreApiWebdavLocksUnlock/unlock.feature:71](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L71)
- [coreApiWebdavLocksUnlock/unlock.feature:72](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlock.feature#L72)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:26](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L26)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L27)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:28](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L28)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:29](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L29)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:34](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L34)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:35](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L35)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:50](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L50)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:51](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L51)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:52](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L52)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:53](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L53)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L58)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:59](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L59)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:74](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L74)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:75](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L75)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:76](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L76)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:77](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L77)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:82](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L82)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:83](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L83)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:98](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L98)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:99](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L99)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:100](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L100)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:101](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L101)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:106](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L106)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:107](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L107)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:122](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L122)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:123](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L123)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:124](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L124)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:125](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L125)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:130](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L130)
- [coreApiWebdavLocksUnlock/unlockSharingToShares.feature:131](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavLocksUnlock/unlockSharingToShares.feature#L131)

### Share

File and sync features in a shared scenario

#### [ocs sharing api always returns an empty exact list while searching for a sharee](https://github.com/owncloud/ocis/issues/4265)

- [coreApiSharees/sharees.feature:98](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L98)
- [coreApiSharees/sharees.feature:99](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L99)
- [coreApiSharees/sharees.feature:118](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L118)
- [coreApiSharees/sharees.feature:119](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L119)
- [coreApiSharees/sharees.feature:138](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L138)
- [coreApiSharees/sharees.feature:139](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L139)
- [coreApiSharees/sharees.feature:158](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L158)
- [coreApiSharees/sharees.feature:159](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L159)
- [coreApiSharees/sharees.feature:178](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L178)
- [coreApiSharees/sharees.feature:179](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharees/sharees.feature#L179)


#### [accepting matching name shared resources from different users/groups sets no serial identifiers on the resource name for the receiver](https://github.com/owncloud/ocis/issues/4289)

- [coreApiShareManagementToShares/acceptShares.feature:285](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L285)
- [coreApiShareManagementToShares/acceptShares.feature:315](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L315)
- [coreApiShareManagementToShares/acceptShares.feature:229](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L229)
- [coreApiShareManagementToShares/acceptShares.feature:494](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L494)
- [coreApiShareManagementToShares/acceptShares.feature:559](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L559)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:128](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L128)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:129](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L129)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:163](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L163)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:164](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L164)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:37](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L37)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:38](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L38)

#### [sharing the shares folder to users exits with different status code than in oc10 backend](https://github.com/owncloud/ocis/issues/2215)

- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:741](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L741)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:742](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L742)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:760](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L760)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:761](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L761)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:776](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L776)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:777](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L777)

#### [file_target of an auto-renamed file is not correct directly after sharing](https://github.com/owncloud/core/issues/32322)

- [coreApiShareManagementToShares/mergeShare.feature:115](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/mergeShare.feature#L115)

#### [File deletion using dav gives unique string in filename in the trashbin](https://github.com/owncloud/product/issues/178)

- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:62](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L62)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:76](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L76)

cannot share a folder with create permission

#### [Resource with share permission create is readable for sharee](https://github.com/owncloud/ocis/issues/4524)

- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:129](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L129)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:142](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L142)


#### [OCS error message for attempting to access share via share id as an unauthorized user is not informative](https://github.com/owncloud/ocis/issues/1233)

- [coreApiShareOperationsToShares1/gettingShares.feature:150](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/gettingShares.feature#L150)
- [coreApiShareOperationsToShares1/gettingShares.feature:151](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/gettingShares.feature#L151)

#### [Listing shares via ocs API does not show path for parent folders](https://github.com/owncloud/ocis/issues/1231)

- [coreApiShareOperationsToShares1/gettingShares.feature:187](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/gettingShares.feature#L187)
- [coreApiShareOperationsToShares1/gettingShares.feature:188](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/gettingShares.feature#L188)

#### [Public link enforce permissions](https://github.com/owncloud/ocis/issues/1269)
- [coreApiSharePublicLink1/createPublicLinkShare.feature:353](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink1/createPublicLinkShare.feature#L353)

#### [download previews of other users file](https://github.com/owncloud/ocis/issues/2071)

- [coreApiWebdavPreviews/previews.feature:92](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavPreviews/previews.feature#L92)

#### [different error message detail for previews of folder](https://github.com/owncloud/ocis/issues/2064)

- [coreApiWebdavPreviews/previews.feature:101](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavPreviews/previews.feature#L101)

#### [copying a folder within a public link folder to folder with same name as an already existing file overwrites the parent file](https://github.com/owncloud/ocis/issues/1232)

- [coreApiSharePublicLink2/copyFromPublicLink.feature:63](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L63)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L89)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:173](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L173)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:174](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L174)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:189](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L189)
- [coreApiSharePublicLink2/copyFromPublicLink.feature:190](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink2/copyFromPublicLink.feature#L190)

#### [Upload-only shares must not overwrite but create a separate file](https://github.com/owncloud/ocis/issues/1267)

- [coreApiSharePublicLink3/uploadToPublicLinkShare.feature:10](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink3/uploadToPublicLinkShare.feature#L10)
- [coreApiSharePublicLink3/uploadToPublicLinkShare.feature:131](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink3/uploadToPublicLinkShare.feature#L131)

#### [Set quota over settings](https://github.com/owncloud/ocis/issues/1290)

- [coreApiSharePublicLink3/uploadToPublicLinkShare.feature:84](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink3/uploadToPublicLinkShare.feature#L84)
- [coreApiSharePublicLink3/uploadToPublicLinkShare.feature:93](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink3/uploadToPublicLinkShare.feature#L93)


#### [deleting a file inside a received shared folder is moved to the trash-bin of the sharer not the receiver](https://github.com/owncloud/ocis/issues/1124)

- [coreApiTrashbin/trashbinSharingToShares.feature:28](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L28)
- [coreApiTrashbin/trashbinSharingToShares.feature:45](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L45)
- [coreApiTrashbin/trashbinSharingToShares.feature:50](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L50)
- [coreApiTrashbin/trashbinSharingToShares.feature:72](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L72)
- [coreApiTrashbin/trashbinSharingToShares.feature:77](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L77)
- [coreApiTrashbin/trashbinSharingToShares.feature:99](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L99)
- [coreApiTrashbin/trashbinSharingToShares.feature:104](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L104)
- [coreApiTrashbin/trashbinSharingToShares.feature:127](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L127)
- [coreApiTrashbin/trashbinSharingToShares.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L132)
- [coreApiTrashbin/trashbinSharingToShares.feature:155](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L155)
- [coreApiTrashbin/trashbinSharingToShares.feature:160](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L160)
- [coreApiTrashbin/trashbinSharingToShares.feature:183](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L183)
- [coreApiTrashbin/trashbinSharingToShares.feature:188](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L188)
- [coreApiTrashbin/trashbinSharingToShares.feature:211](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L211)
- [coreApiTrashbin/trashbinSharingToShares.feature:235](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinSharingToShares.feature#L235)

#### [changing user quota gives ocs status 103 / Cannot set quota](https://github.com/owncloud/product/issues/247)

- [coreApiShareOperationsToShares2/uploadToShare.feature:209](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/uploadToShare.feature#L209)
- [coreApiShareOperationsToShares2/uploadToShare.feature:210](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/uploadToShare.feature#L210)

#### [not possible to move file into a received folder](https://github.com/owncloud/ocis/issues/764)

- [coreApiShareOperationsToShares1/changingFilesShare.feature:24](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/changingFilesShare.feature#L24)
- [coreApiShareOperationsToShares1/changingFilesShare.feature:25](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/changingFilesShare.feature#L25)
- [coreApiShareOperationsToShares1/changingFilesShare.feature:68](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/changingFilesShare.feature#L68)
- [coreApiShareOperationsToShares1/changingFilesShare.feature:69](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/changingFilesShare.feature#L69)
- [coreApiShareOperationsToShares1/changingFilesShare.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/changingFilesShare.feature#L90)
- [coreApiShareOperationsToShares1/changingFilesShare.feature:91](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares1/changingFilesShare.feature#L91)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:532](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L532)
- [coreApiWebdavMove2/moveShareOnOcis.feature:30](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L30)
- [coreApiWebdavMove2/moveShareOnOcis.feature:32](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L32)
- [coreApiWebdavMove2/moveShareOnOcis.feature:98](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L98)
- [coreApiWebdavMove2/moveShareOnOcis.feature:100](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L100)
- [coreApiWebdavMove2/moveShareOnOcis.feature:169](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L169)
- [coreApiWebdavMove2/moveShareOnOcis.feature:170](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveShareOnOcis.feature#L170)

#### [Expiration date for shares is not implemented](https://github.com/owncloud/ocis/issues/1250)

#### Expiration date of user shares

- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:33](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L33)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:34](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L34)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:85](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L85)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:86](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L86)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:142](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L142)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:143](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L143)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:200](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L200)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:201](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L201)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:202](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L202)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:203](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L203)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:286](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L286)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:287](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L287)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:317](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L317)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:318](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L318)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:319](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L319)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:320](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L320)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:378](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L378)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:379](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L379)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:380](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L380)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:381](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L381)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:382](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L382)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:383](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L383)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:412](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L412)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:413](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L413)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:414](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L414)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:415](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L415)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:443](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L443)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:444](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L444)

#### Expiration date of group shares

- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:59](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L59)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:60](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L60)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:115](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L115)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:116](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L116)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:171](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L171)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:172](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L172)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:231](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L231)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:232](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L232)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:233](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L233)
- [coreApiShareReshareToShares3/reShareWithExpiryDate.feature:234](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareReshareToShares3/reShareWithExpiryDate.feature#L234)

#### [Cannot move folder/file from one received share to another](https://github.com/owncloud/ocis/issues/2442)

- [coreApiShareUpdateToShares/updateShare.feature:195](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareUpdateToShares/updateShare.feature#L195)
- [coreApiShareUpdateToShares/updateShare.feature:159](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareUpdateToShares/updateShare.feature#L159)
- [coreApiShareManagementToShares/mergeShare.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/mergeShare.feature#L135)

#### [Sharing folder and sub-folder with same user but different permission,the permission of sub-folder is not obeyed ](https://github.com/owncloud/ocis/issues/2440)

- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:260](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L260)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:295](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L295)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:404](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L404)
- [coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature:439](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareReceivedInMultipleWays.feature#L439)

#### [Empty OCS response for a share create request using a disabled user](https://github.com/owncloud/ocis/issues/2212)

- [coreApiShareCreateSpecialToShares2/createShareWithDisabledUser.feature:19](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareWithDisabledUser.feature#L19)
- [coreApiShareCreateSpecialToShares2/createShareWithDisabledUser.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareCreateSpecialToShares2/createShareWithDisabledUser.feature#L22)


#### [Edit user share response has a "name" field](https://github.com/owncloud/ocis/issues/1225)

- [coreApiShareUpdateToShares/updateShare.feature:241](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareUpdateToShares/updateShare.feature#L241)
- [coreApiShareUpdateToShares/updateShare.feature:242](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareUpdateToShares/updateShare.feature#L242)

#### [user can access version metadata of a received share before accepting it](https://github.com/owncloud/ocis/issues/760)

- [coreApiVersions/fileVersions.feature:487](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiVersions/fileVersions.feature#L487)

#### [Share lists deleted user as 'user'](https://github.com/owncloud/ocis/issues/903)

- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:676](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L676)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:677](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L677)

#### [deleting a share with wrong authentication returns OCS status 996 / HTTP 500](https://github.com/owncloud/ocis/issues/1229)

- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:225](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L225)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:226](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L226)

### User Management

User and group management features

#### [incorrect ocs(v2) status value when getting info of share that does not exist should be 404, gives 998](https://github.com/owncloud/product/issues/250)

_ocs: api compatibility, return correct status code_

- [coreApiShareOperationsToShares2/shareAccessByID.feature:47](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L47)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:48](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L48)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:49](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L49)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:50](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L50)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:51](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L51)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:52](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L52)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:53](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L53)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:54](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L54)

### Other

API, search, favorites, config, capabilities, not existing endpoints, CORS and others

#### [Ability to return error messages in Webdav response bodies](https://github.com/owncloud/ocis/issues/1293)

- [coreApiAuthOcs/ocsDELETEAuth.feature:8](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsDELETEAuth.feature#L8)
- [coreApiAuthOcs/ocsGETAuth.feature:8](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsGETAuth.feature#L8)
- [coreApiAuthOcs/ocsGETAuth.feature:42](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsGETAuth.feature#L42)
- [coreApiAuthOcs/ocsGETAuth.feature:73](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsGETAuth.feature#L73)
- [coreApiAuthOcs/ocsGETAuth.feature:104](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsGETAuth.feature#L104)
- [coreApiAuthOcs/ocsGETAuth.feature:121](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsGETAuth.feature#L121)
- [coreApiAuthOcs/ocsPOSTAuth.feature:8](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsPOSTAuth.feature#L8)
- [coreApiAuthOcs/ocsPUTAuth.feature:8](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsPUTAuth.feature#L8)
- [coreApiSharePublicLink1/createPublicLinkShare.feature:343](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiSharePublicLink1/createPublicLinkShare.feature#L343)

#### [sending MKCOL requests to another user's webDav endpoints as normal user gives 404 instead of 403 ](https://github.com/owncloud/ocis/issues/3872)

_ocdav: api compatibility, return correct status code_

- [coreApiAuthWebDav/webDavMKCOLAuth.feature:52](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavMKCOLAuth.feature#L52)
- [coreApiAuthWebDav/webDavMKCOLAuth.feature:63](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavMKCOLAuth.feature#L63)

#### [trying to lock file of another user gives http 200](https://github.com/owncloud/ocis/issues/2176)

- [coreApiAuthWebDav/webDavLOCKAuth.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavLOCKAuth.feature#L56)
- [coreApiAuthWebDav/webDavLOCKAuth.feature:68](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavLOCKAuth.feature#L68)

#### [send (MOVE,COPY) requests to another user's webDav endpoints as normal user gives 400 instead of 403](https://github.com/owncloud/ocis/issues/3882)

_ocdav: api compatibility, return correct status code_

- [coreApiAuthWebDav/webDavMOVEAuth.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavMOVEAuth.feature#L55)
- [coreApiAuthWebDav/webDavMOVEAuth.feature:64](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavMOVEAuth.feature#L64)
- [coreApiAuthWebDav/webDavCOPYAuth.feature:59](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavCOPYAuth.feature#L59)
- [coreApiAuthWebDav/webDavCOPYAuth.feature:68](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavCOPYAuth.feature#L68)

#### [send POST requests to another user's webDav endpoints as normal user](https://github.com/owncloud/ocis/issues/1287)

_ocdav: api compatibility, return correct status code_

- [coreApiAuthWebDav/webDavPOSTAuth.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavPOSTAuth.feature#L56)
- [coreApiAuthWebDav/webDavPOSTAuth.feature:65](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavPOSTAuth.feature#L65)

#### Another users space literally does not exist because it is not listed as a space for him, 404 seems correct, expects 403

- [coreApiAuthWebDav/webDavPUTAuth.feature:56](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavPUTAuth.feature#L56)
- [coreApiAuthWebDav/webDavPUTAuth.feature:68](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavPUTAuth.feature#L68)

#### [Using double slash in URL to access a folder gives 501 and other status codes](https://github.com/owncloud/ocis/issues/1667)

- [coreApiAuthWebDav/webDavSpecialURLs.feature:13](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L13)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:24](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L24)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:55](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L55)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:66](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L66)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:76](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L76)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:88](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L88)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:100](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L100)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:111](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L111)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:121](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L121)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:132](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L132)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:142](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L142)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:153](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L153)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:163](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L163)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:174](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L174)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:184](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L184)
- [coreApiAuthWebDav/webDavSpecialURLs.feature:195](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavSpecialURLs.feature#L195)

#### [Difference in response content of status.php and default capabilities](https://github.com/owncloud/ocis/issues/1286)

- [coreApiCapabilities/capabilitiesWithNormalUser.feature:11](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiCapabilities/capabilitiesWithNormalUser.feature#L11)

#### [[old/new/spaces] In ocis and oc10, REPORT request response differently](https://github.com/owncloud/ocis/issues/4712)

- [coreApiWebdavOperations/search.feature:208](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L208)
- [coreApiWebdavOperations/search.feature:209](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L209)
- [coreApiWebdavOperations/search.feature:214](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L214)
- [coreApiWebdavOperations/search.feature:240](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L240)
- [coreApiWebdavOperations/search.feature:241](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L241)
- [coreApiWebdavOperations/search.feature:246](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/search.feature#L246)

#### [Support for favorites](https://github.com/owncloud/ocis/issues/1228)

- [coreApiFavorites/favorites.feature:115](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L115)
- [coreApiFavorites/favorites.feature:116](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L116)
- [coreApiFavorites/favorites.feature:142](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L142)
- [coreApiFavorites/favorites.feature:143](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L143)
- [coreApiFavorites/favorites.feature:219](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L219)
- [coreApiFavorites/favorites.feature:220](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L220)

And other missing implementation of favorites

- [coreApiFavorites/favorites.feature:167](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L167)
- [coreApiFavorites/favorites.feature:168](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L168)
- [coreApiFavorites/favorites.feature:173](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L173)
- [coreApiFavorites/favorites.feature:200](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L200)
- [coreApiFavorites/favorites.feature:201](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L201)
- [coreApiFavorites/favorites.feature:206](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L206)
- [coreApiFavorites/favoritesSharingToShares.feature:66](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favoritesSharingToShares.feature#L66)
- [coreApiFavorites/favoritesSharingToShares.feature:67](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favoritesSharingToShares.feature#L67)


#### [WWW-Authenticate header for unauthenticated requests is not clear](https://github.com/owncloud/ocis/issues/2285)

- [coreApiWebdavOperations/refuseAccess.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L22)
- [coreApiWebdavOperations/refuseAccess.feature:23](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L23)

#### [App Passwords/Tokens for legacy WebDAV clients](https://github.com/owncloud/ocis/issues/197)

- [coreApiAuthWebDav/webDavDELETEAuth.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavDELETEAuth.feature#L135)
- [coreApiAuthWebDav/webDavDELETEAuth.feature:149](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavDELETEAuth.feature#L149)
- [coreApiAuthWebDav/webDavDELETEAuth.feature:161](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavDELETEAuth.feature#L161)
- [coreApiAuthWebDav/webDavDELETEAuth.feature:175](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthWebDav/webDavDELETEAuth.feature#L175)

#### [Request to edit non-existing user by authorized admin gets unauthorized in http response](https://github.com/owncloud/core/issues/38423)

- [coreApiAuthOcs/ocsPUTAuth.feature:24](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiAuthOcs/ocsPUTAuth.feature#L24)

#### [Sharing a same file twice to the same group](https://github.com/owncloud/ocis/issues/1710)

- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:724](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L724)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:725](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L725)

#### [PATCH request for TUS upload with wrong checksum gives incorrect response](https://github.com/owncloud/ocis/issues/1755)

- [coreApiWebdavUploadTUS/checksums.feature:84](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L84)
- [coreApiWebdavUploadTUS/checksums.feature:85](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L85)
- [coreApiWebdavUploadTUS/checksums.feature:86](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L86)
- [coreApiWebdavUploadTUS/checksums.feature:87](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L87)
- [coreApiWebdavUploadTUS/checksums.feature:92](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L92)
- [coreApiWebdavUploadTUS/checksums.feature:93](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L93)
- [coreApiWebdavUploadTUS/checksums.feature:173](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L173)
- [coreApiWebdavUploadTUS/checksums.feature:174](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L174)
- [coreApiWebdavUploadTUS/checksums.feature:179](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L179)
- [coreApiWebdavUploadTUS/checksums.feature:226](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L226)
- [coreApiWebdavUploadTUS/checksums.feature:227](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L227)
- [coreApiWebdavUploadTUS/checksums.feature:228](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L228)
- [coreApiWebdavUploadTUS/checksums.feature:229](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L229)
- [coreApiWebdavUploadTUS/checksums.feature:234](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L234)
- [coreApiWebdavUploadTUS/checksums.feature:235](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L235)
- [coreApiWebdavUploadTUS/checksums.feature:282](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L282)
- [coreApiWebdavUploadTUS/checksums.feature:283](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L283)
- [coreApiWebdavUploadTUS/checksums.feature:284](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L284)
- [coreApiWebdavUploadTUS/checksums.feature:285](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L285)
- [coreApiWebdavUploadTUS/checksums.feature:290](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L290)
- [coreApiWebdavUploadTUS/checksums.feature:291](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/checksums.feature#L291)
- [coreApiWebdavUploadTUS/optionsRequest.feature:8](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L8)
- [coreApiWebdavUploadTUS/optionsRequest.feature:22](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L22)
- [coreApiWebdavUploadTUS/optionsRequest.feature:35](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L35)
- [coreApiWebdavUploadTUS/optionsRequest.feature:49](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L49)
- [coreApiWebdavUploadTUS/uploadToShare.feature:175](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L175)
- [coreApiWebdavUploadTUS/uploadToShare.feature:176](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L176)
- [coreApiWebdavUploadTUS/uploadToShare.feature:194](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L194)
- [coreApiWebdavUploadTUS/uploadToShare.feature:195](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L195)
- [coreApiWebdavUploadTUS/uploadToShare.feature:213](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L213)
- [coreApiWebdavUploadTUS/uploadToShare.feature:214](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L214)
- [coreApiWebdavUploadTUS/uploadToShare.feature:252](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L252)
- [coreApiWebdavUploadTUS/uploadToShare.feature:253](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L253)
- [coreApiWebdavUploadTUS/uploadToShare.feature:294](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L294)
- [coreApiWebdavUploadTUS/uploadToShare.feature:295](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/uploadToShare.feature#L295)

#### [TUS OPTIONS requests do not reply with TUS headers when invalid password](https://github.com/owncloud/ocis/issues/1012)

- [coreApiWebdavUploadTUS/optionsRequest.feature:62](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L62)
- [coreApiWebdavUploadTUS/optionsRequest.feature:76](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L76)
- [coreApiWebdavUploadTUS/optionsRequest.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L89)
- [coreApiWebdavUploadTUS/optionsRequest.feature:104](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUploadTUS/optionsRequest.feature#L104)

#### [Trying to accept a share with invalid ID gives incorrect OCS and HTTP status](https://github.com/owncloud/ocis/issues/2111)

- [coreApiShareOperationsToShares2/shareAccessByID.feature:84](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L84)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:85](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L85)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:86](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L86)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:87](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L87)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:88](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L88)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:89](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L89)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:90](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L90)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:91](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L91)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:103](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L103)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:104](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L104)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:135](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L135)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:136](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L136)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:137](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L137)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:138](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L138)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:139](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L139)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:140](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L140)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:141](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L141)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:142](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L142)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:154](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L154)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:155](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L155)


#### [Shares to deleted group listed in the response](https://github.com/owncloud/ocis/issues/2441)

- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:528](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L528)
- [coreApiShareManagementBasicToShares/createShareToSharesFolder.feature:529](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/createShareToSharesFolder.feature#L529)

#### [copying the file inside Shares folder returns 404](https://github.com/owncloud/ocis/issues/3874)

- [coreApiWebdavProperties1/copyFile.feature:409](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L409)
- [coreApiWebdavProperties1/copyFile.feature:410](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L410)
- [coreApiWebdavProperties1/copyFile.feature:415](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L415)
- [coreApiWebdavProperties1/copyFile.feature:435](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L435)
- [coreApiWebdavProperties1/copyFile.feature:436](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L436)
- [coreApiWebdavProperties1/copyFile.feature:441](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L441)

### Won't fix

Not everything needs to be implemented for ocis. While the oc10 testsuite covers these things we are not looking at them right now.

- _The `OC-LazyOps` header is [no longer supported by the client](https://github.com/owncloud/client/pull/8398), implementing this is not necessary for a first production release. We plan to have an upload state machine to visualize the state of a file, see https://github.com/owncloud/ocis/issues/214_
- _Blacklisted ignored files are no longer required because ocis can handle `.htaccess` files without security implications introduced by serving user provided files with apache._

#### [Blacklist files extensions](https://github.com/owncloud/ocis/issues/2177)

- [coreApiWebdavProperties1/copyFile.feature:119](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L119)
- [coreApiWebdavProperties1/copyFile.feature:120](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L120)
- [coreApiWebdavProperties1/copyFile.feature:125](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/copyFile.feature#L125)
- [coreApiWebdavProperties1/createFileFolder.feature:98](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/createFileFolder.feature#L98)
- [coreApiWebdavProperties1/createFileFolder.feature:99](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/createFileFolder.feature#L99)
- [coreApiWebdavProperties1/createFileFolder.feature:104](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavProperties1/createFileFolder.feature#L104)
- [coreApiWebdavUpload1/uploadFile.feature:181](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUpload1/uploadFile.feature#L181)
- [coreApiWebdavUpload1/uploadFile.feature:182](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUpload1/uploadFile.feature#L182)
- [coreApiWebdavUpload1/uploadFile.feature:187](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavUpload1/uploadFile.feature#L187)
- [coreApiWebdavMove2/moveFile.feature:176](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L176)
- [coreApiWebdavMove2/moveFile.feature:177](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L177)
- [coreApiWebdavMove2/moveFile.feature:182](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L182)

#### [Preview of text file with UTF content does not render correctly](https://github.com/owncloud/ocis/issues/2570)

- [coreApiWebdavPreviews/previews.feature:136](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavPreviews/previews.feature#L136)

#### [Share path in the response is different between share states](https://github.com/owncloud/ocis/issues/2540)

- [coreApiShareManagementToShares/acceptShares.feature:64](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L64)
- [coreApiShareManagementToShares/acceptShares.feature:87](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L87)
- [coreApiShareManagementToShares/acceptShares.feature:165](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L165)
- [coreApiShareManagementToShares/acceptShares.feature:188](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L188)
- [coreApiShareManagementToShares/acceptShares.feature:226](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L226)
- [coreApiShareManagementToShares/acceptShares.feature:260](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L260)
- [coreApiShareManagementToShares/acceptShares.feature:480](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/acceptShares.feature#L480)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:122](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L122)
- [coreApiShareOperationsToShares2/shareAccessByID.feature:123](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareOperationsToShares2/shareAccessByID.feature#L123)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:176](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L176)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:177](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L177)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:178](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L178)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:179](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L179)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:195](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L195)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:196](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L196)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:197](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L197)
- [coreApiShareManagementBasicToShares/deleteShareFromShares.feature:198](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementBasicToShares/deleteShareFromShares.feature#L198)

#### [Content-type is not multipart/byteranges when downloading file with Range Header](https://github.com/owncloud/ocis/issues/2677)

- [coreApiWebdavOperations/downloadFile.feature:184](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/downloadFile.feature#L184)
- [coreApiWebdavOperations/downloadFile.feature:185](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/downloadFile.feature#L185)
- [coreApiWebdavOperations/downloadFile.feature:190](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/downloadFile.feature#L190)

#### [moveShareInsideAnotherShare behaves differently on oCIS than oC10](https://github.com/owncloud/ocis/issues/3047)

- [coreApiShareManagementToShares/moveShareInsideAnotherShare.feature:24](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/moveShareInsideAnotherShare.feature#L24)
- [coreApiShareManagementToShares/moveShareInsideAnotherShare.feature:44](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/moveShareInsideAnotherShare.feature#L44)
- [coreApiShareManagementToShares/moveShareInsideAnotherShare.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/moveShareInsideAnotherShare.feature#L58)

#### [Renaming resource to banned name is allowed in spaces webdav](https://github.com/owncloud/ocis/issues/3099)

- [coreApiWebdavMove1/moveFolder.feature:27](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L27)
- [coreApiWebdavMove1/moveFolder.feature:45](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L45)
- [coreApiWebdavMove1/moveFolder.feature:63](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove1/moveFolder.feature#L63)
- [coreApiWebdavMove2/moveFile.feature:138](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavMove2/moveFile.feature#L138)

#### [REPORT method on spaces returns an incorrect d:href response](https://github.com/owncloud/ocis/issues/3111)

- [coreApiFavorites/favorites.feature:121](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L121)
- [coreApiFavorites/favorites.feature:148](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L148)
- [coreApiFavorites/favorites.feature:225](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiFavorites/favorites.feature#L225)

#### [Cannot disable the dav propfind depth infinity for resources](https://github.com/owncloud/ocis/issues/3720)

- [coreApiWebdavOperations/listFiles.feature:364](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/listFiles.feature#L364)
- [coreApiWebdavOperations/listFiles.feature:365](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/listFiles.feature#L365)
- [coreApiWebdavOperations/listFiles.feature:370](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/listFiles.feature#L370)
- [coreApiWebdavOperations/listFiles.feature:384](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/listFiles.feature#L384)
- [coreApiWebdavOperations/listFiles.feature:389](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/listFiles.feature#L389)
- [coreApiWebdavOperations/listFiles.feature:403](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/listFiles.feature#L403)
- [coreApiWebdavOperations/listFiles.feature:404](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/listFiles.feature#L404)
- [coreApiWebdavOperations/listFiles.feature:409](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/listFiles.feature#L409)

### [graph/users: enable/disable users](https://github.com/owncloud/ocis/issues/3064)

- [coreApiWebdavOperations/refuseAccess.feature:35](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L35)
- [coreApiWebdavOperations/refuseAccess.feature:36](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L36)
- [coreApiWebdavOperations/refuseAccess.feature:41](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiWebdavOperations/refuseAccess.feature#L41)

#### [OCS status code zero](https://github.com/owncloud/ocis/issues/3621)

- [coreApiShareManagementToShares/moveReceivedShare.feature:14](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiShareManagementToShares/moveReceivedShare.feature#L14)

#### [HTTP status code differ while deleting file of another user's trash bin](https://github.com/owncloud/ocis/issues/3544)

- [coreApiTrashbin/trashbinDelete.feature:105](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L105)

#### [Problem accessing trashbin with personal space id](https://github.com/owncloud/ocis/issues/3639)

- [coreApiTrashbin/trashbinDelete.feature:35](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L35)
- [coreApiTrashbin/trashbinDelete.feature:36](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L36)
- [coreApiTrashbin/trashbinDelete.feature:58](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L58)
- [coreApiTrashbin/trashbinDelete.feature:85](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L85)
- [coreApiTrashbin/trashbinDelete.feature:128](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L128)
- [coreApiTrashbin/trashbinDelete.feature:151](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L151)
- [coreApiTrashbin/trashbinDelete.feature:176](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L176)
- [coreApiTrashbin/trashbinDelete.feature:201](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L201)
- [coreApiTrashbin/trashbinDelete.feature:238](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L238)
- [coreApiTrashbin/trashbinDelete.feature:275](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L275)
- [coreApiTrashbin/trashbinDelete.feature:324](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinDelete.feature#L324)
- [coreApiTrashbin/trashbinFilesFolders.feature:25](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L25)
- [coreApiTrashbin/trashbinFilesFolders.feature:41](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L41)
- [coreApiTrashbin/trashbinFilesFolders.feature:60](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L60)
- [coreApiTrashbin/trashbinFilesFolders.feature:81](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L81)
- [coreApiTrashbin/trashbinFilesFolders.feature:100](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L100)
- [coreApiTrashbin/trashbinFilesFolders.feature:136](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L136)
- [coreApiTrashbin/trashbinFilesFolders.feature:159](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L159)
- [coreApiTrashbin/trashbinFilesFolders.feature:322](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L322)
- [coreApiTrashbin/trashbinFilesFolders.feature:323](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L323)
- [coreApiTrashbin/trashbinFilesFolders.feature:324](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L324)
- [coreApiTrashbin/trashbinFilesFolders.feature:329](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L329)
- [coreApiTrashbin/trashbinFilesFolders.feature:330](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L330)
- [coreApiTrashbin/trashbinFilesFolders.feature:331](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L331)
- [coreApiTrashbin/trashbinFilesFolders.feature:332](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L332)
- [coreApiTrashbin/trashbinFilesFolders.feature:333](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L333)
- [coreApiTrashbin/trashbinFilesFolders.feature:334](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L334)
- [coreApiTrashbin/trashbinFilesFolders.feature:351](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L351)
- [coreApiTrashbin/trashbinFilesFolders.feature:371](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L371)
- [coreApiTrashbin/trashbinFilesFolders.feature:425](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L425)
- [coreApiTrashbin/trashbinFilesFolders.feature:462](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiTrashbin/trashbinFilesFolders.feature#L462)

#### [Default capabilities for normal user and admin user not same as in oC-core](https://github.com/owncloud/ocis/issues/1285)
- [coreApiCapabilities/capabilities.feature:8](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiCapabilities/capabilities.feature#L8)
- [coreApiCapabilities/capabilities.feature:133](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiCapabilities/capabilities.feature#L133)
- [coreApiCapabilities/capabilities.feature:172](https://github.com/owncloud/ocis/blob/master/tests/acceptance/features/coreApiCapabilities/capabilities.feature#L172)
Note: always have an empty line at the end of this file.
The bash script that processes this file requires that the last line has a newline on the end.
